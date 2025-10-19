#!/usr/bin/env python3
"""
get_access_token.py

Small utility implementing OAuth2 Authorization Code Flow with PKCE for local development.

Features:
- Loads `.local_env` (dotenv) and expands variable references like AUTH_URL=$IDP_BASE_URL/...
- Generates PKCE verifier/challenge
- Attempts to open a private/incognito browser window (Chrome/Chromium/Firefox)
- Attempts to bind a local HTTP listener to capture the authorization code
- Falls back to manual paste mode if binding fails or when --redirected-url is provided

Dependencies:
  pip install requests python-dotenv

Usage:
  python idp/get_access_token.py --auth-url AUTH_URL --token-url TOKEN_URL --client-id CLIENT_ID
  python idp/get_access_token.py --env-file idp/.local_env

"""

from __future__ import annotations

import base64
import hashlib
import json
import os
import secrets
import shutil
import subprocess
import sys
import threading
import urllib.parse
import webbrowser
from http.server import BaseHTTPRequestHandler, HTTPServer
from typing import Optional

import requests
from dotenv import dotenv_values


DEFAULT_REDIRECT = "http://127.0.0.1:8080/callback"
DEFAULT_SCOPE = "openid profile email"


KNOWN_BROWSERS = [
    ("google-chrome", ["--incognito"]),
    ("google-chrome-stable", ["--incognito"]),
    ("chromium", ["--incognito"]),
    ("chromium-browser", ["--incognito"]),
    ("firefox", ["-private-window"]),
]


def expand_env_file(path: str) -> dict:
    raw = dotenv_values(path)
    resolved = dict(raw)
    # Expand using current environment plus resolved values
    for k, v in list(raw.items()):
        if v is None:
            continue
        resolved[k] = os.path.expandvars(v)
    return resolved


def generate_pkce_pair() -> tuple[str, str]:
    verifier = secrets.token_urlsafe(64)
    digest = hashlib.sha256(verifier.encode("utf-8")).digest()
    challenge = base64.urlsafe_b64encode(digest).rstrip(b"=").decode("utf-8")
    return verifier, challenge


class CodeReceiver(BaseHTTPRequestHandler):
    server_version = "CodeReceiver/1.0"

    def do_GET(self):
        parsed = urllib.parse.urlparse(self.path)
        qs = urllib.parse.parse_qs(parsed.query)
        if "code" in qs:
            self.server.code = qs["code"][0]
            self.server.state = qs.get("state", [""])[0]
            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.end_headers()
            self.wfile.write(b"<html><body><h1>Authorization received. You may close this window.</h1></body></html>")
            threading.Thread(target=self.server.shutdown).start()
        else:
            self.send_response(400)
            self.end_headers()


def open_in_incognito(url: str, preferred_bin: Optional[str] = None) -> None:
    candidates = []
    if preferred_bin:
        candidates.append((preferred_bin, ["--incognito"]))
    candidates.extend(KNOWN_BROWSERS)

    for bin_name, flags in candidates:
        path = shutil.which(bin_name)
        if path:
            try:
                subprocess.Popen([path, *flags, url], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, close_fds=True)
                return
            except Exception:
                continue

    webbrowser.open(url)


def extract_code_from_url(redirected_url: str) -> tuple[Optional[str], Optional[str]]:
    parsed = urllib.parse.urlparse(redirected_url)
    qs = urllib.parse.parse_qs(parsed.query)
    return qs.get("code", [None])[0], qs.get("state", [None])[0]


def main(argv: Optional[list[str]] = None) -> int:
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("--auth-url", default=os.environ.get("AUTH_URL"))
    parser.add_argument("--token-url", default=os.environ.get("TOKEN_URL"))
    parser.add_argument("--client-id", default=os.environ.get("CLIENT_ID"))
    parser.add_argument("--redirect-uri", default=os.environ.get("REDIRECT_URL", DEFAULT_REDIRECT))
    parser.add_argument("--scope", default=os.environ.get("OAUTH2_SCOPE", DEFAULT_SCOPE))
    parser.add_argument("--no-browser", action="store_true")
    parser.add_argument("--redirected-url", help="Provide redirected URL (non-interactive)")
    parser.add_argument("--output", help="Save token JSON to file")
    parser.add_argument("--env-file", default=".local_env")
    parser.add_argument("--browser-bin", help="Preferred browser binary name")
    args = parser.parse_args(argv)

    # Load env file
    if args.env_file and os.path.exists(args.env_file):
        resolved = expand_env_file(args.env_file)
        os.environ.update({k: v for k, v in resolved.items() if v is not None})

    auth_url = args.auth_url
    token_url = args.token_url
    client_id = args.client_id
    redirect_uri = args.redirect_uri
    scope = args.scope

    if not auth_url or not token_url or not client_id:
        parser.print_help()
        print("\nRequired: --auth-url, --token-url, --client-id (or set env in .local_env)")
        return 1

    verifier, challenge = generate_pkce_pair()
    state = secrets.token_urlsafe(16)

    params = {
        "response_type": "code",
        "client_id": client_id,
        "redirect_uri": redirect_uri,
        "scope": scope,
        "state": state,
        "code_challenge": challenge,
        "code_challenge_method": "S256",
    }

    auth_full = auth_url + ("?" + urllib.parse.urlencode(params))
    print("Authorization URL:\n", auth_full)

    if not args.no_browser and not args.redirected_url:
        open_in_incognito(auth_full, args.browser_bin)

    code = None
    rec_state = None

    if args.redirected_url:
        code, rec_state = extract_code_from_url(args.redirected_url)
    else:
        parsed = urllib.parse.urlparse(redirect_uri)
        host = parsed.hostname or "127.0.0.1"
        port = parsed.port or 8080
        path = parsed.path or "/"

        try:
            server = HTTPServer((host, port), CodeReceiver)
            print(f"Listening for redirect on {host}:{port}{path} ...")
            server.handle_request()
            code = getattr(server, "code", None)
            rec_state = getattr(server, "state", None)
        except OSError as e:
            print(f"Could not bind to {host}:{port} ({e}). Falling back to manual paste mode.")
            print("Authenticate in the browser, then paste the final redirect URL here.")
            if not args.no_browser:
                open_in_incognito(auth_full, args.browser_bin)
            redirected = input("Paste redirect URL: ").strip()
            code, rec_state = extract_code_from_url(redirected)

    if not code:
        print("Failed to obtain authorization code.", file=sys.stderr)
        return 2

    print("Received code:", code)

    data = {
        "grant_type": "authorization_code",
        "code": code,
        "redirect_uri": redirect_uri,
        "client_id": client_id,
        "code_verifier": verifier,
    }

    print("Exchanging code for token at:", token_url)
    resp = requests.post(token_url, data=data, headers={"Accept": "application/json"})

    try:
        token_json = resp.json()
    except Exception:
        print("Failed to parse token response:", resp.text, file=sys.stderr)
        resp.raise_for_status()

    print(json.dumps(token_json, indent=2))

    if args.output:
        with open(args.output, "w") as f:
            json.dump(token_json, f)
        print("Saved token response to", args.output)

    access_token = token_json.get("access_token")
    if access_token:
        print("\nAccess token:\n", access_token)

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
