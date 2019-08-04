//shows how to watch for new devices and list them
package main

import (
	"time"

	"github.com/muka/go-bluetooth/api"
	log "github.com/sirupsen/logrus"
)

var devicesMap map[string]bool

const logLevel = log.DebugLevel
const machineID = "30:95:E3:C8:33:49"
const adapterID = "hci0"

var isLocked bool

func main() {

	log.SetLevel(logLevel)

	//clean up connection on exit
	devicesMap = map[string]bool{
		"30:95:E3:C8:33:49": true,
		"78:00:9E:88:B6:13": false,
	}

	for true {
		time.Sleep(5 * time.Second)
		startService()
	}
}

func startService() {
	api.ClearDevices()
	devices, err := api.GetDevices()
	// log.Debug(devices)
	// log.Debug("Error ", err)
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		props, err := device.GetProperties()
		if err != nil {
			log.Errorf("%s: Failed to get properties: %s", err.Error())
			return
		}
		// log.Debug("Found RSSI ", props.RSSI)
		if devicesMap[props.Address] {
			isDeviceAvailable, err := api.GetDeviceByAddress(props.Address)
			if err != nil {
				panic(err)
			}
			// log.Debug("Connected : ", isDeviceAvailable.IsConnected())
			if !isDeviceAvailable.IsConnected() {
				// Device not found
				if !isLocked {
					isLocked = true
				}
				isDeviceAvailable.Connect()
				log.Debug("Phone is disconnected -- trying to connect ", props.Address)
				return
			}
			log.Infof("name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
			if isLocked {
				//
				isLocked = false

			}
		}

	}
}
