package read_env

import (
	"bytes"
	"os"
)

func ReadEnv(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lFile := len(file)
	if lFile > 0 {
		if file[lFile-1] != '\n' {
			file = append(file, '\n')
		}
	} else {
		return nil
	}
	buf := bytes.Buffer{}
	field, value := "", ""
	for i := 0; i < len(file); i++ {
		if file[i] == '\n' {
			value = buf.String()
			buf.Reset()
			if field != "" && value != "" {
				err := os.Setenv(field, value)
				if err != nil {
					return err
				}
			}
			field, value = "", ""
			continue
		}
		if file[i] == '=' {
			field = buf.String()
			buf.Reset()
			continue
		}
		buf.WriteByte(file[i])
	}
	return nil
}
