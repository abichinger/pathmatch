package pathmatch_test

import (
	"testing"

	pm "github.com/abichinger/pathmatch"
	"github.com/stretchr/testify/assert"
)

func TestFindSubmatch(t *testing.T) {
	tests := []struct {
		path     string
		str      string
		expected pm.Match
	}{
		{"/", "/", map[string]string{}},
		{"/foo", "/foo", map[string]string{}},
		{"/foo", "/bar", nil},
		{"/foo/:id", "/foo/1", map[string]string{"id": "1"}},
		{"/foo/:id", "/foo/1/bar", nil},
		{"/foo/:id/bar/:name", "/foo/1/bar/tom", map[string]string{"id": "1", "name": "tom"}},
		{"/foo/:id/bar/:id", "/foo/1/bar/2", map[string]string{"id": "2"}},
		{"/foo/:id/bar/:id", "/foo/1/bar/1", map[string]string{"id": "1"}},
		{"/*", "/foo/bar", map[string]string{"$0": "foo/bar"}},
		{"/foo/:id/bar/*", "/foo/1/bar/2/baz/3", map[string]string{"id": "1", "$0": "2/baz/3"}},
		{"/*/bar/:id", "/foo/1/bar/2", map[string]string{"$0": "foo/1", "id": "2"}},
		{"/*/foo/bar", "/api/foo/baz/foo/bar", map[string]string{"$0": "api/foo/baz"}},
		{"/*/a/b/*/c/d", "/x/a/x/a/b/x/a/b/x/c/d", map[string]string{"$0": "x/a/x", "$1": "x/a/b/x"}},
	}

	for _, test := range tests {
		p, err := pm.Compile(test.path)
		if err != nil {
			t.Errorf(err.Error())
		}
		actual := p.FindSubmatch(test.str)
		assert.Equalf(t, test.expected, actual, "path: %s, str %s", test.path, test.str)
		actualBool := p.Match(test.str)
		assert.Equalf(t, test.expected != nil, actualBool, "path: %s, str %s", test.path, test.str)
	}
}

func TestOptions(t *testing.T) {
	tests := []struct {
		path     string
		options  []pm.Option
		str      string
		expected pm.Match
	}{
		{"foo.{{name}}.**", []pm.Option{pm.SetSeperator("."), pm.SetPrefix("{{"), pm.SetSuffix("}}"), pm.SetWildcard("**")}, "foo.bar.baz", map[string]string{"name": "bar", "$0": "baz"}},
		{"/foo/:id/bar/:id", []pm.Option{pm.EnableEqualityCheck(true)}, "/foo/1/bar/2", nil},
		{"/foo/:id/bar/:id", []pm.Option{pm.EnableEqualityCheck(true)}, "/foo/1/bar/1", map[string]string{"id": "1"}},
	}

	for _, test := range tests {
		p, err := pm.Compile(test.path, test.options...)
		if err != nil {
			t.Errorf(err.Error())
		}
		actual := p.FindSubmatch(test.str)
		assert.Equalf(t, test.expected, actual, "path: %s, str %s", test.path, test.str)
		actualBool := p.Match(test.str)
		assert.Equalf(t, test.expected != nil, actualBool, "path: %s, str %s", test.path, test.str)
	}
}

func TestIsStatic(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/", true},
		{"/foo/bar", true},
		{"/foo/:id", false},
		{"/*", false},
		{"/foo/:id/*", false},
	}

	for _, test := range tests {
		p, err := pm.Compile(test.path)
		if err != nil {
			t.Errorf(err.Error())
		}
		actual := p.IsStatic()
		assert.Equalf(t, test.expected, actual, "path: %s", test.path)
	}
}
