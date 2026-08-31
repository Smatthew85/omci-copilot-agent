# Extended OMCI Frame Format

ITU-T G.988 defines an **extended** OMCI frame format identified by Device Identifier `0x0B`. It provides a much larger Message Contents field (up to 1904 bytes) compared to the 32-byte limit of the baseline format.

---

## When Extended Frames Are Used

Extended frames are required whenever the payload of a single OMCI operation exceeds 32 bytes. Common scenarios include:

- **Large `SetTable` operations** (MT = 0x1D) — writing multi-row table data in a single frame.
- **MEs with many attributes** — `Get` and `Get Response` frames for MEs that have more attributes than can fit in 32 bytes (e.g., Extended VLAN Tagging Operation Config Data with many rules).
- **Large software download sections** — `Download Section` (MT = 0x14) frames carrying large firmware chunks.
- **`MIB Upload Next` responses** — when a single ME's attribute snapshot exceeds 32 bytes.

---

## Overall Structure

| Offset (bytes) | Length (bytes) | Field | Description |
|---|---|---|---|
| 0 | 2 | Transaction Correlation ID (TCID) | Same semantics as baseline: big-endian, echoed in response. |
| 2 | 1 | Message Type | Same AR/AK/MT encoding as baseline (see [`baseline-frame.md`](baseline-frame.md)). |
| 3 | 1 | Device Identifier | **`0x0B`** — signals an extended frame. |
| 4 | 2 | ME Class ID | Big-endian ME class identifier. |
| 6 | 2 | ME Instance ID | Big-endian ME instance identifier. |
| 8 | 2 | Message Contents Length | Big-endian unsigned integer: number of bytes in the Message Contents field that follow (0–1904). |
| 10 | 0–1904 | Message Contents | Variable-length action-specific payload. |
| 10+N | 4 | CRC-32 | IEEE 802.3 CRC-32 over all preceding bytes. Big-endian. |

> The total frame length is `10 + N + 4` bytes, where N is the value of the Message Contents Length field.

> **There is no fixed OMCI Trailer in extended frames.** The length information is carried in the Message Contents Length field at offset 8.

---

## Comparison: Baseline vs. Extended

| Attribute | Baseline | Extended |
|---|---|---|
| Device Identifier | `0x0A` | `0x0B` |
| Frame length | Fixed 48 bytes | Variable (14–1918 bytes) |
| Message Contents length | Fixed 32 bytes | Variable 0–1904 bytes (2-byte length field) |
| OMCI Trailer | Present at offset 40 (4 bytes) | Absent |
| CRC-32 position | Fixed at offset 44 | Immediately after Message Contents |
| Maximum payload | 32 bytes | ~1904 bytes |
| ONU support requirement | Mandatory | Optional (negotiated via ONU-G Extended OMCI capability attribute) |

---

## Negotiating Extended Frame Support

Before using extended frames the OLT must confirm that the ONU supports them. The `ONU-G` ME (class 0x0002) contains an **Extended OMCI** attribute. If the ONU advertises support, the OLT may use Device Identifier `0x0B` frames. If the attribute is absent or zero, the OLT must use baseline frames only.

---

## Worked Example: Decoding an Extended Get Response

### Scenario

The OLT sent a `Get` request for all attributes of the **Extended VLAN Tagging Operation Configuration Data** ME (class `0x0171` = 369). The ME has a large table attribute. The ONU replies with an extended `Get Response` frame.

### Raw Frame (hex)

```
0041 29 0B 0171 0001 001A  <20 bytes of attribute data>  <CRC>
```

Expanded with a 26-byte Message Contents for this example:

```
Bytes 0–1:   00 41       TCID = 0x0041 = 65
Byte  2:     29          Message Type = 0x29 → AR=0, AK=1, MT=0x09 (Get response)
Byte  3:     0B          Device Identifier = 0x0B (extended)
Bytes 4–5:   01 71       ME Class ID = 0x0171 = 369 → Extended VLAN Tagging Op Config Data
Bytes 6–7:   00 01       ME Instance ID = 0x0001 = instance 1
Bytes 8–9:   00 1A       Message Contents Length = 0x001A = 26 bytes
Bytes 10–35: <26 bytes>  Message Contents (result code + attribute mask + attribute values)
Bytes 36–39: <CRC-32>    CRC over bytes 0–35
```

### Field-by-Field Breakdown

| Offset | Hex | Value | Interpretation |
|---|---|---|---|
| 0–1 | `00 41` | 65 | TCID — matches the OLT's request TCID 65 |
| 2 | `29` | — | AK=1 (response), AR=0, MT=9 (Get) |
| 3 | `0B` | — | Extended frame |
| 4–5 | `01 71` | 369 | Extended VLAN Tagging Op Config Data |
| 6–7 | `00 01` | 1 | ME instance 1 |
| 8–9 | `00 1A` | 26 | 26 bytes of Message Contents follow |
| 10 | `00` | 0 | Result code: Success |
| 11–12 | `FF FF` | — | Attribute presence mask (all 16 attribute bits set = all attributes present) |
| 13–35 | `...` | — | Encoded attribute values for the requested attributes |
| 36–39 | `<CRC>` | — | CRC-32 over bytes 0–35 |

### Key Observations

- The `0x0B` Device Identifier at offset 3 immediately signals that this is an extended frame.
- The 2-byte Message Contents Length at offsets 8–9 tells the receiver exactly how many bytes to read before expecting the CRC.
- The result code is the first byte of Message Contents (offset 10) for Get responses — same convention as baseline.
- There is **no OMCI Trailer** field; the length is self-describing.

---

## Notes

- Extended frames are backwards-incompatible: a baseline-only ONU will silently discard or reject a frame with Device Identifier `0x0B`.
- When a `Get` response for a large ME must be split across multiple frames, the OLT uses **Get Next** (MT = `0x1A`) with a sequence counter to retrieve subsequent chunks; each chunk may independently be baseline or extended.
- CRC-32 computation covers all bytes from the start of the frame up to (but not including) the CRC field itself, regardless of frame type.
