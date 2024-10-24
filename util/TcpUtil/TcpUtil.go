package TcpUtil

import "net"

// WriteAll net.Conn写入所有数据
func WriteAll(tcp net.Conn, data []uint8) error {

	//要发送的数据总长度
	total := len(data)

	//已经大宋的数据长度
	sendedTotal := 0

	for {
		//n:当前发送的数据长度
		n, err := tcp.Write(data[sendedTotal:])
		if err != nil {
			return err
		}
		sendedTotal += n
		if sendedTotal == total {
			break
		}
	}
	return nil
}

// ReadNByte 读取指定长度数据
func ReadNByte(tcp net.Conn, n int) ([]uint8, error) {

	//记录已经读取到的数据大小
	var readLen = 0
	data := make([]uint8, n)
	for {
		buffer := make([]uint8, n-readLen)
		le, err := tcp.Read(buffer)
		if err != nil {
			return nil, err
		}
		copy(data[readLen:readLen+le], buffer[:le])
		readLen += le
		if readLen == n {
			break
		}
	}
	return data, nil
}
