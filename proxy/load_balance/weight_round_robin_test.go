package load_balance

import (
	"fmt"
	"testing"
)

func TestWeightRoundRobinBalance_Add(t *testing.T) {
	wrb := WeightRoundRobinBalance{}

	wrb.Add("127.0.0.1:2003", "4")
	wrb.Add("127.0.0.1:2004", "3")
	wrb.Add("127.0.0.1:2005", "2")

	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
	fmt.Println(wrb.Next())
}
