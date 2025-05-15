package outis

import "time"

type Indicator struct {
	key       string
	value     float64
	createdAt time.Time
}

// NewIndicator creates a new indicator.
func (ctx *ContextImpl) NewIndicator(key string) *Indicator {
	indicator := &Indicator{key: key, value: 0, createdAt: time.Now()}
	ctx.indicator = append(ctx.indicator, indicator)
	return indicator
}

// GetKey get the key value of an indicator.
func (i *Indicator) GetKey() string { return i.key }

// GetValue get the value of an indicator.
func (i *Indicator) GetValue() float64 { return i.value }

// GetCreatedAt get the creation date of an indicator.
func (i *Indicator) GetCreatedAt() time.Time { return i.createdAt }

// Inc increments the indicator data.
func (i *Indicator) Inc() { i.value++ }

// Add add a value to the indicator.
func (i *Indicator) Add(value float64) { i.value += value }
