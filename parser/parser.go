package parser

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strconv"
	"strings"
)

func (e *Config) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if err := e.Viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// 	do nothing
		default:
			return err
		}
	}
	_ = e.Viper.Unmarshal(rawVal, opts...)
	e.readEnvs(rawVal)
	e.RewriteEnv(rawVal)
	return nil
}

func (e *Config) RewriteEnv(rawVal any) {
	vl := reflect.ValueOf(rawVal).Elem()
	for i := 0; i < vl.NumField(); i++ {
		data := e.Get("test1.0")
		field := vl.Field(i)
		fmt.Println(data)
		if field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				data := e.Get(fmt.Sprintf("%s"))
				if data == nil {
					continue
				}
				field.Index(j).Set(reflect.ValueOf("10"))
			}
		}
	}
}

func (e *Config) readEnvs(rawVal interface{}) {
	e.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	e.bindEnvs(rawVal)
}

func (e *Config) bindEnvs(in interface{}, prev ...string) {
	ifv := reflect.ValueOf(in)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	for i := 0; i < ifv.NumField(); i++ {
		fv := ifv.Field(i)
		if fv.Kind() == reflect.Ptr {
			if fv.IsZero() {
				fv = reflect.New(fv.Type().Elem()).Elem()
			} else {
				fv = fv.Elem()
			}
		}
		t := ifv.Type().Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if ok {
			if tv == ",squash" {
				e.bindEnvs(fv.Interface(), prev...)
				continue
			}
		} else {
			tv = t.Name
		}
		switch fv.Kind() {
		case reflect.Struct:
			e.bindEnvs(fv.Interface(), append(prev, tv)...)
		case reflect.Map:
			iter := fv.MapRange()
			for iter.Next() {
				if key, ok := iter.Key().Interface().(string); ok {
					e.bindEnvs(iter.Value().Interface(), append(prev, tv, key)...)
				}
			}
		case reflect.Slice:
			env := strings.Join(append(prev, tv), ".")
			_ = e.Viper.BindEnv(env)
			for i := 0; i < fv.Len(); i++ {
				data := fv.Index(i)
				if data.Kind() != reflect.Struct && data.Kind() != reflect.Map && data.Kind() != reflect.Slice {
					env := strings.Join(append(prev, tv, strconv.Itoa(i)), ".")
					_ = e.Viper.BindEnv(env)
					continue
				}
				e.bindEnvs(data.Interface(), append(prev, tv, strconv.Itoa(i))...)
			}
		default:
			env := strings.Join(append(prev, tv), ".")
			_ = e.Viper.BindEnv(env)
		}
	}
}
