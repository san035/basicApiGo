package dbtarantool

import (
	"encoding/json"
	"github.com/san035/basicApiGo/pkg/logger"
	"time"
)

type CustomTime uint64

const FormatData = `02.01.2006`

// UnmarshalJSON convert the bytes into an interface
// Пример использования:
// user = new(db.UserProfile)
// err = ctx.BodyParser(user)
func (customTime *CustomTime) UnmarshalJSON(b []byte) error {
	//this will help us check the type of our value
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		if v == "" {
			return nil
		}

		data, err := time.Parse(FormatData, v)
		if err != nil {
			return err
		}
		*customTime = CustomTime(data.Unix())
	default:
		return logger.New(ParseDataUser)
	}
	return nil
}

// MarshalJSON переопределение форматы при вызовове Marshal
func (customTime *CustomTime) MarshalJSON() ([]byte, error) {
	timeStr := `"` + customTime.GetData() + `"`
	return []byte(timeStr), nil
}

func (customTime *CustomTime) GetData() string {
	return time.Unix(int64(*customTime), 0).Format(FormatData)
}
