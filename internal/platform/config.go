package platform

import (
	"crypto/tls"
	"flag"
)

type config struct {
	Addr      string
	Dsn       string
	TestDsn   string
	TlsDir    string
	TlsConfig *tls.Config
}

func (c *config) Init() {
	flag.StringVar(&c.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&c.TlsDir, "tls-dir", "./tls", "Path to tls directory")
	flag.StringVar(&c.Dsn, "dsn", "user:password@/db?parseTime=true", "MySQL data source name")
	flag.StringVar(&c.TestDsn, "test-dsn", "test:pass@/test?parseTime=true", "MySQL test data source name")

	c.TlsConfig = &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
}

var Config = &config{}
