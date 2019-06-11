package structil

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type FinderKeys struct {
	keys []string
}

type confKeys struct {
	Keys []interface{}
}

func NewFinderKeysFromConf(dir string, baseName string) (*FinderKeys, error) {
	viper.SetConfigName(baseName)
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var ck confKeys
	if err := viper.Unmarshal(&ck); err != nil {
		return nil, err
	}

	return newFinderKeysFromConf(ck)
}

func newFinderKeysFromConf(ck confKeys) (*FinderKeys, error) {
	if len(ck.Keys) == 0 {
		return nil, fmt.Errorf("failed to parse or no keys exist in file")
	}

	fks := &FinderKeys{keys: make([]string, 0, len(ck.Keys)+1)}

	var err error
	for _, ckk := range ck.Keys {
		err = fks.addRecursive(ckk, "")
		if err != nil {
			return nil, err
		}
	}

	return fks, nil
}

func (fks *FinderKeys) addRecursive(key interface{}, prefix string) error {
	var res string

	switch t := key.(type) {
	case string:
		res = t
		if prefix != "" {
			res = prefix + defaultSep + res
		}
		fks.keys = append(fks.keys, res) // set here
	case map[string]interface{}:
		var nk string
		var nd interface{}
		for key, value := range t {
			nk = key
			if prefix != "" {
				nk = prefix + "." + nk
			}
			nd = value
			break
		}
		err := fks.addRecursive(nd, nk)
		if err != nil {
			return err
		}
	case map[interface{}]interface{}:
		var nk string
		var nd interface{}
		for key, value := range t {
			nk = key.(string)
			if prefix != "" {
				nk = prefix + "." + nk
			}
			nd = value
			break
		}
		err := fks.addRecursive(nd, nk)
		if err != nil {
			return err
		}
	case []interface{}:
		var err error
		for _, value := range t {
			err = fks.addRecursive(value, prefix)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported type: %#v, prefix: %s", t, prefix)
	}

	return nil
}

func (fks *FinderKeys) Len() int {
	return len(fks.keys)
}

func (fks *FinderKeys) Keys() []string {
	return fks.keys
}

func (fks *FinderKeys) intoAndFindNames(i int) (string, string) {
	s := strings.Split(fks.keys[i], defaultSep)
	if len(s) == 1 {
		return topLevelKey, s[0]
	}

	return strings.Join(s[0:len(s)-1], defaultSep), s[len(s)-1]
}
