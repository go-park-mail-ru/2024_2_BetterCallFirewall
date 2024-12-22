CREATE ROLE admin WITH LOGIN PASSWORD 'admin_password';
ALTER ROLE admin SUPERUSER;

CREATE ROLE profile_service WITH LOGIN PASSWORD 'profile_password';
GRANT CONNECT ON DATABASE mydbvk TO profile_service;
GRANT SELECT, INSERT, DELETE, UPDATE, REFERENCES, TRUNCATE, TRIGGER, MAINTAIN ON profile, friend TO profile_service;

CREATE ROLE community_service WITH LOGIN PASSWORD 'community_password';
GRANT CONNECT ON DATABASE mydbvk TO community_service;
GRANT USAGE ON SCHEMA public TO community_service;
GRANT SELECT, INSERT, DELETE, UPDATE, REFERENCES, TRUNCATE, TRIGGER, MAINTAIN ON community, community_profile, admin TO community_service;

CREATE ROLE post_service WITH LOGIN PASSWORD 'post_password';
GRANT CONNECT ON DATABASE mydbvk TO post_service;
GRANT SELECT, INSERT, DELETE, UPDATE, REFERENCES, TRUNCATE, TRIGGER, MAINTAIN ON post, comment, reaction TO post_service;

CREATE ROLE message_service WITH LOGIN PASSWORD 'message_password';
GRANT CONNECT ON DATABASE mydbvk TO message_service;
GRANT SELECT, INSERT, DELETE, UPDATE, REFERENCES, TRUNCATE, TRIGGER, MAINTAIN ON message TO message_service;