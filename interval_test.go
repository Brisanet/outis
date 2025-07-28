package outis

import (
	"testing"
	"time"
)

func TestCalculateNextExecutionTime(t *testing.T) {
	tests := []struct {
		name                   string
		interval               Interval
		now                    time.Time
		isFirstScriptExecution bool
		expectedNextExecution  time.Time
		expectedWaitDuration   time.Duration
	}{
		// {
		// 	name:                   "First execution without any interval set",
		// 	interval:               Interval{},
		// 	now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
		// 	isFirstScriptExecution: true,
		// 	expectedNextExecution:  time.Date(0, 1, 1, 12, 30, 0, 0, time.UTC),
		// 	expectedWaitDuration:   0,
		// },
		// {
		// 	name:                   "Subsequent execution without any interval set",
		// 	interval:               Interval{},
		// 	now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
		// 	isFirstScriptExecution: false,
		// 	expectedNextExecution:  time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
		// 	expectedWaitDuration:   0,
		// },
		// Minute interval - first executions
		{
			name: "First execution with minute interval (current time within interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			expectedWaitDuration:   0,
		},
		{
			name: "First execution with minute interval (current time before interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 15, 0, 0, time.UTC),
			expectedWaitDuration:   5 * time.Minute,
		},
		{
			name: "First execution with minute interval (current time after interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 50, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 15, 0, 0, time.UTC),
			expectedWaitDuration:   25 * time.Minute,
		},
		// Minute interval - subsequent executions
		{
			name: "Subsequent execution with minute interval (current time within interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 15, 0, 0, time.UTC),
			expectedWaitDuration:   45 * time.Minute,
		},
		{
			name: "Subsequent execution with minute interval (current time before interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 15, 0, 0, time.UTC),
			expectedWaitDuration:   5 * time.Minute,
		},
		{
			name: "Subsequent execution with minute interval (current time after interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 50, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 15, 0, 0, time.UTC),
			expectedWaitDuration:   25 * time.Minute,
		},
		// Hour interval - first executions
		{
			name: "First execution with hour interval (current time within interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			expectedWaitDuration:   0,
		},
		{
			name: "First execution with hour interval (current time before interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
			},
			now:                    time.Date(2023, 1, 1, 8, 30, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			expectedWaitDuration:   2 * time.Hour,
		},
		{
			name: "First execution with hour interval (current time after interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
			},
			now:                    time.Date(2023, 1, 1, 16, 30, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 2, 10, 30, 0, 0, time.UTC),
			expectedWaitDuration:   18 * time.Hour,
		},
		// Hour interval - subsequent executions
		{
			name: "Subsequent execution with hour interval (current time within interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 2, 10, 00, 0, 0, time.UTC),
			expectedWaitDuration:   21*time.Hour + 30*time.Minute,
		},
		{
			name: "Subsequent execution with hour interval (current time before interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
			},
			now:                    time.Date(2023, 1, 1, 8, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			expectedWaitDuration:   2 * time.Hour,
		},
		{
			name: "Subsequent execution with hour interval (current time after interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
			},
			now:                    time.Date(2023, 1, 1, 16, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 2, 10, 30, 0, 0, time.UTC),
			expectedWaitDuration:   18 * time.Hour,
		},
		// Both hour and minute
		{
			name: "First execution with both hour and minute intervals",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     14,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 15, 0, 0, time.UTC),
			expectedWaitDuration:   5 * time.Minute,
		},
		{
			name: "Subsequent execution with both hour and minute intervals",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     14,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
			},
			now:                    time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 15, 0, 0, time.UTC),
			expectedWaitDuration:   5 * time.Minute,
		},
		{
			name: "First execution with 'every' interval",
			interval: Interval{
				everySet: true,
				every:    30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: true,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			expectedWaitDuration:   0,
		},
		{
			name: "Subsequent execution with 'every' interval",
			interval: Interval{
				everySet: true,
				every:    30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC),
			expectedWaitDuration:   30 * time.Minute,
		},
		// Minute + Every combinations
		{
			name: "Subsequent execution with minute interval and every",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 20, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 20, 0, 0, time.UTC).Add(30 * time.Minute),
			expectedWaitDuration:   30 * time.Minute,
		},
		{
			name: "Subsequent execution with minute interval and every (before minute interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 15, 0, 0, time.UTC),
			expectedWaitDuration:   5 * time.Minute,
		},
		{
			name: "Subsequent execution with minute interval and every (within minute interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 20, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 50, 0, 0, time.UTC),
			expectedWaitDuration:   30 * time.Minute,
		},
		{
			name: "Subsequent execution with minute interval and every (after minute interval)",
			interval: Interval{
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 50, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 15, 0, 0, time.UTC),
			expectedWaitDuration:   25 * time.Minute,
		},
		// Hour + Every combinations
		{
			name: "Subsequent execution with hour interval and every",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
				everySet:  true,
				every:     1 * time.Hour,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 30, 0, 0, time.UTC),
			expectedWaitDuration:   1 * time.Hour,
		},
		// Hour + Every combinations - before, within and after hour interval
		{
			name: "Subsequent execution with hour interval and every (before hour interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
				everySet:  true,
				every:     1 * time.Hour,
			},
			now:                    time.Date(2023, 1, 1, 8, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 10, 30, 0, 0, time.UTC),
			expectedWaitDuration:   2 * time.Hour,
		},
		{
			name: "Subsequent execution with hour interval and every (within hour interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
				everySet:  true,
				every:     1 * time.Hour,
			},
			now:                    time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 30, 0, 0, time.UTC),
			expectedWaitDuration:   1 * time.Hour,
		},
		{
			name: "Subsequent execution with hour interval and every (after hour interval)",
			interval: Interval{
				hourSet:   true,
				startHour: 10,
				endHour:   14,
				everySet:  true,
				every:     1 * time.Hour,
			},
			now:                    time.Date(2023, 1, 1, 15, 30, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 2, 10, 30, 0, 0, time.UTC),
			expectedWaitDuration:   19 * time.Hour,
		},
		// Hour + Minute + Every combinations - before, within and after intervals
		{
			name: "Subsequent execution with hour, minute and every (before both intervals)",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     14,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 8, 10, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 10, 15, 0, 0, time.UTC),
			expectedWaitDuration:   2*time.Hour + 5*time.Minute,
		},
		{
			name: "Subsequent execution with hour, minute and every (within hour, before minute)",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     14,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 10, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 15, 0, 0, time.UTC),
			expectedWaitDuration:   5 * time.Minute,
		},
		{
			name: "Subsequent execution with hour, minute and every (within both intervals)",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     14,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 20, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 12, 50, 0, 0, time.UTC),
			expectedWaitDuration:   30 * time.Minute,
		},
		{
			name: "Subsequent execution with hour, minute and every (within hour, after minute)",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     14,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 12, 50, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 1, 13, 15, 0, 0, time.UTC),
			expectedWaitDuration:   25 * time.Minute,
		},
		{
			name: "Subsequent execution with hour, minute and every (after both intervals)",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     14,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 15, 50, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 2, 10, 15, 0, 0, time.UTC),
			expectedWaitDuration:   18*time.Hour + 25*time.Minute,
		},
		{
			name: "Complex case with all intervals set",
			interval: Interval{
				hourSet:     true,
				startHour:   10,
				endHour:     13,
				minuteSet:   true,
				startMinute: 15,
				endMinute:   45,
				everySet:    true,
				every:       30 * time.Minute,
			},
			now:                    time.Date(2023, 1, 1, 13, 50, 0, 0, time.UTC),
			isFirstScriptExecution: false,
			expectedNextExecution:  time.Date(2023, 1, 2, 10, 15, 0, 0, time.UTC),
			expectedWaitDuration:   20*time.Hour + 25*time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextExec, waitDuration := tt.interval.calculateNextExecutionTime(tt.now, tt.isFirstScriptExecution)

			if !nextExec.Equal(tt.expectedNextExecution) {
				t.Errorf("expected next execution time %v, got %v", tt.expectedNextExecution, nextExec)
			}

			if waitDuration != tt.expectedWaitDuration {
				t.Errorf("expected wait duration %v, got %v", tt.expectedWaitDuration, waitDuration)
			}
		})
	}
}
