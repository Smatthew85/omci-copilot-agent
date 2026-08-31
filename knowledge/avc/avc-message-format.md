# AVC Message Format (MT 17 / `0x11`)

> **Source:** ITU-T G.988 §11.3.4 (Attribute value change notification)

---

## Overview

An **Attribute Value Change (AVC)** message is an autonomous notification sent by the ONU
to the OLT to report that one or more AVC-capable attributes of an ME instance have
changed their values.

| Property | Value |
|---|---|
| Message Type byte | `0x11` (bit 7–5 = `000`, bit 4 = `1`, bit 3–0 = `0001` → MT 17) |
| AR bit | 0 (no acknowledgement requested) |
| AK bit | 0 (not an acknowledgement) |
| Direction | ONU → OLT |
| Requires response | No |

---

## Baseline Frame Layout (48 bytes)

```
Byte(s)  Field                   Size   Notes
──────────────────────────────────────────────────────────────────────
0–1      Transaction Correlation   2    0x0000 for autonomous messages
         ID (TCID)
2        Message Type (MT)         1    0x11
3        Device Identifier (DevID) 1    0x0A = baseline OMCI
4–5      ME Class ID               2    Big-endian
6–7      ME Instance ID            2    Big-endian
──────── Content (32 bytes) ─────────────────────────────────────────
8–9      Attribute mask            2    MSB-first; bit N (counting from bit 15) = attribute N
                                        Bit set → attribute N has a new value in this message
10–39    New attribute values     30    Attribute values packed in mask order (attribute index
                                        ascending), each at its natural G.988 size, zero-padded
                                        to fill unused content bytes
──────── OMCI Trailer (4 bytes) ─────────────────────────────────────
40–41    OMCI Message Length       2    0x0028 (40 decimal) for baseline
42       CPCS-UU / CPI             1    0x00
43–47    CRC-32                    4    ITU-T CRC-32 over bytes 0–43
──────────────────────────────────────────────────────────────────────
```

---

## Extended Frame Layout

For extended OMCI (Device Identifier = `0x0B`):

```
Byte(s)  Field                   Size   Notes
──────────────────────────────────────────────────────────────────────
0–1      TCID                      2    0x0000 for autonomous messages
2        MT                        1    0x11
3        DevID                     1    0x0B = extended OMCI
4–5      ME Class ID               2    Big-endian
6–7      ME Instance ID            2    Big-endian
8–9      Content Length            2    Length of content that follows (in bytes)
10–11    Attribute mask            2    Same semantics as baseline
12+      New attribute values      —    Same packing as baseline; size bounded by Content Length
(last 4) CRC-32                    4
──────────────────────────────────────────────────────────────────────
```

---

## Attribute Mask Semantics

The 2-byte attribute mask identifies which attributes have changed and whose new values
appear in the content area.

- **Bit numbering (G.988 §11.2.2):** Bit 15 (MSB of the 2-byte mask) = attribute 1;
  bit 14 = attribute 2; … bit 0 (LSB) = attribute 16.
- **Bit set (1):** The corresponding attribute value has changed; the new value is packed
  into the content area in ascending attribute-index order.
- **AVC-capable restriction:** Only attributes whose ME definition marks them as
  AVC-capable can appear. Consult `knowledge/me-catalog/NNN-slug.json` (once populated)
  for per-attribute AVC capability flags. If a bit is set for a non-AVC-capable attribute,
  treat it as a protocol violation.
- **Value encoding:** Each attribute value is encoded at its natural size as defined in
  G.988 for that ME (e.g., 1-byte operational state, 2-byte pointer). Values are packed
  consecutively with no gaps; remaining content bytes are zero-padded.

### Attribute Mask Extraction Example

Mask bytes = `0x60 0x00`:

```
Binary: 0110 0000 0000 0000
         ││
         │└── bit 14 → attribute 2 changed
         └─── bit 13 → attribute 3 changed
```

Attributes 2 and 3 are present in the content area, in that order.

---

## Decoded Example: ANI-G Optical Power AVC

The following is a structurally valid fabricated example. ANI-G (Class 263, `0x0107`)
instance 0 reports a change in **Received optical signal level** (attribute 9, a 2-byte
value in units of 0.002 dBm, twos-complement).

```
Raw hex (48 bytes, baseline):
0000 11 0A 0107 0000
00 80 F9 8C 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 28 00 00 AABBCCDD
```

Decoded:

```
  Bytes 0–1   TCID           = 0x0000  (autonomous)
  Byte  2     MT             = 0x11    (AVC)
  Byte  3     DevID          = 0x0A    (baseline)
  Bytes 4–5   ME Class       = 0x0107  = 263 (ANI-G)
  Bytes 6–7   ME Instance    = 0x0000  (instance 0)
  Bytes 8–9   Attribute mask = 0x0080  → bit 8 set → attribute 9 changed
  Bytes 10–11 Attribute 9 value = 0xF98C = -1652 → -1652 × 0.002 = -3.304 dBm
  Bytes 12–39 Padding        = 0x00…
  Bytes 40–41 OMCI Length    = 0x0028
  Byte  42    CPCS-UU/CPI    = 0x00
  Bytes 43–47 CRC-32         = 0xAABBCCDD (illustrative)
```

Interpretation: ANI-G instance 0 has autonomously reported that its received optical
signal level changed to approximately −3.30 dBm. No OLT action was required to trigger
this notification; it occurred because the attribute is AVC-capable and the ONU detected
the value had crossed an internal reporting threshold.

---

## References

- ITU-T G.988 §11.3.4 — Attribute value change notification
- ITU-T G.988 §11.2.2 — OMCI baseline frame format
- [`common-avc-triggers.md`](common-avc-triggers.md) — per-ME AVC trigger scenarios
- [`interpretation-workflow.md`](interpretation-workflow.md) — agent decoding recipe
