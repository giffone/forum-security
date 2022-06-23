package read_env

import (
	"errors"
	"os"
	"testing"
)

func TestReadEnv(t *testing.T) {
	t.Parallel()
	tData := map[string]struct {
		path, key, expect string
		err               error
	}{
		"env1 normal":       {"test_env1.env", "test", "abcdefg abcdefg", nil},
		"env2 normal":       {"test_env2.env", "test", "aaaa bbbb", nil},
		"env3 missed value": {"test_env3.env", "test", "", nil},
		"env4 only = value": {"test_env4.env", "test", "", nil},
		"env5 missing file": {"test_env5.env", "", "", errors.New("no file")},
	}

	for name, tcase := range tData {
		v := tcase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := ReadEnv(v.path)
			if err != nil && v.err == nil {
				t.Errorf("[%s] error happend while not expected: %s", name, err.Error())
			}
			got := os.Getenv(v.key)
			os.Clearenv()
			if got != v.expect {
				t.Errorf("[%s] expected \"%s\", but got \"%s\"", name, v.expect, got)
			}
		})
	}
}
