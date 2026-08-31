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
