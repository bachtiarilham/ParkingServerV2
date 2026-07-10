package auth

type LoginRequestDto struct {
	Identity   string `json:"identity"`
	Password   string `json:"password"`
	DeviceId   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	FcmToken   string `json:"fcm_token"`
}
