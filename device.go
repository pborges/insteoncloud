package insteoncloud

type Device struct {
	HouseID          int    `json:"HouseID"`
	DeviceID         int    `json:"DeviceID"`
	DeviceName       string `json:"DeviceName"`
	//IconID           int    `json:"IconID"`
	//AlertOff         int    `json:"AlertOff"`
	//AlertOn          int    `json:"AlertOn"`
	//AlertsEnabled    bool   `json:"AlertsEnabled"`
	//AutoStatus       bool   `json:"AutoStatus"`
	//BeepOnPress      bool   `json:"BeepOnPress"`
	//BlinkOnTraffic   bool   `json:"BlinkOnTraffic"`
	//ConfiguredGroups int    `json:"ConfiguredGroups"`
	//CustomOff        string `json:"CustomOff"`
	//CustomOn         string `json:"CustomOn"`
	//DayMask          int    `json:"DayMask"`
	DevCat           int    `json:"DevCat"`
	DeviceType       int    `json:"DeviceType"`
	//DimLevel         int    `json:"DimLevel"`
	//EnableCustomOff  bool   `json:"EnableCustomOff"`
	//EnableCustomOn   bool   `json:"EnableCustomOn"`
	//Favorite         bool   `json:"Favorite"`
	FirmwareVersion  int    `json:"FirmwareVersion"`
	//Group            int    `json:"Group"`
	//Humidity         bool   `json:"Humidity"`
	//InsteonEngine    int    `json:"InsteonEngine"`
	InsteonID        string `json:"InsteonID"`
	//LEDLevel         int    `json:"LEDLevel"`
	//LinkWithHub      int    `json:"LinkWithHub"`
	//LocalProgramLock bool   `json:"LocalProgramLock"`
	//OffTime          string `json:"OffTime"`
	//OnTime           string `json:"OnTime"`
	//OperationFlags   int    `json:"OperationFlags"`
	//RampRate         int    `json:"RampRate"`
	//SerialNumber     string `json:"SerialNumber"`
	SubCat           int    `json:"SubCat"`
	//TimerEnabled     bool   `json:"TimerEnabled"`
	//DeviceTypeTraits struct {
	//	SecurityDevice  bool   `json:"SecurityDevice"`
	//	TypeDescription string `json:"TypeDescription"`
	//} `json:"DeviceTypeTraits"`
	//GroupList []interface{} `json:"GroupList"`
}
