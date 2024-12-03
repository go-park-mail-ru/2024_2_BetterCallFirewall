package models

type Post struct {
	ID           uint32  `json:"id"`
	Header       Header  `json:"header"`
	PostContent  Content `json:"post_content"`
	LikesCount   uint32  `json:"likes_count"`
	IsLiked      bool    `json:"is_liked"`
	CommentCount uint32  `json:"comment_count"`
}

type Header struct {
	AuthorID    uint32  `json:"author_id"`
	CommunityID uint32  `json:"community_id"`
	Author      string  `json:"author"`
	Avatar      Picture `json:"avatar"`
}

type Comment struct {
	ID         uint32  `json:"id"`
	Header     Header  `json:"header"`
	Content    Content `json:"content"`
	LikesCount uint32  `json:"likes_count"`
	IsLiked    bool    `json:"is_liked"`
}
