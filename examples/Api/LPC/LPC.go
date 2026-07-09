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

// This example connects a HEMS with an EVSE. Once the connection is established it sends
// an active power consumption limit and monitors the EVSE's limit state until it expires.

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {

	// reset simulation session by deleting all components (entities can only be
	// added or removed while the simulation is stopped)
	simulationReset()

	// Get HEMS Ski (the HEMS is a built-in singleton created on startup)
	resp := sendRequest("GET", "/hems", nil)
	hemsSKI := resp.(map[string]any)["ski"]

	// Add EVSE
	evseInfo := map[string]any{"deviceName": "Coretech EVSE WLBX", "deviceCode": "037d42e1", "deviceModel": "Wallbox", "brandName": "Coretech", "vendor": map[string]any{"name": "Coretech", "code": "60745"}, "softwareRev": "0", "hardwareRev": "0", "Manufacturer": map[string]any{"label": "Coretech", "description": "Charging Station"}, "serialNumber": "de07c278", "failsafeValue": 4320, "failsafeDuration": 2, "failSafeDurationMax": 24, "nominalPower": map[string]any{"min": 4320, "max": 23000}, "contractualPowerMax": 23000, "nominalCurrent": map[string]any{"min": 6, "max": 32}, "approveWriteLimit": true}

	// sending request to add EVSE
	resp = sendRequest("POST", "/evse/add", evseInfo)
	evseID := resp.(map[string]any)["id"]

	fmt.Printf("A new EVSE Device is added with ID %d\n", int(evseID.(float64)))

	// start the simulation session (pairing is only allowed while it is running)
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 100,
	}
	sendRequest("POST", "/sim", simulationData)

	// pair the EVSE with the HEMS
	endPoint := fmt.Sprintf("/evse/%d/trust", int(evseID.(float64)))
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": hemsSKI,
	})
	evseSKI := resp.(map[string]any)["ski"]
	fmt.Println(evseSKI)
	// trusting the EVSE from the HEMS side
	sendRequest("POST", "/hems/trust", map[string]any{
		"remoteSKI": evseSKI,
	})
	time.Sleep(3 * time.Second)
	for {
		endPoint = fmt.Sprintf("/evse/%d/lpcCurrentState", int(evseID.(float64)))
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.(map[string]any))
		endPoint = fmt.Sprintf("/evse/%d/lpcActivePowerLimit", int(evseID.(float64)))
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.(map[string]any))
		endPoint = fmt.Sprintf("/evse/%d/lpcFailsafeValues", int(evseID.(float64)))
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.(map[string]any))
		endPoint = fmt.Sprintf("/evse/%d/lpcPowerConstraints", int(evseID.(float64)))
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.(map[string]any))
		time.Sleep(2 * time.Second)
		endPoint = "/hems/ActivePowerConsumptionLimit"
		resp = sendRequest("POST", endPoint, map[string]any{
			"active":          true,
			"value":           10000,
			"durationSeconds": 30,
		})
		time.Sleep(2 * time.Second)
		endPoint = fmt.Sprintf("/evse/%d/lpcCurrentState", int(evseID.(float64)))
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.(map[string]any))
		endPoint = fmt.Sprintf("/evse/%d/lpcActivePowerLimit", int(evseID.(float64)))
		resp = sendRequest("GET", endPoint, nil)
		fmt.Println(resp.(map[string]any))
		time.Sleep(30 * time.Second)
	}
}

func simulationReset() {
	// stop the simulation first — entities can only be removed while it is stopped
	sendRequest("POST", "/sim", map[string]any{"action": "stop"})
	// delete all the EVSEs
	resp := sendRequest("GET", "/evse/list", nil)
	evses := resp.([]map[string]any)
	for _, evse := range evses {
		endpoint := fmt.Sprintf("/evse/%d", int(evse["evseId"].(float64)))
		sendRequest("DELETE", endpoint, nil)
	}
	// delete all the EVs
	resp = sendRequest("GET", "/ev/list", nil)
	evs := resp.([]map[string]any)
	for _, ev := range evs {
		endpoint := fmt.Sprintf("/ev/%d", int(ev["evId"].(float64)))
		sendRequest("DELETE", endpoint, nil)
	}
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
