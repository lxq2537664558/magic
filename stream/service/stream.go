package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aiyun/openapm/mecury/misc"
	"github.com/aiyun/openapm/proto"

	_ "net/http/pprof"

	"github.com/boltdb/bolt"
	"github.com/uber-go/zap"
)

var VLogger zap.Logger

type StreamConfig struct {
	InputerQueue          int
	WriterNum             int
	DisruptorBuffersize   int64
	DisruptorBuffermask   int64
	DisruptorReservations int64
	Dbname                string
	Bucketname            string
	GrpcAddr              string
	PprofAddr             string
}

func (sc *StreamConfig) Show() {
	VLogger.Debug("Show", zap.Int("@InputerQueue", sc.InputerQueue))
	VLogger.Debug("Show", zap.Int("@WriterNum", sc.WriterNum))
	VLogger.Debug("Show", zap.Int64("@DisruptorBuffersize", sc.DisruptorBuffersize))
	VLogger.Debug("Show", zap.Int64("@DisruptorBuffermask", sc.DisruptorBuffermask))
	VLogger.Debug("Show", zap.Int64("@DisruptorReservations", sc.DisruptorReservations))
	VLogger.Debug("Show", zap.String("@Dbname", sc.Dbname))
	VLogger.Debug("Show", zap.String("@Bucketname", sc.Bucketname))
}

// Stream struct
type Stream struct {
	stopPluginsChan chan bool
	metricChan      chan Metrics
	writer          *Writer
	controller      *Controller
	alarmer         *Alarmer
	db              *DB
	groups          *Groups
	hostsTogroups   *HostsToGroup
	grpc            *Grpc
}

var streamer *Stream

// New get new stream struct
func New() *Stream {
	stream := &Stream{}
	streamer = stream
	return stream
}

// Init init stream
func (s *Stream) Init() {
	s.stopPluginsChan = make(chan bool, 1)
	s.metricChan = make(chan Metrics, 1)

	// init disruptor
	s.controller = NewController()
	s.controller.Init(Conf.Stream.DisruptorBuffersize, Conf.Stream.DisruptorBuffermask, Conf.Stream.DisruptorReservations)

	// init alarmer
	s.alarmer = NewAlarm()
	s.alarmer.Init()

	// init db
	s.db = NewDB(Conf.Stream.Dbname, Conf.Stream.Bucketname)
	s.db.Init()

	// init groups
	s.groups = NewGroups()

	// init hostsTogroups
	s.hostsTogroups = NewHostsToGroup()

	// load alerts from db
	s.LoadGroupsAlert()

	// init grpc
	s.grpc = NewGrpc()
	s.grpc.Init(Conf.Stream.GrpcAddr)
}

func (s *Stream) ShowGroups() {
	s.groups.RLock()
	defer s.groups.RUnlock()

	for k, v := range s.groups.groups {
		log.Println("----------------------------------------------------------- start")
		log.Println(k)
		v.Show()
		log.Println("----------------------------------------------------------- end")
	}
}

func (s *Stream) AddGroup(group *proto.Group) error {
	return s.groups.AddGroup(group)
}

func (s *Stream) AddAlerts(alerts *proto.Alerts) error {
	return s.groups.AddAlerts(alerts)
}

func (s *Stream) AddUsers(users *proto.Users) error {
	return s.groups.AddUsers(users)
}
func (s *Stream) AddHosts(hosts *proto.Hosts) error {
	return s.groups.AddHosts(hosts)
}

func (s *Stream) ShowHostsToGroup() {
	s.hostsTogroups.Show()
}

func (s *Stream) LoadGroupsAlert() error {
	VLogger.Debug("DB", zap.String("@Dbname", s.db.dbname))
	ok, err := s.db.isExist()
	if !ok || err != nil {
		return err
	}

	s.db.db.View(func(tx *bolt.Tx) error {
		groupRaws := make(map[string]*GroupRaw)
		b := tx.Bucket(s.db.bucketName)
		c := b.Cursor()
		// loading history data
		for k, v := c.First(); k != nil; k, v = c.Next() {
			g := &GroupRaw{}
			json.Unmarshal(v, g)
			VLogger.Info("LoadGroupsAlert", zap.String("@key", string(k)), zap.Object("@GroupRaw", g))
			groupRaws[string(k)] = g
		}

		// load rawGroups into allGroup
		for gid, graw := range groupRaws {
			// VLogger.Info("Loading", zap.String("@gid", gid), zap.Object("@graw", graw))
			// 初始化该group和parent线上的所有group
			firstG, ok := s.groups.groups[gid] //sg.AllGroup[k]
			// 当前group一旦初始化过,直接返回
			if ok {
				continue
			}

			firstG = NewGroup()
			// 初始化alert静态数据变量
			for rule, alertstatic := range graw.AlertStatics {
				VLogger.Info("LoadGroupsAlert", zap.String("@Rule", rule), zap.Object("@Alertstatic", alertstatic))
				alert := NewAlert()
				alert.AlertSt = alertstatic
				firstG.Alerts[rule] = alert
			}
			firstG.Users = graw.Users
			firstG.ID = graw.ID
			firstG.Hosts = graw.Hosts

			// 添加host和group关联关系
			for hostname, _ := range graw.Hosts {
				s.hostsTogroups.Add(hostname, graw.ID)
			}
			s.groups.groups[gid] = firstG
			parent := graw.Parent
			g := firstG
			// 设置父节点信息
			for {
				// 如果没有parent，则返回
				if parent == "" {
					break
				}
				// 寻找父亲
				pr1, ok := groupRaws[parent]
				if !ok {
					log.Fatal("can find parent, parent name is ", parent)
					return fmt.Errorf("can find parent, parent name is %s ", parent)
				}

				pp1, ok := s.groups.groups[parent]
				// 若第一个父亲已经存在，指针赋值后直接返回
				if ok {
					// 增加子节点
					pp1.AddChild(g.ID)
					g.Parent = pp1
					// 这里要删除父节点和主机之间的关联
					break
				}
				pp1 = NewGroup()
				for rule, alertstatic := range pr1.AlertStatics {
					alert := NewAlert()
					alert.AlertSt = alertstatic
					pp1.Alerts[rule] = alert
				}

				// 添加host和group关联关系
				for hostname, _ := range pr1.Hosts {
					s.hostsTogroups.Add(hostname, pr1.ID)
				}

				pp1.Users = pr1.Users
				pp1.ID = pr1.ID
				g.Parent = pp1
				s.groups.groups[parent] = pp1
				parent = pr1.Parent
				// 增加子节点
				pp1.AddChild(g.ID)
				g = pp1
			}
		}
		// free mem
		// load rawGroups into allGroup
		for k, _ := range groupRaws {
			delete(groupRaws, k)
		}

		return nil
	})

	return nil
}

// Start start stream server
func (s *Stream) Start(shutdown chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			misc.PrintStack(false)
			VLogger.Fatal("Stream fatal error ", zap.Error(err.(error)))
		}
	}()

	pprofStart()

	s.controller.Start()

	s.alarmer.Start()

	go s.grpc.Start()

	// start plugins service
	for _, c := range Conf.Inputs {
		c.Start(s.stopPluginsChan, s.metricChan)
	}

	for _, c := range Conf.Outputs {
		if err := c.Output.Start(); err != nil {
			VLogger.Panic("Output", zap.String("@Name", c.Name), zap.Error(err))
		}
	}

	for _, c := range Conf.Chains {
		c.Start(s.stopPluginsChan)
	}

	for _, c := range Conf.MetricOutputs {
		c.Start(s.stopPluginsChan)
	}
}

// Close close stream server
func (s *Stream) Close() error {
	log.Println("Stream close!")
	close(s.stopPluginsChan)
	close(s.metricChan)

	if err := s.controller.Close(); err != nil {
		VLogger.Error("Close", zap.String("@controller", err.Error()))
	}

	if err := s.alarmer.Close(); err != nil {
		VLogger.Error("Close", zap.String("@alarmer", err.Error()))
	}

	if err := s.grpc.Close(); err != nil {
		VLogger.Error("Close", zap.String("@grpc", err.Error()))
	}

	return nil
}

// pprof
func pprofStart() {
	flag.Parse()
	go func() {
		log.Println(http.ListenAndServe(Conf.Stream.PprofAddr, nil))
	}()
}
