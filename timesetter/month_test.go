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
	updFlagNameMonthsetter     = "upd-gf-monthsetter"
	keepBadFlagNameMonthsetter = "keep-bad-monthsetter"
)

var commonMonthsetterGFC = testhelper.GoldenFileCfg{
	DirNames:               []string{"testdata", "MonthSetter"},
	Pfx:                    "gf",
	Sfx:                    "txt",
	UpdFlagName:            updFlagNameMonthsetter,
	KeepBadResultsFlagName: keepBadFlagNameMonthsetter,
}

func init() {
	commonMonthsetterGFC.AddUpdateFlag()
	commonMonthsetterGFC.AddKeepBadResultsFlag()
}

func TestMonthSetter(t *testing.T) {
	const dfltParamName = "set-month"

	month := time.January

	testCases := []paramtest.Setter{
		{
			ID:      testhelper.MkID("bad-setter-no-value"),
			PSetter: timesetter.MonthSetter{},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName +
					": timesetter.MonthSetter Check failed" +
					": the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("bad-setter-bad-locale"),
			PSetter: timesetter.MonthSetter{
				Value: &month,
			},
			ExpPanic: testhelper.MkExpPanic(
				dfltParamName +
					": timesetter.MonthSetter Check failed" +
					": the Setter is improperly constructed" +
					": the Locale has not been set"),
		},
		{
			ID: testhelper.MkID("good-setter-bad-value"),
			PSetter: timesetter.MonthSetter{
				Value:  &month,
				Locale: &tempus.LocEnglish,
			},
			ParamVal: "nonesuch",
			SetWithValErr: testhelper.MkExpErr(
				`unknown month: "nonesuch"`),
		},
		{
			ID: testhelper.MkID("good-setter-bad-but-close-value"),
			PSetter: timesetter.MonthSetter{
				Value:  &month,
				Locale: &tempus.LocEnglish,
			},
			ParamVal: "Fubruary",
			SetWithValErr: testhelper.MkExpErr(
				`unknown month: "Fubruary",`,
				` did you mean "February" or "february"?`),
		},
		{
			ID: testhelper.MkID("good-setter-good-value"),
			PSetter: timesetter.MonthSetter{
				Value:  &month,
				Locale: &tempus.LocEnglish,
			},
			ParamVal: "February",
		},
		{
			ID: testhelper.MkID("good-setter-good-value-fails-checks"),
			PSetter: timesetter.MonthSetter{
				Value:  &month,
				Locale: &tempus.LocEnglish,
				Checks: []check.ValCk[time.Month]{
					func(val time.Month) error {
						if val == time.May {
							return errors.New("month must not be May")
						}

						return nil
					},
				},
			},
			ParamVal: "May",
			SetWithValErr: testhelper.MkExpErr(
				`month must not be May`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.GFC = commonMonthsetterGFC
			if tc.ParamName == "" {
				tc.ParamName = dfltParamName
			}

			tc.SetVR(param.Mandatory)

			month = time.January

			tc.Test(t)
		})
	}
}
