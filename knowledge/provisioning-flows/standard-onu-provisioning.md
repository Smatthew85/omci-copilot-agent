# Standard ONU OMCI Provisioning Flow

This document describes the standard ONU provisioning sequence using OMCI as defined by
ITU-T G.988. Use this as a reference when diagnosing where in the sequence a failure
occurred.

---

## Overview

A typical GPON/XGS-PON ONU provisioning sequence proceeds in these phases:

1. **PLOAM Activation** (outside OMCI scope)
2. **MIB Synchronization** — Reset + Upload
3. **ONU Capability Discovery** — Read ONU-G, ONU2-G, ANI-G, etc.
4. **T-CONT and GEM Port Configuration**
5. **Bridge and UNI Configuration**
6. **VLAN Configuration**
7. **Service Activation** — final Set/Create to enable traffic

---

## Phase 1: PLOAM Activation

Not OMCI. The OLT and ONU exchange PLOAM messages to establish synchronization,
authentication (SERIAL_NUM_ONU, RANGING_TIME), and grant the ONU a valid ONU-ID.
OMCI cannot proceed until PLOAM activation completes.

---

## Phase 2: MIB Synchronization

```
OLT -> ONU: MIB Reset (MT 15)          -- clear ONU MIB to known state
ONU -> OLT: MIB Reset response         -- result 0x00 = success
OLT -> ONU: MIB Upload (MT 13)         -- request MIB snapshot; ONU returns count N
ONU -> OLT: MIB Upload response        -- content = number of MIB Upload Next messages needed
OLT -> ONU: MIB Upload Next (MT 14, seq=0)
ONU -> OLT: MIB Upload Next response   -- ME instances in content
... (repeat until all N instances returned)
```

**Failure modes**:
- MIB Reset fails: ONU is busy or has a firmware issue.
- MIB Upload count mismatch: restart from MIB Upload.

---

## Phase 3: ONU Capability Discovery

The OLT reads pre-existing MEs created by the ONU autonomously:

- **ONU-G** (ME 256) — ONU identity, vendor ID, version, serial number
- **ONU2-G** (ME 257) — extended ONU capabilities
- **ANI-G** (ME 263) — optical line parameters per PON port
- **UNI-G** (ME 264) — UNI port descriptor
- **PPTP Ethernet UNI** (ME 11) — Ethernet UNI physical port
- **T-CONT** (ME 262) — traffic containers (pre-created by ONU)
- **Priority Queue** (ME 277) — queues associated with T-CONTs and UNIs

These are read-only or SetByCreate from the OLT's perspective at this phase.

---

## Phase 4: T-CONT and GEM Port Configuration

```
OLT -> ONU: Set T-CONT (ME 262) alloc_id = <assigned alloc-ID>
OLT -> ONU: Create GEM Port Network CTP (ME 268) for each GEM port
OLT -> ONU: Create GEM Interworking TP (ME 266) linking GEM port to bridge/T-CONT
```

**ME dependencies**: GEM Port Network CTP → T-CONT (via Alloc-ID);
GEM Interworking TP → GEM Port Network CTP + IW TP (bridge or 802.1p mapper).

---

## Phase 5: Bridge and UNI Configuration

```
OLT -> ONU: Create MAC Bridge Service Profile (ME 45)
OLT -> ONU: Create MAC Bridge Port Config Data (ME 47) for each bridge port
             (one port per GEM Interworking TP, one port per UNI)
OLT -> ONU: Create 802.1p Mapper Service Profile (ME 130) if 802.1p mapping needed
OLT -> ONU: Set PPTP Ethernet UNI (ME 11) admin_state = unlocked
```

**Common failure**: Create ME 45 or ME 47 returns `0x07 Instance Exists` → see
[`knowledge/examples/01-stale-mib-instance-exists/`](../examples/01-stale-mib-instance-exists/).

---

## Phase 6: VLAN Configuration

```
OLT -> ONU: Create/Set Extended VLAN Tagging Operation Config Data (ME 171)
             -- one instance per PPTP Ethernet UNI or VEIP
OLT -> ONU: Set ME 171 attribute 7 -- add VLAN rule table entries (16-byte tuples)
OLT -> ONU: Create VLAN Tagging Filter Data (ME 84) if upstream filtering needed
```

**Common failure**: Set ME 171 returns `0x03 Parameter Error` → see
[`knowledge/examples/02-extended-vlan-parameter-error/`](../examples/02-extended-vlan-parameter-error/).

---

## Phase 7: Service Activation

```
OLT -> ONU: Set ONU-G (ME 256) admin_state = 0 (unlock)
OLT -> ONU: Set ANI-G (ME 263) as needed
-- Traffic can now flow
```

---

## ME Dependency Summary

```
T-CONT (262)
  └── GEM Port Network CTP (268)
        └── GEM Interworking TP (266)
              └── MAC Bridge Port Config Data (47)
                    └── MAC Bridge Service Profile (45)
                          └── Extended VLAN Tagging Op Config Data (171)

PPTP Ethernet UNI (11)
  └── MAC Bridge Port Config Data (47)
        └── (same bridge as above)

802.1p Mapper Service Profile (130)
  └── GEM Interworking TP (266)  [alternative to direct bridge port]
```
