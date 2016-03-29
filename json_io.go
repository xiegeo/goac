package goac

import (
	"encoding/json"
	"io"
)

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

func (v Vertex) EncodeJson(w io.Writer) error {
	bs, err := json.MarshalIndent(v, "", "/t")
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	return err
}

func DecodeJsonToVertex(r io.Reader) (Vertex, error) {
	v := Vertex{}
	dec := json.NewDecoder(r)
	err := dec.Decode(&v)
	return v, err
}
