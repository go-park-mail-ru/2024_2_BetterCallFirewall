package models

type Post struct {
	ID          uint32  `json:"id"`
	Header      Header  `json:"header"`
	PostContent Content `json:"post_content"`
	AuthorID    uint32  `json:"user_id,omitempty"`
}

type Header struct {
	Author string  `json:"author"`
	Avatar Picture `json:"avatar"`
}
