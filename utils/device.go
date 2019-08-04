package device

import (
	"log"
	"os/exec"
)

// DeviceUtils ...
type DeviceUtils struct {
	lockState            bool
	reconnectTryCount    int
	maxreconnectTryCount int
}

// NewDeviceUtils ...
func NewDeviceUtils() *DeviceUtils {
	return &DeviceUtils{}
}

// Lock ...
func (c *DeviceUtils) Lock() {
	c.lockUnlock("lock")
}

// Unlock ...
func (c *DeviceUtils) Unlock() {
	c.lockUnlock("unlock")
}

func (c *DeviceUtils) lockUnlock(typelock string) {
	var err error
	if typelock == "lock" {
		err = exec.Command("loginctl", "lock-session").Run()
	} else {
		err = exec.Command("loginctl", "unlock-session").Run()
	}

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

// SaveLockState ...
func (c *DeviceUtils) SaveLockState(state bool) {
	c.lockState = state
}

// GetLockState ...
func (c *DeviceUtils) GetLockState() bool {
	return c.lockState
}

// IncreRetryCount ...
func (c *DeviceUtils) IncreRetryCount() {
	if c.reconnectTryCount >= c.maxreconnectTryCount {
		return
	}
	c.reconnectTryCount++
}

// ResetRetryCount ...
func (c *DeviceUtils) ResetRetryCount() {
	c.reconnectTryCount = 0
}

// GetRetryCount ...
func (c *DeviceUtils) GetRetryCount() int {
	return c.reconnectTryCount
}

// SetMaxRetryCount ...
func (c *DeviceUtils) SetMaxRetryCount(max int) {
	c.maxreconnectTryCount = max
}
