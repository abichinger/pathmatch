package pathmatch_test

import (
	"testing"

	"github.com/abichinger/pathmatch"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern  string
		path     string
		expected bool
	}{
		{"/", "/", true},
		{"/foo", "/foo", true},
		{"/foo", "/bar", false},
		{"/foo/:id", "/foo/1", true},
		{"/foo/:id", "/foo/1/bar", false},
		{"/foo/:id/bar/:name", "/foo/1/bar/tom", true},
		{"/foo/:id/bar/:id", "/foo/1/bar/tom", false},
		{"/foo/:id/bar/:id", "/foo/1/bar/1", true},
		{"/*", "/foo/bar", true},
		{"foo/:id/bar/*", "foo/1/bar/2/baz/3", true},
	}

	for _, test := range tests {
		p, err := pathmatch.Compile(test.pattern)
		if err != nil {
			t.Errorf(err.Error())
		}
		m := p.MatchString(test.path)
		actual := m != nil
		assert.Equalf(t, test.expected, actual, "pattern: %s, path %s", test.pattern, test.path)
	}
}
