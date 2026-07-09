package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// in this use story we are adding HEMS Device, Heatpump and connecting them with each other
// After the connection is done, we create an OHPCF announcement with power limits and timing constraints:
// - power(goodApproximation): 15000 W
// - powerMaximum: 16000 W
// - activeDurationMinimum: 5 seconds
// - pauseDurationMinimum: 5 seconds
// - startTime: 5 seconds

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {

	fmt.Println("Starting the simulation session example for OHPCF use case")
	// reset simulation session by deleting all components
	simulationReset()

	time.Sleep(5 * time.Second)

	// Get HEMS Ski
	resp := sendRequest("GET", "/hems", nil)
	hemsSKI := resp.(map[string]any)["ski"]

	// Add EVSE
	HeatPumpInfo := map[string]any{
		"deviceName":  "Coretech Heat Pump",
		"deviceCode":  "a7171fc2",
		"deviceModel": "HeatPump-1",
		"brandName":   "Coretech",
		"vendor": map[string]any{
			"name": "Coretech",
			"code": "CT",
		},
		"softwareRev": "0",
		"hardwareRev": "0",
		"manufacturer": map[string]any{
			"label":       "Coretech Innovations",
			"description": "",
		},
		"serialNumber":        "2a7ae968",
		"approveWriteLimit":   true,
		"failsafeValue":       5000,
		"failsafeDuration":    2,
		"failSafeDurationMax": 24,
		"nominalPower": map[string]any{
			"min": 0,
			"max": 23000,
		},
		"nominalCurrent": map[string]any{
			"min": 0,
			"max": 32,
		},
		"manufacturerDescription": "Heat Pump",
	}

	// sending request to add Heat Pump
	resp = sendRequest("POST", "/heatpump/add", HeatPumpInfo)
	heatPumpID := resp.(map[string]any)["id"]
	fmt.Printf("A new Heat Pump Device is added with ID %d\n", int(heatPumpID.(float64)))

	// start the simulation session
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 1,
	}
	sendRequest("POST", "/sim", simulationData)

	endPoint := fmt.Sprintf("/heatpump/%d/trust", int(heatPumpID.(float64)))
	// Running the Heat Pump to connect with the HEMS
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": hemsSKI,
	})
	heatPumpSKI := resp.(map[string]any)["ski"]

	// trusting the Heat Pump from the HEMS Side
	sendRequest("POST", "/hems/trust", map[string]any{
		"remoteSKI": heatPumpSKI,
	})
	fmt.Println("Ski = ", heatPumpSKI, " is trusted by the HEMS")

	var announcmentInfo map[string]any = map[string]any{
		"isPausable":               true,
		"isStoppable":              true,
		"power":                    15000,
		"powerMax":                 16000,
		"powerType":                "both",
		"activeDurationMinSeconds": 5,
		"pauseDurationMinSeconds":  5,
	}

	time.Sleep(5 * time.Second)

	// send request to create an OHPCF announcement
	endPoint = fmt.Sprintf("/heatpump/%d/announceOptional", int(heatPumpID.(float64)))
	resp = sendRequest("POST", endPoint, announcmentInfo)
	fmt.Println("Announcement is created with following parameters :")
	for key, value := range announcmentInfo {
		fmt.Printf("%s = %v\n", key, value)
	}

	resp = sendRequest("GET", "/heatpump/list", nil)
	state := firstHeatPumpState(resp)
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state)

	// send request to make HEMS send a start time to the Heat Pump
	resp = sendRequest("GET", "/heatpump/list", nil)
	deviceAddress := firstHeatPumpAddress(resp)

	time.Sleep(3 * time.Second)

	startime := 5
	endPoint = fmt.Sprintf("/hems/SendOhpcfStartTime/%s", deviceAddress)
	resp = sendRequest("POST", endPoint, map[string]any{
		"startTime": float64(startime),
	})

	fmt.Println("HEMS sent start time with value 5 seconds to the Heat Pump")
	fmt.Println("Heatpump OHPCF state is : Scheduled")

	for i := 0; i < startime; i++ {
		resp = sendRequest("GET", "/heatpump/list", nil)
		state := firstHeatPumpState(resp)
		fmt.Printf("Heat Pump OHPCF current state: %s, %d seconds passed\n", state, i+1)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)

	resp = sendRequest("GET", "/heatpump/list", nil)
	state = firstHeatPumpState(resp)
	fmt.Printf("Start time is reached, Heat Pump OHPCF current state: %s\n", state)

	// send request to pause the Heat Pump OHPCF session
	fmt.Println("HEMS can not pause the Heat Pump OHPCF session because minimum active duration is not reached")
	for i := 0; i < 5; i++ {
		fmt.Println("Waiting minimum active duration to be reached...", i+1, "seconds passed")
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)

	endPoint = fmt.Sprintf("/hems/ChangeOhpcfState/%s", deviceAddress)
	resp = sendRequest("PATCH", endPoint, map[string]any{
		"state": "paused",
	})
	fmt.Println("HEMS sent pause command to the Heat Pump")

	resp = sendRequest("GET", "/heatpump/list", nil)
	state = firstHeatPumpState(resp)
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state)

	// send request to resume the Heat Pump OHPCF session
	fmt.Println("HEMS can not resume the Heat Pump OHPCF session because minimum pause duration is not reached")
	for i := 0; i < 5; i++ {
		fmt.Println("Waiting minimum pause duration to be reached...", i+1, "seconds passed")
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)

	endPoint = fmt.Sprintf("/hems/ChangeOhpcfState/%s", deviceAddress)
	resp = sendRequest("PATCH", endPoint, map[string]any{
		"state": "running",
	})

	fmt.Println("HEMS sent resume command to the Heat Pump")

	resp = sendRequest("GET", "/heatpump/list", nil)
	state = firstHeatPumpState(resp)
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state)

	// send request to abort the Heat Pump OHPCF session
	fmt.Println("HEMS can not abort the Heat Pump OHPCF session because minimum active duration is not reached")
	for i := 0; i < 5; i++ {
		fmt.Println("Waiting minimum active duration to be reached...", i+1, "seconds passed")
		time.Sleep(1 * time.Second)
	}
	time.Sleep(2 * time.Second)

	endPoint = fmt.Sprintf("/hems/ChangeOhpcfState/%s", deviceAddress)
	resp = sendRequest("PATCH", endPoint, map[string]any{
		"state": "invalid",
	})

	fmt.Println("HEMS sent abort command to the Heat Pump")
	resp = sendRequest("GET", "/heatpump/list", nil)
	state = firstHeatPumpState(resp)
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state)

	// wait for the heat pump to clear its power sequence and transition to inactive state
	fmt.Println("Waiting the Heat Pump to clear its power sequence and transition to inactive state...")

	for {
		resp = sendRequest("GET", "/heatpump/list", nil)
		state = firstHeatPumpState(resp)
		if state == "inactive" {
			break
		}
	}
	fmt.Printf("Heat Pump OHPCF current state: %s\n", state)
}

func simulationReset() {
	// clear simulation
	simulationData := map[string]any{
		"action":      "reset",
		"speedFactor": 0,
	}
	sendRequest("POST", "/sim", simulationData)
}

// firstHeatPumpState returns the OHPCF state of the first heat pump in a
// /heatpump/list response, or "" if the list is momentarily empty (which can
// happen while the heat pump is transitioning between states).
func firstHeatPumpState(resp any) string {
	list, ok := resp.([]map[string]any)
	if !ok || len(list) == 0 {
		return ""
	}
	info, ok := list[0]["DeviceInfo"].(map[string]any)
	if !ok {
		return ""
	}
	state, _ := info["state"].(string)
	return state
}

// firstHeatPumpAddress returns the device address of the first heat pump in a
// /heatpump/list response, or "" if the list is momentarily empty.
func firstHeatPumpAddress(resp any) string {
	list, ok := resp.([]map[string]any)
	if !ok || len(list) == 0 {
		return ""
	}
	address, _ := list[0]["deviceAddress"].(string)
	return address
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
	if len(body) == 0 {
		fmt.Printf("Empty response body (status=%d) for %s %s\n", resp.StatusCode, method, url)
		return nil
	}
	var result []map[string]any

	err = json.Unmarshal(body, &result)
	if err == nil {
		return result
	} else {
		var result2 map[string]any
		err = json.Unmarshal(body, &result2)
		if err != nil {
			fmt.Printf("Non-JSON response (status=%d) for %s %s: %s\n", resp.StatusCode, method, url, string(body))
			return nil
		}
		return result2
	}

}
