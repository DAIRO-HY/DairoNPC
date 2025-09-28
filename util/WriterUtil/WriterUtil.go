package WriterUtil

import (
	"io"
)

// WriteFull 尝试将 buf 全部写入 w，直到写完或出错
func WriteFull(w io.Writer, buf []byte) error {
	total := 0
	for total < len(buf) {
		n, err := w.Write(buf[total:])
		if err != nil {
			return err
		}
		total += n
	}
	return nil
}
