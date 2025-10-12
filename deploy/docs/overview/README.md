# Deployment & Operations Overview

## Responsibility
- Orchestrate continuous delivery for backend, frontend, and mobile artifacts with environment parity.
- Maintain observability, incident response procedures, and operational automation.

## High-Level Architecture
- GitHub Actions triggers build pipelines that publish containers to an internal registry and create release candidates.
- ArgoCD promotes workloads across staging and production clusters with progressive delivery gates.
- Terraform manages cloud infrastructure, while Crossplane provisions managed services declaratively.

## Conventions
- State management for configuration uses GitOps repositories with locked main branches and promotion PRs.
- Super comments using `//1.-`, `//2.-` belong in Terraform modules or pipeline scripts explaining sequential automation steps.
- Infrastructure or pipeline files should not exceed 350 lines; break apart modules and reusable actions proactively.
