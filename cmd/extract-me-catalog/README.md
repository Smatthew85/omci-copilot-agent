# extract-me-catalog

`cmd/extract-me-catalog` exports one JSON file per OMCI Managed Entity into `knowledge/me-catalog/`.

## Resolved `omci-lib-go/v2` API

The extractor uses public APIs from `github.com/opencord/omci-lib-go/v2/generated`:

- `GetSupportedClassIDs()` to enumerate known ME class IDs.
- `LoadManagedEntityDefinition(classID)` to load each ME definition.
- `ManagedEntity` accessors (`GetName()`, `GetClassID()`, `GetMessageTypes()`, `GetAttributeDefinitions()`).
- `AttributeDefinition` metadata (`GetIndex()`, `GetName()`, `GetSize()`, `IsTableAttribute()`, `Optional`, and access checks via `SupportsAttributeAccess`).

## Usage

```bash
go run ./cmd/extract-me-catalog [-out DIR] [-pretty=false] [-index=false]
```

Flags:

- `-out string` output directory (default `knowledge/me-catalog`)
- `-pretty` pretty-print JSON (default `true`)
- `-index` write `INDEX.md` (default `true`)

## Updating upstream dependency

```bash
go get github.com/opencord/omci-lib-go/v2@latest && go mod tidy
```

## Exit codes

- `0`: success
- non-zero: extraction or write failure
