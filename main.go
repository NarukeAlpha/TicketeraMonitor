package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"TMv1/Monitor"
	"github.com/playwright-community/playwright-go"
)

var MonitorW struct {
	Webhook string `json:"webhook"`
}

var SetUp struct {
	Completed bool `json:"completed"`
}

var mw io.Writer

func init() {
	err := playwright.Install()
	// load json data file
	// load data into variables
	// check if data is valid
	// if data is not valid, start a setup process
	_, err = os.Stat("data.json")
	if os.IsNotExist(err) {
		_, err = os.Create("data.json")
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = os.Stat("ProxyList.csv")
	if os.IsNotExist(err) {
		log.Fatalln("PoxyList.csv not found, please provide a csv file named ProxyList in the same directory as the exe")

	}
	file, err := os.Open("data.json")
	if err != nil {
		log.Panicf("Error opening data.json: %v", err)

	}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&SetUp); err != nil {
		log.Panicf("Error decoding data.json: %v", err)
	}

	if SetUp.Completed {
		err = decoder.Decode(&MonitorW)
		if err != nil {
			log.Panicf("Error decoding data.json: %v", err)

		}
		Monitor.AssertErrorToNil("failed to close file: %v", file.Close())

	} else {
		Monitor.AssertErrorToNil("failed to close file: %v", file.Close())

		// Ask for Monitor configuration value
		fmt.Println("\nEnter Monitor configuration:")
		fmt.Print("Webhook: ")
		_, err := fmt.Scanln(&MonitorW.Webhook)
		if err != nil {
			log.Panicf("Error scanning input: %v", err)
		}

		SetUp.Completed = true
		file, err = os.OpenFile("data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Panicf("Error opening data.json: %v", err)
		}

		encoder := json.NewEncoder(file)
		if err = encoder.Encode(SetUp); err != nil {
			log.Panicf("Error encoding data.json: %v", err)
		}
		if err = encoder.Encode(MonitorW); err != nil {
			log.Panicf("Error encoding data.json: %v", err)
		}
		Monitor.AssertErrorToNil("Error closing data.json: %v", file.Close())

	}

	_, err = os.Stat("log.txt")
	if os.IsNotExist(err) {
		file, err = os.Create("log.txt")
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err = os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw = io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	log.Println("This is a log message")
	//this is a commit test

}

func main() {

	proxies := Monitor.ProxyLoad()
	log.Println("Proxies loaded")
	go Monitor.TaskInit(proxies, MonitorW.Webhook)
	log.Println("started monitor task")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

}
