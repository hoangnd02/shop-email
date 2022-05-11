package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	return fmt.Sprintf("%v", rand.Intn(max-min+1)+min)
}
