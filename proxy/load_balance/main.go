package load_balance

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

var (
	addr      = "127.0.0.1:2002"
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, //连接超时
			KeepAlive: 30 * time.Second, //长连接超时时间
		}).DialContext,
		MaxIdleConns:          100,              //最大空闲连接
		IdleConnTimeout:       90 * time.Second, //空闲超时时间
		TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  //100-continue状态码超时时间
	}
)

func LoadBalanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	case LbConsistentHash:
		return &ConsistentHashBalance{}
	default:
		return &RandomBalance{}
	}

}

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	// Update 后期服务发现补充
	Update()
}

func main() {
	lb := LoadBalanceFactory(LbConsistentHash)
	err := lb.Add("http://127.0.0.1:2003/base", "10")
	if err != nil {
		fmt.Println(err)
	}
	err = lb.Add("http://127.0.0.1:2004/base", "10")
	if err != nil {
		fmt.Println(err)
	}

	proxy := NewMultipleHostsReverseProxy(lb)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))

}

func NewMultipleHostsReverseProxy(lb LoadBalance) *httputil.ReverseProxy {
	// 请求协助者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil {
			fmt.Println("get next err", err)
		}

		target, err := url.Parse(nextAddr)
		if err != nil {
			fmt.Println("parse next err", err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}

	}

	// 更改内容
	modifyFunc := func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			// 获取内容
			oldPayload, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			// 追加内容
			newPayload := []byte("StatusCode error : " + string(oldPayload))

			resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
			resp.ContentLength = int64(len(newPayload))
			resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(newPayload)), 10))
		}
		return nil
	}

	// 错误回调 : 关闭真实real_server 时测试,错误回调
	// 范围 : transport.RoundTrip发生的错误，以及ModifyResponse发生的错误

	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "ErrorHandler error :"+err.Error(), 500)
	}
	return &httputil.ReverseProxy{Director: director, Transport: transport, ModifyResponse: modifyFunc, ErrorHandler: errFunc}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// 反向代理 目的 是为将请求转发到各个服务器上面
// 负载均衡 目的 获取服务器地址  是为负载（工作任务，访问请求）进行平衡、分摊到多个操作单元（服务器，组件）上进行执行

// 负载均衡 读取监听列表
