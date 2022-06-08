package pathmatch

import (
	"strconv"
	"strings"
)

type SegType int

type MatchDraft struct {
	capture     bool
	unnamed     int
	match       Match
	segments    []ISegment
	strSegments []string
}

func NewMatchDraft(capture bool, match Match, segments []ISegment, strSegments []string) *MatchDraft {
	if capture == false {
		return &MatchDraft{capture, 0, match, segments, strSegments}
	}
	return &MatchDraft{capture, 0, match, segments, strSegments}
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
	if len(m.strSegments) == 0 || m.strSegments[0] != seg.value {
		return nil
	}
	m.strSegments = m.strSegments[1:]
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
	if len(m.strSegments) == 0 || m.strSegments[0] == "" {
		return nil
	}
	/*if value, ok := m.match[seg.key]; ok && value != m.strSegments[0] {
		return nil
	}*/
	m.Add(seg.key, m.strSegments[0])
	m.strSegments = m.strSegments[1:]
	return m
}

type WildcardSegment struct {
	seperator string
}

func NewWildcardSegment(seperator string) *WildcardSegment {
	return &WildcardSegment{seperator}
}

func (seg *WildcardSegment) Type() SegType {
	return Wildcard
}

func (seg *WildcardSegment) Match(m *MatchDraft) *MatchDraft {
	if len(m.strSegments) == 0 {
		return nil
	}
	if len(m.segments) == 0 {
		m.AddUnnamed(strings.Join(m.strSegments, seg.seperator))
		m.strSegments = m.strSegments[:0]
		return m
	}
	for i := 1; i < len(m.strSegments); i++ {
		if nextMatch := m.segments[0].Match(NewMatchDraft(false, nil, m.segments[1:], m.strSegments[i:])); nextMatch != nil {
			m.AddUnnamed(strings.Join(m.strSegments[:i], seg.seperator))
			m.strSegments = m.strSegments[i:]
			return m
		}
	}
	return nil
}
