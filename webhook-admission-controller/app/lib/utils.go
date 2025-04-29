package lib

import (
	"log"
	"os"
	"strings"
)

var debug = false

func Debug(v ...interface{}) {
	if debug {
		log.Println(v...)
	}
}

func Init() {
	if strings.ToLower(os.Getenv("IS_DEBUG")) == "true" {
		debug = true
	}

	log.Printf("💡 Debug enabled: %v", debug)
}
