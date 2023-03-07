package dbtarantool

import (
	"fmt"
	"github.com/san035/basicApiGo/pkg/logger"
)

// ListInterfaceToStruct конвертация списка интерфейсов listInterface в другую структуру
// mapNameColumnToData - map
// ключи id из списка listInterface
// значения - ссылки на значения в структуре
func ListInterfaceToStruct(listInterface *[]interface{}, mapNameColumnToData map[int]interface{}) error {
	countColumnData := len((*listInterface))
	for idColumn, distValue := range mapNameColumnToData {
		if idColumn >= countColumnData {
			continue
		}

		switch valueTo := distValue.(type) {
		case *string:
			*valueTo = (*listInterface)[idColumn].(string)
		case *int64:
			// приемник с типом int64
			switch valueFrom := (*listInterface)[idColumn].(type) {
			case uint64:
				*valueTo = int64(valueFrom)
			case int64:
				*valueTo = valueFrom
			default:
				return logger.New(ParseDataUser, idColumn)
			}
		case *uint64:
			*valueTo = (*listInterface)[idColumn].(uint64)
		case *CustomTime: // uint64
			*valueTo = CustomTime((*listInterface)[idColumn].(uint64))
		default:
			return logger.New(ParseDataUser, map[string]interface{}{"Колонка:": idColumn, "Тип": fmt.Sprintf("%T", distValue)})
		}
	}
	return nil
}
