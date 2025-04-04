package epay

import "strings"

func (dt DeviceType) IsValid() bool {
	switch dt {
	case PC, MOBILE, QQ, WECHAT, ALIPAY:
		return true
	default:
		return false
	}
}

// 转换函数
func ParseDeviceType(device string) DeviceType {
	dt := DeviceType(strings.ToLower(device))
	if dt.IsValid() {
		return dt
	}
	return PC
}
