# Context: [Short Title]

## Provisioning Step

_Describe where in the standard ONU provisioning sequence this failure occurred (e.g., "During initial service profile creation, after MIB Upload Next completed")._

## Preceding Operations

_List the OMCI operations (or VOLTHA state machine transitions) that preceded the failing frame. Paste anonymized log excerpts if available._

```
# Example VOLTHA log excerpt (anonymized):
# INFO  omci-cc: tx OMCI msg  tcid=0x0001 mt=Create meClass=45 meInst=1
# WARN  omci-cc: rx OMCI msg  tcid=0x0001 mt=Create meClass=45 meInst=1 result=0x07
```

## Environment

- **OLT vendor/version**: _replace with actual values or leave as "vendor-neutral"_
- **ONU vendor/model**: _replace with actual values or leave as "vendor-neutral"_
- **OMCI baseline or extended frame**: _Baseline (48 bytes) / Extended_
- **Additional notes**: _any other relevant environmental context_
