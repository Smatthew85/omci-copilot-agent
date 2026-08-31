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
| 4   | 0x04 | Create Complete   | ONU->OLT    | No           |
| 5   | 0x05 | Delete            | OLT->ONU    | Yes          |
| 6   | 0x06 | Set               | OLT->ONU    | Yes          |
| 7   | 0x07 | Get               | OLT->ONU    | Yes          |
| 8   | 0x08 | Get All Alarms    | OLT->ONU    | Yes          |
| 9   | 0x09 | Get All Alarms Next | OLT->ONU  | Yes          |
| 10  | 0x0A | Create            | OLT->ONU    | Yes          |
| 11  | 0x0B | Delete            | OLT->ONU    | Yes          |
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
