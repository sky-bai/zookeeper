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
}

// 添加服务器列表
func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params is empty")
	}

	r.rss = append(r.rss, params...)
	return nil
}

// 获得随机服务器IP
func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	r.curIndex = rand.Intn(len(r.rss))

	return r.rss[r.curIndex]
}
