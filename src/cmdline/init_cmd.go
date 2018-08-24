package cmd

import (
	"flag"
	"schema"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

func init() {
	configPath := flag.String("config", "conf/dev.toml", "specific config file")
	if _, err := toml.DecodeFile(*configPath, &schema.Config); err != nil {
		glog.Fatalf("Parser config error : %+v", err)
	}
	flag.Parse()
}
