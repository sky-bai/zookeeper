package main

import (
	"fmt"
	"github.com/e421083458/gateway_demo/proxy/zookeeper"
	"log"
	"time"
)

func main() {
	// connect zk
	zkManager := zookeeper.NewZkManager([]string{"121.196.163.8:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()

	// 获取某个路径下的服务列表
	zlist, err := zkManager.GetServerListByPath("/go")
	fmt.Println("server node:")
	fmt.Println(zlist)
	if err != nil {
		log.Println(err)
	}

	// 动态监听节点变化 获取某个路径下的服务列表
	chanList, chanErr := zkManager.WatchServerListByPath("/go")
	fmt.Println("节点下游服务器列表", chanList)
	go func() {
		for {
			select {
			case changeErr := <-chanErr:
				// 获取子节点失败直接就让监听程序
				panic(changeErr)

			case changeList := <-chanList:
				fmt.Println("chanList")
				fmt.Println(changeList)
			}

		}
	}()
	time.Sleep(time.Second * 10)

	// 获取到节点到下游服务器列表之后，就可以做负载均衡了
}
