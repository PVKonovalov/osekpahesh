package currency

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strconv"
	"strings"
)

type Currency struct {
	Cent int64
}

type Rate struct {
	UCent int64
}

func New(cents int64) *Currency {
	return &Currency{cents}
}

func (c *Currency) String() string {
	return fmt.Sprintf("%d.%02d", c.Cent/100, c.Cent%100)
}

func (c *Currency) Add(other Currency) {
	c.Cent += other.Cent
}

func (c *Currency) Format(format string) string {
	return humanize.FormatFloat(format, float64(c.Cent)/100.0)
}

// Rate multiplies the currency amount by the rate, divides by 1,000,000 (micro-cents), and rounds the result to the nearest cent.
func (c *Currency) Rate(rate Rate) *Currency {
	return &Currency{
		Cent: (c.Cent*rate.UCent + 500_000) / 1_000_000,
	}
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Currency.
func (c *Currency) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Cent = 0

	var asString string
	if err := unmarshal(&asString); err != nil {
		return fmt.Errorf("invalid currency format: '%s' %v", asString, err)
	}

	parts := strings.Split(asString, ".")

	if len(parts) > 0 {
		cents, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}
		c.Cent += cents * 100
	}

	if len(parts) == 1 {
		return nil
	}

	if len(parts) == 2 {
		cents, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return err
		}
		c.Cent += cents
		return nil
	}

	return fmt.Errorf("invalid currency format: '%s' (expected format 'X.XX')", asString)
}

func (c *Rate) String() string {
	return fmt.Sprintf("%d.%d", c.UCent/1_000_000, c.UCent%1_000_000)
}
func (c *Rate) StringLen(length int) string {
	return fmt.Sprintf("%d.%d", c.UCent/1_000_000, c.UCent%1_000_000)[:length-1]
}
func (c *Rate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.UCent = 0

	var asString string
	if err := unmarshal(&asString); err != nil {
		return fmt.Errorf("invalid currency format: '%s' %v", asString, err)
	}

	parts := strings.Split(asString, ".")

	if len(parts) > 0 {
		cents, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}
		c.UCent += cents * 1_000_000 // Assuming rate is in micro-cents
	}

	if len(parts) == 1 {
		return nil
	}

	if len(parts) == 2 {
		cents, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return err
		}
		c.UCent += cents
		return nil
	}

	return fmt.Errorf("invalid currency format: '%s' (expected format 'X.XX')", asString)
}
