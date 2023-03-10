package util

import (
	"fmt"
	"reflect"
)

// TypeOfStruct 直接获取struct的type
func TypeOfStruct(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// ValueOfStruct 直接获取struct的value
func ValueOfStruct(i interface{}) reflect.Value {
	t := reflect.ValueOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// StructName 直接获取struct的name
func StructName(i interface{}) string {
	return TypeOfStruct(i).Name()
}

// StructRange 循环遍历struct的属性，使用handler处理相对应的属性
func StructRange(i interface{}, handler func(t reflect.StructField, v reflect.Value) error) error {
	v := ValueOfStruct(i)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		var (
			fv = v.Field(i)
			ft = t.Field(i)
		)
		if fv.CanInterface() { //过滤不可访问的属性
			if err := handler(ft, fv); err != nil {
				return err
			}
		}
	}
	return nil
}

// Empty 判断某一属性是否为空值
func IsEmpty(i interface{}) bool {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		v := reflect.ValueOf(i)
		return v.IsNil()
	}
	zero := reflect.New(t).Elem()
	return reflect.DeepEqual(i, zero.Interface())
}

// AOrB 重构了一下赋值方向的代码，设置相同名称的属性，优先参数1的赋值，其次参数2的赋值
// 找不到属性会直接弹出错误
// 属性名为FieldName 注意大小写与结构体保持一直
//
//	if a != nil && a.field != zero {
//		res.set(a)
//	} else if b != nil {
//
//		res.set(b)
//	}
func AOrB(fields []string, val, a, b interface{}) error {
	var (
		va = ValueOfStruct(a)
		vb = ValueOfStruct(b)
		vv = ValueOfStruct(val)
	)
	zero := reflect.Value{}
	wrong := func(fName string, val reflect.Value) error {
		return fmt.Errorf("can not find field %s by %s", fName, val.Type().Name())
	}
	for _, v := range fields {
		vvf := vv.FieldByName(v)
		if vvf == zero {
			return wrong(v, vv)
		}
		vaf := va.FieldByName(v)
		vbf := vb.FieldByName(v)
		switch {
		case vaf == zero && vbf == zero:
			continue
		case vaf == zero || IsEmpty(vaf.Interface()):
			vvf.Set(vbf)
		default:
			vvf.Set(vaf)
		}
	}
	return nil
}
