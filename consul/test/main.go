package main

import (
	"fmt"
	"os"

	"net/http"

	"github.com/hktalent/goutils/consul"
)

func consulApp() {
	healthHost := ":8081"
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		http.ListenAndServe(healthHost, nil)
	}()
	app, err := consul.NewApp(healthHost, "")
	if err != nil {
		fmt.Println("NEWAPP ERROR:", err)
		return
	}
	app.Wait(func(sig os.Signal) {
		fmt.Println("QUIT:", sig)
	})
}

type CFG1 struct {
}

type CFG2 struct {
	HealthHost string `yaml:"health"`
}

func consulAppWithConfig() {
	healthHost := ":8081"
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		http.ListenAndServe(healthHost, nil)
	}()
	var cfg CFG2
	app, err := consul.NewAppWithCfg(&cfg, healthHost)
	if err != nil {
		fmt.Println("NEWAPP ERROR:", err)
		return
	}
	app.Wait(func(sig os.Signal) {
		fmt.Println("QUIT:", sig)
	})
}

func main() {
	consulApp()
}
