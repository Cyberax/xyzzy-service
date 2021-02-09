#!/usr/bin/env python3
from pathlib import Path
import os, sys

if len(sys.argv) != 2:
    print("Usage: renamer.py <service-name>")
    sys.exit(1)

srv_name = sys.argv[1]

capitalized = srv_name.capitalize().encode('utf-8')
lower = srv_name.lower().encode('utf-8')
upper = srv_name.upper().encode('utf-8')
print(f"Creating a {capitalized.decode('utf-8')} service!")

for path in Path('.').rglob('*xyzzy*'):
    src = str(path)
    dst = src.replace('xyzzy', srv_name)
    os.rename(src, dst)

skipPaths = ['.git/', 'renamer.py']

for fileName in Path('.').rglob('*'):
    skipIt = False
    for p in skipPaths:
        if p in str(fileName):
            skipIt = True
            break
    if fileName.is_dir() or skipIt:
        continue
    with open(fileName, "rb") as fl:
        content = fl.read()
        content = content.replace(b'xyzzy', lower)
        content = content.replace(b'Xyzzy', capitalized)
        content = content.replace(b'XYZZY', upper)
    with open(fileName, "wb") as fl:
        fl.write(content)

print("Service created. You can delete the renamer.py and commit the service into its own repo.")
