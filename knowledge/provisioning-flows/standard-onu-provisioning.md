# Standard ONU Provisioning Sequence

This document describes the typical OMCI provisioning sequence that an OLT uses to bring an ONU into service. The Copilot agent uses this sequence to correlate a failing OMCI message with the provisioning step at which it occurred.

---

## Overview

ONU provisioning is performed over the OMCC (ONU Management and Control Channel) using OMCI messages. The OLT drives the sequence; the ONU responds. The sequence below represents the vendor-neutral baseline defined in ITU-T G.988. Vendor-specific extensions (additional MEs, ordering variations) should be documented in separate files alongside this one.

---

## Provisioning Steps

### Step 1 — MIB Reset

**Action:** OLT sends MIB Reset (MT=15) to ONU-G (ME 256, Instance 0).  
**Purpose:** Clear any stale MIB state before starting a fresh provisioning cycle.  
**Expected result:** `0x00` (success).  
**Failure:** `0x01` (command processing error) — ONU internal fault; retry or reboot ONU.

---

### Step 2 — MIB Upload / MIB Upload Next

**Actions:** MIB Upload (MT=13) followed by repeated MIB Upload Next (MT=14) until the ONU reports no more entries.  
**Purpose:** OLT learns the ONU's current MIB — which MEs and instances already exist (factory defaults, auto-created instances).  
**Key check:** The MIB Upload response includes an ME instance count. The OLT must issue exactly that many MIB Upload Next requests. A mismatch indicates MIB desync.  
**Failure:** Count mismatch → issue another MIB Reset and repeat the upload.

---

### Step 3 — ONU Discovery (Read Mandatory MEs)

The OLT issues Get (MT=9) requests to read key auto-created MEs:

| ME Class ID | ME Name | Purpose |
|---|---|---|
| 256 | ONU-G | ONU identity, vendor ID, equipment ID, OMCI version |
| 257 | ONU2-G | Extended ONU capabilities (QoS, battery backup, etc.) |
| 263 | ANI-G | ANI (PON) port parameters: optical level, upstream FEC, etc. |
| 5 / 6 | Cardholder (ME 5) / Circuit Pack (ME 6) | Physical slot/port inventory |
| 11 | PPTP Ethernet UNI | UNI-side Ethernet port; one instance per physical UNI port |

**Expected result:** `0x00` for each Get.  
**Failure:**
- `0x04` Unknown ME — ONU firmware is very old or non-compliant.
- `0x05` Unknown ME instance — instance was deleted or never created; check MIB upload output.

---

### Step 4 — T-CONT Assignment

**ME:** T-CONT (ME 262).  
**Action:** Set (MT=8) the Alloc-ID attribute on an auto-created T-CONT instance.  
**Purpose:** Associate an ONU T-CONT with an OLT-assigned Alloc-ID for upstream bandwidth allocation.  
**Dependency:** T-CONT instances are auto-created by the ONU. The OLT must read them via MIB upload before setting.  
**Expected result:** `0x00`.  
**Failures:**
- `0x05` Unknown ME instance — T-CONT instance not yet discovered; re-run MIB upload.
- `0x08` Attribute failed — Alloc-ID value out of range or already in use.

---

### Step 5 — GEM Port Network CTP Creation

**ME:** GEM Port Network CTP (ME 268).  
**Action:** Create (MT=4).  
**Key attributes:** Port ID, T-CONT pointer, Direction, Traffic Management Pointer.  
**Dependency:** The referenced T-CONT instance (Step 4) must have a valid Alloc-ID set.  
**Expected result:** `0x00`.  
**Failures:**
- `0x03` Parameter error — Port ID out of range or T-CONT pointer invalid.
- `0x07` Instance exists — duplicate Port ID; issue Delete first or reset MIB.

---

### Step 6 — GEM Interworking TP Creation

**ME:** GEM Interworking TP (ME 266).  
**Action:** Create (MT=4).  
**Key attributes:** GEM Port Network CTP pointer, Interworking option, Service Profile pointer.  
**Dependency:** GEM Port Network CTP instance (Step 5) must exist.  
**Expected result:** `0x00`.  
**Failures:**
- `0x05` Unknown ME instance — GEM Port CTP pointer references a non-existent instance.
- `0x03` Parameter error — invalid interworking option or service profile pointer.

---

### Step 7 — 802.1p Mapper Service Profile Creation

**ME:** 802.1p Mapper Service Profile (ME 130).  
**Action:** Create (MT=4).  
**Purpose:** Map 802.1p priority bits to GEM Interworking TP instances for traffic classification.  
**Dependency:** GEM Interworking TP instances (Step 6) must exist to be referenced.  
**Expected result:** `0x00`.

---

### Step 8 — MAC Bridge Service Profile Creation

**ME:** MAC Bridge Service Profile (ME 45).  
**Action:** Create (MT=4).  
**Purpose:** Defines a MAC bridge instance. Acts as the parent for bridge port configuration data.  
**Expected result:** `0x00`.  
**Failure:** `0x07` Instance exists — MAC Bridge already provisioned from a prior cycle; issue Delete or MIB Reset.

---

### Step 9 — MAC Bridge Port Configuration Data Creation

**ME:** MAC Bridge Port Configuration Data (ME 47).  
**Action:** Create (MT=4) for each logical port (UNI side, ANI/GEM side).  
**Key attributes:** Bridge ID pointer (→ ME 45), Port number, TP type, TP pointer.  
**Dependency:** MAC Bridge Service Profile (Step 8) must exist.  
**Expected result:** `0x00`.  
**Failure:**
- `0x05` Unknown ME instance — Bridge ID pointer references a non-existent ME 45 instance.
- `0x03` Parameter error — invalid TP type / TP pointer combination.

---

### Step 10 — VLAN Configuration

Two alternative MEs depending on ONU capability:

#### Option A — VLAN Tagging Filter Data (ME 84)

**Action:** Create (MT=4) and Set (MT=8).  
**Purpose:** Simple VLAN filtering on a bridge port.  
**Key attributes:** VLAN filter list (up to 12 VIDs), Forward operation, Number of entries.

#### Option B — Extended VLAN Tagging Operation Configuration Data (ME 171)

**Action:** Create (MT=4), then Set Table (MT=29) to add rules.  
**Purpose:** Full VLAN tag manipulation (translate, push, pop, double-tag).  
**Key attributes:** Association type, Association pointer, Input/output VLAN operation tables (16-byte rule entries).  
**Failure:**
- `0x03` Parameter error — malformed 16-byte VLAN rule tuple (check filter/treatment fields).
- `0x02` Command not supported — ONU does not support Set Table; use sequential Set instead.

---

### Step 11 — Multicast GEM Interworking TP (Optional)

**ME:** Multicast GEM Interworking TP (ME 281).  
**Action:** Create (MT=4).  
**Purpose:** Extends GEM Interworking for multicast traffic.  
**Dependency:** GEM Port Network CTP instance must exist and be configured for multicast.  
**Note:** Only required when multicast service is being provisioned.

---

## Dependency Graph

```
ONU-G (256) ─────────────────────────────┐
ANI-G (263)                               │
PPTP Ethernet UNI (11)                    │
                                          ▼
T-CONT (262) ──► GEM Port Network CTP (268)
                         │
                         ▼
               GEM Interworking TP (266)
                         │
                         ▼
            802.1p Mapper Service Profile (130)
                         │
                         ▼
            MAC Bridge Service Profile (45)
                         │
                         ▼
        MAC Bridge Port Configuration Data (47)
                         │
                         ▼
    VLAN Tagging Filter Data (84)
    or Extended VLAN Tagging Op Config Data (171)
                         │
                         ▼
    [Multicast GEM Interworking TP (281)] — optional
```

---

## Order-of-Operations Notes

1. **Always MIB Reset before reprovisioning.** Do not attempt to re-create MEs without first resetting to avoid `0x07 Instance exists` errors.
2. **Read before writing.** Use MIB Upload to discover auto-created instances (T-CONT, ANI-G, PPTP) before referencing them.
3. **Parent before child.** Create parent MEs (MAC Bridge, GEM Port CTP, T-CONT) before creating any ME that holds a pointer to them.
4. **Pointer validity.** Every pointer attribute must reference an existing instance. A dangling pointer causes `0x05 Unknown ME instance`.
5. **Attribute mask precision.** In a Set message, only include attributes that are writable for that ME. Including read-only attributes causes `0x08 Attribute failed`.

---

## Vendor-Specific Extensions

Vendor-specific provisioning steps (e.g., additional ME classes, proprietary ordering requirements) should be documented in separate files within this directory, named `<vendor>-onu-provisioning.md`. Keep this file vendor-neutral.

---

## Further Reading

- `knowledge/result-codes/README.md` — full result code reference
- `knowledge/failure-patterns/README.md` — symptom-to-cause table
- `knowledge/me-catalog/` — JSON ME definitions for attribute-level detail
- ITU-T G.988 Section 9 — ME definitions and relationships
