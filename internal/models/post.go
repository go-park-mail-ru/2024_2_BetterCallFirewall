package models

type Post struct {
	ID           uint32  `json:"id"`
	Header       Header  `json:"header"`
	PostContent  Content `json:"post_content"`
	LikesCount   uint32  `json:"likes_count"`
	IsLiked      bool    `json:"is_liked"`
	CommentCount uint32  `json:"comment_count"`
}

func (p *Post) ToDto() PostDto {
	return PostDto{
		ID:           p.ID,
		Header:       p.Header,
		PostContent:  p.PostContent.ToDto(),
		LikesCount:   p.LikesCount,
		IsLiked:      p.IsLiked,
		CommentCount: p.CommentCount,
	}
}

type PostDto struct {
	ID           uint32
	Header       Header
	PostContent  ContentDto
	LikesCount   uint32
	IsLiked      bool
	CommentCount uint32
}

func (p *PostDto) FromDto() Post {
	return Post{
		ID:           p.ID,
		Header:       p.Header,
		PostContent:  p.PostContent.FromDto(),
		LikesCount:   p.LikesCount,
		IsLiked:      p.IsLiked,
		CommentCount: p.CommentCount,
	}
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

func (c *Comment) ToDto() CommentDto {
	return CommentDto{
		ID:         c.ID,
		Header:     c.Header,
		Content:    c.Content.ToDto(),
		LikesCount: c.LikesCount,
		IsLiked:    c.IsLiked,
	}
}

type CommentDto struct {
	ID         uint32
	Header     Header
	Content    ContentDto
	LikesCount uint32
	IsLiked    bool
}

func (c *CommentDto) FromDto() Comment {
	return Comment{
		ID:         c.ID,
		Header:     c.Header,
		Content:    c.Content.FromDto(),
		LikesCount: c.LikesCount,
		IsLiked:    c.IsLiked,
	}
}
