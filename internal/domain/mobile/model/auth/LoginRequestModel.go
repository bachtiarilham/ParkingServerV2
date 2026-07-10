package auth

type LoginRequestModel struct {
	Identity   string `json:"identity"`
	Password   string `json:"password"`
	DeviceId   string
	DeviceName string
	FcmToken   string
	IpAdrress  string
	UserAgent  string
}
