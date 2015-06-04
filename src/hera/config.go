package hera

import (
	"errors"
	"io/ioutil"

	"hera/yaml"
)

var SERVER = make(map[interface{}]string)

var run_mode string

var need_mode = map[string]bool{
	"dev":    true,
	"demo":   true,
	"beta":   true,
	"online": true,
}

type Config struct {
	confPath string
	data     map[interface{}]interface{}
}

func NewConfig(filename string) *Config {
	if filename == "" {
		panic("config file is empty")
	}
	config := &Config{
		confPath: filename,
		data:     make(map[interface{}]interface{}),
	}
	err := config.Init(filename)
	if err != nil {
		panic("config init fail")
	}
	return config
}

func (this *Config) Init(filename string) error {
	stream, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("load config file has error")
	}
	return yaml.Unmarshal(stream, this.data)
}

func MakeServerVar(config *Config) error {
	if config == nil {
		return errors.New("config is empty")
	}
	dict := config.data
	mode, ok := dict["__mode"]
	run_mode = mode.(string)

	if ok == false || need_mode[run_mode] != true {
		panic("mode is illega")
	}

	CpMapValue(dict, &SERVER)
	return nil
}

func CpMapValue(from interface{}, to *map[interface{}]string) {
	switch fromVal := from.(type) {
	case map[interface{}]interface{}:
		for key, value := range fromVal {
			if need_mode[key.(string)] == true && key != run_mode {
				continue
			}
			_, ok_map := value.(map[interface{}]interface{})
			_, ok_slice := value.([]interface{})
			if !ok_map && !ok_slice {
				(*to)[key] = value.(string)
			} else {
				CpMapValue(value, to)
			}
		}
	case []interface{}:
		for _, value := range fromVal {
			CpMapValue(value, to)
		}
	default:
		return
	}
}
