# Context: Create MAC Bridge Service Profile — Instance Exists after ONU Reboot

## Provisioning Step

This failure occurred during the initial service profile creation phase, immediately after
the OLT detected that the ONU had come back online following an unplanned reboot. The OLT
attempted to resume provisioning from its cached MIB state without first performing a MIB
Reset, expecting the ONU MIB to be empty.

## Preceding Operations

The ONU rebooted (power cycle or firmware reset) while provisioned. The OLT's OMCI state
machine detected the ONU coming back online and proceeded directly to service profile
creation, skipping the MIB Reset / MIB Upload sequence.

```
# Anonymized event sequence:
# EVENT  onu-state: ONU detected online  onu=<anonymized>
# INFO   omci-cc: skipping MIB reset (cached MIB assumed valid)
# INFO   omci-cc: tx OMCI  tcid=0x0001 mt=Create  meClass=45 meInst=1
# WARN   omci-cc: rx OMCI  tcid=0x0001 mt=Create  meClass=45 meInst=1 result=0x07
# ERROR  onu-state: provisioning failed — Instance Exists
```

## Environment

- **OLT vendor/version**: vendor-neutral (behavior is standard per G.988)
- **ONU vendor/model**: vendor-neutral
- **OMCI frame type**: Baseline (48 bytes), Device Identifier 0x0A
- **Additional notes**: The ONU retained its MIB across the reboot (non-volatile storage).
  The OLT's MIB cache was stale because no MIB re-sync was triggered on ONU re-discovery.
  This pattern is common when OLT software optimizes for fast re-provisioning but does not
  account for ONU-side MIB persistence.
