package pathmatch_test

import (
	"testing"

	"github.com/abichinger/pathmatch"
	"github.com/stretchr/testify/assert"
)

func TestFindSubmatch(t *testing.T) {
	tests := []struct {
		path     string
		str      string
		expected pathmatch.Match
	}{
		{"/", "/", map[string]string{}},
		{"/foo", "/foo", map[string]string{}},
		{"/foo", "/bar", nil},
		{"/foo/:id", "/foo/1", map[string]string{"id": "1"}},
		{"/foo/:id", "/foo/1/bar", nil},
		{"/foo/:id/bar/:name", "/foo/1/bar/tom", map[string]string{"id": "1", "name": "tom"}},
		//{"/foo/:id/bar/:id", "/foo/1/bar/tom", false},
		{"/foo/:id/bar/:id", "/foo/1/bar/1", map[string]string{"id": "1"}},
		{"/*", "/foo/bar", map[string]string{"$0": "foo/bar"}},
		{"/foo/:id/bar/*", "/foo/1/bar/2/baz/3", map[string]string{"id": "1", "$0": "2/baz/3"}},
		{"/*/bar/:id", "/foo/1/bar/2", map[string]string{"$0": "foo/1", "id": "2"}},
	}

	for _, test := range tests {
		p, err := pathmatch.Compile(test.path)
		if err != nil {
			t.Errorf(err.Error())
		}
		actual := p.FindSubmatch(test.str)
		assert.Equalf(t, test.expected, actual, "path: %s, str %s", test.path, test.path)
		actualBool := p.Match(test.str)
		assert.Equalf(t, test.expected != nil, actualBool, "path: %s, str %s", test.path, test.path)
	}
}
