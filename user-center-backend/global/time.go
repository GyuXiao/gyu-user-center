package global

import (
	"fmt"
	"time"
)

// 为了使 time.Time 的时间格式支持 RFC3339，在这里实现 MarshalJSON 和 UnmarshalJSON 方法

const TimeFormat = "2006-01-02 15:04:05"

type JsonTime time.Time

func (jt *JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, time.Time(*jt).Format(TimeFormat))
	return []byte(formatted), nil
}

func (jt *JsonTime) UnmarshalJSON(data []byte) error {
	if string(data) == `""` {
		return nil
	}
	t, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*jt = JsonTime(t)
	return err
}
