# Implementation Plan: Refactor Infrastructure into API and Update EKS Deployment

## Phase 1: Research and Analysis [checkpoint: bdedd19]
- [x] Task: Analyze current EKS deployment workflow and secret retrieval. [3f0d21a]
- [x] Task: Review external infrastructure repository configurations for secret naming conventions. [3f0d21a]
- [x] Task: Conductor - User Manual Verification 'Phase 1: Research and Analysis' (Protocol in workflow.md) [bdedd19]

## Phase 2: Cleanup and Removal [checkpoint: eb2509b]
- [x] Task: Remove the `infra/` directory from the project. [5c87597]
- [x] Task: Delete the `.github/workflows/build-infra.yml` file. [a06d927]
- [x] Task: Conductor - User Manual Verification 'Phase 2: Cleanup and Removal' (Protocol in workflow.md) [eb2509b]

## Phase 3: Workflow and Secret Configuration [checkpoint: 695ef7a]
- [x] Task: Update `.github/workflows/deploy.yml` to include Lab Role ARN and secret sync. [04895e8]
- [x] Task: Configure the application to prioritize environment variables (from GitHub Secrets) in production. [04895e8]
- [x] Task: Verify that local `.env` loading remains the fallback for development. [04895e8]
- [x] Task: Conductor - User Manual Verification 'Phase 3: Workflow and Secret Configuration' (Protocol in workflow.md) [695ef7a]

## Phase 4: Validation and Testing [checkpoint: e9ec9c7]
- [x] Task: Perform a dry run of the updated deployment pipeline. [695ef7a]
- [x] Task: Verify that the application can connect to the database in a simulated production environment. [695ef7a]
- [x] Task: Ensure local development environment still works correctly. [695ef7a]
- [x] Task: Conductor - User Manual Verification 'Phase 4: Validation and Testing' (Protocol in workflow.md) [e9ec9c7]
