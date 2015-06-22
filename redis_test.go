package hera

import (
	"fmt"
	"testing"
)

func TestDofunc(t *testing.T) {
	NewRedisSvc()
	tmp, _ := Redis.DoCmd("smembers", "userid")
	fmt.Printf("%s", tmp)
}
