package model

type Post struct {
	ID    string `json:"id"`
	Body  string `json:"body"`
	Owner string `json:"owner"`
}

func NewPost(id, body, owner string) *Post {
	return &Post{
		ID:    id,
		Body:  body,
		Owner: owner,
	}
}
