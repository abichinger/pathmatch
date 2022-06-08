package pathmatch

import "strings"

type Match map[string]string

type Option func(path *Path) error

type Path struct {
	Seperator string
	Prefix    string
	Suffix    string
	Wildcard  string
	Segments  []ISegment
}

func Compile(path string, options ...Option) (*Path, error) {
	p := &Path{"/", ":", "", "*", []ISegment{}}
	for _, option := range options {
		if err := option(p); err != nil {
			return nil, err
		}
	}

	strSegments := strings.Split(path, p.Seperator)
	for _, strSeg := range strSegments {
		if strSeg == p.Wildcard {
			p.Segments = append(p.Segments, NewWildcardSegment(p.Seperator))
		} else if (p.Prefix == "" || strings.HasPrefix(strSeg, p.Prefix)) && (p.Suffix == "" || strings.HasSuffix(strSeg, p.Suffix)) {
			key := strSeg[len(p.Prefix) : len(strSeg)-len(p.Suffix)]
			p.Segments = append(p.Segments, NewParamSegment(key))
		} else {
			p.Segments = append(p.Segments, NewStaticSegment(strSeg))
		}
	}

	return p, nil
}

func (p *Path) MatchString(s string) Match {
	draft := NewMatchDraft()
	var segments []ISegment
	segments = append(segments, p.Segments...)
	strSegments := strings.Split(s, p.Seperator)

	for draft != nil && len(segments) > 0 {
		draft, segments, strSegments = segments[0].Match(draft, segments[1:], strSegments)
	}
	if len(strSegments) > 0 || draft == nil {
		return nil
	}
	return draft.match
}

func (p *Path) IsStatic() bool {
	for _, seg := range p.Segments {
		if seg.Type() != Static {
			return false
		}
	}
	return true
}

func MatchString(s string, options ...Option) (Match, error) {
	p, err := Compile(s, options...)
	if err != nil {
		return nil, err
	}
	return p.MatchString(s), nil
}

// /foo/*/bar/baz /foo/bar/foo/bar/baz
