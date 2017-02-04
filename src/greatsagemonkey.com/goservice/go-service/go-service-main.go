/*
 *
 */

package main

import (
	//"encoding/json"
	"github.com/gorilla/mux"
	"greatsagemonkey.com/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const kernelSleepInterval=5*time.Minute

const configName="go-service-config"
const runtimeVersion="v1"

const pathPrefix = "/go-service/v1"

const listeningPortKey="ListeningPort"
const noTLSKey="NO_TLS"
const serverCertKey="ServerCert"
const serverKeyKey="ServerKey"

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}


func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("received ping request")
}

func echo(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()
	//decoder := json.NewDecoder(body)
	log.Println(body)
	w.Header().Set("Content-Type", "application/json")
	//w.Write(body)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Println("request for invalid endpoint: ", r.URL.String())
}

func setUpRoutes() (*mux.Router, error) {
	router := mux.NewRouter().PathPrefix(pathPrefix).Subrouter()

	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/echo", echo).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)

	return router, nil
}

func runServer(config map[string]string) {
	listeningPort := config[listeningPortKey]
	servAddr := fmt.Sprintf(":%v", listeningPort)

	router, error := setUpRoutes()
	if(error != nil) {
		log.Println("error creating http router", error)
	}

	http.Handle("/", router)

	log.Println("listening on port ", servAddr)

	if(config[noTLSKey] == "true" && os.Getenv("STAGE") == "dev") {
		error = http.ListenAndServe(servAddr, nil)
	} else {
		serverCert := config[serverCertKey]
		serverKey := config[serverKeyKey]
		error = http.ListenAndServeTLS(servAddr, serverCert, serverKey, nil)
	}

	if(error != nil) {
		log.Println("error creating http server: ", error)
	}
}

func main() {
	startTime := time.Now()
	log.Println("Service starting at ", startTime)

	config, err := config.SetupConfig(configName, runtimeVersion)
	if(err != nil) {
		log.Println("error reading config")
		os.Exit(1)
	}
	log.Println(config)

	go runServer(config)
	
	for(true) {
		log.Println("Service has been running for ", time.Since(startTime));
		time.Sleep(kernelSleepInterval)
	}
}
