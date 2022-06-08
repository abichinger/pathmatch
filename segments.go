package pathmatch

import (
	"fmt"
	"strings"
)

type SegType int

type MatchDraft struct {
	unnamed int
	match   Match
}

func NewMatchDraft() *MatchDraft {
	return &MatchDraft{0, Match{}}
}

func (md *MatchDraft) AddUnnamed(value string) {
	md.match[fmt.Sprintf("$%d", md.unnamed)] = value
	md.unnamed++
}

const (
	Static SegType = iota
	Parameterized
	Wildcard
)

type ISegment interface {
	Match(res *MatchDraft, nextSegments []ISegment, strSegments []string) (*MatchDraft, []ISegment, []string)
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

func (seg *StaticSegment) Match(res *MatchDraft, nextSegments []ISegment, strSegments []string) (*MatchDraft, []ISegment, []string) {
	if len(strSegments) == 0 || strSegments[0] != seg.value {
		return nil, nil, nil
	}
	return res, nextSegments, strSegments[1:]
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

func (seg *ParamSegment) Match(res *MatchDraft, nextSegments []ISegment, strSegments []string) (*MatchDraft, []ISegment, []string) {
	if len(strSegments) == 0 || strSegments[0] == "" {
		return nil, nil, nil
	}
	if value, ok := res.match[seg.key]; ok && value != strSegments[0] {
		return nil, nil, nil
	}
	res.match[seg.key] = strSegments[0]
	return res, nextSegments, strSegments[1:]
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

func (seg *WildcardSegment) Match(res *MatchDraft, nextSegments []ISegment, strSegments []string) (*MatchDraft, []ISegment, []string) {
	if len(strSegments) == 0 {
		return nil, nil, nil
	}
	if len(nextSegments) == 0 {
		res.AddUnnamed(strings.Join(strSegments, seg.seperator))
		return res, nextSegments, []string{}
	}
	for i := 1; i < len(strSegments); i++ {
		if m, _, _ := nextSegments[0].Match(NewMatchDraft(), nextSegments[1:], strSegments[i:]); m != nil {
			res.AddUnnamed(strings.Join(strSegments[:i], seg.seperator))
			return res, nextSegments, strSegments[i:]
		}
	}
	return nil, nil, nil
}
