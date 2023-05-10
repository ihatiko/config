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
	e.readEnvs(rawVal)
	_ = e.Viper.Unmarshal(rawVal, opts...)
	e.RewriteEnv(rawVal)
	return nil
}

func (e *Config) RewriteEnv(rawVal any) {
	vl := reflect.ValueOf(rawVal).Elem()
	tp := reflect.TypeOf(rawVal).Elem()
	e.rewriteEnv(vl, tp, "")
}

func (e *Config) rewriteEnv(vl reflect.Value, tp reflect.Type, prev string) {
	for i := 0; i < vl.NumField(); i++ {
		fieldName := tp.Field(i).Name
		field := vl.Field(i)
		if field.Kind() == reflect.Ptr {
			if field.IsZero() {
				field = reflect.New(field.Type().Elem()).Elem()
			} else {
				field = field.Elem()
			}
		}
		switch field.Kind() {
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				key := fmt.Sprintf("%s.%s", fieldName, strconv.Itoa(j))
				if prev != "" {
					key = fmt.Sprintf("%s.%s", prev, key)
				}
				fieldData := e.Get(key)
				if fieldData == nil {
					continue
				}
				if _, ok := fieldData.(map[string]interface{}); ok {
					e.rewriteEnv(field.Index(j), field.Index(j).Type(), key)
					continue
				}
				field.Index(j).Set(reflect.ValueOf(fieldData))
			}
		default:
			key := fieldName
			if prev != "" {
				key = fmt.Sprintf("%s.%s", prev, key)
			}
			fieldData := e.getCorrectEnv(field, key)
			if fieldData == nil {
				continue
			}
			if _, ok := fieldData.(map[string]interface{}); ok {
				e.rewriteEnv(field, field.Type(), key)
				continue
			}
			field.Set(reflect.ValueOf(fieldData))
		}
	}
}

func (e *Config) getCorrectEnv(valType reflect.Value, key string) any {
	switch valType.Type().Kind() {
	case reflect.Bool:
		return e.GetBool(key)
	case reflect.Int16:
		return e.GetInt(key)
	case reflect.Int:
		return e.GetInt(key)
	case reflect.Int32:
		return e.GetInt32(key)
	case reflect.Int64:
		return e.GetInt64(key)
	case reflect.String:
		return e.GetString(key)
	default:
		return e.Get(key)
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
