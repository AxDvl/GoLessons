package auxilaries

import (
	"encoding/json"
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

func GetBodyAsJson(r io.ReadCloser, obj any) error {
	bodyText, err := GetStringFromBody(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(bodyText), obj)
	if err != nil {
		return err
	}
	return nil
}
