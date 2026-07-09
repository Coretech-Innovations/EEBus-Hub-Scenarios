package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Coretech-Innovations/EEBUS-Hub-Core/api"
)

// in this use story we are adding HEMS Device and bunch of uncontrollable loads.
// This user story demonstrates the HEMS's capability to receive instantenous power consumption and total energy consumption
// from a smart meter gateway implied in the grid connection point.
// This demo simulates an energy meter that acts as the GCP in which it sends to the HEMS the various parameters supported
// by the MGCP use case.

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println("go run ./examples/Api/MGCP/grid-point-monitoring <gcpRemoteSki>")
		return
	}
	var gcpSki string = os.Args[1]

	// start the simulation session
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 1,
	}
	sendRequest("POST", "/sim", simulationData)
	// The HEMS is a built-in singleton created on startup, so fetch the existing
	// one instead of adding a new device.
	resp := sendRequest("GET", "/hems", nil)
	hemsSKI := resp.(map[string]any)["ski"]

	fmt.Printf("Using the built-in HEMS, ski:%v\n", hemsSKI)

	// Pair with GCP

	resp = sendRequest("POST", "/hems/trust", struct {
		RemoteSki string `json:"remoteSki"`
	}{
		RemoteSki: gcpSki,
	})

	if resp.(map[string]any)["status"] != "OK" {
		fmt.Println("Error trusting", resp.(map[string]any)["err"])
		return
	} else {
		fmt.Println("Succesfully paired HEMS with GCP")
	}

	for i := 0; i < 5; i++ {
		resp = sendRequest("POST", "/uncontrollabledevice", api.UncontrollableDevice{
			DeviceName: "Device " + strconv.Itoa(i),
			Current: struct {
				PhaseA float32 "json:\"a,omitempty\""
				PhaseB float32 "json:\"b,omitempty\""
				PhaseC float32 "json:\"c,omitempty\""
			}{
				PhaseA: 2,
				PhaseB: 2,
				PhaseC: 2,
			},
			PowerFactor:  1,
			PowerStateOn: true,
		})
		if resp.(map[string]any)["status"] != "OK" {
			fmt.Println("Error adding device", i, " - ", resp.(map[string]any)["err"])
		} else {
			fmt.Println("Successfully added uncontrollable device ", resp.(map[string]any)["deviceId"])
			//uncontrollabledevices = append(uncontrollabledevices, resp.(map[string]any)["deviceId"].(string))
		}
	}

	go func() {
		for {
			resp = sendRequest("GET", "/hems/gridmeasurement", nil)
			resp, _ := resp.(map[string]any)
			fmt.Println("======Grid Measurement Data======")
			fmt.Println("Time:", time.Now())
			meas, _ := resp["measurement"].(map[string]any)
			fmt.Println("Voltage:", meas["Voltage"])
			fmt.Println("Current:", meas["Current"])
			fmt.Println("Power:", meas["Power"])
			fmt.Println("Total Energy Consumption:", resp["totalConsumptionKWh"], "kWh")
			time.Sleep(2 * time.Second)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("Quitting.")
	// ---------------------------------

}

func sendRequest(method string, url string, payload any) any {
	client := &http.Client{}
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, BaseIPAddress+url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result []map[string]any

	err = json.Unmarshal(body, &result)
	if err == nil {
		return result
	} else {
		var result2 map[string]any
		err = json.Unmarshal(body, &result2)
		if err != nil {
			log.Fatal("Can not convert to map")
			return nil
		}
		return result2
	}

}
