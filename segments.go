package pathmatch

import (
	"strconv"
	"strings"
)

type SegType int

type MatchDraft struct {
	capture  bool
	unnamed  int
	match    Match
	segments []ISegment
	str      string
	sep      string
}

func NewMatchDraft(capture bool, match Match, segments []ISegment, str string, sep string) *MatchDraft {
	if capture == false {
		return &MatchDraft{capture, 0, match, segments, str, sep}
	}
	return &MatchDraft{capture, 0, make(Match), segments, str, sep}
}

func (m *MatchDraft) nextSeg(start int) (next string, done bool) {
	i := strings.Index(m.str[start:], m.sep)
	if i == -1 {
		return m.str, true
	}
	return m.str[:i+start], false
}

func (m *MatchDraft) sliceNext(next string, done bool) string {
	if done {
		return m.str[len(next):]
	}
	return m.str[len(next)+1:]
}

func (m *MatchDraft) AddUnnamed(value string) {
	if m.capture == false {
		return
	}
	m.match["$"+strconv.Itoa(m.unnamed)] = value
	m.unnamed++
}

func (m *MatchDraft) Add(key, value string) {
	if m.capture == false {
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
	Match(m *MatchDraft) *MatchDraft
	Type() SegType
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

func (seg *StaticSegment) Match(m *MatchDraft) *MatchDraft {
	next, done := m.nextSeg(0)
	if done && len(m.segments) > 0 {
		return nil
	}
	if next != seg.value {
		return nil
	}
	m.str = m.sliceNext(next, done)
	return m
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

func (seg *ParamSegment) Match(m *MatchDraft) *MatchDraft {
	next, done := m.nextSeg(0)
	if done && len(m.segments) > 0 {
		return nil
	}
	m.Add(seg.key, next)
	m.str = m.sliceNext(next, done)
	return m
}

type WildcardSegment struct {
}

func NewWildcardSegment() *WildcardSegment {
	return &WildcardSegment{}
}

func (seg *WildcardSegment) Type() SegType {
	return Wildcard
}

func (seg *WildcardSegment) Match(m *MatchDraft) *MatchDraft {
	next, done := m.nextSeg(0)
	if done && len(m.segments) > 0 {
		return nil
	}

	if len(m.segments) == 0 {
		m.AddUnnamed(m.str)
		m.str = m.str[:0]
		return m
	}
	for !done {
		if nextMatch := m.segments[0].Match(NewMatchDraft(false, nil, m.segments[1:], m.sliceNext(next, done), m.sep)); nextMatch != nil {
			m.AddUnnamed(next)
			m.str = m.sliceNext(next, done)
			return m
		}
		next, done = m.nextSeg(len(next) + len(m.sep))
	}
	return nil
}
