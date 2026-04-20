#!/usr/bin/env python3
"""
quick_install.py -- Batch install APKs from a folder or URL
Usage: python3 quick_install.py ~/apks/
       python3 quick_install.py --url https://example.com/app.apk
"""
import subprocess, os, sys, argparse
from pathlib import Path

def adb_install(apk_path):
    r = subprocess.run(f"adb install -r {apk_path}", shell=True, capture_output=True, text=True)
    return "Success" in r.stdout

parser = argparse.ArgumentParser()
parser.add_argument("path", nargs="?")
parser.add_argument("--url", help="Install from URL")
args = parser.parse_args()

if args.url:
    name = args.url.split("/")[-1]
    print(f"Downloading {name}...")
    subprocess.run(f"wget {args.url} -O /tmp/{name}", shell=True)
    apk_install(f"/tmp/{name}")
elif args.path:
    folder = Path(args.path)
    apks = list(folder.glob("*.apk"))
    print(f"Found {len(apks)} APKs\n")
    for apk in apks:
        print(f"Installing {apk.name}...", end=" ")
        if adb_install(str(apk)):
            print("✓")
        else:
            print("✗")
else:
    print("Usage: python3 quick_install.py <folder>")
