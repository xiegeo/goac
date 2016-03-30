package goac

import (
	"bytes"
	"testing"
)

func TestVertexJSON(t *testing.T) {
	v := Vertex{
		Name: "tester",
		FullAssignments: []FullAssignment{
			FullAssignment{
				Elevate: "a",
				Over:    "b",
				Comments: map[string]string{
					"generated-by": "test",
				},
			},
		},
	}
	json := `{
	"name": "tester",
	"fullAssignments": [
		{
			"elevate": "a",
			"over": "b",
			"comments": {
				"generated-by": "test"
			}
		}
	]
}`
	buf := bytes.NewBuffer(nil)
	v.EncodeJson(buf)
	if json != buf.String() {
		t.Fatalf("encode not expected:%v", buf.String())
	}
}
