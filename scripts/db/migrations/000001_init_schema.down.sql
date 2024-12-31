-- Drop triggers
DROP TRIGGER IF EXISTS update_api_tokens_updated_at ON api_tokens;
DROP TRIGGER IF EXISTS update_licenses_updated_at ON licenses;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients;
DROP TRIGGER IF EXISTS update_license_types_updated_at ON license_types;
DROP TRIGGER IF EXISTS update_applications_updated_at ON applications;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS license_activities;
DROP TABLE IF EXISTS api_tokens;
DROP TABLE IF EXISTS licenses;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS license_types;
DROP TABLE IF EXISTS applications;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";