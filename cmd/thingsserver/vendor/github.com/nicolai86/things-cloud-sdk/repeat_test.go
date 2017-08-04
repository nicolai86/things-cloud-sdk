package thingscloud

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var (
	rcEveryDay          = []byte(`{"ia":1504396800,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":16,"sr":1499644800,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rcEveryDayEndDate   = []byte(`{"ia":1519776000,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":16,"sr":1519776000,"fa":1,"rc":0,"ts":0,"ed":1519862400}`)
	rcEveryDayEndRepeat = []byte(`{"ia":1519776000,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":16,"sr":1519776000,"fa":1,"rc":2,"ts":0}`)
	rcEvery2ndDay       = []byte(`{"ia":1504396800,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":16,"sr":1499644800,"fa":2,"rc":0,"ts":0,"ed":64092211200}`)

	rcEveryWeekOnMonday              = []byte(`{"ia":1504483200,"rrv":4,"tp":0,"of":[{"wd":1}],"fu":256,"sr":1499644800,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rcEveryWeekOnMondayEndDate       = []byte(`{"ia":1520208000,"rrv":4,"tp":0,"of":[{"wd":1}],"fu":256,"sr":1519776000,"fa":1,"rc":0,"ts":0,"ed":1521331200}`)
	rcEveryWeekOnMondayEndRepeat     = []byte(`{"ia":1520208000,"rrv":4,"tp":0,"of":[{"wd":1}],"fu":256,"sr":1520208000,"fa":1,"rc":2,"ts":0}`)
	rcEveryWeekOnMondayAndTuesday    = []byte(`{"ia":1504483200,"rrv":4,"tp":0,"of":[{"wd":1},{"wd":2}],"fu":256,"sr":1517356800,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rcEvery2ndWeekOnMonday           = []byte(`{"ia":1504483200,"rrv":4,"tp":0,"of":[{"wd":1}],"fu":256,"sr":1499644800,"fa":2,"rc":0,"ts":0,"ed":64092211200}`)
	rcEvery2ndWeekOnMondayAndTuesday = []byte(`{"ia":1504483200,"rrv":4,"tp":0,"of":[{"wd":1},{"wd":2}],"fu":256,"sr":1499644800,"fa":2,"rc":0,"ts":0,"ed":64092211200}`)

	rc1stDayEveryMonth             = []byte(`{"ia":1506816000,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":8,"sr":1499644800,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rc1stDayEveryMonthEndDate      = []byte(`{"ia":1517443200,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":8,"sr":1517443200,"fa":1,"rc":0,"ts":0,"ed":1519948800}`)
	rc1stDayEveryMonthEndRepeat    = []byte(`{"ia":1517443200,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":8,"sr":1517443200,"fa":1,"rc":2,"ts":0}`)
	rc1stDayAnd3rdDayEveryMonth    = []byte(`{"ia":1506816000,"rrv":4,"tp":0,"of":[{"dy":0},{"dy":2}],"fu":8,"sr":1502323200,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rc1stDayAnd2ndMondayEveryMonth = []byte(`{"ia":1504224000,"rrv":4,"tp":0,"of":[{"dy":0},{"wdo":2,"wd":1}],"fu":8,"sr":1504224000,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rc1stDayEvery2ndMonth          = []byte(`{"ia":1504224000,"rrv":4,"tp":0,"of":[{"dy":0}],"fu":8,"sr":1499644800,"fa":2,"rc":0,"ts":0,"ed":64092211200}`)
	rc1stAndLastDayEveryMonth      = []byte(`{"ia":1501459200,"rrv":4,"tp":0,"of":[{"dy":0},{"dy":-1}],"fu":8,"sr":1499472000,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rcLastDayEvery2ndMonth         = []byte(`{"ia":1506729600,"rrv":4,"tp":0,"of":[{"dy":-1}],"fu":8,"sr":1501545600,"fa":2,"rc":0,"ts":0,"ed":64092211200}`)
	rcLastMondayEvery2ndMonth      = []byte(`{"ia":1506297600,"rrv":4,"tp":0,"of":[{"wdo":-1,"wd":1}],"fu":8,"sr":1504137600,"fa":2,"rc":0,"ts":0,"ed":64092211200}`)
	rcFirstMondayEvery2ndMonth     = []byte(`{"ia":1502064000,"rrv":4,"tp":0,"of":[{"wdo":1,"wd":1}],"fu":8,"sr":1509321600,"fa":2,"rc":0,"ts":0,"ed":64092211200}`)

	rc1stDayJanuaryEveryYear                     = []byte(`{"ia":1514764800,"rrv":4,"tp":0,"of":[{"dy":0,"mo":0}],"fu":4,"sr":1512345600,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rcLastDayJanuaryEveryYear                    = []byte(`{"ia":1517356800,"rrv":4,"tp":0,"of":[{"dy":-1,"mo":0}],"fu":4,"sr":1514764800,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rcLastDayJanuaryEveryYearEndDate             = []byte(`{"ia":1517356800,"rrv":4,"tp":0,"of":[{"dy":-1,"mo":0}],"fu":4,"sr":1499472000,"fa":1,"rc":0,"ts":0,"ed":1551225600}`)
	rcLastDayJanuaryEveryYearEndRepeat           = []byte(`{"ia":1517356800,"rrv":4,"tp":0,"of":[{"dy":-1,"mo":0}],"fu":4,"sr":1517356800,"fa":1,"rc":2,"ts":0}`)
	rcLastDayFebuaryEveryYear                    = []byte(`{"ia":1519776000,"rrv":4,"tp":0,"of":[{"dy":-1,"mo":1}],"fu":4,"sr":1517443200,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rc1stAndLastDayFebuaryEveryYear              = []byte(`{"ia":1517443200,"rrv":4,"tp":0,"of":[{"dy":0,"mo":1},{"dy":-1,"mo":1}],"fu":4,"sr":1499472000,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rc1stJanuaryAnd1stMarchEveryYear             = []byte(`{"ia":1514764800,"rrv":4,"tp":0,"of":[{"dy":0,"mo":0},{"dy":0,"mo":2}],"fu":4,"sr":1504224000,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rc1stJanuaryAndLastWednesdayFebuaryEveryYear = []byte(`{"ia":1514764800,"rrv":4,"tp":0,"of":[{"dy":0,"mo":0},{"wdo":-1,"wd":3,"mo":1}],"fu":4,"sr":1514764800,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
	rcLastWednesdayFebuaryEveryYear              = []byte(`{"ia":1519776000,"rrv":4,"tp":0,"of":[{"wdo":-1,"wd":3,"mo":1}],"fu":4,"sr":1519776000,"fa":1,"rc":0,"ts":0,"ed":64092211200}`)
)

func TestRepeaterConfiguration_IsNeverending(t *testing.T) {
	ts := &Timestamp{}
	ts.UnmarshalJSON([]byte(`64092211200`))
	rc := RepeaterConfiguration{LastScheduledAt: ts}

	if !rc.IsNeverending() {
		t.Fatal("Expect repeater to be never ending, but isn't")
	}
}

func Test_nthDayOfMonthOfYear(t *testing.T) {
	testCases := []struct {
		Date         string
		Month        int
		Day          int
		ExpectedDate string
	}{
		{"2017-01-10", 0, 0, "2017-01-01"},
		{"2017-01-10", 0, 30, "2017-01-31"},
		// {"2017-02-10", "2017-02-01"},
		// {"2017-02-01", "2017-02-01"},
		// {"2017-03-20", "2017-03-01"},
		// {"2018-02-10", "2018-02-01"},
		// {"2018-02-01", "2018-02-01"},
		// {"2018-03-20", "2018-03-01"},
	}
	for _, testCase := range testCases {
		ts, err := time.Parse("2006-01-02", testCase.Date)
		if err != nil {
			t.Error(err)
		}
		nt := nthDayOfMonthOfYear(ts, testCase.Month, testCase.Day)
		if nt.Format("2006-01-02") != testCase.ExpectedDate {
			t.Errorf("Expected %d day of %d month of %q to be %q, but got %q", testCase.Day, testCase.Month, testCase.Date, testCase.ExpectedDate, nt.Format("2006-01-02"))
		}
	}
}

func Test_firstDayOfMonth(t *testing.T) {
	testCases := []struct {
		Date         string
		ExpectedDate string
	}{
		{"2017-01-10", "2017-01-01"},
		{"2017-02-10", "2017-02-01"},
		{"2017-02-01", "2017-02-01"},
		{"2017-03-20", "2017-03-01"},
		{"2018-02-10", "2018-02-01"},
		{"2018-02-01", "2018-02-01"},
		{"2018-03-20", "2018-03-01"},
	}
	for _, testCase := range testCases {
		ts, err := time.Parse("2006-01-02", testCase.Date)
		if err != nil {
			t.Error(err)
		}
		if firstDayOfMonth(ts).Format("2006-01-02") != testCase.ExpectedDate {
			t.Errorf("Expected first day of month of %q to be %q, but got %q", testCase.Date, testCase.ExpectedDate, firstDayOfMonth(ts).Format("2006-01-02"))
		}
	}
}

func Test_lastDayOfMonth(t *testing.T) {
	testCases := []struct {
		Date         string
		ExpectedDate string
	}{
		{"2017-01-10", "2017-01-31"},
		{"2017-02-10", "2017-02-28"},
		{"2017-02-01", "2017-02-28"},
		{"2017-03-20", "2017-03-31"},
		{"2020-02-10", "2020-02-29"},
		{"2020-02-01", "2020-02-29"},
		{"2020-03-20", "2020-03-31"},
	}
	for _, testCase := range testCases {
		ts, err := time.Parse("2006-01-02", testCase.Date)
		if err != nil {
			t.Error(err)
		}
		if lastDayOfMonth(ts).Format("2006-01-02") != testCase.ExpectedDate {
			t.Errorf("Expected first day of month of %q to be %q, but got %q", testCase.Date, testCase.ExpectedDate, lastDayOfMonth(ts).Format("2006-01-02"))
		}
	}
}

func Test_lastWeekdayOfMonth(t *testing.T) {
	testCases := []struct {
		Date         string
		Weekday      time.Weekday
		ExpectedDate string
	}{
		{"2018-02-02", time.Wednesday, "2018-02-28"},
		{"2018-02-02", time.Thursday, "2018-02-22"},
		{"2018-02-02", time.Saturday, "2018-02-24"},
	}
	for _, testCase := range testCases {
		ts, err := time.Parse("2006-01-02", testCase.Date)
		if err != nil {
			t.Error(err)
		}
		nt := lastWeekdayOfMonth(ts, testCase.Weekday)
		if nt.Format("2006-01-02") != testCase.ExpectedDate {
			t.Errorf("Expected last %s of month of %q to be %q, but got %q", testCase.Weekday.String(), testCase.Date, testCase.ExpectedDate, nt.Format("2006-01-02"))
		}
	}
}

func Test_nthWeekdayOfMonth(t *testing.T) {
	testCases := []struct {
		Date         string
		Nth          int
		Weekday      time.Weekday
		ExpectedDate string
	}{
		{"2018-02-02", 1, time.Wednesday, "2018-02-07"},
		{"2018-02-02", 2, time.Thursday, "2018-02-08"},
		{"2018-02-02", 3, time.Saturday, "2018-02-17"},
	}
	for _, testCase := range testCases {
		ts, err := time.Parse("2006-01-02", testCase.Date)
		if err != nil {
			t.Error(err)
		}
		nt := nthWeekdayOfMonth(ts, testCase.Weekday, testCase.Nth)
		if nt.Format("2006-01-02") != testCase.ExpectedDate {
			t.Errorf("Expected %dth %s of month of %q to be %q, but got %q", testCase.Nth, testCase.Weekday.String(), testCase.Date, testCase.ExpectedDate, nt.Format("2006-01-02"))
		}
	}
}

func TestRepeaterConfiguration_NextScheduledAtEndDate(t *testing.T) {
	testCases := []struct {
		Title             string
		Data              []byte
		ExpectedNextDates []string
	}{
		{"Every day w/ date", rcEveryDayEndDate, []string{"2018-02-28", "2018-03-01", "0001-01-01"}},
		{"Every day w/ count", rcEveryDayEndRepeat, []string{"2018-02-28", "2018-03-01", "0001-01-01"}},
		{"Every week on monday w/ date", rcEveryWeekOnMondayEndDate, []string{"2018-03-05", "2018-03-12", "0001-01-01"}},
		{"Every week on monday w/ count", rcEveryWeekOnMondayEndRepeat, []string{"2018-03-05", "2018-03-12", "0001-01-01"}},
		{"Every 1st day of every month w/ date", rc1stDayEveryMonthEndDate, []string{"2018-02-01", "2018-03-01", "0001-01-01"}},
		{"Every 1st day of every month w/ count", rc1stDayEveryMonthEndRepeat, []string{"2018-02-01", "2018-03-01", "0001-01-01"}},
		{"Every last day of january every year w/ date", rcLastDayJanuaryEveryYearEndDate, []string{"2018-01-31", "2019-01-31", "0001-01-01"}},
		{"Every last day of january every year w/ date", rcLastDayJanuaryEveryYearEndRepeat, []string{"2018-01-31", "2019-01-31", "0001-01-01"}},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase %q", testCase.Title), func(t *testing.T) {
			var rc RepeaterConfiguration
			err := json.Unmarshal(testCase.Data, &rc)
			if err != nil {
				t.Fatalf("Failed to deserialize repeater configuration: %v", err)
			}

			for i, date := range testCase.ExpectedNextDates {
				nts := rc.NextScheduledAt(i)
				if nts.Format("2006-01-02") != date {
					t.Errorf("Expected %q for next %d date, but got %q", date, i+1, nts.Format("2006-01-02"))
				}
			}
		})
	}
}

func TestRepeaterConfiguration_ComputeFirstScheduledAt(t *testing.T) {
	testCases := []struct {
		Title            string
		Data             []byte
		StartAt          string
		FirstScheduledAt string
	}{
		{"Every week on monday", rcEveryWeekOnMonday, "2017-09-03", "2017-09-04"},
		{"Every week on monday next week", rcEveryWeekOnMonday, "2017-09-05", "2017-09-11"},
		{"Every week on monday same day", rcEveryWeekOnMonday, "2017-09-04", "2017-09-04"},

		{"Every 1st day every month", rc1stDayEveryMonth, "2017-09-03", "2017-10-01"},
		{"Every last day every 2nd month", rcLastDayEvery2ndMonth, "2017-09-03", "2017-09-30"},
		{"Every last monday every 2nd month", rcLastMondayEvery2ndMonth, "2017-09-03", "2017-09-25"},
		{"Every 1st and 3rd day every month 1st", rc1stDayAnd3rdDayEveryMonth, "2017-09-01", "2017-09-01"},
		{"Every 1st and 3rd day every month 3rd", rc1stDayAnd3rdDayEveryMonth, "2017-09-03", "2017-09-03"},
		{"Every 1st and 3rd day every month", rc1stDayAnd3rdDayEveryMonth, "2017-09-02", "2017-09-03"},
		{"Every 1st and 3rd day every month next 1st", rc1stDayAnd3rdDayEveryMonth, "2017-08-30", "2017-09-01"},
		{"Every 1st and last day every month", rc1stAndLastDayEveryMonth, "2017-08-15", "2017-08-31"},
		{"Every 1st and last day every month", rc1stAndLastDayEveryMonth, "2017-08-31", "2017-08-31"},
		{"Every 1st and last day every month", rc1stAndLastDayEveryMonth, "2017-09-01", "2017-09-01"},

		{"Every 1st January and last Wednesday of Febuary every year", rc1stJanuaryAndLastWednesdayFebuaryEveryYear, "2017-12-22", "2018-01-01"},
		{"Every 1st January and last Wednesday of Febuary every year 1st", rc1stJanuaryAndLastWednesdayFebuaryEveryYear, "2018-01-01", "2018-01-01"},
		{"Every 1st January and last Wednesday of Febuary every year last", rc1stJanuaryAndLastWednesdayFebuaryEveryYear, "2018-01-02", "2018-02-28"},
		{"Every 1st January and last Wednesday of Febuary every year last", rc1stJanuaryAndLastWednesdayFebuaryEveryYear, "2018-03-02", "2019-01-01"},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase %q", testCase.Title), func(t *testing.T) {
			var rc RepeaterConfiguration
			err := json.Unmarshal(testCase.Data, &rc)
			if err != nil {
				t.Fatalf("Failed to deserialize repeater configuration: %v", err)
			}

			tt, err := time.Parse("2006-01-02", testCase.StartAt)
			if err != nil {
				t.Fatalf("Failed to parse date: %v", err)
			}
			nts := rc.ComputeFirstScheduledAt(tt)
			if nts.Format("2006-01-02") != testCase.FirstScheduledAt {
				t.Errorf("Expected start date for %s to be %s, but got %q", testCase.StartAt, testCase.FirstScheduledAt, nts.Format("2006-01-02"))
			}

		})
	}
}

func TestRepeaterConfiguration_NextScheduledAt(t *testing.T) {
	testCases := []struct {
		Title                      string
		Data                       []byte
		ExpectedFrequencyUnit      FrequencyUnit
		ExpectedFrequencyAmplitude int64
		ExpectedNextDates          []string
	}{
		{"Every day", rcEveryDay, FrequencyUnitDaily, 1, []string{"2017-09-03", "2017-09-04", "2017-09-05"}},
		{"Every 2nd day", rcEvery2ndDay, FrequencyUnitDaily, 2, []string{"2017-09-03", "2017-09-05", "2017-09-07"}},

		{"Every week on monday", rcEveryWeekOnMonday, FrequencyUnitWeekly, 1, []string{"2017-09-04", "2017-09-11"}},
		{"Every week on monday and tuesday", rcEveryWeekOnMondayAndTuesday, FrequencyUnitWeekly, 1, []string{"2017-09-04", "2017-09-05", "2017-09-11", "2017-09-12"}},
		{"Every 2n week on monday", rcEvery2ndWeekOnMonday, FrequencyUnitWeekly, 2, []string{"2017-09-04", "2017-09-18"}},
		{"Every 2n week on monday and tuesday", rcEvery2ndWeekOnMondayAndTuesday, FrequencyUnitWeekly, 2, []string{"2017-09-04", "2017-09-05", "2017-09-18", "2017-09-19"}},

		{"Every first day of every month", rc1stDayEveryMonth, FrequencyUnitMonthly, 1, []string{"2017-10-01", "2017-11-01"}},
		{"Every first and third day of every month", rc1stDayAnd3rdDayEveryMonth, FrequencyUnitMonthly, 1, []string{"2017-10-01", "2017-10-03", "2017-11-01", "2017-11-03"}},
		{"Every first day and 2nd monday of every month", rc1stDayAnd2ndMondayEveryMonth, FrequencyUnitMonthly, 1, []string{"2017-09-01", "2017-09-11", "2017-10-01", "2017-10-09"}},
		{"Every first day of every 2nd month", rc1stDayEvery2ndMonth, FrequencyUnitMonthly, 2, []string{"2017-09-01", "2017-11-01", "2018-01-01"}},
		{"Every last day of every 2nd month", rcLastDayEvery2ndMonth, FrequencyUnitMonthly, 2, []string{"2017-09-30", "2017-11-30", "2018-01-31", "2018-03-31"}},
		{"Every last Monday of every 2nd month", rcLastMondayEvery2ndMonth, FrequencyUnitMonthly, 2, []string{"2017-09-25", "2017-11-27", "2018-01-29", "2018-03-26"}},
		// {"Every last and first day of every month", rc1stAndLastDayEveryMonth, FrequencyUnitMonthly, 1, []string{"2017-07-31", "2017-08-01", "2017-08-31", "2017-09-01", "2017-09-30"}},
		{"Every first Monday of every 2nd month", rcFirstMondayEvery2ndMonth, FrequencyUnitMonthly, 2, []string{"2017-08-07", "2017-10-02", "2017-12-04"}},

		{"Every first day of january of every year", rc1stDayJanuaryEveryYear, FrequencyUnitYearly, 1, []string{"2018-01-01", "2019-01-01", "2020-01-01"}},
		{"Every last day of january of every year", rcLastDayJanuaryEveryYear, FrequencyUnitYearly, 1, []string{"2018-01-31", "2019-01-31", "2020-01-31"}},
		{"Every last day of february of every year", rcLastDayFebuaryEveryYear, FrequencyUnitYearly, 1, []string{"2018-02-28", "2019-02-28", "2020-02-29", "2021-02-28"}},
		{"Every first and last day of february of every year", rc1stAndLastDayFebuaryEveryYear, FrequencyUnitYearly, 1, []string{"2018-02-01", "2018-02-28", "2019-02-01", "2019-02-28"}},
		{"Every first day of january and first day of march of every year", rc1stJanuaryAnd1stMarchEveryYear, FrequencyUnitYearly, 1, []string{"2018-01-01", "2018-03-01", "2019-01-01", "2019-03-01"}},
		{"Every first day of january and last Wednesday of febuary of every year", rc1stJanuaryAndLastWednesdayFebuaryEveryYear, FrequencyUnitYearly, 1, []string{"2018-01-01", "2018-02-28", "2019-01-01", "2019-02-27", "2020-01-01", "2020-02-26"}},
		{"Every last Wednesday of febuary of every year", rcLastWednesdayFebuaryEveryYear, FrequencyUnitYearly, 1, []string{"2018-02-28", "2019-02-27", "2020-02-26", "2021-02-24", "2022-02-23"}},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("testCase %q", testCase.Title), func(t *testing.T) {
			var rc RepeaterConfiguration
			err := json.Unmarshal(testCase.Data, &rc)
			if err != nil {
				t.Fatalf("Failed to deserialize repeater configuration: %v", err)
			}
			if rc.FrequencyAmplitude != testCase.ExpectedFrequencyAmplitude {
				t.Fatalf("Expected fa of %d but got %d", testCase.ExpectedFrequencyAmplitude, rc.FrequencyAmplitude)
			}
			if rc.FrequencyUnit != testCase.ExpectedFrequencyUnit {
				t.Fatalf("Expected fu of %d but got %d", testCase.ExpectedFrequencyUnit, rc.FrequencyUnit)
			}

			for i, date := range testCase.ExpectedNextDates {
				nts := rc.NextScheduledAt(i)
				if nts.Format("2006-01-02") != date {
					t.Errorf("Expected %q for next %d date, but got %q", date, i+1, nts.Format("2006-01-02"))
				}
			}
		})
	}
}
