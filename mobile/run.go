package frp

import (
	"encoding/json"
	"fmt"
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/client/proxy"
	"github.com/fatedier/frp/pkg/auth"
	"github.com/fatedier/frp/pkg/config"
	"github.com/fatedier/frp/pkg/consts"
	"github.com/fatedier/frp/pkg/util/log"
	"time"
)

var (
	globalService *client.Service = nil
)

type FrpCallback client.FrpCallback

type FrpConfig struct {
	ServerAddr   string
	ServerPort   int
	ServerToken  string
	RemotePort   int
	HttpUser     string
	HttpPassword string
}

func GetStatus() string {
	if globalService != nil {
		if out, err := json.Marshal(struct {
			Status []*proxy.WorkingStatus `json:"status"`
		}{Status: globalService.GetController().Status()}); err == nil {
			return string(out)
		} else {
			return "ERROR"
		}
	}
	return "SERVICE IS NIL"
}

func Stop() {
    if globalService != nil {
        globalService.GracefulClose(time.Millisecond * 250)
    }
    //os.Exit(1)
}

func RunFRP(cb FrpCallback, conf *FrpConfig) error {
	log.InitLog("console", "", "error", 0, true)
	c := config.GetDefaultClientConf()
	c.ServerAddr = conf.ServerAddr
	c.ServerPort = conf.ServerPort
	c.AuthenticationMethod = "token"
	c.DialServerTimeout = 6
	c.TokenConfig = auth.TokenConfig{Token: conf.ServerToken}
	name := fmt.Sprintf("androfrp_port_%d", conf.RemotePort)
	service, err := client.NewService(c,
		map[string]config.ProxyConf{name: &config.TCPProxyConf{
			RemotePort: conf.RemotePort,
			BaseProxyConf: config.BaseProxyConf{
				ProxyName: name,
				ProxyType: consts.TCPProxy,
				LocalSvrConf: config.LocalSvrConf{
					LocalIP: "127.0.0.1",
					Plugin:  "http_proxy",
					PluginParams: map[string]string{
						"plugin_http_user":   conf.HttpUser,
						"plugin_http_passwd": conf.HttpPassword,
					},
				},
			},
		}}, map[string]config.VisitorConf{}, "", cb)
	if err != nil {
		return err
	}
	globalService = service
	return service.Run()
}
