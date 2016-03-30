package goac

import (
	"encoding/json"
	"io"
)

type Vertex struct {
	Name            string           `json:"name,omitempty"`
	FullAssignments []FullAssignment `json:"fullAssignments,omitempty"`
}

type FullAssignment struct {
	Elevate  string            `json:"elevate,omitempty"`
	Over     string            `json:"over,omitempty"`
	Comments map[string]string `json:"comments,omitempty"`
}

func (v Vertex) EncodeJson(w io.Writer) error {
	bs, err := json.MarshalIndent(v, "", "\t")
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
