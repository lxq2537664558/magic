package log_collect

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"sync"

	"github.com/corego/vgo/mecury/agent"
	"github.com/corego/vgo/mecury/misc/errchan"
	"github.com/corego/vgo/mecury/misc/globpath"
	"github.com/corego/vgo/mecury/plugins/input/log_collect/parsers/raw"
	"github.com/hpcloud/tail"
	"github.com/syndtr/goleveldb/leveldb"
)

type LogParser interface {
	ParseLine(line string) (agent.Metric, error)
	Compile() error
}

type LogParserPlugin struct {
	Files         []string
	FromBeginning bool

	tailers []*tail.Tail
	lines   chan string
	done    chan struct{}
	wg      sync.WaitGroup
	acc     agent.Accumulator
	parsers []LogParser

	sync.Mutex

	RawParser *raw.Parser `toml:"raw"`
	localdb   *leveldb.DB
}

const sampleConfig = `
`

func (l *LogParserPlugin) SampleConfig() string {
	return sampleConfig
}

func (l *LogParserPlugin) Description() string {
	return "Stream and parse log file(s)."
}

func (l *LogParserPlugin) Gather(acc agent.Accumulator) error {
	return nil
}

func (l *LogParserPlugin) Start(acc agent.Accumulator) error {
	l.Lock()
	defer l.Unlock()

	l.localdb = initLocalDB()

	l.acc = acc
	l.lines = make(chan string, 1000)
	l.done = make(chan struct{})

	// Looks for fields which implement LogParser interface
	l.parsers = []LogParser{}
	s := reflect.ValueOf(l).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if !f.CanInterface() {
			continue
		}

		if lpPlugin, ok := f.Interface().(LogParser); ok {
			if reflect.ValueOf(lpPlugin).IsNil() {
				continue
			}
			l.parsers = append(l.parsers, lpPlugin)
		}
	}

	if len(l.parsers) == 0 {
		return fmt.Errorf("ERROR: logparser input plugin: no parser defined.")
	}

	// compile log parser patterns:
	errChan := errchan.New(len(l.parsers))
	for _, parser := range l.parsers {
		if err := parser.Compile(); err != nil {
			errChan.C <- err
		}
	}
	if err := errChan.Error(); err != nil {
		return err
	}

	var seek tail.SeekInfo
	if !l.FromBeginning {
		seek.Whence = 2
		seek.Offset = 0
	}

	l.wg.Add(1)
	go l.parser()

	// Create a "tailer" for each file
	for _, filepath := range l.Files {
		// get files by filepath
		g, err := globpath.Compile(filepath)
		if err != nil {
			log.Printf("ERROR Glob %s failed to compile, %s", filepath, err)
			continue
		}
		files := g.Match()

		errChan = errchan.New(len(files))

		for file := range files {
			// load last offset for each file
			seek.Offset = offset(file, l)

			tailer, err := tail.TailFile(file,
				tail.Config{
					ReOpen:    true,
					Follow:    true,
					Location:  &seek,
					MustExist: true,
				})
			errChan.C <- err

			// create a goroutine for each "tailer"
			l.wg.Add(1)
			go l.receiver(tailer)
			l.tailers = append(l.tailers, tailer)
		}
	}

	return errChan.Error()
}

func (l *LogParserPlugin) Stop() {
	l.Lock()
	defer l.Unlock()

	for _, t := range l.tailers {
		setOffset(t, l)
		err := t.Stop()
		if err != nil {
			log.Printf("ERROR stopping tail on file %s\n", t.Filename)
		}
		t.Cleanup()
	}
	close(l.done)
	l.wg.Wait()
}

func offset(file string, l *LogParserPlugin) int64 {
	var offset int64
	offsetB, err := l.localdb.Get([]byte(file), nil)
	if err != nil {
		offset = 0
	} else {
		oi, err := strconv.ParseInt(string(offsetB), 10, 64)
		if err != nil {
			log.Printf("[ERROR] conv %v's offset %v to int error:%v", file, offsetB, err)
		}
		offset = oi
	}

	return offset
}

func setOffset(t *tail.Tail, l *LogParserPlugin) {
	// record the newest offset of each tailer
	offset, _ := t.Tell()
	offsetS := strconv.FormatInt(offset, 10)
	l.localdb.Put([]byte(t.Filename), []byte(offsetS), nil)
}

func initLocalDB() *leveldb.DB {
	lvdb, err := leveldb.OpenFile("data.db", nil)
	if err != nil {
		log.Fatal("[FATAL] init local db error: ", err)
	}
	return lvdb
}

// receiver is launched as a goroutine to continuously watch a tailed logfile
// for changes and send any log lines down the l.lines channel.
func (l *LogParserPlugin) receiver(tailer *tail.Tail) {
	defer l.wg.Done()

	var line *tail.Line
	for line = range tailer.Lines {
		if line.Err != nil {
			log.Printf("ERROR tailing file %s, Error: %s\n",
				tailer.Filename, line.Err)
			continue
		}

		select {
		case <-l.done:
		case l.lines <- line.Text:
		}
	}
}

// parser is launched as a goroutine to watch the l.lines channel.
// when a line is available, parser parses it and adds the metric(s) to the
// accumulator.
func (l *LogParserPlugin) parser() {
	defer l.wg.Done()

	var m agent.Metric
	var err error
	var line string
	for {
		select {
		case <-l.done:
			return
		case line = <-l.lines:
			if line == "" || line == "\n" {
				continue
			}
		}

		for _, parser := range l.parsers {
			m, err = parser.ParseLine(line)
			if err == nil {
				if m != nil {
					l.acc.AddFields(m.Name(), m.Fields(), m.Tags(), m.Time())
				}
			}
		}
	}
}

func init() {
	agent.AddInput("log_collect", &LogParserPlugin{})
}
