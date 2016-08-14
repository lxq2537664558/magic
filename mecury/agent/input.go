package agent

import "fmt"

type Inputer interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns a one-sentence description on the Input
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
	Name []string
	name Filter

	Field []string
	field Filter

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
	f.name, err = CompileFilter(f.Name)
	if err != nil {
		return fmt.Errorf("Error compiling 'namedrop', %s", err)
	}

	f.field, err = CompileFilter(f.Field)
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

func (f *InputFilter) ShouldMetricPass(metric Metric) bool {
	if f.ShouldNamePass(metric.Name()) && f.ShouldTagsPass(metric.Tags()) {
		return true
	}
	return false
}

// ShouldFieldsPass returns true if the metric should pass, false if should drop
// based on the drop/pass filter parameters
func (f *InputFilter) ShouldNamePass(key string) bool {
	if f.name != nil {
		if f.name.Match(key) {
			return false
		}
	}
	return true
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
	if f.field != nil {
		if f.field.Match(key) {
			return false
		}
	}
	return true
}
