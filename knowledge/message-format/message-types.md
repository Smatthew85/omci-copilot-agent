# OMCI Message Types

ITU-T G.988 defines a set of standardised Message Type (MT) values that fit in the 5 low-order bits (bits 4–0) of the Message Type byte. The full Message Type byte also encodes the AR and AK flags; see [`baseline-frame.md`](baseline-frame.md) for the bit layout.

---

## AR and AK Bit Semantics

| Flag | Name | Meaning |
|---|---|---|
| **AR** | Acknowledge Request | When `1`, the sender requires the peer to send an acknowledgement frame. The OLT sets AR=1 for almost all commands. The ONU sets AR=1 only for autonomous notifications that require a response (none in standard G.988). |
| **AK** | Acknowledgement | When `1`, this frame is the acknowledgement (response) to a previous AR=1 frame. A frame must not have both AR=1 and AK=1 simultaneously. |

Autonomous notifications from the ONU (Alarm Notification, Attribute Value Change, Test Result) set AR=0 and AK=0.

---

## Message Type Reference Table

| MT (decimal) | MT (hex) | Name | Direction | AR | Notes |
|---|---|---|---|---|---|
| 4 | `0x04` | Create | OLT → ONU | 1 | Instantiate a new ME instance. Response contains a result code. |
| 6 | `0x06` | Delete | OLT → ONU | 1 | Remove an existing ME instance. Response contains a result code. |
| 8 | `0x08` | Set | OLT → ONU | 1 | Modify one or more attributes of an existing ME instance. |
| 9 | `0x09` | Get | OLT → ONU | 1 | Read one or more attributes. Response carries the attribute values. |
| 11 | `0x0B` | Get All Alarms | OLT → ONU | 1 | Request the full alarm status bitmap for all MEs. Followed by Get All Alarms Next. |
| 12 | `0x0C` | Get All Alarms Next | OLT → ONU | 1 | Fetch the next entry in the alarm list started by Get All Alarms. |
| 13 | `0x0D` | MIB Upload | OLT → ONU | 1 | Initiate a MIB synchronisation. Response includes the number of subsequent MIB Upload Next frames to expect. |
| 14 | `0x0E` | MIB Upload Next | OLT → ONU | 1 | Retrieve the next ME snapshot during a MIB upload sequence. |
| 15 | `0x0F` | MIB Reset | OLT → ONU | 1 | Reset the ONU MIB to its factory-default state. Use with caution — clears all provisioned MEs. |
| 16 | `0x10` | Alarm Notification | ONU → OLT | 0 | Autonomous: ONU reports an alarm condition change. AR=0, AK=0. |
| 17 | `0x11` | Attribute Value Change | ONU → OLT | 0 | Autonomous: ONU reports a spontaneous attribute change (e.g., optical Rx power drift). AR=0, AK=0. |
| 18 | `0x12` | Test | OLT → ONU | 1 | Initiate a self-test or optical diagnostic (BERT, loopback, etc.). |
| 19 | `0x13` | Start Software Download | OLT → ONU | 1 | Begin a firmware download session. Negotiates window size and image size. |
| 20 | `0x14` | Download Section | OLT → ONU | 1 (last section) | Carry a firmware image section. AR is set only on the last section of a window; intermediate sections use AR=0. |
| 21 | `0x15` | End Software Download | OLT → ONU | 1 | Signal the end of a firmware download; ONU validates the image CRC. |
| 22 | `0x16` | Activate Software | OLT → ONU | 1 | Instruct the ONU to boot from the downloaded image. |
| 23 | `0x17` | Commit Software | OLT → ONU | 1 | Mark the active image as the committed (fallback) image. |
| 24 | `0x18` | Synchronize Time | OLT → ONU | 1 | Set the ONU's real-time clock. |
| 25 | `0x19` | Reboot | OLT → ONU | 1 | Instruct the ONU to perform a soft reboot. |
| 26 | `0x1A` | Get Next | OLT → ONU | 1 | Retrieve the next chunk of a multi-part Get response or table attribute. Carries a sequence number. |
| 27 | `0x1B` | Test Result | ONU → OLT | 0 | Autonomous: ONU delivers the result of a previously requested Test (MT 18). AR=0, AK=0. |
| 28 | `0x1C` | Get Current Data | OLT → ONU | 1 | Read current PM counter values (as opposed to interval snapshots read with Get). |
| 29 | `0x1D` | Set Table | OLT → ONU | 1 | Write one or more rows to a table attribute in bulk. Typically requires extended frames for large tables. |

---

## Constructing the Full Message Type Byte

Given an MT value and the desired AR/AK state, the Message Type byte is:

```
MT_byte = (DB << 7) | (AR << 6) | (AK << 5) | (MT & 0x1F)
```

### Examples

| Operation | DB | AR | AK | MT | MT_byte (hex) |
|---|---|---|---|---|---|
| Create request (OLT→ONU) | 0 | 1 | 0 | 4 | `0x44` |
| Create response (ONU→OLT) | 0 | 0 | 1 | 4 | `0x24` |
| MIB Reset request | 0 | 1 | 0 | 15 | `0x4F` |
| MIB Reset response | 0 | 0 | 1 | 15 | `0x2F` |
| Alarm Notification (autonomous) | 1 | 0 | 0 | 16 | `0x90` |
| Attribute Value Change (autonomous) | 1 | 0 | 0 | 17 | `0x91` |
| Set Table request | 0 | 1 | 0 | 29 | `0x5D` |
| Set Table response | 0 | 0 | 1 | 29 | `0x3D` |

---

## Notes

- MT values 0–3 and 30–31 are reserved in G.988 and should not appear in conformant implementations.
- The **Download Section** MT (20) is unusual in that AR is set to `0` for all sections within a window except the last, which uses AR=1 to trigger an acknowledgement.
- Autonomous notifications (MT 16, 17, 27) originate from the ONU and always set DB=1, AR=0, AK=0. They do not carry a TCID that correlates to a previous OLT request.
