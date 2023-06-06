package utils

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type Timestamp time.Time

// MarshalJSON implements json.Marshaler.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	//entity your serializing here
	stamp := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	var ts int64
	ts, err = strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	theTime := time.Unix(ts, 0)
	*t = Timestamp(theTime)
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *Timestamp) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Timestamp(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t Timestamp) GetTime() time.Time {
	return time.Time(t)
}

// GetUnixTimeSql 获取unix时间戳sql
func GetUnixTimeSql(unixTime int64) string {
	return fmt.Sprintf("FROM_UNIXTIME(%d)", unixTime)
}
