package pusher

import (
	"github.com/pusher/pusher-http-go/v5"
	"github.com/stelofinance/api/tools"
)

var PusherClient pusher.Client

func ConnectClient() {
	PusherClient = pusher.Client{
		Host:    tools.EnvVars.Pusher.Host,
		AppID:   tools.EnvVars.Pusher.AppId,
		Key:     tools.EnvVars.Pusher.AppKey,
		Secret:  tools.EnvVars.Pusher.AppSecret,
		Secure:  tools.EnvVars.ProductionEnv,
		Cluster: "mt1",
	}
}
