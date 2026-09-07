package timesetter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramtest"
	"github.com/nickwells/tempus.mod/tempus"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
	"github.com/nickwells/timesetter.mod/timesetter"
)

const (
	updFlagNameWeekdaysetter     = "upd-gf-weekdaysetter"
	keepBadFlagNameWeekdaysetter = "keep-bad-weekdaysetter"
)

var commonWeekdaysetterGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "WeekdaySetter"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameWeekdaysetter,
	KeepBadResultsFlagName: keepBadFlagNameWeekdaysetter,
}

func init() {
	commonWeekdaysetterGFC.AddUpdateFlag()
	commonWeekdaysetterGFC.AddKeepBadResultsFlag()
}

func TestWeekdaySetter(t *testing.T) {
	const dfltParamName = "set-weekday"

	weekday := time.Monday

	testCases := []paramtest.Setter{
		{
			ID:      testhelper.MkID("bad-setter-no-value"),
			PSetter: timesetter.WeekdaySetter{},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName +
					": timesetter.WeekdaySetter Check failed" +
					": the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("bad-setter-bad-locale"),
			PSetter: timesetter.WeekdaySetter{
				Value: &weekday,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName +
					": timesetter.WeekdaySetter Check failed" +
					": the Setter is improperly constructed" +
					": the Locale has not been set"),
		},
		{
			ID: testhelper.MkID("good-setter-bad-value"),
			PSetter: timesetter.WeekdaySetter{
				Value:  &weekday,
				Locale: &tempus.LocEnglish,
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				`unknown day of week: "nonesuch"`),
		},
		{
			ID: testhelper.MkID("good-setter-bad-but-close-value"),
			PSetter: timesetter.WeekdaySetter{
				Value:  &weekday,
				Locale: &tempus.LocEnglish,
			},
			ParamVal: "Wudnesday",
			SetWithValErr: testhelper.MkExpErr(
				`unknown day of week: "Wudnesday",`,
				` did you mean "Wednesday" or "wednesday"?`),
		},
		{
			ID: testhelper.MkID("good-setter-good-value"),
			PSetter: timesetter.WeekdaySetter{
				Value:  &weekday,
				Locale: &tempus.LocEnglish,
			},
			ParamVal: "Wednesday",
		},
		{
			ID: testhelper.MkID("good-setter-good-value-fails-checks"),
			PSetter: timesetter.WeekdaySetter{
				Value:  &weekday,
				Locale: &tempus.LocEnglish,
				Checks: []check.ValCk[time.Weekday]{
					func(val time.Weekday) error {
						if val == time.Wednesday {
							return errors.New(
								"day of week must not be a Wednesday")
						}

						return nil
					},
				},
			},
			ParamVal: "Wednesday",
			SetWithValErr: testhelper.MkExpErr(
				`day of week must not be a Wednesday`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.GFC = commonWeekdaysetterGFC
			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			tc.SetVR(param.Mandatory)

			weekday = time.Monday

			tc.Test(t)
		})
	}
}
