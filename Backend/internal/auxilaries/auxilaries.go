package auxilaries

import (
	"io"
)

func GetStringFromBody(r io.ReadCloser) (string, error) {
	buf := make([]byte, 100)
	s := ""
	for {
		n, err := r.Read(buf)
		s += string(buf[:n])

		if err == io.EOF {
			break
		}

		if err != nil {
			return s, err
		}
	}
	return s, nil

}
