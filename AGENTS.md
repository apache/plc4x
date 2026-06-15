# Agent guidance

This file is read by automated agents (security scanners, code analyzers,
AI assistants) operating on this repository. It points them at the
human-authored references they should consult before producing output.

## Security

Security model: [SECURITY.md](./SECURITY.md), which links to the project's
threat model at [draft-THREAT-MODEL.md](./draft-THREAT-MODEL.md).

Agents that scan this repository should consult `draft-THREAT-MODEL.md` for
the project's in-scope / out-of-scope declarations, adversary model, security
properties, and known non-findings before reporting issues. Apache PLC4X is a
client library for industrial protocols (most of which are unauthenticated /
unencrypted by design); the threat model is about which threats the library
takes on versus which are left to the operator and the OT network.
