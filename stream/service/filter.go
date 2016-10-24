package service

import (
	"strings"

	"github.com/gobwas/glob"
)

type Filter interface {
	Match(string) bool
}

type GlobalFilter struct {
	InputDrop []string
	inputDrop Filter

	AlarmDrop []string
	alarmDrop Filter

	Metric_OutputDrop []string
	metric_OutputDrop Filter

	ChainDrop []string
	chainDrop Filter
}

// ShouldFieldsPass returns true if the metric should pass, false if should drop
// based on the drop/pass filter parameters
func (f *GlobalFilter) ShouldInputPass(key string) bool {
	if f.inputDrop != nil {
		if f.inputDrop.Match(key) {
			return false
		}
	}
	return true
}

func (f *GlobalFilter) ShouldAlarmDropPass(key string) bool {
	if f.alarmDrop != nil {
		if f.alarmDrop.Match(key) {
			return false
		}
	}
	return true
}

func (f *GlobalFilter) ShouldMetric_OutputDropPass(key string) bool {
	if f.metric_OutputDrop != nil {
		if f.metric_OutputDrop.Match(key) {
			return false
		}
	}
	return true
}

func (f *GlobalFilter) ShouldChainDropPass(key string) bool {
	if f.chainDrop != nil {
		if f.chainDrop.Match(key) {
			return false
		}
	}
	return true
}

// CompileFilter takes a list of string filters and returns a Filter interface
// for matching a given string against the filter list. The filter list
// supports glob matching too, ie:
//
//   f, _ := CompileFilter([]string{"cpu", "mem", "net*"})
//   f.Match("cpu")     // true
//   f.Match("network") // true
//   f.Match("memory")  // false
//
func CompileFilter(filters []string) (Filter, error) {
	// return if there is nothing to compile
	if len(filters) == 0 {
		return nil, nil
	}

	// check if we can compile a non-glob filter
	noGlob := true
	for _, filter := range filters {
		if hasMeta(filter) {
			noGlob = false
			break
		}
	}

	switch {
	case noGlob:
		// return non-globbing filter if not needed.
		return compileFilterNoGlob(filters), nil
	case len(filters) == 1:
		return glob.Compile(filters[0])
	default:
		return glob.Compile("{" + strings.Join(filters, ",") + "}")
	}
}

// hasMeta reports whether path contains any magic glob characters.
func hasMeta(s string) bool {
	return strings.IndexAny(s, "*?[") >= 0
}

type filter struct {
	m map[string]struct{}
}

func (f *filter) Match(s string) bool {
	_, ok := f.m[s]
	return ok
}

type filtersingle struct {
	s string
}

func (f *filtersingle) Match(s string) bool {
	return f.s == s
}

func compileFilterNoGlob(filters []string) Filter {
	if len(filters) == 1 {
		return &filtersingle{s: filters[0]}
	}
	out := filter{m: make(map[string]struct{})}
	for _, filter := range filters {
		out.m[filter] = struct{}{}
	}
	return &out
}
