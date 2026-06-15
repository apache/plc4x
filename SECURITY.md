# Security Policy

## Reporting a Vulnerability

Please report suspected security vulnerabilities in Apache PLC4X privately to
the Apache Security Team at <security@apache.org>, following the ASF process
at <https://www.apache.org/security/>. Do not open public GitHub issues or
pull requests for security reports.

## Threat Model

Apache PLC4X's security threat model — what is in and out of scope, the
security properties the project provides and disclaims, the adversary model,
the environmental assumptions, and how findings are triaged — is documented in
[draft-THREAT-MODEL.md](./draft-THREAT-MODEL.md).

PLC4X speaks industrial protocols (Modbus, S7, OPC-UA, ADS, EtherNet/IP, …),
most of which are unauthenticated and unencrypted by design. The threat model
covers the parser/driver trust boundary (responses from the device/wire) and
draws the line on what is the operator's and the OT network's responsibility.
