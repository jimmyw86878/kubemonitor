package gateway

import (
	"fmt"
	"kubemonitor/internal/kubeutil"
	"kubemonitor/internal/logger"
	"kubemonitor/internal/util"
	"os"
	"os/signal"
	"time"
)

var (
	//Ticker define
	Ticker *time.Ticker
	//TickerDone define
	TickerDone = make(chan bool)
	//Store define, this is main storage for pod status
	Store *kubeutil.Store
)

//StartServeK8S for demonstration of K8S
func StartServeK8S(c chan os.Signal) {
	//init logger
	logger.NewLogger()
	//init store
	Store = kubeutil.InitStore()
	if Store == nil {
		os.Exit(1)
	}
	Ticker = time.NewTicker(time.Duration(util.LoadInt64Env("checkPeriodSec", 30)) * time.Second)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("Shutting down...")
			Ticker.Stop()
			TickerDone <- true
		}
	}()
	//Start ticker
	startTicker()
}

func startTicker() {
	for {
		select {
		case <-TickerDone:
			return
		case t := <-Ticker.C:
			fmt.Println("-------------------------")
			fmt.Println("Check at", t)
			errArr := Store.CheckAndUpdatePodStat()
			if len(errArr) == 0 {
				fmt.Println("Check all pods status successfully")
			}
			fmt.Println("-------------------------")

		}
	}
}
