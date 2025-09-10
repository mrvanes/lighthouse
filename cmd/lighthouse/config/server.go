package config

import (
	"github.com/go-oidfed/lighthouse"
)

var defaultServerConf = lighthouse.ServerConf{
	Address:           "127.0.0.1:7672",
	ForwardedIPHeader: "X-Forwarded-For",
}
