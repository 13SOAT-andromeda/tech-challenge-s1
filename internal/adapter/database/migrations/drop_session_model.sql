-- Migration: drop_session_model
-- Description: Drops the sessions table now that in-app session management
--              has been removed. Authentication is handled by the Lambda Authorizer.
-- Run manually before or during the deployment of the remove-jwt-auth-middleware change.

DROP TABLE IF EXISTS session_model;
