#!/usr/bin/env python3
"""Generate the embedded HL7 FHIR R4 base catalog from definitions.json.zip."""

from __future__ import annotations

import json
import sys
import urllib.request
import zipfile
from io import BytesIO
from pathlib import Path

DEFINITIONS_URL = "https://hl7.org/fhir/R4/definitions.json.zip"
ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "pkg/registry/internal/bundles/r4"
SD_DIR = OUT / "structure-definitions"
SP_DIR = OUT / "search-parameters"


def main() -> int:
    print(f"Downloading {DEFINITIONS_URL}")
    with urllib.request.urlopen(DEFINITIONS_URL, timeout=120) as resp:
        payload = resp.read()

    with zipfile.ZipFile(BytesIO(payload)) as zf:
        profiles = json.loads(zf.read("profiles-resources.json"))
        search = json.loads(zf.read("search-parameters.json"))

    if SD_DIR.exists():
        for path in SD_DIR.glob("*.json"):
            path.unlink()
    else:
        SD_DIR.mkdir(parents=True, exist_ok=True)

    if SP_DIR.exists():
        for path in SP_DIR.glob("*.json"):
            path.unlink()
    else:
        SP_DIR.mkdir(parents=True, exist_ok=True)

    sd_count = 0
    for entry in profiles.get("entry", []):
        resource = entry.get("resource", {})
        if resource.get("resourceType") != "StructureDefinition":
            continue
        if resource.get("kind") != "resource" or not resource.get("type"):
            continue
        name = resource.get("id") or resource["type"]
        (SD_DIR / f"{name}.json").write_text(json.dumps(resource, separators=(",", ":")))
        sd_count += 1

    sp_count = 0
    for entry in search.get("entry", []):
        resource = entry.get("resource", {})
        if resource.get("resourceType") != "SearchParameter":
            continue
        name = resource.get("id") or resource.get("code") or str(sp_count)
        (SP_DIR / f"{name}.json").write_text(json.dumps(resource, separators=(",", ":")))
        sp_count += 1

    print(f"Wrote {sd_count} StructureDefinitions and {sp_count} SearchParameters to {OUT}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
