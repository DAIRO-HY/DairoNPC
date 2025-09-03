package TcpUtil

import "net"

// WriteAll net.Conn写入所有数据
func WriteAll(tcp net.Conn, data []uint8) error {

	//要发送的数据总长度
	total := len(data)

	//已经发送的数据长度
	sendTotal := 0

	for {
		//n:当前发送的数据长度
		n, err := tcp.Write(data[sendTotal:])
		if err != nil {
			return err
		}
		sendTotal += n
		if sendTotal == total {
			break
		}
	}
	return nil
}

// ReadNByte 读取指定长度数据
func ReadNByte(tcp net.Conn, len int) ([]uint8, error) {

	//记录已经读取到的数据大小
	var readLen = 0
	data := make([]uint8, len)
	for {
		buffer := make([]uint8, len-readLen)
		n, readErr := tcp.Read(buffer)
		if n > 0 {
			copy(data[readLen:readLen+n], buffer[:n])
			readLen += n
			if readLen == len {
				break
			}
		}
		if readErr != nil {
			return nil, readErr
		}
	}
	return data, nil
}
