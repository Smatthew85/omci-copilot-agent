# Alarm Message Format (MT 16 / `0x10`)

> **Source:** ITU-T G.988 §11.3.3 (Alarm notification)

---

## Overview

An **Alarm** message is an autonomous notification sent by the ONU to the OLT to
report that a fault condition has been raised or cleared on a specific ME instance.

| Property | Value |
|---|---|
| Message Type byte | `0x10` (bit 7–5 = `000`, bit 4 = `1`, bit 3–0 = `0000` → MT 16) |
| AR bit | 0 (no acknowledgement requested) |
| AK bit | 0 (not an acknowledgement) |
| Direction | ONU → OLT |
| Requires response | No |

---

## Baseline Frame Layout (48 bytes)

```
Byte(s)  Field                   Size   Notes
──────────────────────────────────────────────────────────────────────
0–1      Transaction Correlation   2    Set to 0x0000 by ONU for autonomous messages
         ID (TCID)
2        Message Type (MT)         1    0x10
3        Device Identifier (DevID) 1    0x0A = baseline OMCI
4–5      ME Class ID               2    Big-endian
6–7      ME Instance ID            2    Big-endian
──────── Content (32 bytes, bytes 8–39) ────────────────────────────
8–35     Alarm bitmap             28    224 bits; bit 7 of byte 8 = alarm 0 (MSB-first)
                                        Bit set → alarm active
                                        Bit clear → alarm cleared / not present
36–38    Reserved                  3    Set to 0x000000 by ONU; ignored by OLT
39       Alarm sequence number     1    Monotonically increasing per ONU (wraps 0xFF→0x00)
──────── OMCI Trailer (8 bytes, bytes 40–47) ───────────────────────
40–41    OMCI Message Length       2    0x0028 (40 decimal) for baseline
42       CPCS-UU                   1    0x00
43       CPI                       1    0x00
44–47    CRC-32 (MIC)              4    ITU-T CRC-32 over bytes 0–43
──────────────────────────────────────────────────────────────────────
```

> **Note (G.988 §11.2.2):** The 32-byte content area occupies bytes 8–39. It contains
> 28 bytes of alarm bitmap + 3 reserved bytes + 1 alarm sequence number byte = 32 bytes.
> The OMCI trailer (8 bytes) occupies bytes 40–47.

---

## Extended Frame Layout

For extended OMCI (Device Identifier = `0x0B`), the same fields appear in the extended
framing defined in G.988 §11.2.3:

```
Byte(s)  Field                   Size   Notes
──────────────────────────────────────────────────────────────────────
0–1      TCID                      2    0x0000 for autonomous messages
2        MT                        1    0x10
3        DevID                     1    0x0B = extended OMCI
4–5      ME Class ID               2    Big-endian
6–7      ME Instance ID            2    Big-endian
8–9      Content Length            2    Length of content that follows (in bytes)
10–37    Alarm bitmap             28    Same semantics as baseline
38–40    Reserved                  3    0x000000
41       Alarm sequence number     1
42+      Padding / extension       —    To OMCI content-length boundary
(last 4) CRC-32                    4
──────────────────────────────────────────────────────────────────────
```

---

## Alarm Bitmap Semantics

The 28-byte (224-bit) alarm bitmap is the core payload of an alarm notification.

- **Bit numbering:** Within each byte, **bit 7 is the most significant** (alarm 0 of that
  byte group). Across the 28 bytes, bit 0 of byte 8 (of the full OMCI frame) = alarm 0
  for the ME.
- **Bit set (1):** The alarm condition is **active** (raised).
- **Bit clear (0):** The alarm condition is **cleared** (not present).
- **ME-specific meaning:** Which alarm number maps to which fault condition is defined
  per ME class. See [`common-alarms.md`](common-alarms.md) for well-known MEs, or consult
  `knowledge/me-catalog/NNN-slug.json` (once populated) for the authoritative mapping.
- **Bits beyond the defined range:** Must be set to 0 by the ONU. The OLT should ignore
  undefined bits rather than treating them as faults.

### Alarm Bit Extraction Example

To check whether alarm bit 5 is set in a 28-byte bitmap:

```
Byte index within bitmap = 5 ÷ 8 = 0 (byte 0 of the bitmap, i.e., frame byte 8)
Bit position within byte = 5 mod 8 = 5 → bit 5 counting from MSB = bit 2 counting from LSB
Mask = 0x04
```

In general: `alarm_active = (bitmap[alarm_n // 8] >> (7 - (alarm_n % 8))) & 1`

---

## Alarm Sequence Number

The 1-byte alarm sequence number (byte 39 in a baseline frame — the last byte of the
32-byte content area) is a monotonically increasing counter maintained by the ONU across
**all** alarm notifications, regardless of ME class or instance.

- Range: `0x00`–`0xFF`, then wraps to `0x00`.
- The OLT uses gaps in this sequence to detect missed alarm notifications.
- Upon detecting a gap the OLT should initiate a **Get All Alarms** (MT 27) resync.
  See [`alarm-synchronization.md`](alarm-synchronization.md) for the full procedure.

---

## Decoded Example

The following is a structurally valid fabricated example of an alarm notification from
an ANI-G ME (Class 263, `0x0107`) instance 0, raising alarm bit 2 (Signal fail):

```
Raw hex (48 bytes, baseline):
0000 10 0A 0107 0000
04 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
00 00 00 00 00 00 00 00 00 00 00 00 00 00 3C 00
28 00 00 00 AABBCCDD

Formatted:
  Bytes 0–1   TCID           = 0x0000  (autonomous)
  Byte  2     MT             = 0x10    (Alarm)
  Byte  3     DevID          = 0x0A    (baseline)
  Bytes 4–5   ME Class       = 0x0107  = 263 (ANI-G)
  Bytes 6–7   ME Instance    = 0x0000  (instance 0)
  Bytes 8–35  Alarm bitmap:
                Byte 8       = 0x04   → binary 0000 0100
  Bytes 36–38 Reserved       = 0x000000
  Byte  39    Seq number     = 0x3C   = 60
  Bytes 40–41 OMCI Length    = 0x0028 = 40
  Bytes 42–43 CPCS-UU/CPI   = 0x0000
  Bytes 44–47 CRC-32         = 0xAABBCCDD (illustrative)
```

Decoding bitmap byte 8 = `0x04` = `0000 0100`:
- Bit 7 (alarm 0) = 0 → alarm 0 not active
- Bit 6 (alarm 1) = 0
- Bit 5 (alarm 2) = 1 → **alarm 2 active**
- Bits 4–0        = 0

For ANI-G (ME 263), alarm 2 corresponds to **Signal fail** (see [`common-alarms.md`](common-alarms.md)).

---

## References

- ITU-T G.988 §11.3.3 — Alarm notification message
- ITU-T G.988 §11.2.2 — OMCI baseline frame format
- [`common-alarms.md`](common-alarms.md) — per-ME alarm bit definitions
- [`alarm-synchronization.md`](alarm-synchronization.md) — sequence number handling
