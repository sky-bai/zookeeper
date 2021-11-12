package load_balance

import (
	"fmt"
	"testing"
)

func TestRoundRobinBalance_Add(t *testing.T) {
	rb := &RoundRobinBalance{}
	rb.Add("127.0.0.1:2003")
	rb.Add("127.0.0.1:2004")
	rb.Add("127.0.0.1:2005")
	rb.Add("127.0.0.1:2006")
	rb.Add("127.0.0.1:2007")
	rb.Add("127.0.0.1:2008")

	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())
	fmt.Println(rb.Next())

}
