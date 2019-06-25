/**
 * 问题：
 * 使用 json 格式化 struct 时，time.Time 被格式化成
 * ”2006-01-02T15:04:05.999999999Z07:00“ 格式
 *
 * golang 的 time.Time 的默认 json 格式化格式叫做 RFC3339。
 * 好像是一种国际标准，被推荐用作 json 时间的标准格式。
 * 但是使用中不需要这种，而且不容易解析。
 */

package utils

import (
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

/**
 * String 格式化输出
 */
func (t Time) String() string {
	return time.Time(t).Format(TimeFormat)
}
