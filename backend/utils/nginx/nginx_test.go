package nginx

import (
	"fmt"
	"testing"
)

func TestNginx(t *testing.T) {
	config, err := GetConfig("/opt/panelx/data/apps/nginx/nginx-new/conf/conf.d/panelx.cloud.conf")
	if err != nil {
		panic(err)
	}

	//server := config.FindServers()[0]
	//fmt.Println(server)
	//serverD := config.FindServers()[0]
	//serverD.AddListen("8989", false)

	fmt.Println(DumpConfig(config, IndentedStyle))
}
