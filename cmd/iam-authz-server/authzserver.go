package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/rose839/IAM/internal/authzserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	authzserver.NewApp("iam-authz-server").Run()
}
