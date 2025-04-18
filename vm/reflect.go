package smanchai

import (
	"fmt"
	"reflect"
)

func Reflect(data any) (*Data, error) {
	return toData(reflect.ValueOf(data).Kind(), reflect.ValueOf(data))
}

func toData(typ reflect.Kind, data reflect.Value) (*Data, error) {
	switch typ {
	case reflect.Struct:
		omap := map[string]*Data{}
		for i := 0; i < data.NumField(); i++ {
			kind := data.Field(i).Kind()
			name := data.Type().Field(i).Name
			value := data.Field(i)
			v, err := toData(kind, value)
			if err != nil {
				return nil, err
			}
			omap[name] = v
		}
		return &Data{
			Type: DTypeObject,
			Value: &DataObjectMap{
				Data: omap,
			},
		}, nil
	case reflect.Array:
		oarr := []*Data{}
		for i := 0; i < data.NumField(); i++ {
			kind := data.Field(i).Kind()
			value := data.Field(i)
			v, err := toData(kind, value)
			if err != nil {
				return nil, err
			}
			oarr = append(oarr, v)
		}
		return &Data{
			Type: DTypeArray,
			Value: &DataObjectArray{
				Data: oarr,
			},
		}, nil
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return &Data{
			Type:  DTypeInt,
			Value: data.Int(),
		}, nil
	case reflect.Float32, reflect.Float64:
		return &Data{
			Type:  DTypeDouble,
			Value: data.Float(),
		}, nil
	case reflect.Chan:
		return &Data{
			Type:  DTypeChar,
			Value: data.Int(),
		}, nil
	case reflect.Bool:
		m := 0
		if data.Bool() {
			m = 1
		}
		return &Data{
			Type:  DTypeBool,
			Value: m,
		}, nil
	case reflect.String:
		return &Data{
			Type:  DTypeString,
			Value: data.String(),
		}, nil
	}
	return nil, fmt.Errorf("unsupported data type")
}
