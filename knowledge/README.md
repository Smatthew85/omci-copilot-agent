# Knowledge Base Index

This directory contains the reference knowledge used by the OMCI Protocol Analyst Copilot agent. All content is grounded in ITU-T G.988 and is intended to be vendor-neutral unless explicitly noted.

---

## Directory Layout

| Folder | Contents |
|---|---|
| `me-catalog/` | JSON ME definitions produced by the existing extractor tooling (sourced from `opencord/omci-lib-go`). Each file covers one ME class and includes the Class ID, name, attributes (index, name, size, access type, mandatory/optional), and supported actions. |
| `result-codes/` | OMCI result/reason code reference — code table with hex values, names, meanings, and typical causes; plus the full G.988 message type reference. |
| `provisioning-flows/` | Standard and (optionally) vendor-specific ONU provisioning sequences. Documents the expected order of Create/Set operations, ME dependency graph, and common failure points at each step. |
| `failure-patterns/` | Symptom → likely cause → suggested action mappings for common OMCI provisioning failures. |

---

## Files

- [`result-codes/README.md`](result-codes/README.md) — G.988 result/reason code reference
- [`provisioning-flows/standard-onu-provisioning.md`](provisioning-flows/standard-onu-provisioning.md) — Standard ONU provisioning sequence
- [`failure-patterns/README.md`](failure-patterns/README.md) — Failure symptom diagnostic table

---

## Adding Content

- **New ME JSON files** — place them in `me-catalog/` following the existing file naming convention.
- **Vendor-specific provisioning flows** — add `<vendor>-onu-provisioning.md` alongside the standard flow.
- **Vendor-specific failure patterns** — add `<vendor>-failure-patterns.md` alongside the baseline patterns.
- **Examples** — create an `examples/` subdirectory with annotated real-world failure cases (anonymized).
