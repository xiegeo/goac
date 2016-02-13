package goac

type CommentFeature struct {
	Comments map[string]string
}

func (c *CommentFeature) GetAllComments() map[string]string {
	return c.Comments
}
func (c *CommentFeature) SetAllComments(cs map[string]string) {
	c.Comments = cs
}

type Vertex struct {
	Name           string
	FullAssigments []FullAssigment
}

type FullAssigment struct {
	CommentFeature
	Elevate string
	Over    string
}
