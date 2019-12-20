package insteoncloud

type Scene struct {
	SceneID int `json:"SceneID"`
	//AutoStatus      bool   `json:"AutoStatus"`
	//CustomOff       string `json:"CustomOff"`
	//CustomOn        string `json:"CustomOn"`
	//DayMask         int64  `json:"DayMask"`
	//EnableCustomOff bool   `json:"EnableCustomOff"`
	//EnableCustomOn  bool   `json:"EnableCustomOn"`
	//Favorite        bool   `json:"Favorite"`
	Group int64 `json:"Group"`
	//IconID          int    `json:"IconID"`
	HouseID int64 `json:"HouseID"`
	//OffTime         string `json:"OffTime"`
	//OnTime          string `json:"OnTime"`
	SceneName    string `json:"SceneName"`
	StatusDevice string `json:"StatusDevice"`
	//TimerEnabled    bool   `json:"TimerEnabled"`
	//Visible         bool   `json:"Visible"`
	DeviceList []struct {
		OnLevel             int `json:"OnLevel"`
		DeviceRoleMask      int `json:"DeviceRoleMask"`
		DeviceGroupDetailID int `json:"DeviceGroupDetailID"`
		RampRate            int `json:"RampRate"`
		DeviceID            int `json:"DeviceID"`
	} `json:"DeviceList"`
}
