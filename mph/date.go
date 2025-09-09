package mph

import (
	"encoding/json"
	"time"

	"braces.dev/errtrace"
)

// Date is a custom type for representing dates in the format YYYYMMDD
type Date struct {
	Time time.Time
}

var _ json.Marshaler = &Date{}
var _ json.Unmarshaler = &Date{}

const dateFormat = "20060102"

// NewDate is used to create a new Date object
func NewDate(year, month, day int) Date {
	return Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)}
}

// NewDatePtr is used to create a pointer to a new Date object
func NewDatePtr(year, month, day int) *Date {
	return &Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)}
}

func (d Date) String() string {
	if d.Time.IsZero() {
		return ""
	}
	return d.Time.Format(dateFormat)
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.String() + `"`), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if v := string(data); v == "null" || v == `""` {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(`"`+dateFormat+`"`, string(data))
	*d = Date{t}
	return errtrace.Wrap(err)
}
