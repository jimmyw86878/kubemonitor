package main

import (
	"kubemonitor/internal/gateway"
	"os"
)

func main() {
	gateway.StartServeK8S(make(chan os.Signal, 1))
}
