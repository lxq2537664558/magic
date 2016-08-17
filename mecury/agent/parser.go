package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/influxdata/toml/ast"
)

// ParseConfig is a struct that covers the data types needed for all parser types,
// and can be used to instantiate _any_ of the parsers.
type ParseConfig struct {
	// Dataformat can be one of: json, influx, graphite, value, nagios
	DataFormat string

	// Separator only applied to Graphite data.
	Separator string
	// Templates only apply to Graphite data.
	Templates []string

	// TagKeys only apply to JSON data
	TagKeys []string
	// MetricName applies to JSON & value. This will be the name of the measurement.
	MetricName string

	// DataType only applies to value, this will be the type to parse value to
	DataType string

	// DefaultTags are the default tags that will be added to all parsed metrics.
	DefaultTags map[string]string
}

// ParserInput is an interface for input plugins that are able to parse
// arbitrary data formats.
type ParserInput interface {
	// SetParser sets the parser function for the interface
	SetParser(parser Parser)
}

// Parser is an interface defining functions that a parser plugin must satisfy.
type Parser interface {
	// Parse takes a byte buffer separated by newlines
	// ie, `cpu.usage.idle 90\ncpu.usage.busy 10`
	// and parses it into telegraf metrics
	Parse(buf []byte) ([]Metric, error)

	// ParseLine takes a single string metric
	// ie, "cpu.usage.idle 90"
	// and parses it into a telegraf metric.
	ParseLine(line string) (Metric, error)

	// SetDefaultTags tells the parser to add all of the given tags
	// to each parsed metric.
	// NOTE: do _not_ modify the map after you've passed it here!!
	SetDefaultTags(tags map[string]string)
}

// NewParser returns a Parser interface based on the given config.
func NewParser(config *ParseConfig) Parser {
	var parser Parser
	var err error
	switch config.DataFormat {
	case "json":
		parser, err = NewJSONParser(config.MetricName,
			config.TagKeys, config.DefaultTags)
	case "influx":
		parser, err = NewInfluxParser()
	default:
		log.Fatalf("Invalid data format: %s", config.DataFormat)
	}

	if err != nil {
		log.Fatalln("[FATAL] create parser : ", err)
	}
	return parser
}

// buildParser grabs the necessary entries from the ast.Table for creating
// a parsers.Parser object, and creates it, which can then be added onto
// an Input object.
func parserInit(name string, tbl *ast.Table) Parser {
	c := &ParseConfig{}

	df, ok := tbl.Fields["data_format"]
	if ok {
		kv, ok := df.(*ast.KeyValue)
		if ok {
			v, ok := kv.Value.(*ast.String)
			if ok {
				c.DataFormat = v.Value
			}
		}
	}
	if c.DataFormat == "" {
		log.Fatalln("[FATAL] a input parse init without data_format found")
	}

	// Legacy support, exec plugin originally parsed JSON by default.
	if name == "exec" && c.DataFormat == "" {
		c.DataFormat = "json"
	} else if c.DataFormat == "" {
		c.DataFormat = "influx"
	}

	if node, ok := tbl.Fields["separator"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				c.Separator = str.Value
			}
		}
	}

	if node, ok := tbl.Fields["templates"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						c.Templates = append(c.Templates, str.Value)
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["tag_keys"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if ary, ok := kv.Value.(*ast.Array); ok {
				for _, elem := range ary.Value {
					if str, ok := elem.(*ast.String); ok {
						c.TagKeys = append(c.TagKeys, str.Value)
					}
				}
			}
		}
	}

	if node, ok := tbl.Fields["data_type"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.String); ok {
				c.DataType = str.Value
			}
		}
	}

	c.MetricName = name

	delete(tbl.Fields, "data_format")
	delete(tbl.Fields, "separator")
	delete(tbl.Fields, "templates")
	delete(tbl.Fields, "tag_keys")
	delete(tbl.Fields, "data_type")

	return NewParser(c)
}

func NewJSONParser(
	metricName string,
	tagKeys []string,
	defaultTags map[string]string,
) (Parser, error) {
	parser := &JSONParser{
		MetricName:  metricName,
		TagKeys:     tagKeys,
		DefaultTags: defaultTags,
	}
	return parser, nil
}

// JSONParser
type JSONParser struct {
	MetricName  string
	TagKeys     []string
	DefaultTags map[string]string
}

func (p *JSONParser) Parse(buf []byte) ([]Metric, error) {
	metrics := make([]Metric, 0)

	var jsonOut map[string]interface{}
	err := json.Unmarshal(buf, &jsonOut)
	if err != nil {
		err = fmt.Errorf("unable to parse out as JSON, %s", err)
		return nil, err
	}

	tags := make(map[string]string)
	for k, v := range p.DefaultTags {
		tags[k] = v
	}

	for _, tag := range p.TagKeys {
		switch v := jsonOut[tag].(type) {
		case string:
			tags[tag] = v
		}
		delete(jsonOut, tag)
	}

	f := JSONFlattener{}
	err = f.FlattenJSON("", jsonOut)
	if err != nil {
		return nil, err
	}

	metric, err := NewMetric(p.MetricName, tags, f.Fields, time.Now().UTC())

	if err != nil {
		return nil, err
	}
	return append(metrics, metric), nil
}

func (p *JSONParser) ParseLine(line string) (Metric, error) {
	metrics, err := p.Parse([]byte(line + "\n"))

	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf("Can not parse the line: %s, for data format: influx ", line)
	}

	return metrics[0], nil
}

func (p *JSONParser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}

type JSONFlattener struct {
	Fields map[string]interface{}
}

// FlattenJSON flattens nested maps/interfaces into a fields map
func (f *JSONFlattener) FlattenJSON(
	fieldname string,
	v interface{},
) error {
	if f.Fields == nil {
		f.Fields = make(map[string]interface{})
	}
	fieldname = strings.Trim(fieldname, "_")
	switch t := v.(type) {
	case map[string]interface{}:
		for k, v := range t {
			err := f.FlattenJSON(fieldname+"_"+k+"_", v)
			if err != nil {
				return err
			}
		}
	case []interface{}:
		for i, v := range t {
			k := strconv.Itoa(i)
			err := f.FlattenJSON(fieldname+"_"+k+"_", v)
			if err != nil {
				return nil
			}
		}
	case float64:
		f.Fields[fieldname] = t
	case bool, string, nil:
		// ignored types
		return nil
	default:
		return fmt.Errorf("JSON Flattener: got unexpected type %T with value %v (%s)",
			t, t, fieldname)
	}
	return nil
}

func NewInfluxParser() (Parser, error) {
	return &InfluxParser{}, nil
}

// InfluxParser is an object for Parsing incoming metrics.
type InfluxParser struct {
	// DefaultTags will be added to every parsed metric
	DefaultTags map[string]string
}

// Parse returns a slice of Metrics from a text representation of a
// metric (in line-protocol format)
// with each metric separated by newlines. If any metrics fail to parse,
// a non-nil error will be returned in addition to the metrics that parsed
// successfully.
func (p *InfluxParser) Parse(buf []byte) ([]Metric, error) {
	// parse even if the buffer begins with a newline
	buf = bytes.TrimPrefix(buf, []byte("\n"))
	points, err := misc.ParsePoints(buf)
	metrics := make([]Metric, len(points))
	for i, point := range points {
		tags := point.Tags()
		for k, v := range p.DefaultTags {
			// Only set tags not in parsed metric
			if _, ok := tags[k]; !ok {
				tags[k] = v
			}
		}
		// Ignore error here because it's impossible that a model.Point
		// wouldn't parse into client.Point properly
		metrics[i], _ = NewMetric(point.Name(), tags,
			point.Fields(), point.Time())
	}
	return metrics, err
}

func (p *InfluxParser) ParseLine(line string) (Metric, error) {
	metrics, err := p.Parse([]byte(line + "\n"))

	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf(
			"Can not parse the line: %s, for data format: influx ", line)
	}

	return metrics[0], nil
}

func (p *InfluxParser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}
