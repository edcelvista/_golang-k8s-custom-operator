package lib

import (
	"math/rand"
)

func GetRandomAge() (randn int) {
	randn = rand.Intn(100)
	return
}
