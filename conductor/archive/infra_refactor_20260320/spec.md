# Specification: Refactor Infrastructure into API and Update EKS Deployment

## Overview
This track aims to simplify the project's infrastructure by removing the internal `@infra/` folder and integrating with external infrastructure repositories (`iac-tech-challenge-infra`, `iac-tech-challenge-gateway`, `iac-tech-challenge-data`). The deployment process to EKS will be updated to assume that the required environment is already provisioned, focusing on syncing secrets and updating the GitHub Action workflows for a more consistent and clean setup.

## Functional Requirements
- **Infrastructure Removal:** Delete the existing `infra/` directory.
- **Secret Management (Production):** Sync database passwords and service tokens (e.g., `MAILTRAP_TOKEN`) from GitHub Secrets.
- **Secret Management (Local):** Ensure local development continues to use `.env` files for configuration.
- **Pipeline Update:** Update `.github/workflows/deploy.yml` to support the new configuration and sync with external repository values.
- **Workflow Cleanup:** Remove the obsolete `.github/workflows/build-infra.yml` pipeline.
- **Lab Role Integration:** Ensure the AWS Lab Role ARN is correctly considered in the configuration for EKS deployment.

## Non-Functional Requirements
- **Consistency:** Maintain a clean and simple setup consistent with other repositories in the organization.
- **Security:** Avoid exposing secrets in the codebase; use secure environment variables and GitHub Secrets.
- **Maintainability:** Document the new deployment process clearly.

## Acceptance Criteria
- [ ] `infra/` directory is removed.
- [ ] `.github/workflows/build-infra.yml` is deleted.
- [ ] `.github/workflows/deploy.yml` is updated and successfully deploys to EKS.
- [ ] Application correctly retrieves secrets from GitHub Secrets in production.
- [ ] Local development setup using `.env` remains functional.
- [ ] Lab Role ARN is used in the EKS deployment configuration.

## Out of Scope
- Modifying the external infrastructure repositories.
- Changing the application's core logic or Hexagonal Architecture.
