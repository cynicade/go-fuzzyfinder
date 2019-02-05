// Package matching provides matching features that find appropriate strings
// by using a passed input string.
package matching

import (
	"sort"
	"strings"
	"unicode"

	"github.com/ktr0731/go-fuzzyfinder/scoring"
)

// Matched represents a result of FindAll.
type Matched struct {
	// Idx is the index of an item of the original slice which was used to
	// search matched strings.
	Idx int
	// Pos is the range of matched position.
	// [2]int represents a closed interval of a position.
	Pos [2]int
	// score is the value that indicates how it similar to the input string.
	// The bigger score, the more similar it is.
	score int
}

type option func(*opt)

type Mode int

const (
	ModeSmart Mode = iota
	ModeCaseSensitive
	ModeCaseInsensitive
)

// opt represents available options and its default values.
type opt struct {
	mode Mode
}

// WithMode specifies a matching mode. The default mode is ModeSmart.
func WithMode(m Mode) option {
	return func(o *opt) {
		o.mode = m
	}
}

// FindAll tries to find out sub-strings from slice that match the passed argument in.
// The returned slice is sorted by similarity scores in descending order.
func FindAll(in string, slice []string, opts ...option) []Matched {
	var opt opt
	for _, o := range opts {
		o(&opt)
	}
	m := match(in, slice, opt)
	sort.Slice(m, func(i, j int) bool {
		return m[i].score > m[j].score
	})
	return m
}

// match iterates each string of slice for check whether it is matched to the input string.
func match(input string, slice []string, opt opt) (res []Matched) {
	if opt.mode == ModeSmart {
		// Find an upper-case rune
		n := strings.IndexFunc(input, unicode.IsUpper)
		if n == -1 {
			opt.mode = ModeCaseInsensitive
			input = strings.ToLower(input)
		} else {
			opt.mode = ModeCaseSensitive
		}
	}

	in := []rune(input)
	for idxOfSlice, s := range slice {
		var from, idx int
		if opt.mode == ModeCaseInsensitive {
			s = strings.ToLower(s)
		}
	LINE_MATCHING:
		for i, r := range []rune(s) {
			if r == in[idx] {
				if idx == 0 {
					from = i
				}
				idx++
				if idx == len(in) {
					res = append(res, Matched{
						Idx: idxOfSlice,
						Pos: [2]int{from, i + 1},
						// TODO: 引数と順番をあわせる
						score: scoring.Calculate(s, input),
					})
					break LINE_MATCHING
				}
			}
		}
	}
	return
}
