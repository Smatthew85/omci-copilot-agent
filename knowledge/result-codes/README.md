# OMCI Result Codes (G.988 Reason Codes)

This document lists the result codes returned in OMCI response messages as defined in
ITU-T G.988. Result codes appear in the first byte of the response content field.

| Hex  | Decimal | Name                          | Description |
|------|---------|-------------------------------|-------------|
| 0x00 | 0       | Success                       | The operation completed successfully. |
| 0x01 | 1       | Command processing error      | The ONU could not process the command; general error not covered by other codes. |
| 0x02 | 2       | Command not supported         | The ONU does not support this action for this ME. |
| 0x03 | 3       | Parameter error               | One or more attribute values are out of the allowed range or structurally invalid. |
| 0x04 | 4       | Unknown managed entity        | The ME Class ID is not recognized by this ONU firmware. |
| 0x05 | 5       | Unknown managed entity instance | The ME Class ID is known but the specified instance does not exist. |
| 0x06 | 6       | Device busy                   | The ONU is temporarily unable to process the request; the OLT may retry. |
| 0x07 | 7       | Instance exists               | A Create action was attempted for an ME instance that already exists in the ONU MIB. |
| 0x08 | 8       | Attribute(s) failed or unknown | One or more attributes in a Set/Get could not be processed; see failed/unsupported mask. |
| 0x09 | 9       | (Reserved)                    | Not defined in baseline G.988; do not use. |

## Notes on 0x08

For result code 0x08, the response content includes:
- Bytes 1–2: **Failed attribute mask** — bitmask of attributes that exist but could not
  be set (e.g., read-only attribute, value out of range for that specific attribute).
- Bytes 3–4: **Unsupported attribute mask** — bitmask of attributes not supported by this
  ONU implementation.

## Message Type (MT) Reference

The MT field occupies the lower 6 bits of the Message Type byte; bits 7–6 are AR (Action
Request) and AK (Acknowledgement).

| MT  | Hex  | Name              | Direction   | Has Response |
|-----|------|-------------------|-------------|--------------|
| 8   | 0x08 | Set               | OLT->ONU    | Yes          |
| 9   | 0x09 | Get               | OLT->ONU    | Yes          |
| 10  | 0x0A | Create            | OLT->ONU    | Yes          |
| 11  | 0x0B | Delete            | OLT->ONU    | Yes          |
| 12  | 0x0C | Get All Alarms    | OLT->ONU    | Yes          |
| 13  | 0x0D | MIB Upload        | OLT->ONU    | Yes          |
| 14  | 0x0E | MIB Upload Next   | OLT->ONU    | Yes          |
| 15  | 0x0F | MIB Reset         | OLT->ONU    | Yes          |
| 16  | 0x10 | Alarm Notification | ONU->OLT   | No           |
| 17  | 0x11 | Attribute Value Change | ONU->OLT | No         |
| 18  | 0x12 | Test              | OLT->ONU    | Yes          |
| 19  | 0x13 | Start Software Download | OLT->ONU | Yes       |
| 20  | 0x14 | Download Section  | OLT->ONU    | Conditional  |
| 21  | 0x15 | End Software Download | OLT->ONU | Yes        |
| 22  | 0x16 | Activate Software | OLT->ONU    | Yes          |
| 23  | 0x17 | Commit Software   | OLT->ONU    | Yes          |
| 24  | 0x18 | Synchronize Time  | OLT->ONU    | Yes          |
| 25  | 0x19 | Reboot            | OLT->ONU    | Yes          |
| 26  | 0x1A | Get Next           | OLT->ONU   | Yes          |
| 27  | 0x1B | Test Result       | ONU->OLT    | No           |
| 28  | 0x1C | Get Current Data  | OLT->ONU    | Yes          |
| 29  | 0x1D | Set Table         | OLT->ONU    | Yes          |
# OMCI Result / Reason Codes

Reference for ITU-T G.988 OMCI baseline message **Result and Reason** codes returned in response messages.

---

## Result / Reason Code Table

| Code (hex) | Name | Meaning | Typical Causes |
|---|---|---|---|
| `0x00` | Command processed successfully | Success — the requested action completed without error | — |
| `0x01` | Command processing error | Generic failure at the ONU | ONU internal error, resource exhaustion, unspecified hardware fault |
| `0x02` | Command not supported | The ME or action is not implemented on this ONU | ONU firmware lacks support for this ME class or this action on that ME |
| `0x03` | Parameter error | Invalid attribute value or malformed request | Bad VLAN rule, out-of-range attribute value, wrong attribute mask |
| `0x04` | Unknown managed entity | ME class ID not recognized by the ONU | ONU does not implement this ME class; check ONU model and firmware version |
| `0x05` | Unknown managed entity instance | ME instance ID not found | Referenced instance was never created or was previously deleted; MIB may be out of sync |
| `0x06` | Device busy | ONU cannot process the request at this time | Retry after backoff; possible MIB sync in progress or ONU resource contention |
| `0x07` | Instance exists | A Create was issued for an instance that already exists | Stale MIB, duplicate provisioning — consider issuing a MIB Reset and re-syncing |
| `0x08` | Attribute(s) failed or unknown | One or more attributes could not be set or retrieved | Read-only attribute included in Set mask, unsupported attribute bit in mask, value out of range |

> **Note:** Result codes are only present in response messages (AR bit = 0 in the Message Type byte). Autonomous notifications (e.g., Alarm Notification, AVC) do not carry result codes.

---

## Message Type Reference

The following table lists standard OMCI message type values as defined in ITU-T G.988. The **Message Type** byte in an OMCI frame encodes both the action and the AR/AK bits (bits 6–7).

| Decimal | Hex | Message Type Name | Direction |
|---|---|---|---|
| 4 | `0x04` | Create | OLT → ONU |
| 6 | `0x06` | Delete | OLT → ONU |
| 8 | `0x08` | Set | OLT → ONU |
| 9 | `0x09` | Get | OLT → ONU |
| 11 | `0x0B` | Get All Alarms | OLT → ONU |
| 12 | `0x0C` | Get All Alarms Next | OLT → ONU |
| 13 | `0x0D` | MIB Upload | OLT → ONU |
| 14 | `0x0E` | MIB Upload Next | OLT → ONU |
| 15 | `0x0F` | MIB Reset | OLT → ONU |
| 16 | `0x10` | Alarm Notification | ONU → OLT (autonomous) |
| 17 | `0x11` | Attribute Value Change (AVC) | ONU → OLT (autonomous) |
| 18 | `0x12` | Test | OLT → ONU |
| 19 | `0x13` | Start Software Download | OLT → ONU |
| 20 | `0x14` | Download Section | OLT → ONU |
| 21 | `0x15` | End Software Download | OLT → ONU |
| 22 | `0x16` | Activate Software | OLT → ONU |
| 23 | `0x17` | Commit Software | OLT → ONU |
| 24 | `0x18` | Synchronize Time | OLT → ONU |
| 25 | `0x19` | Reboot | OLT → ONU |
| 26 | `0x1A` | Get Next | OLT → ONU |
| 27 | `0x1B` | Test Result | ONU → OLT (autonomous) |
| 28 | `0x1C` | Get Current Data | OLT → ONU |
| 29 | `0x1D` | Set Table | OLT → ONU |

> **AR/AK encoding:** In the raw Message Type byte, bit 6 (AR) = 1 indicates the sender requests a response; bit 7 (AK) = 1 indicates this message is a response. For example, a Create request has MT byte `0x44` (AR=1, MT=4) and the Create response has `0x24` (AK=1, MT=4).

---

## Further Reading

- ITU-T G.988 Section 11 — OMCI message format and result codes
- `knowledge/failure-patterns/README.md` — symptom-to-cause mappings using these codes
- `knowledge/provisioning-flows/standard-onu-provisioning.md` — which result codes to expect at each provisioning step
