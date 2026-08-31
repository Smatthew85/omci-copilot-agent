# Copilot Instructions

## Authoritative OMCI ME Source

- `knowledge/me-catalog/*.json` is the authoritative source for OMCI Managed Entity definitions, including class IDs, attributes, and access rules.
- When referencing ME details in analysis or diagnostics, cite the specific catalog file path (for example, `knowledge/me-catalog/045-mac-bridge-service-profile.json`).
- Do not infer unsupported attributes or access permissions if they are not present in the relevant ME catalog JSON.
