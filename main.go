package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"

	"github.com/benjaminbear/xmltv-proxy/config"
)

var Version = "undefined"

func main() {
	fmt.Println("Version:", Version)

	// Parse config from environment
	conf, err := config.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	runtime := &RunTime{
		Config: conf,
	}

	// Start cron once
	err = runtime.EPGCron()
	if err != nil {
		fmt.Println(err)
	}

	// Add cron process
	c := cron.New(cron.WithLocation(runtime.Config.TimeZone))
	c.AddFunc("12 6 * * *", func() {
		err := runtime.EPGCron()
		if err != nil {
			fmt.Println(err)
		}
	})
	c.Start()

	// Run Server
	http.HandleFunc("/", runtime.XMLTvServer)
	fmt.Println("Webserver ready.")

	err = http.ListenAndServe(fmt.Sprintf(":%d", runtime.Config.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
