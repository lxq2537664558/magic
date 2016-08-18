package misc

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strconv"
	"time"
)

const alphanum string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Duration just wraps time.Duration1
type Duration struct {
	Duration time.Duration
}

// UnmarshalTOML parses the duration from the TOML config file
func (d *Duration) UnmarshalTOML(b []byte) error {
	var err error
	// Parse string duration, ie, "1s"
	d.Duration, err = time.ParseDuration(string(b[1 : len(b)-1]))
	if err == nil {
		return nil
	}

	// First try parsing as integer seconds
	sI, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		d.Duration = time.Second * time.Duration(sI)
		return nil
	}
	// Second try parsing as float seconds
	sF, err := strconv.ParseFloat(string(b), 64)
	if err == nil {
		d.Duration = time.Second * time.Duration(sF)
		return nil
	}

	return nil
}

func PrintStack(all bool) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, all)

	log.Println("[FATAL] catch a panic,stack is: ", string(buf[:n]))
}

// RandomString returns a random string of alpha-numeric characters
func RandomString(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

// GetTLSConfig gets a tls.Config object from the given certs, key, and CA files.
// you must give the full path to the files.
// If all files are blank and InsecureSkipVerify=false, returns a nil pointer.
func GetTLSConfig(
	SSLCert, SSLKey, SSLCA string,
	InsecureSkipVerify bool,
) (*tls.Config, error) {
	if SSLCert == "" && SSLKey == "" && SSLCA == "" && !InsecureSkipVerify {
		return nil, nil
	}

	t := &tls.Config{
		InsecureSkipVerify: InsecureSkipVerify,
	}

	if SSLCA != "" {
		caCert, err := ioutil.ReadFile(SSLCA)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not load TLS CA: %s",
				err))
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		t.RootCAs = caCertPool
	}

	if SSLCert != "" && SSLKey != "" {
		cert, err := tls.LoadX509KeyPair(SSLCert, SSLKey)
		if err != nil {
			return nil, errors.New(fmt.Sprintf(
				"Could not load TLS client key/certificate from %s:%s: %s",
				SSLKey, SSLCert, err))
		}

		t.Certificates = []tls.Certificate{cert}
		t.BuildNameToCertificate()
	}

	// will be nil by default if nothing is provided
	return t, nil
}
