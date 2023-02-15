package univers

import "reflect"

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
