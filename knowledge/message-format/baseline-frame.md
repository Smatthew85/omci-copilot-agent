# Baseline OMCI Frame Format

ITU-T G.988 defines two OMCI frame formats. The **baseline** format is the original 48-byte fixed-length frame that all compliant ONUs must support. The extended format (Device Identifier `0x0B`) is covered in [`extended-frame.md`](extended-frame.md).

---

## Overall Structure

Every baseline OMCI frame is exactly **48 bytes**.

| Offset (bytes) | Length (bytes) | Field | Description |
|---|---|---|---|
| 0 | 2 | Transaction Correlation ID (TCID) | Big-endian unsigned integer. The OLT assigns a unique TCID per request; the ONU echoes the same TCID in the response. Used to match request/response pairs. |
| 2 | 1 | Message Type | Encodes the AR/AK flags and the MT value (see [Message Type Byte](#message-type-byte) below). |
| 3 | 1 | Device Identifier | Always `0x0A` for baseline frames. |
| 4 | 2 | ME Class ID | Big-endian unsigned integer identifying the Managed Entity class (e.g., `0x0002` = ONU-G). |
| 6 | 2 | ME Instance ID | Big-endian unsigned integer identifying the specific ME instance (e.g., `0x0000` for the single ONU-G instance). |
| 8 | 32 | Message Contents | Action-specific payload. Unused bytes are padded with `0x00`. |
| 40 | 4 | OMCI Trailer | Upper 2 bytes encode the length of the frame excluding the trailer and CRC (bytes 0–39), always `0x0028` = 40 for baseline frames. Lower 2 bytes are padding (`0x0000`). |
| 44 | 4 | CRC-32 | IEEE 802.3 CRC-32 over bytes 0–43. Big-endian. Used for frame-integrity validation. |

> **Note:** ME Class ID and ME Instance ID are both big-endian 16-bit values.

---

## Message Type Byte

Byte offset 2 packs three logical fields into a single octet:

```
  Bit 7 (MSB)   Bit 6   Bit 5   Bits 4–0 (LSB)
  ┌──────────┬───────┬───────┬──────────────────┐
  │    DB    │  AR   │  AK   │     MT[4:0]      │
  └──────────┴───────┴───────┴──────────────────┘
```

| Bit(s) | Name | Description |
|---|---|---|
| 7 | **DB** (Destination Bit) | `0` = OLT→ONU (command), `1` = ONU→OLT (autonomous notification). For standard request/response the OLT sets `0`; alarm and AVC notifications from the ONU set `1`. |
| 6 | **AR** (Acknowledge Request) | `1` = sender expects an acknowledgement. OLT sets AR=1 for most commands. |
| 5 | **AK** (Acknowledgement) | `1` = this frame IS an acknowledgement. ONU sets AK=1 in response frames. A frame cannot have both AR=1 and AK=1. |
| 4–0 | **MT** | Message Type value (5 bits, 0–31). See [`message-types.md`](message-types.md) for the full table. |

### Common Message Type Byte Values

| Hex | Binary | Meaning |
|---|---|---|
| `0x44` | `0100 0100` | AR=1, AK=0, MT=4 → **Create request** |
| `0x24` | `0010 0100` | AR=0, AK=1, MT=4 → **Create response** |
| `0x48` | `0100 1000` | AR=1, AK=0, MT=8 → **Set request** |
| `0x28` | `0010 1000` | AR=0, AK=1, MT=8 → **Set response** |
| `0x49` | `0100 1001` | AR=1, AK=0, MT=9 → **Get request** |
| `0x29` | `0010 1001` | AR=0, AK=1, MT=9 → **Get response** |
| `0x4F` | `0100 1111` | AR=1, AK=0, MT=15 → **MIB Reset request** |
| `0x2F` | `0010 1111` | AR=0, AK=1, MT=15 → **MIB Reset response** |

---

## Worked Example: Decoding a Create Request

### Raw Frame (hex)

```
0023 44 0A 00AB 0001 0000000000000000000000000000000000000000000000000000000000000000 0028 0000 A1B2C3D4
```

Formatted as a single hex string (48 bytes = 96 hex characters):

```
002344 0A 00AB 0001 00000000000000000000000000000000 00000000000000000000000000000000 00280000 A1B2C3D4
```

Let's lay it out byte by byte:

```
Offset  Hex Bytes       Field
------  -----------     -----
0–1     00 23           TCID = 0x0023 = 35
2       44              Message Type = 0x44 → AR=1, AK=0, MT=0x04 (Create)
3       0A              Device Identifier = 0x0A (baseline)
4–5     00 AB           ME Class ID = 0x00AB = 171 decimal → MAC Bridge Service Profile
6–7     00 01           ME Instance ID = 0x0001 = instance 1
8–39    00 00 ... 00    Message Contents (32 bytes) — Create attributes (all zeros in this example)
40–43   00 28 00 00     OMCI Trailer — length = 0x0028 = 40 bytes, padding = 0x0000
44–47   A1 B2 C3 D4     CRC-32
```

### Interpretation

- **TCID 35**: The OLT will match the ONU's response by looking for a frame with TCID = 35 and AK=1.
- **MT = Create (0x04)**: The OLT is asking the ONU to instantiate a new ME.
- **ME Class 0x00AB (171)**: MAC Bridge Service Profile — a Layer-2 bridging entity.
- **ME Instance 0x0001**: Instance 1 of the MAC Bridge Service Profile.
- **Message Contents**: The 32-byte payload carries the initial attribute values for the new ME instance. In this example all bytes are zero, which means the ONU will use default values for every attribute.
- **Trailer length 0x0028 = 40**: Counts bytes 0–39 (everything before the trailer itself).
- **CRC-32 `A1B2C3D4`**: Computed over bytes 0–43.

### Expected Response

The ONU should reply with TCID=35, Message Type=`0x24` (AK=1, MT=Create), and a 1-byte result code in the Message Contents:

| Result | Meaning |
|---|---|
| `0x00` | Success |
| `0x07` | Instance exists (already provisioned) |
| `0x04` | Unknown ME (ONU firmware does not support class 0x00AB) |
| `0x03` | Parameter error (invalid attribute value in the Create payload) |

---

## Notes

- The 32-byte Message Contents field is sufficient for most MEs. For MEs with large table attributes or many attributes, the **extended frame** format should be used instead.
- The OMCI Trailer length field always counts bytes 0 through the last content byte (offset 39), giving a constant value of `0x0028` (40) for baseline frames.
- CRC-32 validation is mandatory for conformant implementations but is sometimes disabled in test environments.
