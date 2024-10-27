package SecurityUtil

// ClientSecurityKey 客户端加密秘钥
var ClientSecurityKey [256]uint8

// Mapping 解密数据
func Mapping(data []uint8, len int) {
	for i := 0; i < len; i++ {
		value := data[i]
		data[i] = ClientSecurityKey[value]
	}
}
