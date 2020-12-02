package yaml

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ardanlabs/conf"
	"gopkg.in/yaml.v3"
)

type YAMLSourcer struct {
	m map[string]string
}

// NewSource returns a conf.Sourcer and, potentially, an error if a
// read error occurs or the Reader contains an invalid YAML document.
func NewSource(r io.Reader) (conf.Sourcer, error) {
	if r == nil {
		return &YAMLSourcer{m: nil}, nil
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tmpMap := make(map[string]interface{})
	err = yaml.Unmarshal(src, &tmpMap)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for key, value := range tmpMap {
		switch v := value.(type) {
		case float64:
			m[key] = strings.TrimRight(fmt.Sprintf("%f", v), "0.")
		case bool:
			m[key] = fmt.Sprintf("%t", v)
		case string:
			m[key] = value.(string)
		}
	}

	return &YAMLSourcer{m: m}, nil
}

func (s *YAMLSourcer) Source(fld conf.Field) (string, bool) {
	if fld.Options.ShortFlagChar != 0 {
		flagKey := fld.Options.ShortFlagChar
		k := strings.ToLower(string(flagKey))
		if val, found := s.m[k]; found {
			return val, found
		}
	}

	k := strings.ToLower(strings.Join(fld.FlagKey, `_`))
	val, found := s.m[k]
	return val, found
}
