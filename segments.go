package pathmatch

type SegType int

type matchDraft struct {
	capture bool
	match   Match
}

func newMatchDraft(capture bool, match Match) *matchDraft {
	if !capture {
		return &matchDraft{capture, match}
	}
	return &matchDraft{capture, make(Match)}
}

func (m *matchDraft) set(key, value string) {
	if !m.capture {
		return
	}
	m.match[key] = value
}

const (
	Static SegType = iota
	Parameterized
	Wildcard
)

type ISegment interface {
	// Match returns m if the segment matches s,
	Match(m *matchDraft, s string) *matchDraft

	// Type returns the segment type
	Type() SegType

	// Multiple returns true, if the segment can match one or more string segments
	Multiple() bool
}

type staticSegment struct {
	value string
}

func newStaticSegment(value string) *staticSegment {
	return &staticSegment{value}
}

func (seg *staticSegment) Type() SegType {
	return Static
}

func (seg *staticSegment) Match(m *matchDraft, s string) *matchDraft {
	if s != seg.value {
		return nil
	}
	return m
}

func (seg *staticSegment) Multiple() bool {
	return false
}

type paramSegment struct {
	key        string
	equalCheck bool
}

func newParamSegment(key string, equalCheck bool) *paramSegment {
	return &paramSegment{key, equalCheck}
}

func (seg *paramSegment) Type() SegType {
	return Parameterized
}

func (seg *paramSegment) Match(m *matchDraft, s string) *matchDraft {
	if value, ok := m.match[seg.key]; seg.equalCheck && ok && s != value {
		return nil
	}
	m.set(seg.key, s)
	return m
}

func (seg *paramSegment) Multiple() bool {
	return false
}

type wildcardSegment struct {
	key string
}

func newWildcardSegment(key string) *wildcardSegment {
	return &wildcardSegment{key}
}

func (seg *wildcardSegment) Type() SegType {
	return Wildcard
}

func (seg *wildcardSegment) Match(m *matchDraft, s string) *matchDraft {
	m.set(seg.key, s)
	return m
}

func (seg *wildcardSegment) Multiple() bool {
	return true
}
