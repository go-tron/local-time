package localTime

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	Layout = "2006-01-02 15:04:05"
	Zone   = "Asia/Shanghai"
)

func Now() Time {
	return Time(time.Now())
}

func ParseLocalLayout(l, val string) (Time, error) {
	t, err := time.ParseInLocation(l, val, time.Local)
	if err != nil {
		return Time{}, err
	}
	return Time(t), nil
}
func ParseLocal(val string) (Time, error) {
	t, err := time.ParseInLocation(Layout, val, time.Local)
	if err != nil {
		return Time{}, err
	}
	return Time(t), nil
}
func ParseDate(val string) (Time, error) {
	t, err := time.ParseInLocation("2006-01-02", val, time.Local)
	if err != nil {
		return Time{}, err
	}
	return Time(t), nil
}
func ParseCompact(val string) (Time, error) {
	t, err := time.ParseInLocation("20060102150405", val, time.Local)
	if err != nil {
		return Time{}, err
	}
	return Time(t), nil
}
func ParseCompactDate(val string) (Time, error) {
	t, err := time.ParseInLocation("20060102", val, time.Local)
	if err != nil {
		return Time{}, err
	}
	return Time(t), nil
}

func FromTimestamppb(ts *timestamppb.Timestamp) Time {
	return Time(ts.AsTime())
}

func Since(t Time) time.Duration {
	return time.Since(time.Time(t))
}

func Until(t Time) time.Duration {
	return time.Until(time.Time(t))
}
func Unix(sec int64, nsec int64) Time {
	return Time(time.Unix(sec, nsec))
}

func StartTimeOfYear(t Time) Time {
	d := time.Time(t)
	return Time(time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location()))
}

func EndTimeOfYear(t Time) Time {
	d := time.Time(t)
	return StartTimeOfYear(Time(d)).AddDate(1, 0, 0).Add(-time.Second)
}

func StartTimeOfMonth(t Time) Time {
	d := time.Time(t)
	return Time(time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location()))
}

func EndTimeOfMonth(t Time) Time {
	d := time.Time(t)
	return StartTimeOfMonth(Time(d)).AddDate(0, 1, 0).Add(-time.Second)
}

func StartTimeOfDate(t Time) Time {
	d := time.Time(t)
	return Time(time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()))
}

func EndTimeOfDate(t Time) Time {
	d := time.Time(t)
	return Time(time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location()))
}

func StartTimeOfWeek(t Time) Time {
	d := time.Time(t)
	offset := int(time.Monday - d.Weekday())
	if offset > 0 {
		offset = -6
	}
	return Time(time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location()).AddDate(0, 0, offset))
}

func EndTimeOfWeek(t Time) Time {
	d := time.Time(t)
	offset := int(7 - d.Weekday())
	if offset == 7 {
		offset = 0
	}
	return Time(time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()).AddDate(0, 0, offset))
}

type Time time.Time

func (t Time) Ptr() *Time {
	return &t
}

func (t Time) AddDate(years int, months int, days int) Time {
	ti := time.Time(t).AddDate(years, months, days)
	return Time(ti)
}

func (t Time) Add(d time.Duration) Time {
	ti := time.Time(t).Add(d)
	return Time(ti)
}
func (t Time) Sub(u Time) time.Duration {
	return time.Time(t).Sub(time.Time(u))
}
func (t Time) Before(u Time) bool {
	return time.Time(t).Before(time.Time(u))
}
func (t Time) After(u Time) bool {
	return time.Time(t).After(time.Time(u))
}
func (t Time) Equal(u Time) bool {
	return time.Time(t).Equal(time.Time(u))
}
func (t Time) Format(l string) string {
	return time.Time(t).Format(l)
}
func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}
func (t Time) UnixNano() int64 {
	return time.Time(t).UnixNano()
}
func (t Time) UnixMillisecond() int64 {
	return time.Time(t).UnixNano() / int64(time.Millisecond)
}
func (t Time) ToTimestamppb() *timestamppb.Timestamp {
	return timestamppb.New(time.Time(t))
}

func (t Time) MarshalJSON() ([]byte, error) {
	if reflect.ValueOf(t).IsZero() {
		return []byte(`""`), nil
	}
	tune := t.Format(`"` + Layout + `"`)
	return []byte(tune), nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || string(data) == "\"\"" {
		return
	}
	now, err := ParseLocalLayout(`"`+Layout+`"`, string(data))
	if err != nil {
		return err
	}
	*t = now
	return
}

func (t Time) IsZero() bool {
	return reflect.ValueOf(t).IsZero()
}
func (t Time) String() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return t.Format(Layout)
}
func (t Time) StringMillisecond() string {
	return t.Format(Layout + ".000")
}
func (t Time) Date() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
func (t Time) Month() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return t.Format("2006-01")
}
func (t Time) Week() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	y, w := time.Time(t).ISOWeek()
	return fmt.Sprintf("%d-%02d", y, w)
}
func (t Time) Year() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return strconv.Itoa(time.Time(t).Year())
}
func (t Time) Time() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return t.Format("15:04:05")
}
func (t Time) TimeMillisecond() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return t.Format("15:04:05.000")
}
func (t Time) RFC3339() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z07:00")
}
func (t Time) Compact() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return t.Format("20060102150405")
}
func (t Time) CompactMillisecond() string {
	if reflect.ValueOf(t).IsZero() {
		return ""
	}
	return strings.Replace(t.Format("20060102150405.000"), ".", "", 1)
}

func (t Time) Value() (driver.Value, error) {
	if reflect.ValueOf(t).IsZero() {
		return nil, nil
	}
	return t.Format(Layout), nil
}

func (t *Time) Scan(src interface{}) error {
	switch v := src.(type) {
	case time.Time:
		*t = Time(v)
	case *time.Time:
		*t = Time(*v)
	case Time:
		*t = v
	case *Time:
		*t = *v
	case string:
		ti, err := ParseLocal(v)
		if err != nil {
			return errors.New("Invalid string for LocalTime")
		}
		*t = ti
	case []uint8:
		ti, err := ParseLocal(string(v))
		if err != nil {
			return errors.New("Invalid string for LocalTime")
		}
		*t = ti
	default:
		return errors.New("Incompatible type for LocalTime")
	}
	return nil
}
