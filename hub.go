package insteoncloud

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const BaseUrl = "https://connect.insteon.com/api/v2"

type Hub struct {
	Username string
	Password string
	ClientId string
	Log      *log.Logger

	connected bool
	houseId   int
	token     Token
	devices   map[string]Device
	scenes    map[int]Scene
	lock      sync.Mutex
}

func (h *Hub) getJson(url string, res interface{}) error {
	return h.reqForJson(http.MethodGet, url, nil, res)
}

func (h *Hub) postJson(url string, req interface{}, res interface{}) error {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(req); err == nil {
		return h.reqForJson(http.MethodPost, url, body, res)
	} else {
		return err
	}
}
func (h *Hub) reqForJson(method string, url string, body io.Reader, res interface{}) error {
	response, err := h.req(method, url, body)
	if err != nil {
		return err
	}
	return json.NewDecoder(response.Body).Decode(res)
}
func (h *Hub) req(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authentication", "APIKey "+h.ClientId)
	req.Header.Add("Authorization", "Bearer "+h.token.AccessToken)
	res, err := http.DefaultClient.Do(req)

	if h.Log != nil {
		h.Log.Println(method, url, "error:", err)
	}

	if err != nil {
		// attempt a refresh
		if err := h.refresh(); err != nil {
			return nil, err
		}
		// try the request again
		return h.req(method, url, body)
	}
	return res, nil
}

func (h *Hub) refresh() error {
	if h.Log == nil {
		h.Log.Println("refresh token")
	}
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("client_id", h.ClientId)
	form.Add("refresh_token", h.token.RefreshToken)

	res, err := http.Post(
		BaseUrl+"/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(&h.token)
}

func (h *Hub) login() error {
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("client_id", h.ClientId)
	form.Add("username", h.Username)
	form.Add("password", h.Password)

	res, err := http.Post(
		BaseUrl+"/oauth2/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(&h.token)
}

func (h *Hub) findHouseId() error {
	var res response
	if err := h.getJson(BaseUrl+"/houses", &res); err != nil {
		return err
	}

	if len(res.HouseList) != 1 {
		return errors.New("idk what to do with this many houses")
	}

	h.houseId = res.HouseList[0].HouseID
	return nil
}

func (h *Hub) queryDevices() error {
	var res response
	if err := h.getJson(BaseUrl+"/devices?properties=all", &res); err != nil {
		return err
	}

	h.lock.Lock()
	for _, dev := range res.DeviceList {
		h.devices[dev.InsteonID] = dev
	}
	h.lock.Unlock()
	return nil
}

func (h *Hub) queryScenes() error {
	var res response
	if err := h.getJson(BaseUrl+"/scenes?properties=all", &res); err != nil {
		return err
	}

	h.lock.Lock()
	for _, dev := range res.SceneList {
		h.scenes[dev.SceneID] = dev
	}
	h.lock.Unlock()
	return nil
}

func (h *Hub) command(cmd interface{}, res interface{}) error {
	if h.Log != nil {
		h.Log.Printf("command: %+v\n", cmd)
	}
	var response commandResponse
	if err := h.postJson(BaseUrl+"/commands", cmd, &response); err != nil {
		if h.Log != nil {
			h.Log.Printf("command: %+v error: %s\n", cmd, err)
		}
		return err
	}

	if h.Log != nil {
		h.Log.Printf("command: %+v res: %+v\n", cmd, response)
	}

	for response.Status == "pending" {
		time.Sleep(1 * time.Second)
		statusUrl := fmt.Sprintf("%s/commands/%d", BaseUrl, response.ID)
		if h.Log != nil {
			h.Log.Printf("command retry: %+v url: %s\n", cmd, statusUrl)
		}
		if err := h.getJson(statusUrl, &response); err != nil {
			return err
		}
	}

	if response.Status == "failed" {
		return errors.New("failed to execute command")
	}
	if response.Response != nil {
		return json.Unmarshal(response.Response, res)
	}
	return nil
}

func (h *Hub) RefreshDevicesAndScenes() error {
	if !h.connected {
		return errors.New("not connected")
	}

	if err := h.queryDevices(); err != nil {
		return err
	}
	if err := h.queryScenes(); err != nil {
		return err
	}
	if h.Log != nil {
		h.Log.Println("query devices, found", len(h.devices), "devices and", len(h.scenes), "scenes")
	}
	return nil
}

func (h *Hub) Subscribe(fn func(dev Device, state string)) error {
	if !h.connected {
		return errors.New("not connected")
	}

	streamUrl := fmt.Sprintf("%s/houses/%d/stream", BaseUrl, h.houseId)
	if h.Log != nil {
		h.Log.Println("subscribe", streamUrl)
	}
	res, err := h.req(http.MethodGet, streamUrl, nil)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("invalid response for listenForEvents")
	}

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if h.Log != nil {
			h.Log.Println(line)
		}
		if strings.HasPrefix(line, "data: ") {
			if h.Log != nil {
				h.Log.Println(line)
			}
			line = strings.TrimPrefix(line, "data: ")
			var data eventMessage
			if err := json.Unmarshal([]byte(line), &data); err != nil {
				return err
			}

			if dev, err := h.Device(data.DeviceInsteonID); err == nil {
				fn(dev, data.Status)
			} else {
				return err
			}
		}
	}
	return scanner.Err()
}

func (h *Hub) Device(insteonId string) (Device, error) {
	if !h.connected {
		return Device{}, errors.New("not connected")
	}

	h.lock.Lock()
	defer h.lock.Unlock()
	for _, dev := range h.devices {
		if dev.InsteonID == insteonId {
			return dev, nil
		}
	}
	return Device{}, errors.New("unknown device")
}

func (h *Hub) Devices() ([]Device, error) {
	if !h.connected {
		return nil, errors.New("not connected")
	}

	h.lock.Lock()
	defer h.lock.Unlock()
	devs := make([]Device, 0, len(h.devices))
	for _, dev := range h.devices {
		devs = append(devs, dev)
	}
	return devs, nil
}

func (h *Hub) Scene(sceneId int) (Scene, error) {
	if !h.connected {
		return Scene{}, errors.New("not connected")
	}

	h.lock.Lock()
	defer h.lock.Unlock()
	for _, scene := range h.scenes {
		if scene.SceneID == sceneId {
			return scene, nil
		}
	}
	return Scene{}, errors.New("unknown scene")
}

func (h *Hub) Scenes() ([]Scene, error) {
	if !h.connected {
		return nil, errors.New("not connected")
	}

	h.lock.Lock()
	defer h.lock.Unlock()
	scenes := make([]Scene, 0, len(h.scenes))
	for _, scene := range h.scenes {
		scenes = append(scenes, scene)
	}
	return scenes, nil
}

func (h *Hub) Connect() error {
	h.lock.Lock()
	if h.connected {
		return errors.New("already connected")
	}
	h.devices = make(map[string]Device)
	h.scenes = make(map[int]Scene)
	h.connected = true
	h.lock.Unlock()

	if err := h.login(); err != nil {
		return err
	}

	if h.Log != nil {
		h.Log.Println("login success")
	}

	if err := h.findHouseId(); err != nil {
		return err
	}

	if h.Log != nil {
		h.Log.Println("found house id", h.houseId)
	}

	return h.RefreshDevicesAndScenes()
}

func (h *Hub) SetDeviceLevel(insteonId string, level int) error {
	if !h.connected {
		return errors.New("not connected")
	}

	if level < 0 || level > 100 {
		return errors.New("level must be between 0 and 100")
	}

	dev, err := h.Device(insteonId)
	if err != nil {
		return err
	}

	var res responseStatus
	req := commandDeviceRequest{
		Command:  "on",
		DeviceID: dev.DeviceID,
		Level:    level,
	}
	if level == 0 {
		req.Command = "off"
	}

	if err := h.command(req, &res); err != nil {
		return err

	}

	return nil
}

func (h *Hub) SetSceneState(sceneId int, state bool) error {
	if !h.connected {
		return errors.New("not connected")
	}

	dev, err := h.Scene(sceneId)
	if err != nil {
		return err
	}

	req := commandSceneRequest{
		Command: "off",
		SceneID: dev.SceneID,
	}

	if state {
		req.Command = "on"
	}

	var res responseStatus
	return h.command(req, &res)
}

func (h *Hub) GetStatus(insteonId string) (int, error) {
	if !h.connected {
		return 0, errors.New("not connected")
	}

	dev, err := h.Device(insteonId)
	if err != nil {
		return 0, err
	}
	var res responseStatus
	if err := h.command(commandDeviceRequest{
		Command:  "get_status",
		DeviceID: dev.DeviceID,
	}, &res); err != nil {
		return 0, err
	}
	return res.Level, nil
}
