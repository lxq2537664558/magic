package mongodb

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/aiyun/openapm/mecury/agent"
	"github.com/aiyun/openapm/mecury/misc/errchan"

	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	Servers          []string
	Ssl              Ssl
	mongos           map[string]*Server
	GatherPerdbStats bool
}

type Ssl struct {
	Enabled bool
	CaCerts []string `toml:"cacerts"`
}

var sampleConfig = `
  ## An array of URI to gather stats about. Specify an ip or hostname
  ## with optional port add password. ie,
  ##   mongodb://user:auth_key@10.10.3.30:27017,
  ##   mongodb://10.10.3.33:18832,
  ##   10.0.0.1:10000, etc.
  servers = ["127.0.0.1:27017"]
  gather_perdb_stats = false
`

func (m *MongoDB) SampleConfig() string {
	return sampleConfig
}

func (*MongoDB) Description() string {
	return "Read metrics from one or many MongoDB servers"
}

var localhost = &url.URL{Host: "127.0.0.1:27017"}

// Reads stats from all configured servers accumulates stats.
// Returns one of the errors encountered while gather stats (if any).
func (m *MongoDB) Gather(acc agent.Accumulator) error {
	if len(m.Servers) == 0 {
		m.gatherServer(m.getMongoServer(localhost), acc)
		return nil
	}

	var wg sync.WaitGroup
	errChan := errchan.New(len(m.Servers))
	for _, serv := range m.Servers {
		u, err := url.Parse(serv)
		if err != nil {
			return fmt.Errorf("Unable to parse to address '%s': %s", serv, err)
		} else if u.Scheme == "" {
			u.Scheme = "mongodb"
			// fallback to simple string based address (i.e. "10.0.0.1:10000")
			u.Host = serv
			if u.Path == u.Host {
				u.Path = ""
			}
		}
		wg.Add(1)
		go func(srv *Server) {
			defer wg.Done()
			errChan.C <- m.gatherServer(srv, acc)
		}(m.getMongoServer(u))
	}

	wg.Wait()
	return errChan.Error()
}

func (m *MongoDB) getMongoServer(url *url.URL) *Server {
	if _, ok := m.mongos[url.Host]; !ok {
		m.mongos[url.Host] = &Server{
			Url: url,
		}
	}
	return m.mongos[url.Host]
}

func (m *MongoDB) gatherServer(server *Server, acc agent.Accumulator) error {
	if server.Session == nil {
		var dialAddrs []string
		if server.Url.User != nil {
			dialAddrs = []string{server.Url.String()}
		} else {
			dialAddrs = []string{server.Url.Host}
		}
		dialInfo, err := mgo.ParseURL(dialAddrs[0])
		if err != nil {
			return fmt.Errorf("Unable to parse URL (%s), %s\n",
				dialAddrs[0], err.Error())
		}
		dialInfo.Direct = true
		dialInfo.Timeout = 5 * time.Second

		if m.Ssl.Enabled {
			tlsConfig := &tls.Config{}
			if len(m.Ssl.CaCerts) > 0 {
				roots := x509.NewCertPool()
				for _, caCert := range m.Ssl.CaCerts {
					ok := roots.AppendCertsFromPEM([]byte(caCert))
					if !ok {
						return fmt.Errorf("failed to parse root certificate")
					}
				}
				tlsConfig.RootCAs = roots
			} else {
				tlsConfig.InsecureSkipVerify = true
			}
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
				if err != nil {
					fmt.Printf("error in Dial, %s\n", err.Error())
				}
				return conn, err
			}
		}

		sess, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			fmt.Printf("error dialing over ssl, %s\n", err.Error())
			return fmt.Errorf("Unable to connect to MongoDB, %s\n", err.Error())
		}
		server.Session = sess
	}
	return server.gatherData(acc, m.GatherPerdbStats)
}

func init() {
	agent.AddInput("mongodb", &MongoDB{mongos: make(map[string]*Server)})
}
