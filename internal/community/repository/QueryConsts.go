package repository

const (
	CreateNewCommunity = `WITH new_community AS (
    INSERT INTO community(name, about) VALUES ($1, $2) RETURNING id
) INSERT INTO community_profile(community_id, profile_id) VALUES ((SELECT id FROM new_community), $3) RETURNING id;`
	GetOne = `
SELECT community.id, name, file_path, about, 
       (SELECT COUNT(*) FROM community_profile WHERE community_id = $1) AS subs
    FROM community
    INNER JOIN file ON avatar = file.id
WHERE community.id = $1;`
	UpdateWithoutAvatar = `UPDATE community SET name = $1, about = $2 WHERE id = $3;`
	UpdateWithAvatar    = ` 
WITH newAvatar AS (
	INSERT INTO file VALUES($2) RETURNING id		
) UPDATE community SET name = $1, avatar = SELECT id FROM newAvatar, about = $3 WHERE id = $4;`
	Delete   = `DELETE FROM community WHERE id = $1;`
	GetBatch = `
SELECT community.id, name, file_path, about 
FROM community 
    JOIN file ON avatar = file.id 
WHERE community.id > $1 
ORDER BY community.id ASC 
LIMIT $2;`
	JoinCommunity  = `INSERT INTO community_profile(community_id, profile_id)  VALUES ($1, $2);`
	LeaveCommunity = `DELETE FROM community_profile WHERE community_id = $1 AND profile_id = $2;`
)
