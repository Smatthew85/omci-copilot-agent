# Context: Set Extended VLAN Tagging Op Config Data — Parameter Error mid-provisioning

## Provisioning Step

This failure occurred mid-provisioning, after MAC Bridge Service Profile (ME 45) and
MAC Bridge Port Configuration Data (ME 47) were successfully created. The OLT was
configuring VLAN translation rules as the next step in the service activation sequence.

## Preceding Operations

The preceding operations completed without error:

```
# Anonymized operation sequence:
# INFO   omci-cc: tx mt=Create  meClass=45  meInst=1  -> result=0x00 (Success)
# INFO   omci-cc: tx mt=Create  meClass=47  meInst=1  -> result=0x00 (Success)
# INFO   omci-cc: tx mt=Set     meClass=171 meInst=1  attrMask=0x0100
# WARN   omci-cc: rx mt=Set     meClass=171 meInst=1  result=0x03 (Parameter Error)
# ERROR  onu-state: VLAN rule configuration failed
```

The ME 171 instance itself was pre-existing (created by the ONU autonomously per G.988
or set by an earlier Create action not shown here).

## Environment

- **OLT vendor/version**: vendor-neutral (encoding issue originates in OLT provisioning
  software that constructs the 16-byte rule tuple)
- **ONU vendor/model**: vendor-neutral (ONU correctly rejects the malformed tuple)
- **OMCI frame type**: Baseline (48 bytes), Device Identifier 0x0A
- **Additional notes**: This failure pattern is commonly triggered when OLT provisioning
  software auto-generates rule tuples from higher-level VLAN policy objects without
  fully validating each field against G.988 allowed ranges. The bug may be latent until
  a specific VLAN policy combination exercises the out-of-range code path.
