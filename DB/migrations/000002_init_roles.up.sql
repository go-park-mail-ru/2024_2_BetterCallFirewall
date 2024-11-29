CREATE ROLE admin WITH LOGIN PASSWORD 'admin_password';
ALTER ROLE admin SUPERUSER;

CREATE ROLE profile_service WITH LOGIN PASSWORD 'profile_password';
GRANT CONNECT ON DATABASE mydbvk TO profile_service;
GRANT SELECT, INSERT, DELETE, UPDATE ON profile, friend TO profile_service;

CREATE ROLE community_service WITH LOGIN PASSWORD 'community_password';
GRANT CONNECT ON DATABASE mydbvk TO community_service;
GRANT SELECT, INSERT, DELETE, UPDATE ON community, community_profile TO community_service;

CREATE ROLE post_service WITH LOGIN PASSWORD 'post_password';
GRANT CONNECT ON DATABASE mydbvk TO post_service;
GRANT SELECT, INSERT, DELETE, UPDATE ON post, comment, reaction TO post_service;

CREATE ROLE message_service WITH LOGIN PASSWORD 'message_password';
GRANT CONNECT ON DATABASE mydbvk TO message_service;
GRANT SELECT, INSERT, DELETE ON message TO message_service;