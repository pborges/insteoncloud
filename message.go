package insteoncloud

import (
	"encoding/json"
	"time"
)

type eventMessage struct {
	HubInsteonID    string    `json:"hub_insteon_id"`
	DeviceInsteonID string    `json:"device_insteon_id"`
	DeviceGroup     int       `json:"device_group"`
	Status          string    `json:"status"`
	ReceivedAt      time.Time `json:"received_at"`
}

type responseStatus struct {
	Level int `json:"level"`
}
type response struct {
	DeviceList []Device `json:"DeviceList"`
	SceneList  []Scene  `json:"SceneList"`
	HouseList  []struct {
		HouseID   int    `json:"HouseID"`
		HouseName string `json:"HouseName"`
		IconID    int    `json:"IconID"`
	} `json:"HouseList"`
}

type commandDeviceRequest struct {
	Command  string `json:"command"`
	DeviceID int    `json:"device_id"`
	Level    int    `json:"level"`
}
type commandResponse struct {
	Status   string          `json:"status"`
	Link     string          `json:"link"`
	ID       int             `json:"id"`
	Response json.RawMessage `json:"response"`
}

type commandSceneRequest struct {
	Command  string `json:"command"`
	SceneID  int    `json:"scene_id"`
}