# Alarm Synchronization

> **Source:** ITU-T G.988 §11.3.3, §11.3.17, §11.3.18

This document covers how the OLT detects missed alarm notifications and how to
resynchronize alarm state using **Get All Alarms** (MT 27) and **Get All Alarms Next**
(MT 28).

---

## 1. Sequence-Number Gap Detection

Every baseline alarm notification (MT 16) carries a 1-byte **alarm sequence number**
(byte 40 of the 48-byte frame) that the ONU increments for every alarm it sends,
wrapping from `0xFF` back to `0x00`.

The OLT tracks the last-seen sequence number per ONU. A gap indicates missed frames:

```
Expected next = (last_seen + 1) mod 256
If received_seq ≠ expected_next → gap detected
```

**Common causes of gaps:**

| Cause | Notes |
|---|---|
| PON upstream burst loss | Transient optical impairment; retransmit not possible (autonomous, no AR) |
| ONU reboot | Sequence number reset; OLT may also see the ONU re-register |
| OLT processing overload | Frame received at PON MAC but dropped before delivery to OMCI stack |

**OLT action on gap detection:** Initiate the Get All Alarms resync procedure
(§3 below) to rebuild a consistent alarm state from scratch.

---

## 2. Get All Alarms — MT 27

The OLT sends a **Get All Alarms** request to instruct the ONU to report all currently
active alarms. The ONU responds with the number of Get All Alarms Next commands required
to retrieve the full alarm set.

### Request (OLT → ONU)

```
Byte(s)  Field              Value
──────────────────────────────────────────────────────
0–1      TCID               OLT-assigned, non-zero
2        MT                 0x5B  (AR=1, AK=0, MT=27)
3        DevID              0x0A
4–5      ME Class           0x0002  (ONU Data, ME 2)
6–7      ME Instance        0x0000
8        Alarm Retrieval    0x00 = do not reset alarm sequence number
         Mode               0x01 = reset alarm sequence number to 0
──────── padding ──────────────────────────────────────
9–39     0x00…
──────── trailer ──────────────────────────────────────
```

### Response (ONU → OLT)

```
Byte(s)  Field              Value
──────────────────────────────────────────────────────
0–1      TCID               Same as request
2        MT                 0x3B  (AR=0, AK=1, MT=27)
3        DevID              0x0A
4–5      ME Class           0x0002
6–7      ME Instance        0x0000
8–9      Number of commands 2-byte count of Get All Alarms Next commands
                            needed to retrieve all active alarms
──────── padding ──────────────────────────────────────
10–39    0x00…
──────── trailer ──────────────────────────────────────
```

**Retrieval mode:**
- `0x00` — retrieve alarms without disturbing the ongoing sequence number.
- `0x01` — reset the ONU's alarm sequence number to 0 before responding; use this to
  definitively resynchronize after a detected gap.

---

## 3. Get All Alarms Next — MT 28

The OLT issues one or more **Get All Alarms Next** requests, using a sequential
**Command Sequence Number** starting at 0, to retrieve each active alarm.

### Request (OLT → ONU)

```
Byte(s)  Field                        Value
──────────────────────────────────────────────────
0–1      TCID                         OLT-assigned
2        MT                           0x5C  (AR=1, AK=0, MT=28)
3        DevID                        0x0A
4–5      ME Class                     0x0002
6–7      ME Instance                  0x0000
8–9      Command Sequence Number      0-based index of this request
──────── padding ───────────────────────────────────
10–39    0x00…
──────── trailer ───────────────────────────────────
```

### Response (ONU → OLT)

```
Byte(s)  Field              Value
──────────────────────────────────────────────────
0–1      TCID               Same as request
2        MT                 0x3C  (AR=0, AK=1, MT=28)
3        DevID              0x0A
4–5      ME Class           0x0002
6–7      ME Instance        0x0000
8–9      Reported ME Class  ME Class ID of the ME with an active alarm
10–11    Reported ME Inst   ME Instance ID
12–39    Alarm bitmap       28 bytes; same semantics as in MT 16
──────── trailer ──────────────────────────────────
```

Each response contains one ME instance with its current full alarm bitmap. If no alarms
are active for the ONU, the ONU responds with Command Sequence Number 0 and all-zero
ME Class and Instance fields.

---

## 4. Resync Workflow

```
OLT detects sequence-number gap (or ONU re-registration)
  │
  ▼
Send Get All Alarms request (MT 27, Retrieval Mode = 0x01)
  │
  ▼
Receive response: Number of commands = N
  │
  ├─ N = 0 → No active alarms; alarm state cleared; done.
  │
  └─ N > 0 → for seq = 0 to N-1:
               Send Get All Alarms Next (MT 28, Seq = seq)
               Receive response: ME Class, ME Instance, alarm bitmap
               Update OLT alarm table for that ME instance
             Done → alarm table is now consistent.
```

**Expected result codes:**

| Step | Expected Result |
|---|---|
| MT 27 response | No result code field; success indicated by valid `Number of commands` |
| MT 28 response | No result code field; success indicated by non-zero ME Class in content |

---

## 5. Common Failure Modes

| Symptom | Likely Cause | Recommended Action |
|---|---|---|
| ONU returns `Number of commands = 0` but alarms are known active | ONU MIB inconsistency or reboot cleared alarm state | Accept as ONU ground truth; update OLT table accordingly |
| MT 28 response ME Class = 0x0000 for Seq 0 | No alarms active despite non-zero N in MT 27 response | Log discrepancy; treat as no-alarm state |
| Sequence number gap persists after resync | Ongoing upstream burst loss | Check optical link; repeat resync after stabilization |
| MT 27 or MT 28 times out | ONU unresponsive or OMCI channel disrupted | Check OMCI channel connectivity; attempt MIB Reset if needed |

---

## 6. Worked Example: 3-Frame Resync Sequence

Scenario: ONU has two active alarms — ANI-G (Class 263) instance 0 alarm 2, and
PPTP Ethernet UNI (Class 11) instance 1 alarm 0.

**Frame 1: Get All Alarms Request (OLT → ONU)**

```
Hex: 0042 5B 0A 0002 0000 01 00 00 00 00 00 00 00
     00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
     00 00 00 00 00 28 00 00 AABBCCDD

  TCID = 0x0042, MT = 0x5B (Get All Alarms, AR=1),
  ME 2 inst 0, Retrieval Mode = 0x01 (reset seq number)
```

**Frame 2: Get All Alarms Response (ONU → OLT)**

```
Hex: 0042 3B 0A 0002 0000 00 02 00 00 00 00 00 00
     00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
     00 00 00 00 00 28 00 00 AABBCCDD

  TCID = 0x0042, MT = 0x3B (Get All Alarms response, AK=1)
  Number of commands = 0x0002 → 2 × Get All Alarms Next to follow
```

**Frame 3: Get All Alarms Next Request Seq=0 (OLT → ONU)**

```
Hex: 0043 5C 0A 0002 0000 00 00 00 00 00 00 00 00
     00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
     00 00 00 00 00 28 00 00 AABBCCDD

  MT = 0x5C, Seq = 0x0000
```

**Frame 4: Get All Alarms Next Response Seq=0 (ONU → OLT)**

```
Hex: 0043 3C 0A 0002 0000 01 07 00 00 04 00 00 00
     00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
     00 00 00 00 00 28 00 00 AABBCCDD

  Reported ME Class = 0x0107 = 263 (ANI-G)
  Reported ME Inst  = 0x0000
  Alarm bitmap byte 0 = 0x04 → bit 5 set → alarm 2 active (Signal fail)
```

**Frame 5: Get All Alarms Next Request Seq=1 (OLT → ONU)**

(Same structure as Frame 3 with Seq = 0x0001)

**Frame 6: Get All Alarms Next Response Seq=1 (ONU → OLT)**

```
  Reported ME Class = 0x000B = 11 (PPTP Ethernet UNI)
  Reported ME Inst  = 0x0001
  Alarm bitmap byte 0 = 0x80 → bit 7 set → alarm 0 active (LAN LOS)
```

OLT now has a consistent alarm table: ANI-G/0 alarm 2 active, PPTP Eth UNI/1 alarm 0
active. Any subsequent MT 16 messages update incremental state from this baseline.

---

## References

- ITU-T G.988 §11.3.3 — Alarm notification
- ITU-T G.988 §11.3.17 — Get All Alarms
- ITU-T G.988 §11.3.18 — Get All Alarms Next
- [`alarm-message-format.md`](alarm-message-format.md)
- [`common-alarms.md`](common-alarms.md)
