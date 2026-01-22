package stdlib

import (
	"fmt"
	"time"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

// RegisterTime registers time core functions in the environment.
func RegisterTime(env *runtime.Environment) error {
	functions := map[string]runtime.CoreFunction{
		"time:now":      runtime.NewCoreFunction(timeNow),
		"time:unix":     runtime.NewCoreFunction(timeUnix),
		"time:year":     runtime.NewCoreFunction(timeYear),
		"time:month":    runtime.NewCoreFunction(timeMonth),
		"time:day":      runtime.NewCoreFunction(timeDay),
		"time:hour":     runtime.NewCoreFunction(timeHour),
		"time:minute":   runtime.NewCoreFunction(timeMinute),
		"time:second":   runtime.NewCoreFunction(timeSecond),
		"time:format":   runtime.NewCoreFunction(timeFormat),
		"time:parse":    runtime.NewCoreFunction(timeParse),
		"time:add":      runtime.NewCoreFunction(timeAdd),
		"time:sub":      runtime.NewCoreFunction(timeSub),
		"time:diff":     runtime.NewCoreFunction(timeDiff),
		"time:is-leap":  runtime.NewCoreFunction(timeIsLeap),
	}

	for name, fn := range functions {
		if _, err := env.Define(name, fn); err != nil {
			return fmt.Errorf("failed to register time function `%s`: %v", name, err)
		}
	}

	return nil
}

// timeNow implements the current time core function.
// Usage: (time:now) => 1737489123
func timeNow(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:now"

	if err := expectArgs(name, 0, args); err != nil {
		return nil, err
	}

	return runtime.NewNumber(float64(time.Now().Unix())), nil
}

// timeUnix implements the Unix timestamp conversion core function.
// Usage: (time:unix 1737489123) => 1737489123
func timeUnix(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:unix"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	return runtime.NewNumber(timestamp.Value), nil
}

// timeYear implements the year extraction core function.
// Usage: (time:year 1737489123) => 2025
func timeYear(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:year"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	return runtime.NewNumber(float64(t.Year())), nil
}

// timeMonth implements the month extraction core function.
// Usage: (time:month 1737489123) => 1
func timeMonth(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:month"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	return runtime.NewNumber(float64(t.Month())), nil
}

// timeDay implements the day extraction core function.
// Usage: (time:day 1737489123) => 21
func timeDay(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:day"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	return runtime.NewNumber(float64(t.Day())), nil
}

// timeHour implements the hour extraction core function.
// Usage: (time:hour 1737489123) => 14
func timeHour(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:hour"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	return runtime.NewNumber(float64(t.Hour())), nil
}

// timeMinute implements the minute extraction core function.
// Usage: (time:minute 1737489123) => 25
func timeMinute(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:minute"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	return runtime.NewNumber(float64(t.Minute())), nil
}

// timeSecond implements the second extraction core function.
// Usage: (time:second 1737489123) => 23
func timeSecond(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:second"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	return runtime.NewNumber(float64(t.Second())), nil
}

// timeFormat implements the time formatting core function.
// Usage: (time:format 1737489123 "2006-01-02") => "2025-01-21"
func timeFormat(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:format"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	layout, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	return runtime.NewString(t.Format(layout.Value)), nil
}

// timeParse implements the time parsing core function.
// Usage: (time:parse "2025-01-21" "2006-01-02") => 1737417600
func timeParse(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:parse"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	value, err := expectString(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	layout, err := expectString(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(layout.Value, value.Value)
	if err != nil {
		return nil, fmt.Errorf("`%s` failed to parse: %v", name, err)
	}

	return runtime.NewNumber(float64(t.Unix())), nil
}

// timeAdd implements the time addition core function.
// Usage: (time:add 1737489123 3600) => 1737492723
func timeAdd(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:add"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	seconds, err := expectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	newTime := t.Add(time.Duration(seconds.Value) * time.Second)
	return runtime.NewNumber(float64(newTime.Unix())), nil
}

// timeSub implements the time subtraction core function.
// Usage: (time:sub 1737489123 3600) => 1737485523
func timeSub(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:sub"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	seconds, err := expectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t := time.Unix(int64(timestamp.Value), 0).UTC()
	newTime := t.Add(-time.Duration(seconds.Value) * time.Second)
	return runtime.NewNumber(float64(newTime.Unix())), nil
}

// timeDiff implements the time difference core function.
// Usage: (time:diff 1737492723 1737489123) => 3600
func timeDiff(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:diff"

	if err := expectArgs(name, 2, args); err != nil {
		return nil, err
	}

	timestamp1, err := expectNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	timestamp2, err := expectNumber(name, 1, args[1])
	if err != nil {
		return nil, err
	}

	t1 := time.Unix(int64(timestamp1.Value), 0).UTC()
	t2 := time.Unix(int64(timestamp2.Value), 0).UTC()
	diff := t1.Sub(t2).Seconds()

	return runtime.NewNumber(diff), nil
}

// timeIsLeap implements the leap year check core function.
// Usage: (time:is-leap 2024) => true
func timeIsLeap(args ...runtime.Value) (runtime.Value, error) {
	const name = "time:is-leap"

	if err := expectArgs(name, 1, args); err != nil {
		return nil, err
	}

	yearNum, err := expectIntegerNumber(name, 0, args[0])
	if err != nil {
		return nil, err
	}

	year := int(yearNum.Value)
	isLeap := (year%4 == 0 && year%100 != 0) || (year%400 == 0)

	return runtime.NewBool(isLeap), nil
}
