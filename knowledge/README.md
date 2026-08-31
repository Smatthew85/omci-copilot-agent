# OMCI Copilot Agent — Knowledge Base

This directory contains the reference material used by the Copilot agent to analyze
OMCI messages and diagnose provisioning failures.

---

## Index

| Directory | Contents |
|-----------|----------|
| [`result-codes/`](./result-codes/) | G.988 OMCI result/reason codes and message type reference |
| [`failure-patterns/`](./failure-patterns/) | Symptom → root cause → remediation diagnostic table |
| [`provisioning-flows/`](./provisioning-flows/) | Standard ONU OMCI provisioning sequence with ME dependencies |
| [`examples/`](./examples/) | Golden set of anonymized OMCI failure cases with expected diagnoses (few-shot references) |
| [`logs/`](./logs/) | VOLTHA and OLT hardware log formats — how to extract OMCI frames from pasted log output |
| [`vendors/`](./vendors/) | Vendor-specific ONU/OLT deviations from G.988 (ME gaps, firmware quirks, workarounds) |

---

## Contributing

- Add new result codes or message types to [`result-codes/README.md`](./result-codes/README.md).
- Add new failure patterns to [`failure-patterns/README.md`](./failure-patterns/README.md).
- Add new provisioning flow documentation under [`provisioning-flows/`](./provisioning-flows/).
- Add new golden examples following the guide in [`examples/README.md`](./examples/README.md).
# Copilot Agent Knowledge Base

This directory contains reference documentation used by the GitHub Copilot agent to analyze ITU-T G.988 OMCI messages and diagnose ONU provisioning failures.

## Contents

| Directory | Description |
|---|---|
| [`message-format/`](message-format/) | OMCI frame formats — baseline (48-byte) and extended layouts, Message Type reference |
| [`result-codes/`](result-codes/) | G.988 result and reason codes with diagnostic guidance |
| [`provisioning-flows/`](provisioning-flows/) | Standard ONU provisioning sequences and ME dependency order |
| [`failure-patterns/`](failure-patterns/) | Symptom → root-cause diagnostic table |

## How to Use

The Copilot agent is instructed (via [`.github/copilot-instructions.md`](../.github/copilot-instructions.md)) to consult this knowledge base when:

1. A user pastes a raw hex OMCI frame or a sequence of frames.
2. A user describes a provisioning failure symptom.
3. A user asks about a specific ME class, attribute, or result code.

Start with [`message-format/`](message-format/) to decode the frame structure, then cross-reference [`result-codes/`](result-codes/) for response interpretation, and consult [`failure-patterns/`](failure-patterns/) for root-cause guidance.
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
