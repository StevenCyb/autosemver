package utils

import (
	"log"
	"strconv"
)

func MustParseUint(s string) uint {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return uint(i)
}
