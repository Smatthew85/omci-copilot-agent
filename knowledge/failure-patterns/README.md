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
# OMCI Provisioning Failure Patterns

Symptom → likely cause → suggested action mappings for common OMCI provisioning failures. Use these patterns alongside `knowledge/result-codes/README.md` and the ME catalog to diagnose issues quickly.

---

## Failure Pattern Table

| Symptom | Likely Cause | Suggested Action |
|---|---|---|
| Create returns `0x07` Instance exists | Stale MIB or duplicate provisioning — the ONU already has an instance with this class + instance ID | Issue MIB Reset (MT=15), re-sync via MIB Upload, then retry Create |
| Set returns `0x08` Attribute(s) failed or unknown | Read-only attribute included in the Set mask, or the attribute is not supported by this ONU | Verify attribute access rules (R/W/SetByCreate) in the ME catalog JSON; remove read-only attributes from the mask |
| Create returns `0x04` Unknown managed entity | ONU firmware does not implement this ME class | Check ONU model/firmware version; use an alternative ME if available (e.g., ME 84 instead of ME 171) |
| Create returns `0x02` Command not supported | The Create action is not in the supported actions list for this ME on this ONU | Verify supported actions in the ME catalog; consult ONU vendor documentation |
| `0x03` Parameter error on Extended VLAN Tagging Operation Config Data (ME 171) | Malformed 16-byte VLAN rule tuple — invalid filter or treatment field values | Validate the 16-byte rule structure per G.988 Table 9.3.13-1; check filter outer/inner VLAN, treatment tags |
| MIB Upload count mismatch (OLT receives fewer/more entries than the count in the MIB Upload response) | ONU/OLT MIB desync — ONU state changed during upload | Issue MIB Reset + full MIB re-upload; ensure no parallel provisioning is in progress |
| TCID timeout / no response from ONU | Transport-layer issue — OMCC channel problem, not OMCI protocol logic | Check OMCC GEM port and T-CONT configuration at the transport layer; verify ONU is registered and optical signal is present |
| GEM Interworking TP (ME 266) Create fails after GEM Port Network CTP (ME 268) Create succeeds | Missing or incorrect T-CONT pointer, or GEM Port CTP instance pointer in ME 266 references wrong instance | Verify the GEM Port CTP pointer attribute in ME 266 matches the just-created ME 268 instance ID |
| MAC Bridge Port Configuration Data (ME 47) Create fails with `0x05` Unknown ME instance | MAC Bridge Service Profile (ME 45) was not created, or the Bridge ID pointer references the wrong instance | Ensure ME 45 is created first (Step 8 of the provisioning sequence); confirm the bridge instance ID |
| Set on T-CONT (ME 262) returns `0x03` Parameter error | Alloc-ID value is out of the valid range or is already assigned to another T-CONT | Verify the Alloc-ID assigned by the OLT is within the valid range and not duplicated |
| `0x05` Unknown ME instance on any pointer attribute | The referenced ME instance was never created, was deleted, or was not yet discovered via MIB Upload | Re-run MIB Upload to discover existing instances; create the parent ME before the child |
| `0x06` Device busy | ONU is occupied (MIB sync in progress, software download, etc.) | Back off and retry; do not send further commands until the ONU becomes available |
| Create returns `0x01` Command processing error | ONU internal error — resource exhaustion, hardware fault | Check ONU logs/alarms; retry; if persistent, reboot ONU and re-provision |
| Alarm Notification received unexpectedly during provisioning | ONU detected a fault condition (LOS, SF, SD, etc.) | Investigate physical layer before continuing OMCI provisioning |
| VLAN Tagging Filter Data (ME 84) Set returns `0x03` | Number of entries field does not match the populated entries in the VLAN filter list | Ensure the `Number of entries` attribute is set to the exact count of valid VIDs in the filter list |

---

## General Diagnostic Checklist

1. **Identify the ME class and instance** from the failing message (Class ID + Instance ID).
2. **Look up the ME** in `knowledge/me-catalog/` to confirm supported actions and attribute access rules.
3. **Map the result code** using `knowledge/result-codes/README.md`.
4. **Locate the step** in `knowledge/provisioning-flows/standard-onu-provisioning.md` to check ordering and dependencies.
5. **Check parent MEs** — ensure all MEs referenced by pointer attributes were created in prior steps.
6. **Consider MIB state** — when in doubt, MIB Reset and re-upload is the safest recovery path.

---

## Vendor-Specific Patterns

Vendor-specific failure patterns (e.g., proprietary ME extensions, ONU-model-specific quirks) should be documented in separate files alongside this one, named `<vendor>-failure-patterns.md`. Keep this file vendor-neutral.

---

## Further Reading

- `knowledge/result-codes/README.md` — full result and reason code reference
- `knowledge/provisioning-flows/standard-onu-provisioning.md` — standard provisioning sequence with per-step failure notes
- `knowledge/me-catalog/` — JSON ME definitions for attribute-level detail
