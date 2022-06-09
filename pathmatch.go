// Package pathmatch allows you to quickly parse and match URL paths.
//
// Path expression are a composition of static, parameterized and wildcard segments.
//
// The wildcard * matches one or more segments.
//
//  /* / 		-> no match
//  /* /foo 	-> * matches foo
//  /* /foo/bar -> * matches foo/bar
//
// Parameterized segments start with a colon and have a name
//
//  /foo/:name 	/foo 		-> no match
//  /foo/:name 	/foo/bar 	-> :name matches bar
//
package pathmatch

import (
	"strconv"
	"strings"
)

type Match map[string]string

type savePoint struct {
	i           int
	sIndex      int
	searchStart int
	valid       bool
}

type Path struct {
	Seperator string
	Prefix    string
	Suffix    string
	Wildcard  string
	Segments  []ISegment
	match     Match
	save      *savePoint
}

// Compile parses a path expression and returns a Path if successful
func Compile(path string, options ...Option) (*Path, error) {
	p := &Path{"/", ":", "", "*", []ISegment{}, make(Match, 0), &savePoint{}}
	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}

	unnamed := 0
	strSegments := strings.Split(path, p.Seperator)
	for _, strSeg := range strSegments {
		if strSeg == p.Wildcard {
			key := "$" + strconv.Itoa(unnamed)
			unnamed++
			p.Segments = append(p.Segments, NewWildcardSegment(key))
		} else if (p.Prefix == "" || strings.HasPrefix(strSeg, p.Prefix)) && (p.Suffix == "" || strings.HasSuffix(strSeg, p.Suffix)) {
			key := strSeg[len(p.Prefix) : len(strSeg)-len(p.Suffix)]
			if key == "" {
				key = "$" + strconv.Itoa(unnamed)
				unnamed++
			}
			p.Segments = append(p.Segments, NewParamSegment(key))
		} else {
			p.Segments = append(p.Segments, NewStaticSegment(strSeg))
		}
	}

	return p, nil
}

// Match returns true if s and p match
func (p *Path) Match(s string) bool {
	m := p.getMatch(s, false)
	return m != nil
}

// FindSubmatch returns a map with the values of parameterized segments, if s and p match
// Otherwise nil is returned
// Wildcard segments are named $0, $1, ...
func (p *Path) FindSubmatch(s string) Match {
	return p.getMatch(s, true)
}

func sliceSegment(s string, sep string, start int, offset int) (string, bool) {
	str := s[start:]
	i := strings.Index(str[offset:], sep)
	if i == -1 {
		return str, true
	}
	return str[:i+offset], false
}

func segmentLen(s string, sep string, done bool) int {
	if done {
		return len(s)
	}
	return len(s) + len(sep)
}

func (p *Path) getMatch(s string, capture bool) Match {
	draft := NewMatchDraft(capture, p.match)

	sIndex := 0
	searchStart := 0

	for i := 0; draft != nil && i < len(p.Segments); i++ {
		seg := p.Segments[i]

		str, done := sliceSegment(s, p.Seperator, sIndex, searchStart)
		if done && len(p.Segments)-1 != i {
			return nil
		}

		if seg.Multiple() {

			if len(p.Segments)-1 == i {
				draft = seg.Match(draft, s[sIndex:])
				sIndex = len(s)
				break
			}

			if p.save.valid && p.save.i == i {
				p.save.searchStart = segmentLen(str, p.Seperator, done)
			} else {
				p.save.i = i
				p.save.sIndex = sIndex
				p.save.searchStart = segmentLen(str, p.Seperator, done)
				p.save.valid = true
			}
		}

		m := seg.Match(draft, str)
		if m == nil && p.save.valid {
			i = p.save.i - 1
			sIndex = p.save.sIndex
			searchStart = p.save.searchStart
			continue
		}

		draft = m
		sIndex += segmentLen(str, p.Seperator, done)
		searchStart = 0
	}
	if draft == nil || len(s) != sIndex {
		return nil
	}
	return draft.match
}

// IsStatic returns true if p only contains static segments
func (p *Path) IsStatic() bool {
	for _, seg := range p.Segments {
		if seg.Type() != Static {
			return false
		}
	}
	return true
}
