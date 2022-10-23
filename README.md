# cronfmt

A multi-platform CLI tool that parses Vixie `cron` expressions and prints their extended format in stdout.

It emits a formatted table with the field name and the times as a space-separated list following it.

```bash
> cronfmt */15 0 1,15 * 1-5 /usr/bin/find

```

It supports special field value to further specify the execution time:

| Special field value |                         Description                         |                       Example                       |
|:-------------------:|:-----------------------------------------------------------:|:---------------------------------------------------:|
| Asterisk            | An asterisk represents every allowed value (first to last). | * (run every hour, month, etc.)                     |
| Range               | A range consists of two numbers separated by a hyphen.      | 0-5 (run from 0th to 5th hour, month, etc.)         |
| List                | A list is a set of numbers or ranges separated by commas.   | 0,1,2,3,4,5 (run from 0th to 5th hour, month, etc.) |
| Step                | A step is used in conjunction with ranges or asterisks.     | */2 (run every second hour, month, etc.)            |

## Usage

All dependencies are statically self-contained when compiled, and the resultant binary file weights a little more than 4MB.

Currently, just the `root` and the `version` commands are implemented leveraging the powerful [Cobra CLI framework](https://github.com/spf13/cobra)

```bash
cronfmt <cron expression>[mm hh dom mon dow command]

```

## Limitations

- The special `?` character is not supported.
- Named months and weekdays are not supported.
- The `@yearly` `@monthly` `@weekly` `@daily` `@hourly` and `@reboot` expressions are not supported.

## Good coding practices

### Variable naming conventions

It sticks to the Golang standard used by thousands of developers in the open-source community. Also documented on the [standard library](https://go.dev/doc/effective_go#names)

### Clear and concise code and comments

### Portability

### Testing
