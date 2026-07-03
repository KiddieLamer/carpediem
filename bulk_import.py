#!/usr/bin/env python3
"""Bulk import auth.json files into 9router."""

import json, glob, hashlib, os, re
from urllib.request import Request, urlopen
from urllib.error import URLError

DATADIR = os.path.expanduser("~/.9router")
API_URL = "http://localhost:20128/api/oauth/codex/import-token"

def get_cli_token():
    machine_id = open(os.path.join(DATADIR, "machine-id")).read().strip()
    cli_secret = open(os.path.join(DATADIR, "auth", "cli-secret")).read().strip()
    h = hashlib.sha256(f"{machine_id}9r-cli-auth{cli_secret}".encode()).hexdigest()
    return h[:16]

def import_token(access_token, name):
    data = json.dumps({"accessToken": access_token, "name": name}).encode()
    req = Request(API_URL, data=data, headers={
        "Content-Type": "application/json",
        "x-9r-cli-token": get_cli_token(),
    })
    try:
        resp = urlopen(req)
        result = json.loads(resp.read())
        if result.get("success"):
            print(f"  OK  {name}")
        else:
            print(f"  FAIL {name}: {result}")
    except URLError as e:
        print(f"  ERROR {name}: {e}")

def load_auth(path):
    with open(path) as f:
        d = json.load(f)
    token = d.get("tokens", {}).get("access_token")
    if not token:
        token = d.get("accessToken")
    return token

def extract_email(access_token):
    parts = access_token.split(".")
    if len(parts) == 3:
        try:
            payload = parts[1]
            padding = 4 - len(payload) % 4
            if padding != 4:
                payload += "=" * padding
            import base64
            data = json.loads(base64.urlsafe_b64decode(payload))
            email = (data.get("https://api.openai.com/profile") or {}).get("email") or data.get("email")
            if email:
                return email
        except:
            pass
    return None

def main():
    files = sorted(glob.glob(os.path.expanduser("auth*.json")))
    if not files:
        print("Tidak ada file auth*.json di direktori ini.")
        return

    print(f"Mengimport {len(files)} file...\n")
    for path in files:
        name = os.path.basename(path)
        token = load_auth(path)
        if not token:
            print(f"  SKIP {name}: access_token tidak ditemukan")
            continue
        email = extract_email(token) or name.replace(".json", "")
        import_token(token, email)
    print("\nSelesai.")

if __name__ == "__main__":
    main()
