package dto

import (
	"fmt"
	"strings"
	"time"
)

const monthYearLayout = "01-2006"

type MonthYear struct {
	time.Time
}

func (m *MonthYear) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}

	t, err := time.Parse(monthYearLayout, s)
	if err != nil {
		return fmt.Errorf("неверный формат даты, ожидается MM-YYYY: %w", err)
	}

	m.Time = t
	return nil
}

func (m MonthYear) MarshalJSON() ([]byte, error) {
	if m.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + m.Time.Format(monthYearLayout) + `"`), nil
}
