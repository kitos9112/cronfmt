package cmd

import (
	"reflect"
	"testing"
)

func TestCronfmtValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		cronfmt *Cronfmt
		errMsg  string
		wantErr bool
	}{
		{
			name:    "noArgs",
			args:    []string{},
			errMsg:  "Invalid number of arguments\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "tooManyArgs",
			args:    []string{"*", "10", "20", "*", "12", "echo \"hello\"", "extra"},
			errMsg:  "Invalid number of arguments\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "fewArgs",
			args:    []string{"1", "2", "3", "4", "5"},
			errMsg:  "Invalid number of arguments\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "emptyArg",
			args:    []string{"", "10", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Empty argument identified\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMinute",
			args:    []string{"_", "10", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: _ (each argument should fall within '^[\\*\\/\\d,-]+$' regex)\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionDayOfWeek",
			args:    []string{"*", "10", "20", "12", "bananas", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: bananas (each argument should fall within '^[\\*\\/\\d,-]+$' regex)\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2A",
			args:    []string{"*/5,5", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: */5,5 Special operator `*` cannot be used in conjunction `-` or `,`\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2A1",
			args:    []string{"**", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: ** Special operator `*` or `/` can only be used once\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2B",
			args:    []string{"5,-5", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: 5,-5 Special operator `-` cannot be used in conjunction `/` `*` or `,`\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2B2",
			args:    []string{"5-1-1", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: 5-1-1 Special operator `-` can only be used once\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidMinute",
			args:    []string{"61", "10", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid minute value: 61\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidHour",
			args:    []string{"*", "25", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid hour value: 25\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidDayOfMonth",
			args:    []string{"*", "10", "32", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid day of month value: 32\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:    "invalidMonth",
			args:    []string{"*", "10", "20", "13", "6", "echo \"hello\""},
			errMsg:  "Invalid month value: 13\n",
			wantErr: true,
		},
		{
			name:    "invalidDayOfWeek",
			args:    []string{"*", "10", "20", "12", "8", "echo \"hello\""},
			errMsg:  "Invalid day of week value: 8\n",
			cronfmt: new(Cronfmt),
			wantErr: true,
		},
		{
			name:   "validArgsExample1",
			args:   []string{"*/18", "*/3", "5,10", "*", "6", "echo \"hello\""},
			errMsg: "",
			cronfmt: &Cronfmt{
				Minute:     "18 36 54",
				Hour:       "0 3 6 9 12 15 18 21",
				DayOfMonth: "5 10",
				Month:      "1 2 3 4 5 6 7 8 9 10 11 12",
				DayOfWeek:  "6",
				Command:    "echo \"hello\"",
			},
			wantErr: false,
		},
		{
			name:   "validArgsExample2",
			args:   []string{"*/18", "*/3", "5,10", "*", "6", "/bin/cat /tmp/myHiddenFile"},
			errMsg: "",
			cronfmt: &Cronfmt{
				Minute:     "18 36 54",
				Hour:       "0 3 6 9 12 15 18 21",
				DayOfMonth: "5 10",
				Month:      "1 2 3 4 5 6 7 8 9 10 11 12",
				DayOfWeek:  "6",
				Command:    "/bin/cat /tmp/myHiddenFile",
			},
			wantErr: false,
		},
		{
			name:   "validArgsExample3",
			args:   []string{"*/30", "*/5", "5,7,9,10", "5", "6", "/usr/bin/call-home -m 'I am alive'"},
			errMsg: "",
			cronfmt: &Cronfmt{
				Minute:     "0 30",
				Hour:       "5 10 15 20",
				DayOfMonth: "5 7 9 10",
				Month:      "5",
				DayOfWeek:  "6",
				Command:    "/usr/bin/call-home -m 'I am alive'",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cronfmt, err := validateAndExtractArgs(tt.args)
			// test whether or not the error flag is set correctly
			if (err != nil) != tt.wantErr {
				t.Log("Test error flag")
				t.Errorf("validateArgs() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test the error message (if any) is output correctly
			if err != nil && err.Error() != tt.errMsg {
				t.Log("Evaluate error message")
				t.Errorf("validateArgs() err = %v, errMsg=%v", err, tt.errMsg)
			}
			// test the cronfmt struct is populated correctly
			if !tt.wantErr && !reflect.DeepEqual(cronfmt, tt.cronfmt) {
				t.Log("Evaluate cronfmt struct")
				t.Errorf("validateArgs() cronfmt = %v, want %v", cronfmt, tt.cronfmt)
			}
		})
	}
}
