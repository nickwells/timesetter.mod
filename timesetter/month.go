package timesetter

import (
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/tempus.mod/tempus"
)

// MonthSetter allows you to set the value of a time.Month variable.
type MonthSetter struct {
	psetter.ValueReqMandatory
	psetter.ValueChecker[time.Month]

	// You must set a Value, the program will panic if not. This is a pointer
	// to the month value that the setter is setting.
	Value *time.Month

	// You must set a Locale, the program will panic if not. This is a
	// pointer to the mappings between month names and time.Month values that
	// the setter will use to convert the parameter value.
	Locale *tempus.Locale

	// If you set a ValDesc this will be used in help messages to label the
	// Value. If you leave this blank then the Value will be labelled
	// "month".
	ValDesc string
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s MonthSetter) suggestAltVal(val string) string {
	names := s.Locale.MonthNames()

	return strdist.SuggestionString(strdist.SuggestedVals(val, names))
}

// SetWithVal (called when a value follows a parameter) parses the paramVal
// as a month using the MonthSetter's Locale. An error is returned if no
// month can be found. Only if the value is parsed successfully is the Value
// set.
func (s MonthSetter) SetWithVal(_ string, paramVal string) error {
	m, err := s.Locale.ToMonth(paramVal)
	if err != nil {
		return fmt.Errorf("%v%s", err, s.suggestAltVal(paramVal))
	}

	if err := s.ApplyChecks(m); err != nil {
		return err
	}

	*s.Value = m

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s MonthSetter) AllowedValues() string {
	names := s.Locale.MonthNames()

	rval := strings.Join(names, ", ")

	return rval
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s MonthSetter) ValDescribe() string {
	if s.ValDesc != "" {
		return s.ValDesc
	}

	return "month"
}

// CurrentValue returns the current setting of the parameter value
func (s MonthSetter) CurrentValue() string {
	return s.Value.String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil, if the Locale is nil or if one of the check functions is
// nil.
func (s MonthSetter) CheckSetter(name string) {
	if s.Value == nil {
		panic(psetter.NilValueMessage(name, fmt.Sprintf("%T", s)))
	}

	if s.Locale == nil {
		panic(psetter.BadSetterMessage(
			name, fmt.Sprintf("%T", s),
			"the Locale has not been set"))
	}

	s.VerifyChecks(name, fmt.Sprintf("%T", s))
}
