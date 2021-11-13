package load_balance

import (
	"fmt"
	"testing"
)

func TestNewConsistentHashBalance(t *testing.T) {
	chb := NewConsistentHashBalance(10, nil)

	chb.Add("127.0.0.1:2003")
	chb.Add("127.0.0.1:2004")
	chb.Add("127.0.0.1:2005")
	chb.Add("127.0.0.1:2006")
	chb.Add("127.0.0.1:2007")

	// url hash
	fmt.Println(chb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(chb.Get("http://127.0.0.1:2002/base/error"))
	fmt.Println(chb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(chb.Get("http://127.0.0.1:2002/base/changepwd"))

	// ip hash
	fmt.Println(chb.Get("127.0.0.1"))
	fmt.Println(chb.Get("192.168.0.1"))
	fmt.Println(chb.Get("127.0.0.1"))

}
