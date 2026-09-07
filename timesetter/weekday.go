package timesetter

import (
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/strdist.mod/v2/strdist"
	"github.com/nickwells/tempus.mod/tempus"
)

// WeekdaySetter allows you to set the value of a time.Weekday variable.
type WeekdaySetter struct {
	psetter.ValueReqMandatory
	psetter.ValueChecker[time.Weekday]

	// You must set a Value, the program will panic if not. This is a pointer
	// to the weekday value that the setter is setting.
	Value *time.Weekday

	// You must set a Locale, the program will panic if not. This is a
	// pointer to the mappings between weekday names and time.Weekday values
	// that the setter will use to convert the parameter value.
	Locale *tempus.Locale

	// If you set a ValDesc this will be used in help messages to label the
	// Value. If you leave this blank then the Value will be labelled
	// "weekday".
	ValDesc string
}

// suggestAltVal will suggest a possible alternative value for the parameter
// value. It will find those strings in the set of possible values that are
// closest to the given value
func (s WeekdaySetter) suggestAltVal(val string) string {
	names := s.Locale.WeekdayNames()

	return strdist.SuggestionString(strdist.SuggestedVals(val, names))
}

// SetWithVal (called when a value follows a parameter) parses the paramVal
// as a weekday using the WeekdaySetter's Locale. An error is returned if no
// weekday can be found. Only if the value is parsed successfully is the
// Value set.
func (s WeekdaySetter) SetWithVal(_ string, paramVal string) error {
	w, err := s.Locale.ToWeekday(paramVal)
	if err != nil {
		return fmt.Errorf("%v%s", err, s.suggestAltVal(paramVal))
	}

	if err := s.ApplyChecks(w); err != nil {
		return err
	}

	*s.Value = w

	return nil
}

// AllowedValues returns a string describing the allowed values
func (s WeekdaySetter) AllowedValues() string {
	names := s.Locale.WeekdayNames()

	rval := strings.Join(names, ", ")

	return rval
}

// ValDescribe returns a string describing the value that can follow the
// parameter
func (s WeekdaySetter) ValDescribe() string {
	if s.ValDesc != "" {
		return s.ValDesc
	}

	return "weekday"
}

// CurrentValue returns the current setting of the parameter value
func (s WeekdaySetter) CurrentValue() string {
	return s.Value.String()
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil, if the Locale is nil or if one of the check functions is
// nil.
func (s WeekdaySetter) CheckSetter(name string) {
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
