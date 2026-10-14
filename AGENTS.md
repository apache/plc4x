# Agent guidance

This file is read by automated agents (security scanners, code analyzers,
AI assistants) operating on this repository. It points them at the
human-authored references they should consult before producing output.

## Security

Security model: [SECURITY.md](./SECURITY.md), which links to the project's
threat model at [THREAT-MODEL.md](THREAT-MODEL.md).

Agents that scan this repository should consult `THREAT-MODEL.md` for
the project's in-scope / out-of-scope declarations, adversary model, security
properties, and known non-findings before reporting issues. Apache PLC4X is a
client library for industrial protocols (most of which are unauthenticated /
unencrypted by design); the threat model is about which threats the library
takes on versus which are left to the operator and the OT network.

## CI workflows and Dependabot

GitHub Actions runners are shared by every ASF project. Agents adding or
editing anything under `.github/` must follow Infra's
[GitHub Actions policy](https://infra.apache.org/github-actions-policy.html)
(mandatory) and its
[guidance on reducing queue load](https://infra.apache.org/services.html#github-actions).
In this repository that means:

- Platform workflows that run on pull requests carry the concurrency group
  already used in `java-platform.yml`, so a new push to a PR cancels its
  earlier runs. Pushes to `develop` are never cancelled:

  ```yaml
  concurrency:
    group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
    cancel-in-progress: ${{ github.event_name == 'pull_request' }}
  ```

- Short jobs that need no Docker run on `ubuntu-slim` (1 vCPU, 15 minute
  limit), not `ubuntu-latest`.
- Each Dependabot ecosystem in `.github/dependabot.yml` groups both version
  updates and `applies-to: security-updates`, so each ecosystem opens one PR
  instead of one per dependency.
- Actions outside `apache/*`, `github/*` and `actions/*` are pinned to a full
  commit SHA and must be on Infra's approved list
  ([apache/infrastructure-actions](https://github.com/apache/infrastructure-actions)).
- No `pull_request_target`, and no matrix above 20 concurrent jobs per
  workflow (15 or fewer recommended).
