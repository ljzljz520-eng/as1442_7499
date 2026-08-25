package store

type Clock interface {
	Now() string
}

type SequenceClock struct {
	value string
	step  int
}

func NewSequenceClock(start string) *SequenceClock {
	return &SequenceClock{value: start, step: 0}
}

func (c *SequenceClock) Now() string {
	c.step++
	return c.value + "." + twoDigits(c.step)
}

func twoDigits(value int) string {
	if value < 10 {
		return "0" + string(rune('0'+value))
	}
	if value < 100 {
		return string(rune('0'+value/10)) + string(rune('0'+value%10))
	}
	return "99"
}
