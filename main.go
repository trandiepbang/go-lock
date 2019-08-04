//shows how to watch for new devices and list them
package main

import (
	"flag"
	device "go-lock/utils"
	"time"

	"github.com/muka/go-bluetooth/api"
	log "github.com/sirupsen/logrus"
)

var configMachineID *string
var configMaxRetryConnectCount *int

const adapterID = "hci0"
const logLevel = log.DebugLevel

var deviceUtils *device.DeviceUtils

func main() {
	log.SetLevel(logLevel)
	deviceUtils = device.NewDeviceUtils()
	configMachineID = flag.String("device", "", "a device mac addres")
	configMaxRetryConnectCount = flag.Int("maxretry", 3, "numbers time try to reconnect device")

	flag.Parse()

	if *configMachineID == "" {
		panic("You need to pass into device mac address")
	}
	log.Info("Device info : ", *configMachineID)

	deviceUtils.SetMaxRetryCount(*configMaxRetryConnectCount)

	for true {
		time.Sleep(3 * time.Second)
		startService()
	}
}

func startService() {
	//Clear device cache ...
	api.ClearDevices()
	devices, err := api.GetDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		props, err := device.GetProperties()
		if err != nil {
			log.Errorf("%s: Failed to get properties: %s", err.Error())
			return
		}
		if props.Address == *configMachineID {
			currDevice, err := api.GetDeviceByAddress(props.Address)
			if err != nil {
				panic(err)
			}
			if !currDevice.IsConnected() {
				// Device not found
				log.Debug("Phone is disconnected trying to connect : ", props.Address)
				if !deviceUtils.GetLockState() && deviceUtils.GetRetryCount() > *configMaxRetryConnectCount {
					deviceUtils.SaveLockState(true)
					deviceUtils.Lock()
				}

				deviceUtils.IncreRetryCount()
				currDevice.Connect()
				return
			}

			log.Infof("name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
			if deviceUtils.GetLockState() {
				deviceUtils.SaveLockState(false)
				deviceUtils.Unlock()
			}

			deviceUtils.ResetRetryCount()
		}
	}
}
