package models

// 登录状态
type LoginStatus struct {
	AccessIP      string `json:"access_ip"`
	AccessPort    string `json:"access_port"`
	AppId         uint32 `json:"appId"`
	UserId        string `json:"userId"`
	ClientIP      string `json:"client_ip"`
	ClientPort    string `json:"client_port"`
	LoginTime     uint64 `json:"login_time"`      // 用户上次登录时间
	HeartBeatTime uint64 `json:"heart_beat_time"` // 用户上次心跳时间
	LogoutTime    uint64 `json:"logout_time"`     // 用户上次退出登录的时间
	DeviceInfo    string `json:"device_info"`     // 设备信息
	Logoff        bool   `json:"logoff"`          // 是否已离线
}
