# Diagnosis: Set Extended VLAN Tagging Op Config Data returns 0x03 Parameter Error

## Summary

The OLT issued a Set request targeting attribute 7 ("Received frame VLAN tagging operation
table") of Extended VLAN Tagging Operation Configuration Data (ME Class 171), instance 1.
The ONU responded with result code **0x03 — Parameter Error**, indicating that one or more
fields in the 16-byte VLAN rule tuple fall outside the allowed ranges defined in G.988.

## Root Cause

Each entry in the "Received frame VLAN tagging operation table" attribute is a **16-byte
rule tuple** encoding filter fields (bytes 0–7) and treatment fields (bytes 8–15). The
ONU validates every field against the allowed ranges specified in G.988, clause 9.3.11.

In this case the malformed tuple contains:

- **Filter outer priority (bits [63:61])**: value `0xF` (binary `1111`) — this value is
  reserved. Valid values are `0–7` (specific priority), `8` (don't care / no tag), or
  `15` (don't care with any tag); `0xF` with conflicting semantics is rejected by
  strict implementations.
- **Filter TPID/DE (bits [47:32])**: value `0xEEEE` — not a valid IEEE TPID
  (`0x8100`, `0x88A8`, `0x9100` are the recognized values per G.988).
- **Treatment outer priority (bits [31:29] of treatment section)**: value `0xFF` exceeds
  the 4-bit field width.

The ONU rejects the entire Set operation before applying any changes to its MIB.

## Evidence

- **Message Type**: MT=8 — Set (request byte `0x48`, response byte `0x28`)
- **ME Class**: 171 — Extended VLAN Tagging Operation Configuration Data
- **ME Instance**: 1
- **Attribute Mask**: `0x0100` — attribute 7 only ("Received frame VLAN tagging
  operation table")
- **Result Code**: `0x03` — Parameter Error (response content byte 0)
- **Suspect bytes in rule tuple** (raw `0FFFEEEEFFFFFFFF0000000000000000`):
  - Bytes 0–1: `0F FF` — filter outer priority nibble `0xF` is out of range
  - Bytes 2–3: `EE EE` — invalid TPID value
  - Bytes 8–9: `FF FF` — treatment priority overflow

See also:
- [`knowledge/result-codes/README.md`](../../result-codes/README.md) — definition of 0x03
- [`knowledge/failure-patterns/README.md`](../../failure-patterns/README.md) — "Set returns Parameter Error" pattern

## Remediation

1. **Validate the 16-byte rule tuple** against G.988, clause 9.3.11 (ME 171 attribute
   table encoding). Specifically:
   - Filter outer priority: `0`–`7` for a specific 802.1p priority; `8` = no outer tag;
     `14` = don't care; `15` = treat as untagged (check the exact semantics in your
     G.988 edition).
   - Filter TPID: use only `0x8100` (C-VLAN), `0x88A8` (S-VLAN), or `0x9100`
     (legacy double-tag); do not use arbitrary values.
   - Treatment priority: must fit in 4 bits (`0`–`15`).
2. **Correct the rule tuple** in the OLT provisioning configuration and reissue the Set.
3. **Verify all table entries** if the provisioner constructs rules dynamically — the same
   field encoding error may affect other rule entries.

## References

- [`knowledge/provisioning-flows/`](../../provisioning-flows/) — standard ONU provisioning
  sequence, VLAN configuration section
- [`knowledge/result-codes/README.md`](../../result-codes/README.md)
- [`knowledge/failure-patterns/README.md`](../../failure-patterns/README.md)
- ITU-T G.988, clause 9.3.11 — Extended VLAN Tagging Operation Configuration Data
- ITU-T G.988, Table 9.3.11-1 — VLAN tagging operation table entry format
