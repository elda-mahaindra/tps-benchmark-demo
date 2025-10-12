\c demo_db;

-- Permission grants

-- Grant permissions to demo schema
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA demo TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA demo TO postgres;
GRANT USAGE ON SCHEMA demo TO postgres;
