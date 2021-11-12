package main

import (
	"fmt"
	"github.com/e421083458/gateway_demo/proxy/zookeeper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs1.Run()

	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	rs2.Run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type RealServer struct {
	Addr string
}

func (r *RealServer) Run() {
	// 1. 创建一个http服务
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", r.PingHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)

	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}
	log.Fatal(server.ListenAndServe())

	// 2.为每一个服务注册zk节点
	go func() {
		// 每个下游服务器需要注册zk节点
		zkManager := zookeeper.NewZkManager([]string{"121.196.163.8:2181"})
		err := zkManager.GetConnect()
		if err != nil {
			fmt.Printf("connect zk error : %s\n", err)
		}
		defer zkManager.Close()

		zkManager.RegistServerPath("/real_server", r.Addr)
		if err != nil {
			fmt.Printf("register node error : %s", err)
		}
		zList, err := zkManager.GetServerListByPath("/real_server")
		fmt.Println(zList)
	}()
}

func (r *RealServer) PingHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("internal server error"))
}
