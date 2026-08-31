# OMCI Copilot Agent — Golden Examples

This directory contains a curated set of anonymized OMCI failure cases with expected
diagnoses. These serve two purposes:

1. **Few-shot references** — the Copilot agent consults these examples when analyzing
   user-provided OMCI captures to produce consistent, structured diagnoses.
2. **Regression cases** — each example's `diagnosis.md` is the expected agent output and
   can be used to validate diagnostic quality as the knowledge base evolves.

---

## Directory Index

| # | Slug | ME Class | Result Code | Summary |
|---|------|----------|-------------|---------|
| 01 | [stale-mib-instance-exists](./01-stale-mib-instance-exists/) | 45 — MAC Bridge Service Profile | 0x07 Instance Exists | Create rejected because prior provisioning attempt left stale MIB state on the ONU. |
| 02 | [extended-vlan-parameter-error](./02-extended-vlan-parameter-error/) | 171 — Extended VLAN Tagging Op Config Data | 0x03 Parameter Error | Set rejected due to malformed 16-byte VLAN rule tuple. |

---

## Naming Convention

Each example lives in a subdirectory named `NN-short-slug/`, where:

- `NN` is a zero-padded sequence number (01, 02, …).
- `short-slug` is a lowercase, hyphen-separated description of the failure mode.

Example: `01-stale-mib-instance-exists/`

---

## Required Files per Example

| File | Required | Description |
|------|----------|-------------|
| `input.txt` | ✅ | Raw OMCI frame(s), one per line as continuous hex. Lines starting with `#` are comments; blank lines are ignored. |
| `decoded.json` | ✅ | Structured decoded fields: TCID, MT, ME Class/Instance, attributes, result code. |
| `diagnosis.md` | ✅ | Expected agent output: Summary, Root Cause, Evidence, Remediation, References. |
| `context.md` | Optional | Surrounding provisioning sequence, VOLTHA log excerpts, vendor/environment info. |

> **Note:** `input.txt` contains hex frames as text. The `.txt` extension is used (rather than `.hex`) so that Copilot Spaces accepts the file.

---

## Adding a New Example

1. Copy the `_template/` folder to a new `NN-short-slug/` directory.
2. Replace all placeholder values in each file.
3. Verify your hex frames are structurally plausible (correct DevID `0x0A` for baseline,
   correct field offsets, plausible TCID; mark CRC placeholders in comments).
4. Add a row to the index table above.
5. Open a PR; the example will be reviewed for accuracy and anonymization.

See [`_template/`](./_template/) for starter files.

---

## Anonymization Checklist

Before committing, confirm you have stripped or replaced:

- [ ] ONU serial numbers
- [ ] MAC addresses
- [ ] IP addresses (management plane references)
- [ ] Vendor identifiers embedded in ME attributes
- [ ] Customer names or site identifiers
- [ ] Any other personally identifiable or commercially sensitive data
