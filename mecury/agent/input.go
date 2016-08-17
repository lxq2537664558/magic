package agent

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/config"
	"github.com/influxdata/toml/ast"
)

type InputConfig struct {
	Name   string
	Prefix string
	Suffix string

	Input Inputer

	Tags     map[string]string
	Filter   InputFilter
	Interval time.Duration
}

type Inputer interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns  a one-sentence description on the Input
	Description() string
	// Gather takes in an accumulator and adds the metrics that the Input
	// gathers. This is called every "interval"
	Gather(Accumulator) error
}

var Inputs = map[string]Inputer{}

func AddInput(name string, input Inputer) {
	Inputs[name] = input
}

// InputFilter containing drop/pass and tagdrop/tagpass rules
type InputFilter struct {
	FieldDrop []string
	fieldDrop Filter

	TagDrop []TagFilter

	IsActive bool
}

// TagFilter is the name of a tag, and the values on which to filter
type TagFilter struct {
	Name   string
	Filter []string
	filter Filter
}

// Compile all Filter lists into filter.Filter objects.
func (f *InputFilter) CompileFilter() error {
	var err error

	f.fieldDrop, err = CompileFilter(f.FieldDrop)
	if err != nil {
		return fmt.Errorf("Error compiling 'fielddrop', %s", err)
	}

	for i := range f.TagDrop {
		f.TagDrop[i].filter, err = CompileFilter(f.TagDrop[i].Filter)
		if err != nil {
			return fmt.Errorf("Error compiling 'tagdrop', %s", err)
		}
	}

	return nil
}

// ShouldTagsPass returns true if the metric should pass, false if should drop
// based on the tagdrop/tagpass filter parameters
func (f *InputFilter) ShouldTagsPass(tags map[string]string) bool {
	if f.TagDrop != nil {
		for _, pat := range f.TagDrop {
			if pat.filter == nil {
				continue
			}
			if tagval, ok := tags[pat.Name]; ok {
				if pat.filter.Match(tagval) {
					return false
				}
			}
		}
		return true
	}

	return true
}

// ShouldFieldsPass returns true if the metric should pass, false if should drop
// based on the drop/pass filter parameters
func (f *InputFilter) ShouldFieldsPass(key string) bool {
	if f.fieldDrop != nil {
		if f.fieldDrop.Match(key) {
			return false
		}
	}
	return true
}

// buildInput parses input specific items from the ast.Table,
// builds the filter and returns a
// models.InputConfig to be inserted into models.RunningInput
func buildInput(name string, tbl *ast.Table) (*InputConfig, error) {
	cp := &InputConfig{Name: name}
	if node, ok := tbl.Fields["interval"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				dur, err := time.ParseDuration(str.Value)
				if err != nil {
					return nil, err
				}

				cp.Interval = dur
			}
		}
	}

	if node, ok := tbl.Fields["name_prefix"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.Prefix = str.Value
			}
		}
	}

	if node, ok := tbl.Fields["name_suffix"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				cp.Suffix = str.Value
			}
		}
	}

	cp.Tags = make(map[string]string)
	if node, ok := tbl.Fields["tags"]; ok {
		if subtbl, ok := node.(*ast.Table); ok {
			if err := config.UnmarshalTable(subtbl, cp.Tags); err != nil {
				log.Printf("Could not parse tags for input %s\n", name)
			}
		}
	}

	delete(tbl.Fields, "name_prefix")
	delete(tbl.Fields, "name_suffix")
	delete(tbl.Fields, "name_override")
	delete(tbl.Fields, "interval")
	delete(tbl.Fields, "tags")
	var err error
	cp.Filter, err = buildFilter(tbl)
	if err != nil {
		return cp, err
	}
	return cp, nil
}
