package main

import (
	"fmt"
	"github.com/e421083458/gateway_demo/proxy/zookeeper"
	"time"
)

func main() {
	zkManager := zookeeper.NewZkManager([]string{"121.196.163.8:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()

	i := 0
	for {
		conf := fmt.Sprintf("{name:" + fmt.Sprint(i) + "}")
		fmt.Println(conf)
		zkManager.SetPathData("/rs_server_conf", []byte(conf), int32(i))
		time.Sleep(5 * time.Second)
		i++
	}
}
