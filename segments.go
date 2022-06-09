package pathmatch

type SegType int

type MatchDraft struct {
	capture bool
	match   Match
}

func NewMatchDraft(capture bool, match Match) *MatchDraft {
	if !capture {
		return &MatchDraft{capture, match}
	}
	return &MatchDraft{capture, make(Match)}
}

func (m *MatchDraft) Set(key, value string) {
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
	//returns m if the segment matches s,
	Match(m *MatchDraft, s string) *MatchDraft
	Type() SegType
	Multiple() bool
}

type StaticSegment struct {
	value string
}

func NewStaticSegment(value string) *StaticSegment {
	return &StaticSegment{value}
}

func (seg *StaticSegment) Type() SegType {
	return Static
}

func (seg *StaticSegment) Match(m *MatchDraft, s string) *MatchDraft {
	if s != seg.value {
		return nil
	}
	return m
}

func (seg *StaticSegment) Multiple() bool {
	return false
}

type ParamSegment struct {
	key string
}

func NewParamSegment(key string) *ParamSegment {
	return &ParamSegment{key}
}

func (seg *ParamSegment) Type() SegType {
	return Parameterized
}

func (seg *ParamSegment) Match(m *MatchDraft, s string) *MatchDraft {
	m.Set(seg.key, s)
	return m
}

func (seg *ParamSegment) Multiple() bool {
	return false
}

type WildcardSegment struct {
	key string
}

func NewWildcardSegment(key string) *WildcardSegment {
	return &WildcardSegment{key}
}

func (seg *WildcardSegment) Type() SegType {
	return Wildcard
}

func (seg *WildcardSegment) Match(m *MatchDraft, s string) *MatchDraft {
	m.Set(seg.key, s)
	return m
}

func (seg *WildcardSegment) Multiple() bool {
	return true
}
