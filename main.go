package main

import (
	"fmt"
	"strconv"
)

func main() {
	userID := 1728362829
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	fmt.Println(userIDStr)
}
