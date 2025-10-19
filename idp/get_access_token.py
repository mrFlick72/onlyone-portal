#!/usr/bin/env python3
"""
get_access_token.py

Performs OAuth2 Authorization Code Flow with PKCE for local development.

Usage:
  python idp/get_access_token.py --auth-url AUTH_URL --token-url TOKEN_URL --client-id CLIENT_ID

Environment variables also supported: OAUTH2_AUTH_URL, OAUTH2_TOKEN_URL, OAUTH2_CLIENT_ID, OAUTH2_REDIRECT_URL, OAUTH2_SCOPE

The script opens the system browser to the authorization URL, starts a local HTTP server to capture the
authorization code at the registered redirect URI, exchanges the code for tokens, prints the JSON response and
optionally writes it to a file.

Requires: requests (pip install requests)
"""

import argparse
import base64
import hashlib
import json
import os
import secrets
import sys
import threading
import urllib.parse
import webbrowser
from http.server import BaseHTTPRequestHandler, HTTPServer
import requests
from dotenv import load_dotenv

load_dotenv(dotenv_path=".local_env")


DEFAULT_SCOPE = "openid profile email"


def generate_pkce_pair():
    verifier = secrets.token_urlsafe(64)
    m = hashlib.sha256()
    m.update(verifier.encode("utf-8"))
    challenge = base64.urlsafe_b64encode(m.digest()).rstrip(b"=").decode("utf-8")
    return verifier, challenge


class CodeReceiver(BaseHTTPRequestHandler):
    server_version = "CodeReceiver/0.1"

    def do_GET(self):
        parsed = urllib.parse.urlparse(self.path)
        qs = urllib.parse.parse_qs(parsed.query)
        if "code" in qs:
            self.server.code = qs["code"][0]
            self.server.state = qs.get("state", [""])[0]
            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.end_headers()
            self.wfile.write(
                b"<html><body><h1>Authorization received. You can close this window.</h1></body></html>"
            )
            # shutdown server in separate thread to avoid blocking
            threading.Thread(target=self.server.shutdown).start()
        else:
            self.send_response(400)
            self.end_headers()


def find_port_from_redirect(uri: str) -> int:
    parsed = urllib.parse.urlparse(uri)
    port = parsed.port
    if port:
        return port
    # default port for http
    return 80 if parsed.scheme == "http" else 443


def main(argv=None):
    AUTH_URL = os.getenv("AUTH_URL")
    TOKEN_URL = os.getenv("TOKEN_URL")
    CLIENT_ID = os.getenv("CLIENT_ID")
    REDIRECT_URL = os.getenv("REDIRECT_URL")

    print("Using Authorization URL:", AUTH_URL)
    print("Using Token URL:", TOKEN_URL)
    print("Using Client ID:", CLIENT_ID)
    print("Using Redirect URI:", REDIRECT_URL)
    verifier, challenge = generate_pkce_pair()

    state = secrets.token_urlsafe(16)

    params = {
        "response_type": "code",
        "client_id": CLIENT_ID,
        "redirect_uri": REDIRECT_URL,
        "scope": DEFAULT_SCOPE,
        "state": state,
        "code_challenge": challenge,
        "code_challenge_method": "S256",
    }

    auth_url = AUTH_URL + ("?" + urllib.parse.urlencode(params))
    print("Opening browser to:", auth_url)

    webbrowser.open(auth_url)

    parsed = urllib.parse.urlparse(REDIRECT_URL)
    host = parsed.hostname
    port = parsed.port
    path = parsed.path

    print("parsed", parsed)
    server = HTTPServer((host, port), CodeReceiver)

    print(f"Listening for redirect on {host}:{port}{path} ...")
    try:
        server.handle_request()  # handle a single request and return
    except Exception as exc:
        print("Server error:", exc, file=sys.stderr)
        sys.exit(2)

    code = getattr(server, "code", None)
    rec_state = getattr(server, "state", None)

    if not code:
        print("Did not receive authorization code", file=sys.stderr)
        sys.exit(3)

    print("Received code:", code)
    if rec_state:
        print("State:", rec_state)

    # Exchange code for token
    data = {
        "grant_type": "authorization_code",
        "code": code,
        "redirect_uri": REDIRECT_URL,
        "client_id": CLIENT_ID,
        "code_verifier": verifier,
    }

    print("Exchanging code for token...")
    resp = requests.post(TOKEN_URL, data=data, headers={"Accept": "application/json"})

    try:
        token_json = resp.json()
    except Exception:
        print("Failed to parse token response:", resp.text, file=sys.stderr)
        resp.raise_for_status()

    print(json.dumps(token_json, indent=2))

    access_token = token_json.get("access_token")
    if access_token:
        print("\nAccess token:\n", access_token)


if __name__ == "__main__":
    main()
