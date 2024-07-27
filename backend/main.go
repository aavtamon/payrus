// PizTec Corporation, 2024. All Rights Reserved.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var cfgFile string
	flag.StringVar(&cfgFile, "cfg", "", "Configuration file")
	flag.Parse()

	if cfgFile == "" {
		Logger.Error("Config file is not specified. Exiting...")
		os.Exit(1)
	}

	err := InitConfig(cfgFile)
	if err != nil {
		Logger.Error("Unable to read config file %s: %v", cfgFile, err)
		os.Exit(1)
	}

	stopChan := make(chan bool, 1)

	// Create the SIGTERM listener.
	stop := make(chan bool, 1)
	go func() {
		termChannel := make(chan os.Signal, 1)
		signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM)
		<-termChannel
		stopChan <- true

		Logger.Info("Payrus is stopping")
		stop <- true
	}()

	err = startWebServer(stopChan)
	if err != nil {
		Logger.Error("Unable to start web server: %v", err)
		os.Exit(1)
	}

	Logger.Info("Payrus is ready")

	// We are up and running so wait until we are
	// told to stop.
	<-stop

	Logger.Info("Payrus exited")
}

func startWebServer(stopChan chan bool) error {
	webRootDir := Config.GetString("server.web_root")
	_, err := os.Stat(webRootDir)
	if err != nil {
		return fmt.Errorf("web server cannot be started - WebRoot directory %s is not found: %v", webRootDir, err)
	}

	webMux := http.NewServeMux()
	fmt.Println(">>> ", http.Dir(webRootDir))
	webMux.Handle("/", http.FileServer(http.Dir(webRootDir)))
	webMux.HandleFunc("/api/", apiHandler)

	go func() {
		port := Config.GetString("server.web_port")
		http.ListenAndServe(":"+port, webMux)

		stopChan <- true
	}()

	return nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
}
