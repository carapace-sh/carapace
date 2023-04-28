package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/rsteube/carapace/pkg/xdg"
)

type configMap map[string]interface{}

func (c configMap) Keys() []string {
	keys := make([]string, 0)
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// func (c configMap) Fields(name string, styled bool) ([]string, error) {
// 	if i, ok := c[name]; ok {
// 		fields := make([]string, 0)
// 		t := reflect.TypeOf(i).Elem()
// 		for index := 0; index < t.NumField(); index++ {
// 			field := t.Field(index)
// 			style := ""
// 			if styled {
// 				if field.Type.Name() != "string" {
// 					return nil, fmt.Errorf("invalid field type [name: '%v', type: '%v']", field.Name, field.Type.Name())
// 				}
// 				v := reflect.ValueOf(i).Elem()
// 				style = v.FieldByName(field.Name).String()
// 			}
// 			fields = append(fields, field.Name, field.Tag.Get("desc"), style)
// 		}
// 		return fields, nil
// 	}
// 	return nil, fmt.Errorf("unknown config: '%v'", name)
// }

func (c configMap) Fields(name string, styled bool) ([]Field, error) {
	if i, ok := c[name]; ok {
		fields := make([]Field, 0)
		t := reflect.TypeOf(i).Elem()
		v := reflect.ValueOf(i).Elem()
		for index := 0; index < t.NumField(); index++ {
			field := t.Field(index)
			if styled && field.Type.Name() != "string" {
				return nil, fmt.Errorf("invalid field type [name: '%v', type: '%v']", field.Name, field.Type.Name())
			}
			fields = append(fields, Field{
				Name:        field.Name,
				Description: field.Tag.Get("desc"),
				Style:       v.FieldByName(field.Name).String(), // TODO only if styled
				Tag:         field.Tag.Get("tag"),
				Type:        field.Type,
			})
		}
		return fields, nil
	}
	return nil, fmt.Errorf("unknown config: '%v'", name)
}

var config = struct {
	Configs configMap
	Styles  configMap
}{
	Configs: make(configMap),
	Styles:  make(configMap),
}

func RegisterConfig(name string, i interface{}) {
	config.Configs[name] = i
}

func RegisterStyle(name string, i interface{}) {
	config.Styles[name] = i
}

func Load() error {
	if err := load("styles", config.Styles); err != nil {
		return err
	}

	// TODO duplicated, ok or improve?
	if err := load("configs", config.Configs); err != nil {
		return err
	}
	return nil
}

func load(name string, c configMap) error {
	if dir, err := xdg.UserConfigDir(); err == nil {
		content, err := os.ReadFile(fmt.Sprintf("%v/carapace/%v.json", dir, name))
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		var unmarshalled map[string]map[string]interface{}
		if err := json.Unmarshal(content, &unmarshalled); err != nil {
			return err
		}

		for key, value := range unmarshalled {
			if s, ok := c[key]; ok {
				elem := reflect.ValueOf(s).Elem()
				for k, v := range value {
					if field := elem.FieldByName(k); field != (reflect.Value{}) {
						field.Set(reflect.ValueOf(v).Convert(field.Type()))
					}
				}
			}
		}
	}
	return nil
}

func SetConfig(key, value string) error {
	return set("configs", key, strings.Replace(value, ",", " ", -1))
}

func GetConfigs() []string                         { return config.Configs.Keys() }
func GetConfigFields(name string) ([]Field, error) { return config.Configs.Fields(name, false) }
func GetConfigMap(name string) interface{}         { return config.Configs[name] }

func GetStyleConfigs() []string                   { return config.Styles.Keys() }
func GetStyleFields(name string) ([]Field, error) { return config.Styles.Fields(name, true) }
func SetStyle(key, value string) error {
	return set("styles", key, strings.Replace(value, ",", " ", -1))
}

func set(name, key, value string) error {
	dir, err := xdg.UserConfigDir()
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%v/carapace/%v.json", dir, name)
	content, err := os.ReadFile(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
			return err
		}
		content = []byte("{}")
	}

	var config map[string]map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		return err
	}

	if splitted := strings.Split(key, "."); len(splitted) != 2 {
		return errors.New("invalid key")
	} else {
		if _, ok := config[splitted[0]]; !ok {
			config[splitted[0]] = make(map[string]interface{}, 0)
		}
		if strings.TrimSpace(value) == "" {
			delete(config[splitted[0]], splitted[1])
		} else {
			switch reflect.TypeOf(config[splitted[0]][splitted[1]]).Kind() {
			case reflect.Int:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				config[splitted[0]][splitted[1]] = intValue

			case reflect.String:
				config[splitted[0]][splitted[1]] = value

			case reflect.Slice:
				// TODO
			}
		}
	}

	marshalled, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	os.WriteFile(file, marshalled, os.ModePerm)

	return nil
}
