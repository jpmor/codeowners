package codeowners

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	testcases := []struct {
		input  string
		output *Entry
	}{{
		input: "#This is a code owners file",
		output: &Entry{
			path:    "",
			comment: "#This is a code owners file",
			owners:  make([]string, 0),
		},
	},
		{
			input: "app/lib/index.html @product",
			output: &Entry{
				path:   "app/lib/index.html",
				suffix: Absolute,
				owners: []string{"@product"},
			},
		},
		{
			input: "app/lib/index.html @product @leadership",
			output: &Entry{
				path:   "app/lib/index.html",
				suffix: Absolute,
				owners: []string{"@product", "@leadership"},
			},
		},
	}
	for _, tc := range testcases {
		p := NewParser(strings.NewReader(tc.input))
		got, err := p.Parse()

		if err != nil {
			t.Errorf(`%s err: %v, want no error`, tc.input, err)
		}
		if !reflect.DeepEqual(got, tc.output) {
			t.Errorf("Input %s: got, \n %#v \n want \n %#v", tc.input, got, tc.output)
		}
	}
}
