package main

import (
	"fmt"
	"github.com/e421083458/gateway_demo/proxy/zookeeper"
	"log"
	"time"
)

func main() {
	zkManager := zookeeper.NewZkManager([]string{"121.196.163.8:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()

	// 第一次同步获取节点内容
	zc, _, err := zkManager.GetPathData("/rs_server_conf")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("get node data :")
	fmt.Println(string(zc))

	// 第二次异步地监听节点内容
	dataChan, dataErrChan := zkManager.WatchPathData("/rs_server_conf")
	go func() {
		for {
			select {
			case changeData := <-dataChan:
				fmt.Println("watchGetData changed :")
				fmt.Println(string(changeData))
			case changeErr := <-dataErrChan:
				fmt.Println("changeErr :")
				fmt.Println(changeErr)
			}
		}
	}()
	time.Sleep(100 * time.Second)
}
