#!/usr/bin/env bash
set -euo pipefail

# get_access_token.sh
#
# Simple helper to perform an OAuth2 Authorization Code Flow with PKCE from the shell.
# - Opens the authorization URL in the default browser
# - Starts a minimal local HTTP listener to capture the redirect with the authorization code
# - Exchanges the code for tokens and prints the access_token (and full JSON response)
#
# Configurable via environment variables (or override on the command line):
#   OAUTH2_AUTH_URL - Authorization endpoint (required)
#   OAUTH2_TOKEN_URL - Token endpoint (required)
#   OAUTH2_CLIENT_ID - Client ID (required)
#   OAUTH2_REDIRECT_URI - Redirect URI (default: http://127.0.0.1:8080/callback)
#   OAUTH2_SCOPE - space-separated scopes (default: openid profile email)
#   OAUTH2_STATE - optional state value (auto-generated if empty)
#   OUTPUT_FILE - optional file to save JSON token response
#
# Requirements: curl, openssl (or head/xxd), python3 (to run a tiny http server and url parsing), xdg-open/open

DEFAULT_REDIRECT_URI="http://127.0.0.1:8080/callback"
DEFAULT_SCOPE="openid profile email"

OAUTH2_AUTH_URL="${OAUTH2_AUTH_URL:-}"
OAUTH2_TOKEN_URL="${OAUTH2_TOKEN_URL:-}"
OAUTH2_CLIENT_ID="${OAUTH2_CLIENT_ID:-}"
OAUTH2_REDIRECT_URI="${OAUTH2_REDIRECT_URI:-$DEFAULT_REDIRECT_URI}"
OAUTH2_SCOPE="${OAUTH2_SCOPE:-$DEFAULT_SCOPE}"
OAUTH2_STATE="${OAUTH2_STATE:-}"
OUTPUT_FILE="${OUTPUT_FILE:-}"

usage(){
	cat <<EOF
Usage: ${0##*/} [--auth-url URL] [--token-url URL] [--client-id ID] [--redirect URI] [--scope SCOPE] [--output FILE]

Environment variables can also be used: OAUTH2_AUTH_URL, OAUTH2_TOKEN_URL, OAUTH2_CLIENT_ID, OAUTH2_REDIRECT_URI, OAUTH2_SCOPE, OUTPUT_FILE
EOF
	exit 1
}

while [[ $# -gt 0 ]]; do
	case "$1" in
		--auth-url) OAUTH2_AUTH_URL="$2"; shift 2;;
		--token-url) OAUTH2_TOKEN_URL="$2"; shift 2;;
		--client-id) OAUTH2_CLIENT_ID="$2"; shift 2;;
		--redirect) OAUTH2_REDIRECT_URI="$2"; shift 2;;
		--scope) OAUTH2_SCOPE="$2"; shift 2;;
		--output) OUTPUT_FILE="$2"; shift 2;;
		-h|--help) usage;;
		*) echo "Unknown arg: $1"; usage;;
	esac
done

if [[ -z "$OAUTH2_AUTH_URL" || -z "$OAUTH2_TOKEN_URL" || -z "$OAUTH2_CLIENT_ID" ]]; then
	echo "ERROR: --auth-url, --token-url and --client-id are required (or set OAUTH2_AUTH_URL/OAUTH2_TOKEN_URL/OAUTH2_CLIENT_ID)" >&2
	usage
fi

# generate code_verifier and code_challenge (PKCE)
code_verifier=$(python3 - <<'PY'
import secrets,base64
v = secrets.token_urlsafe(64)
print(v)
PY
)

code_challenge=$(python3 - <<'PY'
import hashlib,base64,sys
v = sys.stdin.read().strip().encode()
ch = hashlib.sha256(v).digest()
print(base64.urlsafe_b64encode(ch).rstrip(b"=").decode())
PY <<EOF
$code_verifier
EOF
)

if [[ -z "$OAUTH2_STATE" ]]; then
	OAUTH2_STATE=$(python3 - <<'PY'
import secrets
print(secrets.token_urlsafe(16))
PY
)
fi

# Build authorization URL
auth_url="${OAUTH2_AUTH_URL}?response_type=code&client_id=$(python3 - <<'PY'
import urllib.parse,sys
print(urllib.parse.quote(sys.argv[1], safe=''))
PY "$OAUTH2_CLIENT_ID")&redirect_uri=$(python3 - <<'PY'
import urllib.parse,sys
print(urllib.parse.quote(sys.argv[1], safe=''))
PY "$OAUTH2_REDIRECT_URI")&scope=$(python3 - <<'PY'
import urllib.parse,sys
print(urllib.parse.quote(sys.argv[1], safe=''))
PY "$OAUTH2_SCOPE")&state=$(python3 - <<'PY'
import urllib.parse,sys
print(urllib.parse.quote(sys.argv[1], safe=''))
PY "$OAUTH2_STATE")&code_challenge=$code_challenge&code_challenge_method=S256"

echo "Opening browser to authorize..."

# Try to open browser; fallback to printing URL
if command -v xdg-open >/dev/null 2>&1; then
	xdg-open "$auth_url" >/dev/null 2>&1 || echo "Open this URL in your browser:\n$auth_url"
elif command -v open >/dev/null 2>&1; then
	open "$auth_url" || echo "Open this URL in your browser:\n$auth_url"
else
	echo "Open this URL in your browser:\n$auth_url"
fi

echo "Waiting for authorization response on $OAUTH2_REDIRECT_URI ..."

# Start a tiny Python HTTP server to capture the redirect and extract code
tmpfile=$(mktemp)
python3 - <<'PY' > "$tmpfile" 2>&1 &
import http.server, socketserver, urllib.parse, sys

PORT = int(urllib.parse.urlparse('%s').port or 8080)
PATH = urllib.parse.urlparse('%s').path

class Handler(http.server.BaseHTTPRequestHandler):
		def do_GET(self):
				qs = urllib.parse.urlparse(self.path).query
				params = urllib.parse.parse_qs(qs)
				if 'code' in params:
						code = params['code'][0]
						state = params.get('state', [''])[0]
						# Print to stdout so parent shell can read it via tmp file
						print('CODE=' + code)
						print('STATE=' + state)
						self.send_response(200)
						self.send_header('Content-type','text/html')
						self.end_headers()
						self.wfile.write(b'<html><body><h1>Authorization received. You can close this window.</h1></body></html>')
						# shut down server after handling
						def shutdown(server):
								server.shutdown()
						import threading
						threading.Thread(target=shutdown, args=(self.server,)).start()
				else:
						self.send_response(400)
						self.end_headers()

with socketserver.TCPServer(('127.0.0.1', PORT), Handler) as httpd:
		httpd.serve_forever()
PY

# Wait for the python server to print the captured code (it writes to tmpfile)
code=""
state=""
for i in {1..300}; do
	if grep -q "CODE=" "$tmpfile"; then
		code=$(grep "CODE=" "$tmpfile" | head -n1 | sed 's/CODE=//')
		state=$(grep "STATE=" "$tmpfile" | head -n1 | sed 's/STATE=//')
		break
	fi
	sleep 0.1
done

if [[ -z "$code" ]]; then
	echo "Failed to obtain authorization code (timeout)" >&2
	cat "$tmpfile" >&2 || true
	rm -f "$tmpfile"
	exit 2
fi

echo "Received code: $code"
if [[ -n "$state" ]]; then
	echo "Received state: $state"
fi

# Exchange code for token
echo "Exchanging code for token..."

token_response=$(curl -s -X POST "$OAUTH2_TOKEN_URL" \
	-H "Content-Type: application/x-www-form-urlencoded" \
	-d "grant_type=authorization_code" \
	-d "code=$code" \
	-d "redirect_uri=$OAUTH2_REDIRECT_URI" \
	-d "client_id=$OAUTH2_CLIENT_ID" \
	-d "code_verifier=$code_verifier")

if [[ -z "$token_response" ]]; then
	echo "Empty response from token endpoint" >&2
	rm -f "$tmpfile"
	exit 3
fi

echo "$token_response" | jq '.' 2>/dev/null || echo "$token_response"

if [[ -n "$OUTPUT_FILE" ]]; then
	echo "$token_response" > "$OUTPUT_FILE"
	echo "Saved token response to $OUTPUT_FILE"
fi

# Print access_token if present
access_token=$(echo "$token_response" | python3 -c "import sys, json
try:
	print(json.load(sys.stdin).get('access_token',''))
except Exception:
	sys.exit(0)
")

if [[ -n "$access_token" ]]; then
	echo "\nAccess token:\n$access_token"
fi

rm -f "$tmpfile"

exit 0
