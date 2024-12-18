package repository

const (
	CreateNewCommunity = `WITH new_community AS (
    INSERT INTO community(name, about) VALUES ($1, $2) RETURNING id
) INSERT INTO community_profile(community_id, profile_id) VALUES ((SELECT id FROM new_community), $3)
RETURNING (SELECT id FROM new_community);`
	CreateNewCommunityWithAvatar = `WITH new_community AS (
    INSERT INTO community(name, about, avatar) VALUES ($1, $2, $3) RETURNING id
) INSERT INTO community_profile(community_id, profile_id) VALUES ((SELECT id FROM new_community), $4)
RETURNING (SELECT id FROM new_community);`
	GetOne = `
SELECT community.id, name, avatar, about, 
       (SELECT COUNT(*) FROM community_profile WHERE community_id = $1) AS subs
    FROM community
WHERE community.id = $1;`
	UpdateWithoutAvatar = `UPDATE community SET name = $1, about = $2 WHERE id = $3;`
	UpdateWithAvatar    = `UPDATE community SET name = $1, avatar = $2, about = $3 WHERE id = $4;`
	Delete              = `DELETE FROM community WHERE id = $1;`
	GetBatch            = `
SELECT community.id, name, avatar, about 
FROM community  
WHERE community.id < $1 
ORDER BY community.id DESC 
LIMIT $2;`
	JoinCommunity  = `INSERT INTO community_profile(community_id, profile_id)  VALUES ($1, $2);`
	LeaveCommunity = `DELETE FROM community_profile WHERE community_id = $1 AND profile_id = $2;`

	Search = `
SELECT community.id, name, avatar, about
FROM community
WHERE 
    (name ILIKE '%' || $1 || '%' OR about ILIKE '%' || $1 || '%')
	AND community.id < $2
ORDER BY community.name ASC
LIMIT $3;`

	GetHeader      = `SELECT id, name, avatar FROM community WHERE id = $1`
	IsFollow       = `SELECT COUNT(*) FROM community_profile WHERE community_id = $1 AND profile_id = $2`
	InsertNewAdmin = `INSERT INTO admin(community_id, admin_id) VALUES ($1, $2)`
	CheckAccess    = `SELECT COUNT(*) FROM admin WHERE community_id = $1 AND admin_id = $2`
	DeleteAdmin    = `DELETE FROM admin WHERE community_id = $1 AND admin_id = $2`
)
