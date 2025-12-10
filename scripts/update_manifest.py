import json
import os
import hashlib
import sys
import re

def compute_sha256(filepath):
    sha256_hash = hashlib.sha256()
    with open(filepath, "rb") as f:
        for byte_block in iter(lambda: f.read(4096), b""):
            sha256_hash.update(byte_block)
    return sha256_hash.hexdigest()

def get_platform_key(filename):
    # Expected format: btxz-{OS}-{ARCH}[.exe]
    # Regex to handle optional .exe and extract OS-ARCH
    match = re.match(r"btxz-([a-z0-9]+-[a-z0-9]+)(?:\.exe)?$", filename)
    if match:
        return match.group(1)
    return None

def main():
    if len(sys.argv) < 3:
        print("Usage: update_manifest.py <version_file> <assets_json_string> <artifacts_dir>")
        sys.exit(1)

    version_file = sys.argv[1]
    assets_json_str = sys.argv[2]
    artifacts_dir = sys.argv[3]

    try:
        assets = json.loads(assets_json_str)
    except json.JSONDecodeError as e:
        print(f"Error decoding assets JSON: {e}")
        sys.exit(1)

    with open(version_file, 'r') as f:
        manifest = json.load(f)

    if 'platforms' not in manifest:
        manifest['platforms'] = {}

    print(f"Updating manifest for version {manifest.get('version', 'unknown')}...")

    for asset in assets:
        filename = asset['name']
        download_url = asset['browser_download_url']
        
        # Skip non-binary files (like sha256sums.txt)
        if filename.endswith(".txt") or filename.endswith(".sh") or filename.endswith(".ps1"):
            continue

        key = get_platform_key(filename)
        if not key:
            print(f"Skipping {filename}: could not determine platform key")
            continue

        local_path = os.path.join(artifacts_dir, filename)
        if not os.path.exists(local_path):
            print(f"Warning: Local file {local_path} not found. Cannot compute checksum.")
            checksum = ""
        else:
            checksum = compute_sha256(local_path)
            print(f"Processed {filename} -> {key} (SHA256: {checksum[:8]}...)")

        manifest['platforms'][key] = {
            "url": download_url,
            "sha256": checksum
        }

    with open(version_file, 'w') as f:
        json.dump(manifest, f, indent=4)
        f.write('\n') # Add trailing newline

    print("Manifest updated successfully.")

if __name__ == "__main__":
    main()
