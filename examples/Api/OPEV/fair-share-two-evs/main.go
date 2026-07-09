package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/enbility/spine-go/model"
)

// in this use story we will add 2 EVs to the system, both max current under the available limits

var BaseIPAddress string = "http://localhost:8080/api/v1"

func main() {

	// reset simulation session by deleting all components
	simulationReset()

	// Get HEMS SKI
	resp := sendRequest("GET", "/hems", nil)
	hemsSKI := resp.(map[string]any)["ski"]

	// adding EVSE (entities can only be added while the simulation is stopped)
	evseInfo := map[string]any{
		"deviceName": "eSystems EVSE WallBox",
		"deviceCode": "0001",
		"vendor": map[string]any{
			"name": "eSystems",
			"code": "60745",
		},
		"softwareRev": "0",
		"hardwareRev": "0",
		"brandName":   "eSystems",
		"Manufacturer": map[string]any{
			"label":       "eSystems",
			"description": "Charging Station",
		},
		"deviceModel":         "EVSE",
		"serialNumber":        "00000002",
		"port":                4712,
		"approveWriteLimit":   true,
		"failsafeValue":       5000,
		"failsafeDuration":    2,
		"failSafeDurationMax": 24,
		"nominalPower":        map[string]any{"min": 0, "max": 23000},
		"nominalCurrent":      map[string]any{"min": 0, "max": 32},
	}
	// sending request to add EVSE
	resp = sendRequest("POST", "/evse/add", evseInfo)
	evseID := int(resp.(map[string]any)["id"].(float64))

	fmt.Printf("A new EVSE Device is added with ID %d\n", evseID)

	// add EV
	var EV1 map[string]any = map[string]any{
		"asymmetricCharging": false,
		"currentLimits": map[string]int{
			"min": 6,
			"max": 18,
		},
		"dischargingEnable": false,
		"chargingEnable":    true,
		"chargingCapacity":  80,
		"charged":           20,
		"batteryHealth":     100,
	}
	resp = sendRequest("POST", "/ev/add", EV1)
	ev1ID := int(resp.(map[string]any)["id"].(float64))
	fmt.Printf("A new EV Device is added with ID %d\n", ev1ID)

	// creating second EV
	var EV2 map[string]any = map[string]any{
		"asymmetricCharging": false,
		"currentLimits": map[string]int{
			"min": 10,
			"max": 20,
		},
		"dischargingEnable": false,
		"chargingEnable":    true,
		"chargingCapacity":  80,
		"charged":           20,
		"batteryHealth":     100,
	}
	resp = sendRequest("POST", "/ev/add", EV2)
	ev2ID := int(resp.(map[string]any)["id"].(float64))
	fmt.Printf("A new EV Device is added with ID %d\n", ev2ID)

	// start the simulation session (pairing is only allowed while it is running)
	simulationData := map[string]any{
		"action":      "start",
		"speedFactor": 100,
	}
	sendRequest("POST", "/sim", simulationData)

	endPoint := fmt.Sprintf("/evse/%d/trust", evseID)
	// Running the EVSE to connect with the HEMS
	resp = sendRequest("POST", endPoint, map[string]any{
		"remoteSKI": hemsSKI,
	})
	evseSKI := resp.(map[string]any)["ski"]
	// trusting the EVSE from the HEMS Side
	sendRequest("POST", "/hems/trust", map[string]any{
		"remoteSKI": evseSKI,
	})

	// connect EV1 to EVSE
	endPoint = fmt.Sprintf("/ev/%d/evse/%d", ev1ID, evseID)
	sendRequest("POST", endPoint, nil)

	// connect EV2 to EVSE
	endPoint = fmt.Sprintf("/ev/%d/evse/%d", ev2ID, evseID)
	sendRequest("POST", endPoint, nil)
	time.Sleep(5 * time.Second)
	for {
		endPoint = fmt.Sprintf("/ev/%d/LoadControlLimit", ev1ID)
		resp = sendRequest("GET", endPoint, nil)
		{
			limits, _ := json.Marshal(resp)
			var limitsCasted []model.LoadControlLimitDataType
			json.Unmarshal(limits, &limitsCasted)
			if len(limitsCasted) == 0 {
				fmt.Printf("per-phase limit for EV%d :  not available yet\n", ev1ID)
			} else {
				val := limitsCasted[0].Value.GetValue()
				fmt.Println(fmt.Sprintf("per-phase limit for EV%d : ", ev1ID), val)
			}
		}
		endPoint = fmt.Sprintf("/ev/%d/LoadControlLimit", ev2ID)
		resp = sendRequest("GET", endPoint, nil)
		{
			limits, _ := json.Marshal(resp)
			var limitsCasted []model.LoadControlLimitDataType
			json.Unmarshal(limits, &limitsCasted)
			if len(limitsCasted) == 0 {
				fmt.Printf("per-phase limit for EV%d :  not available yet\n", ev2ID)
			} else {
				val := limitsCasted[0].Value.GetValue()
				fmt.Println(fmt.Sprintf("per-phase limit for EV%d : ", ev2ID), val)
			}
		}
		time.Sleep(2 * time.Second)
	}

}

func simulationReset() {
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
