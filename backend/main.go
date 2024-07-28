// PizTec Corporation, 2024. All Rights Reserved.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

const API_ROOT = "/api/"
const AUTH_TOKEN = "Payrus-Auth-Token"

var loggedAccounts map[string]*Account = make(map[string]*Account)

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
	webMux.Handle("/", http.FileServer(http.Dir(webRootDir)))
	webMux.HandleFunc(API_ROOT, apiHandler)

	go func() {
		port := Config.GetString("server.web_port")
		http.ListenAndServe(":"+port, webMux)

		stopChan <- true
	}()

	return nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2) // simulated delay in response

	if r.Method == http.MethodPost && r.URL.Path == path.Join(API_ROOT, "login") {
		processLoginRequest(w, r)

		return
	}
	if r.Method == http.MethodPost && r.URL.Path == path.Join(API_ROOT, "logout") {
		processLogoutRequest(w, r)

		return
	}
	if r.Method == http.MethodPost && r.URL.Path == path.Join(API_ROOT, "create_account") {
		processCreateAccountRequest(w, r)

		return
	}
	if r.Method == http.MethodPost && r.URL.Path == path.Join(API_ROOT, "create_card") {
		processCreateCardRequest(w, r)

		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(fmt.Sprintf("Unsupported request: method=%s, path=%s", r.Method, r.RequestURI)))
}

func processLoginRequest(w http.ResponseWriter, r *http.Request) {
	creds, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read login credentials"))
		return
	}

	loginCreds := &AccountCredentials{}
	err = json.Unmarshal(creds, loginCreds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Login credentials are not provided"))
		return
	}

	account := LoadAccount(loginCreds)

	if account == nil {
		Logger.Warning("Login request for %s - account NOT found", loginCreds.Email)

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Login credentials are incorrect"))
		return
	}

	Logger.Info("Login request for %s - account located", loginCreds.Email)

	sessionToken := createSessionToken()
	loggedAccounts[sessionToken] = account

	accountAsJson, _ := json.Marshal(account)
	w.Header().Set(AUTH_TOKEN, sessionToken)
	w.WriteHeader(http.StatusOK)
	w.Write(accountAsJson)
}

func processLogoutRequest(w http.ResponseWriter, r *http.Request) {
	account, token := findLoggedAccount(w, r)
	if account != nil {
		delete(loggedAccounts, token)

		Logger.Info("Logout request for %s", account.Credentials.Email)
		w.WriteHeader(http.StatusOK)
	}
}

func processCreateAccountRequest(w http.ResponseWriter, r *http.Request) {
	creds, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read account credentials"))
		return
	}

	accountCreds := &AccountCredentials{}
	err = json.Unmarshal(creds, accountCreds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Account credentials are not provided"))
		return
	}

	account := LoadAccount(accountCreds)
	if account != nil {
		Logger.Warning("Account creation request for an existing account %s", accountCreds.Email)

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Account already exists"))
		return
	}

	account, err = CreateAccount(accountCreds)
	if err != nil {
		Logger.Warning("Account creation failed for %s: %v", accountCreds.Email, err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Account creation failed for an unknown reason"))

		return
	}

	Logger.Info("Account %s is successfully created", accountCreds.Email)

	accountAsJson, _ := json.Marshal(account)
	w.WriteHeader(http.StatusOK)
	w.Write(accountAsJson)
}

func processCreateCardRequest(w http.ResponseWriter, r *http.Request) {
	cardInfo, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read card details"))
		return
	}

	cardRequest := &Card{}
	err = json.Unmarshal(cardInfo, cardRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Card details are not provided"))
		return
	}

	if cardRequest.FirstName == "" || cardRequest.LastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("First or last name is not provided"))
		return
	}

	account, _ := findLoggedAccount(w, r)
	if account == nil {
		return
	}

	Logger.Info("Processing card creation request for account %s", account.Credentials.Email)

	err = CreateCard(account, cardRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to create a card"))

		Logger.Error("Failure creating card for account %s: %v", account.Credentials.Email, err)
		return
	}

	accountAsJson, _ := json.Marshal(account)
	w.WriteHeader(http.StatusOK)
	w.Write(accountAsJson)

	Logger.Info("Card created for account %s", account.Credentials.Email)
}

func findLoggedAccount(w http.ResponseWriter, r *http.Request) (*Account, string) {
	token := r.Header.Get(AUTH_TOKEN)
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Auth token is not present"))
		return nil, ""
	}

	account := loggedAccounts[token]
	if account == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid or expired auth token"))
		return nil, ""
	}

	return account, token
}

func createSessionToken() string {
	return fmt.Sprintf("token_%d", time.Now().UnixMilli())
}
