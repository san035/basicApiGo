package common

import (
	"reflect"
)

// Заполнение структуры valueStruct значениями из тега default
func SetValueDefaultByReflect(valueStruct interface{}) {
	itemsType := reflect.TypeOf(valueStruct).Elem()
	itemsValue := reflect.ValueOf(valueStruct)
	for i := 0; i < itemsType.NumField(); i++ {
		itemValue := itemsValue.Elem().Field(i)
		if itemValue.CanSet() {
			newValue := itemsType.Field(i).Tag.Get("default")
			itemValue.SetString(newValue)
		}
	}
	return
}

// StrListStructure Возаращает строку списока имен структуры
func StrListStructure(s interface{}) (rez string) {
	a := reflect.ValueOf(s)
	if a.Kind() != reflect.Ptr {
		return ""
	}
	cointField := reflect.ValueOf(s).Elem().NumField()
	for x := 0; x < cointField; x++ {
		rez += reflect.TypeOf(s).Elem().Field(x).Name + " "
	}
	return
}
