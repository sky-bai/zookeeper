package load_balance

import (
	"errors"
	"math/rand"
	"zookeeper/load_balance"
)

type RandomBalance struct {
	curIndex int
	rss      []string

	// 这里我需要维护下游服务器列表
	conf load_balance.LoadBalanceConf
	// 这里定义了一个关于负载均衡配置的接口 它可以获得1.服务器配置 2.更新服务器配置 也就是说这里抽象出一组方法 不管具体实现 让其他负载均衡配置来实现这个接口
	// 每个负载均衡配置方法内部都有一个负载均衡的配置
	// 它可以获取服务器列表
}

func (r *RandomBalance) Get(s string) (string, error) {
	panic("implement me")
}

func (r *RandomBalance) Update() {
	panic("implement me")
}

// Add 添加服务器列表
func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params is empty")
	}

	r.rss = append(r.rss, params...)
	return nil
}

// Next 获得随机服务器IP
func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss)) //nolint:gosec
	return r.rss[r.curIndex]
}
