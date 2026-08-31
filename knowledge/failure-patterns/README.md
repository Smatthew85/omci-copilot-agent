# OMCI Failure Patterns

This document maps observed OMCI failure symptoms to likely root causes and remediation
actions. Use this as a quick-reference when analyzing OMCI captures.

---

## Quick Reference Table

| Symptom | Likely Root Cause | Remediation |
|---------|-------------------|-------------|
| `Create` → `0x07 Instance Exists` | Stale MIB state from prior provisioning; no MIB Reset performed | MIB Reset (MT 15) + MIB Upload (MT 13/14), then reissue Create |
| `Set` → `0x08 Attribute(s) Failed` | Attempting to Set a read-only attribute, or attribute value out of range | Check attribute access (R/W/SetByCreate) in ME catalog; correct the value |
| `Create` → `0x04 Unknown ME` | ONU firmware does not implement this ME Class | Verify ONU capability; use a different ME or skip this step for this ONU |
| `Set` → `0x03 Parameter Error` on ME 171 | Malformed 16-byte VLAN rule tuple (invalid priority/TPID/treatment fields) | Validate rule encoding per G.988 clause 9.3.11; fix field values |
| `Get` → `0x05 Unknown Instance` | Attempting to Get an ME instance that does not exist | Check ONU MIB with MIB Upload; create the instance first if required |
| MIB Upload count mismatch | OLT and ONU MIB state diverged (e.g., after ONU reboot without re-sync) | MIB Reset (MT 15) + full MIB Upload (MT 13/14) to resync |
| No response / TCID timeout | Transport-layer issue (ONU offline, PLOAM not established, queue congestion) | Verify PLOAM / activation state before diagnosing OMCI messages |
| `Create` → `0x06 Device Busy` | ONU is processing a prior operation or performing internal maintenance | Retry after a short delay; if persistent, check ONU health |
| `Set` → `0x02 Command Not Supported` | ONU firmware does not support this action for the ME | Check supported actions in ME catalog; use an alternative approach |
| `Delete` → `0x05 Unknown Instance` | Attempting to delete an ME instance that is already absent | Verify MIB state; the instance may have been deleted in a prior cycle |
| `Alarm Notification` (MT 16) with unexpected alarm | ONU hardware or link condition has changed | Correlate with physical layer (optical power, BER) and PLOAM state |

---

## Detailed Patterns

### Pattern 1: Create Returns 0x07 Instance Exists

**Symptom**: OLT sends a Create request; ONU responds with result code `0x07`.

**Causes**:
1. The ONU MIB retained the instance from a previous provisioning cycle (e.g., after
   ONU reboot without OLT-side MIB re-sync).
2. Duplicate Create issued by the OLT (bug in provisioning state machine).

**Remediation**:
- Issue MIB Reset (MT 15) to clear ONU MIB, then MIB Upload (MT 13 + MT 14) to resync,
  then reissue Create.
- Alternatively, issue Delete (MT 11) for the specific instance, then reissue Create.

**Reference**: [`knowledge/examples/01-stale-mib-instance-exists/`](../examples/01-stale-mib-instance-exists/)

---

### Pattern 2: Set Returns 0x03 Parameter Error

**Symptom**: OLT sends a Set request; ONU responds with result code `0x03`.

**Causes**:
1. An attribute value falls outside the range specified in G.988 for that ME.
2. A structured attribute (e.g., VLAN rule table entry) has fields at invalid offsets or
   with out-of-range values.
3. Conflicting attribute combination (two attributes whose values are mutually exclusive).

**Remediation**:
- Identify the failing attribute from the request's attribute mask.
- Validate the attribute value against the G.988 ME definition.
- For ME 171 (Extended VLAN Tagging Op Config Data), validate each 16-byte rule tuple
  field-by-field per G.988 clause 9.3.11.

**Reference**: [`knowledge/examples/02-extended-vlan-parameter-error/`](../examples/02-extended-vlan-parameter-error/)

---

### Pattern 3: MIB Upload Count Mismatch

**Symptom**: The count returned by MIB Upload (MT 13) does not match the number of ME
instances enumerated by subsequent MIB Upload Next (MT 14) exchanges.

**Causes**:
1. ONU MIB changed between the MIB Upload command and the MIB Upload Next sequence
   (race condition).
2. ONU firmware bug causing incorrect count reporting.

**Remediation**:
- Restart the MIB Upload sequence from MT 13.
- If the mismatch persists, perform MIB Reset (MT 15) first.

---

### Pattern 4: No Response / TCID Timeout

**Symptom**: The OLT transmits an OMCI request but receives no response within the
configured TCID timeout window.

**Causes**:
1. ONU is offline (power loss, fiber cut, PLOAM not established).
2. ONU OMCI message queue is full (ONU is overloaded).
3. Transport-layer issue (GEM port misconfigured for OMCI channel).

**Remediation**:
- Verify PLOAM activation state before diagnosing at the OMCI layer.
- Check GEM port configuration for the OMCI channel (GEM Port Network CTP pointing to
  the ONU-G managed entity).
- Retry the OMCI request; if the problem persists, investigate the physical layer.
