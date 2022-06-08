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
	match     Match
}

func Compile(path string, options ...Option) (*Path, error) {
	p := &Path{"/", ":", "", "*", []ISegment{}, make(Match, 0)}
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

func (p *Path) Match(s string) bool {
	m := p.getMatch(s, false)
	return m != nil
}

func (p *Path) FindSubmatch(s string) Match {
	return p.getMatch(s, true)
}

func (p *Path) getMatch(s string, capture bool) Match {
	strSegments := strings.Split(s, p.Seperator)
	draft := NewMatchDraft(capture, p.match, p.Segments, strSegments)

	for draft != nil && len(draft.segments) > 0 {
		seg := draft.segments[0]
		draft.segments = draft.segments[1:]
		draft = seg.Match(draft)
	}
	if draft == nil || len(draft.strSegments) > 0 {
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

func MatchString(s string, options ...Option) (bool, error) {
	p, err := Compile(s, options...)
	if err != nil {
		return false, err
	}
	return p.Match(s), nil
}

// /foo/*/bar/baz /foo/bar/foo/bar/baz
