-- Native PostgreSQL setup for TaskFlow (Stage 1)
-- Run as postgres superuser:
--   sudo -u postgres psql -p 5433 -f 001-setup.sql

CREATE ROLE task_service
WITH LOGIN PASSWORD 'devpassword';

CREATE DATABASE taskflow
OWNER task_service;
