package outis

import (
	"context"
	"fmt"
	"time"
)

type Interval struct {
	hourSet                bool
	startHour, endHour     uint
	minuteSet              bool
	startMinute, endMinute uint
	everySet               bool
	every                  time.Duration
	onlyOnce               bool
	executeFirstTimeNow    bool
	logger                 ILogger
}

func (i Interval) validate() (err error) {
	// if !i.everySet || !i.onlyOnce {
	// 	return errors.New("interval frequency is required")
	// }
	// todo: validar hora e minuto
	return
}

func NewInterval(opts ...IntervalOpts) (interval Interval) {
	for _, opt := range opts {
		opt(&interval)
	}
	return
}

type IntervalOpts func(interval *Interval)

// WithHours sets the start and end time of script execution
func WithHours(start, end uint) IntervalOpts {
	return func(interval *Interval) {
		interval.hourSet, interval.startHour, interval.endHour = true, start, end
	}
}

// WithMinutes sets the start and end minutes of script execution
func WithMinutes(start, end uint) IntervalOpts {
	return func(interval *Interval) {
		interval.minuteSet, interval.startMinute, interval.endMinute = true, start, end
	}
}

func WithEvery(every time.Duration) IntervalOpts {
	return func(interval *Interval) {
		interval.everySet = true
		interval.every = every
	}
}

func WithExecuteFirstTimeNow() IntervalOpts {
	return func(interval *Interval) {
		interval.executeFirstTimeNow = true
	}
}

func WithOnlyOnce() IntervalOpts {
	return func(interval *Interval) {
		interval.onlyOnce = true
	}
}

func (i Interval) mustWait(time int, start, end uint) bool {
	if start <= end {
		return !(time >= int(start) && time <= int(end))
	}

	return !(time >= int(start) || time <= int(end))
}

// func (i Interval) nextTime(now time.Time, hour, minute uint) time.Time {
// 	today := time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(minute), 0, 0, now.Location())
// 	if now.Before(today) || now.Equal(today) {
// 		return today
// 	}
// 	return today.Add(24 * time.Hour)
// }

func (i Interval) Wait(ctx context.Context, now time.Time, isFirstScriptExecution bool) {
	if i.executeFirstTimeNow && isFirstScriptExecution {
		return
	}
	nextExecution, waitDuration := i.calculateNextExecutionTime(now, isFirstScriptExecution)
	if waitDuration > 0 {
		i.logger.Info(fmt.Sprintf("Waiting until next window at %s (in %s)", nextExecution.Format("02/01/2006 15:04:05"), waitDuration))
	}

	select {
	case <-time.After(waitDuration):
	case <-ctx.Done():
	}
}

func (i Interval) calculateNextExecutionTime(now time.Time, isFirstScriptExecution bool) (time.Time, time.Duration) {
	var (
		next              = now
		forceNextInterval bool
	)
	if isFirstScriptExecution {
		hourOK := !i.hourSet || !isOutsideWindow(uint(now.Hour()), i.startHour, i.endHour)
		minuteOK := !i.minuteSet || !isOutsideWindow(uint(now.Minute()), i.startMinute, i.endMinute)
		if hourOK && minuteOK {
			return now, 0
		}
	} else if i.everySet {
		next = now.Add(i.every)
	} else {
		forceNextInterval = true
	}

	if i.hourSet {
		currentHour := uint(next.Hour())
		if isOutsideWindow(currentHour, i.startHour, i.endHour) || forceNextInterval {
			nextHour := int(i.startHour)
			nextDay := next.Day()
			nextMinute := 0
			// If today's window start is already past, aim for the next day.
			if i.startHour <= i.endHour && (currentHour > i.startHour || (currentHour == i.startHour && uint(next.Minute()) > i.startMinute)) {
				nextDay += 1
			}
			// When we jump to a new hour, minute must be reset to its window start.
			if i.minuteSet {
				nextMinute = int(i.startMinute)
			}
			next = time.Date(next.Year(), next.Month(), nextDay, nextHour, nextMinute, 0, 0, next.Location())
		}
	}

	if i.minuteSet {
		currentMinute := uint(next.Minute())
		if isOutsideWindow(currentMinute, i.startMinute, i.endMinute) || forceNextInterval {
			nextHour := next.Hour()
			nextMinute := int(i.startMinute)
			// If this hour's window start is already past, aim for the next hour.
			if i.startMinute <= i.endMinute && currentMinute >= i.startMinute {
				nextHour += 1
			}
			next = time.Date(next.Year(), next.Month(), next.Day(), nextHour, nextMinute, 0, 0, next.Location())
		}
	}

	waitDuration := next.Sub(now)
	if waitDuration < 0 {
		// TODO: add warning
		return now, 0
	}

	return next, waitDuration
}

// backup v.1
// func (i Interval) calculateNextExecutionTime(now time.Time, isFirstScriptExecution bool) (time.Time, time.Duration) {
// 	var (
// 		nextExecutionHour          uint = uint(now.Hour())
// 		nextExecutionMinute        uint = uint(now.Minute())
// 		nextExecutionDay           uint = uint(now.Day())
// 		isFirstExecutionOnInterval      = isFirstScriptExecution
// 	)

// 	if i.minuteSet {
// 		if isOutsideWindow(nextExecutionMinute, i.startMinute, i.endMinute) {
// 			nextExecutionMinute = i.startMinute
// 			isFirstExecutionOnInterval = true
// 			if i.startMinute <= nextExecutionMinute {
// 				nextExecutionHour++
// 			}
// 		}
// 	}

// 	if i.hourSet {
// 		if isOutsideWindow(nextExecutionHour, i.startHour, i.endHour) {
// 			nextExecutionHour = i.startHour
// 			isFirstExecutionOnInterval = true
// 			if i.startHour <= nextExecutionHour {
// 				nextExecutionDay++
// 			}
// 		}
// 	}

// 	var nextExecution = time.Date(now.Year(), now.Month(), int(nextExecutionDay), int(nextExecutionHour), int(nextExecutionMinute), now.Second(), now.Nanosecond(), now.Location())
// 	if !isFirstExecutionOnInterval {
// 		if i.everySet {
// 			nextExecution = nextExecution.Add(i.every)
// 		} else {
// 			nextExecutionMinute = 0
// 			if i.hourSet {
// 				nextExecutionHour = i.startHour
// 				nextExecutionDay = nextExecutionDay + 1
// 			}
// 			if i.minuteSet {
// 				nextExecutionMinute = i.startMinute
// 			}

// 			nextExecution = time.Date(now.Year(), now.Month(), int(nextExecutionDay), int(nextExecutionHour), int(nextExecutionMinute), now.Second(), now.Nanosecond(), now.Location())
// 		}
// 	}

// 	// garantindo que não está no passado
// 	if nextExecution.Before(now) {
// 		// TODO: add warning
// 		nextExecution = now
// 	}

// 	waitDuration := nextExecution.Sub(now)
// 	return nextExecution, waitDuration
// }

// isOutsideWindow checks if a value is outside a [start, end) range, handling wrap-around.
// The end is exclusive.
func isOutsideWindow(value, start, end uint) bool {
	if start <= end {
		// Normal interval, e.g., 10-14. Outside if value < 10 or value >= 14.
		return value < start || value >= end
	}
	// Wrapped interval, e.g., 22-02. Outside if value >= 2 and value < 22.
	return value >= end && value < start
}

// func (i Interval) WaitHourAndMinute(now time.Time) {
// 	startHour := now.Hour()

// 	if i.hourSet {
// 		startHour = int(i.startHour)
// 		if i.mustWait(now.Hour(), i.startHour, i.endHour) {
// 			nextTime := i.nextTime(now, startHour, 0)
// 			i.logger.Info("Waiting until " + nextTime.Format("02/01/2006 15:04:05"))
// 			time.Sleep(nextTime.Sub(now))
// 		}
// 	}

// 	if i.minuteSet {
// 		if i.mustWait(now.Minute(), i.startMinute, i.endMinute) {
// 			nextTime := i.nextTime(now, startHour, int(i.startMinute))
// 			i.logger.Info("Waiting until " + nextTime.Format("02/01/2006 15:04:05"))
// 			time.Sleep(nextTime.Sub(now))
// 		}

// 	}
// }
