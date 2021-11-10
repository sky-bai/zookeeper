package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

var host = []string{"121.196.163.8:2181"}

func main() {
	conn, _, err := zk.Connect(host, time.Second*5)
	if err != nil {
		fmt.Println("连接zookeeper失败", err)
	}
	defer conn.Close()
	fmt.Println("连接zookeeper成功")

	// 创建节点
	_, err = conn.Create("/go", []byte("hello"), 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Println("create err", err)
	}
	fmt.Println("创建节点成功")
	// 验证节点是否存在
	exists, _, err := conn.Exists("/go")
	if err != nil {
		fmt.Println("exists err", err)
	}
	if exists {
		fmt.Println("节点存在")
	} else {
		fmt.Println("节点不存在")
	}

	// 获取节点
	data, _, err := conn.Get("/go")
	if err != nil {
		fmt.Println("获取节点 err", err)
	}
	fmt.Println("获取节点 内容为", string(data))

	// 删除节点
	err = conn.Delete("/go", -1)
	if err != nil {
		fmt.Println("删除节点 err", err)
	}
	fmt.Println("删除节点成功")
}
