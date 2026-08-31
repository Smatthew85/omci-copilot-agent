# ME Catalog

This folder contains the authoritative, machine-readable OMCI Managed Entity (ME) catalog used by the Copilot agent.

## Source

Generated from [`github.com/opencord/omci-lib-go/v2`](https://github.com/opencord/omci-lib-go/tree/master/generated).

## Regeneration

```bash
go run ./cmd/extract-me-catalog
```

The extractor writes one JSON file per ME and regenerates `INDEX.md`.

## JSON schema

Each `NNN-name.json` file includes:

- `class_id` (number)
- `name` (string)
- `message_types` (optional string array)
- `attributes` (array), where each attribute contains:
  - `index` (number)
  - `name` (string)
  - `size_bytes` (number; `0` for table/variable-size)
  - `access` (string array of any: `Read`, `Write`, `SetByCreate`, `Delete` when exposed)
  - `mandatory` (boolean)
  - `table` (boolean)
- `source` object:
  - `library` (`github.com/opencord/omci-lib-go/v2`)
  - `version` (resolved module version at extraction time)

## Automation

`.github/workflows/sync-me-catalog.yml` re-runs this extractor weekly.
