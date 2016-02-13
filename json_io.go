package goac

type CommentFeature struct {
	Comments map[string]string `json:"comments,omitempty"`
}

func (c *CommentFeature) GetAllComments() map[string]string {
	return c.Comments
}
func (c *CommentFeature) SetAllComments(cs map[string]string) {
	c.Comments = cs
}

type Vertex struct {
	Name            string           `json:"name,omitempty"`
	FullAssignments []FullAssignment `json:"fullAssignments,omitempty"`
}

type FullAssignment struct {
	CommentFeature
	Elevate string `json:"elevate,omitempty"`
	Over    string `json:"over,omitempty"`
}
