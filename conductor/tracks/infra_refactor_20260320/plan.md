# Implementation Plan: Refactor Infrastructure into API and Update EKS Deployment

## Phase 1: Research and Analysis
- [ ] Task: Analyze current EKS deployment workflow and secret retrieval.
- [ ] Task: Review external infrastructure repository configurations for secret naming conventions.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Research and Analysis' (Protocol in workflow.md)

## Phase 2: Cleanup and Removal
- [ ] Task: Remove the `infra/` directory from the project.
- [ ] Task: Delete the `.github/workflows/build-infra.yml` file.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Cleanup and Removal' (Protocol in workflow.md)

## Phase 3: Workflow and Secret Configuration
- [ ] Task: Update `.github/workflows/deploy.yml` to include Lab Role ARN and secret sync.
- [ ] Task: Configure the application to prioritize environment variables (from GitHub Secrets) in production.
- [ ] Task: Verify that local `.env` loading remains the fallback for development.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Workflow and Secret Configuration' (Protocol in workflow.md)

## Phase 4: Validation and Testing
- [ ] Task: Perform a dry run of the updated deployment pipeline.
- [ ] Task: Verify that the application can connect to the database in a simulated production environment.
- [ ] Task: Ensure local development environment still works correctly.
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Validation and Testing' (Protocol in workflow.md)
