package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type Cronfmt struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Command    string
}

var (
	// VERSION is the version of the application
	VERSION string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "cronfmt",
		Short: "Cron expression parser",
		Long:  "Parses cron expressions and prints their extended space-separated format in stdout",
		Example: `
$ cronfmt "*/15" 0 1,15 "*" 1-5 "/usr/bin/find / -type f .terraform"
+-----------------+------------------------------------+
| CRON EXPRESSION | EXTENDED FORMAT                    |
+-----------------+------------------------------------+
| minute          | 0 15 30 45                         |
| hour            | 0                                  |
| day of month    | 1 15                               |
| month           | 1 2 3 4 5 6 7 8 9 10 11 12         |
| day of week     | 1 2 3 4 5                          |
| command         | /usr/bin/find / -type f .terraform |
+-----------------+------------------------------------+

$ cronfmt "*/18" "*/3" 5,15 "*" 1-5 "/usr/bin/call-home -m 'I am alive'"
+-----------------+------------------------------------+
| CRON EXPRESSION | EXTENDED FORMAT                    |
+-----------------+------------------------------------+
| minute          | 18 36 54                           |
| hour            | 0 3 6 9 12 15 18 21                |
| day of month    | 5 15                               |
| month           | 1 2 3 4 5 6 7 8 9 10 11 12         |
| day of week     | 1 2 3 4 5                          |
| command         | /usr/bin/call-home -m 'I am alive' |
+-----------------+------------------------------------+

$ cronfmt "*/5" 1,2,3,4,5 1 "*/4" 1-5 "/bin/cat /tmp/myHiddenFile"
+-----------------+-----------------------------------+
| CRON EXPRESSION | EXTENDED FORMAT                   |
+-----------------+-----------------------------------+
| minute          | 0 5 10 15 20 25 30 35 40 45 50 55 |
| hour            | 1 2 3 4 5                         |
| day of month    | 1                                 |
| month           | 4 8 12                            |
| day of week     | 1 2 3 4 5                         |
| command         | /bin/cat /tmp/myHiddenFile        |
+-----------------+-----------------------------------+
`,
		Args: cobra.MaximumNArgs(6),
		Run: func(cmd *cobra.Command, args []string) {
			cronfmt, err := validateAndExtractArgs(args)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Cron Expression", "Extended Format"})
			t.AppendSeparator()
			t.AppendRows([]table.Row{
				{"minute", cronfmt.Minute},
				{"hour", cronfmt.Hour},
				{"day of month", cronfmt.DayOfMonth},
				{"month", cronfmt.Month},
				{"day of week", cronfmt.DayOfWeek},
				{"command", cronfmt.Command},
			})
			t.AppendSeparator()
			t.Render()
		},
		// All flags are disabled and passed to root command as arguments
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// validateAndExtractArgs ensures the cron format suitability of the arguments passed to the root command
func validateAndExtractArgs(args []string) (*Cronfmt, error) {
	if len(args) != 6 {
		return nil, errors.New("Invalid number of arguments\n")
	}

	// Initialise an empty cron struct
	cronfmt := Cronfmt{}

	// Iterate through each argument, validate it and extract its extended value
	for index, arg := range args {
		if arg == "" {
			return nil, errors.New("Empty argument identified\n")
		}
		switch index {
		case 0:
			// validate minute positional argument
			if err := validateCronExpr(arg, 0, 59, "minute"); err != nil {
				return nil, err
			}
			minuteExpanded, err := expandCronExpr(arg, 0, 59, "minute")
			if err != nil {
				return nil, err
			}
			cronfmt.Minute = minuteExpanded
		case 1:
			// validate hour positional argument
			if err := validateCronExpr(arg, 0, 23, "hour"); err != nil {
				return nil, err
			}
			hourExpanded, err := expandCronExpr(arg, 0, 23, "hour")
			if err != nil {
				return nil, err
			}
			cronfmt.Hour = hourExpanded
		case 2:
			// validate day of month positional argument
			if err := validateCronExpr(arg, 1, 31, "day of month"); err != nil {
				return nil, err
			}
			dayOfMonthExpanded, err := expandCronExpr(arg, 1, 31, "day of month")
			if err != nil {
				return nil, err
			}
			cronfmt.DayOfMonth = dayOfMonthExpanded
		case 3:
			// validate month positional argument
			if err := validateCronExpr(arg, 1, 12, "month"); err != nil {
				return nil, err
			}
			monthExpanded, err := expandCronExpr(arg, 1, 12, "month")
			if err != nil {
				return nil, err
			}
			cronfmt.Month = monthExpanded
		case 4:
			// validate day of week positional argument
			if err := validateCronExpr(arg, 0, 6, "day of week"); err != nil {
				return nil, err
			}
			dayOfWeekExpanded, err := expandCronExpr(arg, 0, 6, "day of week")
			if err != nil {
				return nil, err
			}
			cronfmt.DayOfWeek = dayOfWeekExpanded
		}

		cronfmt.Command = args[5]
	}

	return &cronfmt, nil
}

// validateCronExpr ensures the cron expression entered is valid and sticks to the rules
func validateCronExpr(expression string, minRange int, maxRange int, position string) error {
	// 1) Verify the cron expression syntax
	cronRegex := regexp.MustCompile(`^[\*\/\d,-]+$`)
	numberRegex := regexp.MustCompile(`\d+`)
	if !cronRegex.MatchString(expression) {
		return errors.New("Invalid cron expression value: " + expression + " (each argument should fall within '^[\\*\\/\\d,-]+$' regex)\n")
	}
	// 2) Re-assure the cron expression only contains a valid especial operator. RE2 does not come with negative lookahead or lookbehind
	// A) Contains `*` but not `,` or `-`
	// A.1) Ensure `*` or `/` are only used once`
	if isSubstring(expression, `\*`) || isSubstring(expression, `/`) {
		if isSubstring(expression, `-`) || isSubstring(expression, `,`) {
			return errors.New("Invalid cron expression value: " + expression + " Special operator `*` cannot be used in conjunction `-` or `,`\n")
		} else if getMatchesLength(expression, `\*`) > 1 || getMatchesLength(expression, `/`) > 1 {
			return errors.New("Invalid cron expression value: " + expression + " Special operator `*` or `/` can only be used once\n")
		}
	}
	// B) Contains `-` but not `,` `*` or `/`
	// B.1) Ensure the `-` is only used once
	if isSubstring(expression, `-`) {
		if isSubstring(expression, `/`) || isSubstring(expression, `,`) || isSubstring(expression, `\*`) {
			return errors.New("Invalid cron expression value: " + expression + " Special operator `-` cannot be used in conjunction `/` `*` or `,`\n")
		} else if getMatchesLength(expression, `-`) > 1 {
			return errors.New("Invalid cron expression value: " + expression + " Special operator `-` can only be used once\n")
		}
	}

	// 3) Match all numbers in the expression argument and ensure they are within the given range
	numbers := numberRegex.FindAllString(expression, -1)
	for _, numberStr := range numbers {
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return errors.New("Unknown error occurred while parsing " + numberStr + " as uint\n" + err.Error() + "\n")
		}
		if number < minRange || number > maxRange {
			return errors.New("Invalid " + position + " value: " + strconv.Itoa(number) + "\n")
		}
	}
	return nil
}

// extendCronExpr extends the cron expression to its full space-separated format
func expandCronExpr(expression string, minRange int, maxRange int, position string) (string, error) {
	var expandedExpr string

	var matchAll = regexp.MustCompile(`^\*$`)
	var matchRange = regexp.MustCompile(`^(?P<minRange>\d+)-(?P<maxRange>\d+)$`)
	var matchAllStep = regexp.MustCompile(`^\*\/(?P<step>\d+)$`)
	var matchList = regexp.MustCompile(`^(\d+)(,\d+)+$`)
	var matchSingle = regexp.MustCompile(`^\d+$`)

	switch {
	// 1) Asterisk (*)
	case matchAll.MatchString(expression):
		expandedExpr = createStringRange(minRange, maxRange, 1)
	// 2) Hyphen-separated Range (min-max) (\d-\d)
	case matchRange.MatchString(expression):
		results := matchRange.FindStringSubmatch(expression)
		if len(results) != 3 {
			return "", errors.New("Unknown error occurred while parsing " + expression + " in " + position + " as range\n")
		}
		minRange, err := strconv.Atoi(results[1])
		if err != nil {
			return "", errors.New("Unknown error occurred while parsing " + results[1] + " in " + position + " as int\n" + err.Error() + "\n")
		}
		maxRange, err := strconv.Atoi(results[2])
		if err != nil {
			return "", errors.New("Unknown error occurred while parsing " + results[2] + " in " + position + " as int\n" + err.Error() + "\n")
		}
		if minRange > maxRange {
			return "", errors.New("Invalid range in " + position + ": " + expression + "\n")
		}

		expandedExpr = createStringRange(minRange, maxRange, 1)
	// 3) Asterisk with Step (*/step)
	case matchAllStep.MatchString(expression):
		results := matchAllStep.FindStringSubmatch(expression)
		if len(results) != 2 {
			return "", errors.New("Unknown error occurred while parsing " + expression + " in " + position + " as match all step\n")
		}
		step, err := strconv.Atoi(results[1])
		if err != nil {
			return "", errors.New("Unknown error occurred while parsing " + results[1] + " in " + position + " as int\n" + err.Error() + "\n")
		}
		expandedExpr = createStringRange(minRange, maxRange, step)
	// 4) List of Values (\d,\d,\d)
	case matchList.MatchString(expression):
		expandedExpr = strings.ReplaceAll(expression, ",", " ")
	// 5) Single value (\d)
	case matchSingle.MatchString(expression):
		expandedExpr = expression
	default:
		return "", errors.New("Invalid " + position + " expression or not yet recognisable\n")
	}
	return expandedExpr, nil
}
