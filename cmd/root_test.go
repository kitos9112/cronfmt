package cmd

import "testing"

func TestCronfmtValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		errMsg  string
		wantErr bool
	}{
		{
			name:    "noArgs",
			args:    []string{},
			errMsg:  "Invalid number of arguments\n",
			wantErr: true,
		},
		{
			name:    "tooManyArgs",
			args:    []string{"*", "10", "20", "*", "12", "echo \"hello\"", "extra"},
			errMsg:  "Invalid number of arguments\n",
			wantErr: true,
		},
		{
			name:    "fewArgs",
			args:    []string{"1", "2", "3", "4", "5"},
			errMsg:  "Invalid number of arguments\n",
			wantErr: true,
		},
		{
			name:    "validArgs",
			args:    []string{"*", "10", "20", "12", "6", "echo \"hello\""},
			errMsg:  "",
			wantErr: false,
		},
		{
			name:    "emptyArg",
			args:    []string{"", "10", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Empty argument identified\n",
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMinute",
			args:    []string{"_", "10", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: _ (each argument should fall within '^[\\*\\/\\d,-]+$' regex)\n",
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionDayOfWeek",
			args:    []string{"*", "10", "20", "12", "bananas", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: bananas (each argument should fall within '^[\\*\\/\\d,-]+$' regex)\n",
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2A",
			args:    []string{"*/5,5", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: */5,5 Special operator `*` cannot be used in conjunction `-` or `,`\n",
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2A1",
			args:    []string{"**", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: ** Special operator `*` or `/` can only be used once\n",
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2B",
			args:    []string{"5,-5", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: 5,-5 Special operator `-` cannot be used in conjunction `/` `*` or `,`\n",
			wantErr: true,
		},
		{
			name:    "invalidCronExpressionMixOperator2B2",
			args:    []string{"5-1-1", "10", "20", "12", "3", "echo \"hello\""},
			errMsg:  "Invalid cron expression value: 5-1-1 Special operator `-` can only be used once\n",
			wantErr: true,
		},
		{
			name:    "invalidMinute",
			args:    []string{"61", "10", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid minute value: 61\n",
			wantErr: true,
		},
		{
			name:    "invalidHour",
			args:    []string{"*", "25", "20", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid hour value: 25\n",
			wantErr: true,
		},
		{
			name:    "invalidDayOfMonth",
			args:    []string{"*", "10", "32", "12", "6", "echo \"hello\""},
			errMsg:  "Invalid day of month value: 32\n",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAndExtractArgs(tt.args)
			// test whether or not the error flag is set correctly
			t.Log("Test error flag")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateArgs() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test the error message (if any) is output correctly
			t.Log("Evaluate error message")
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("validateArgs() err = %v, errMsg=%v", err, tt.errMsg)
			}
		})
	}
}
