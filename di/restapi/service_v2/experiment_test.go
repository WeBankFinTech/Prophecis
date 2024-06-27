package service_v2

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	fmt.Println(nextVersion("v1"))
	fmt.Println(nextVersion("v1.1"))
	fmt.Println(nextVersion("v1.3.5"))
}
