package response

import (
	"strconv"
	"strings"
	"time"
)

type Date time.Time
type DateTime time.Time

const (
	dateLayout     = "2006-01-02 15:04:05"
	dateTimeLayout = "2006-01-02 15:04:05.9999999"
)

func (d *Date) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(strings.Trim(string(b), "\""), 10, 64)
	if err != nil {
		return err
	}
	*d = Date(time.Unix(i, 0))
	return nil
}

func (d Date) String() string {
	return time.Time(d).Format(dateTimeLayout)
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(`"`+dateTimeLayout+`"`, string(b))
	if err != nil {
		return err
	}
	*d = DateTime(t)
	return nil
}

func (d DateTime) String() string {
	return time.Time(d).Format(dateTimeLayout)
}
