package stdlib

import (
	"fmt"
	"strings"
	"time"

	"github.com/danielspk/tatu-lang/pkg/core"
	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterTime registers time functions.
func RegisterTime(natives map[string]runtime.NativeFunction) {
	natives["time:now"] = runtime.NewNativeFunction(timeNow)
	natives["time:unix"] = runtime.NewNativeFunction(timeUnix)
	natives["time:year"] = runtime.NewNativeFunction(timeYear)
	natives["time:month"] = runtime.NewNativeFunction(timeMonth)
	natives["time:day"] = runtime.NewNativeFunction(timeDay)
	natives["time:hour"] = runtime.NewNativeFunction(timeHour)
	natives["time:minute"] = runtime.NewNativeFunction(timeMinute)
	natives["time:second"] = runtime.NewNativeFunction(timeSecond)
	natives["time:format"] = runtime.NewNativeFunction(timeFormat)
	natives["time:parse"] = runtime.NewNativeFunction(timeParse)
	natives["time:add"] = runtime.NewNativeFunction(timeAdd)
	natives["time:sub"] = runtime.NewNativeFunction(timeSub)
	natives["time:diff"] = runtime.NewNativeFunction(timeDiff)
	natives["time:is-leap"] = runtime.NewNativeFunction(timeIsLeap)
}

// timeNow implements the current time function.
// Usage: (time:now) => 1737489123
func timeNow(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:now"

	if err := core.ExpectArgs(name, 0, args); err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(time.Now().Unix())), nil
}

// timeUnix implements the Unix timestamp conversion function.
// Usage: (time:unix 1737489123) => 1737489123
func timeUnix(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:unix"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(timestamp.Value), nil
}

// timeYear implements the year extraction function.
// Usage: (time:year 1737489123) => 2025
func timeYear(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:year"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()

	return runtime.NewNumber(float64(t.Year())), nil
}

// timeMonth implements the month extraction function.
// Usage: (time:month 1737489123) => 1
func timeMonth(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:month"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()

	return runtime.NewNumber(float64(t.Month())), nil
}

// timeDay implements the day extraction function.
// Usage: (time:day 1737489123) => 21
func timeDay(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:day"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()

	return runtime.NewNumber(float64(t.Day())), nil
}

// timeHour implements the hour extraction function.
// Usage: (time:hour 1737489123) => 14
func timeHour(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:hour"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()

	return runtime.NewNumber(float64(t.Hour())), nil
}

// timeMinute implements the minute extraction function.
// Usage: (time:minute 1737489123) => 25
func timeMinute(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:minute"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()

	return runtime.NewNumber(float64(t.Minute())), nil
}

// timeSecond implements the second extraction function.
// Usage: (time:second 1737489123) => 23
func timeSecond(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:second"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()

	return runtime.NewNumber(float64(t.Second())), nil
}

// timeFormat implements the time formatting function.
// Usage: (time:format 1737489123 "YYYY-MM-DD") => "2025-01-21"
func timeFormat(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:format"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	layout, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()

	return runtime.NewString(t.Format(translateLayout(layout.Value))), nil
}

// timeParse implements the time parsing function.
// Usage: (time:parse "2025-01-21" "YYYY-MM-DD") => 1737417600
func timeParse(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:parse"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	value, err := core.ExpectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	layout, err := core.ExpectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(translateLayout(layout.Value), value.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to parse: %w", name, err)
	}

	return runtime.NewNumber(float64(t.Unix())), nil
}

// timeAdd implements the time addition function.
// Usage: (time:add 1737489123 3600) => 1737492723
func timeAdd(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:add"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	seconds, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	newTime := t.Add(time.Duration(seconds.Value) * time.Second)

	return runtime.NewNumber(float64(newTime.Unix())), nil
}

// timeSub implements the time subtraction function.
// Usage: (time:sub 1737489123 3600) => 1737485523
func timeSub(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:sub"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	seconds, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	newTime := t.Add(-time.Duration(seconds.Value) * time.Second)

	return runtime.NewNumber(float64(newTime.Unix())), nil
}

// timeDiff implements the time difference function.
// Usage: (time:diff 1737492723 1737489123) => 3600
func timeDiff(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:diff"

	if err := core.ExpectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp1, err := core.ExpectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	timestamp2, err := core.ExpectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t1 := time.Unix(int64(timestamp1.Value), 0).UTC()
	t2 := time.Unix(int64(timestamp2.Value), 0).UTC()
	diff := t1.Sub(t2).Seconds()

	return runtime.NewNumber(diff), nil
}

// timeIsLeap implements the leap year check function.
// Usage: (time:is-leap 2024) => true
func timeIsLeap(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:is-leap"

	if err := core.ExpectArgs(name, 1, args); err != nil {
		return nil, err
	}

	yearNum, err := core.ExpectIntegerNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	year := int(yearNum.Value)
	isLeap := (year%4 == 0 && year%100 != 0) || (year%400 == 0)

	return runtime.NewBool(isLeap), nil
}

// translateLayout converts common date patterns to go time format.
func translateLayout(layout string) string {
	replacements := []struct {
		from, to string
	}{
		{"YYYY", "2006"},
		{"YY", "06"},
		{"MMMM", "January"},
		{"MMM", "Jan"},
		{"MM", "01"},
		{"DD", "02"},
		{"HH", "15"},
		{"hh", "03"},
		{"mm", "04"},
		{"ss", "05"},
		{"SSS", "000"},
		{"A", "PM"},
		{"dddd", "Monday"},
		{"ddd", "Mon"},
	}

	result := layout
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.from, r.to)
	}

	return result
}
