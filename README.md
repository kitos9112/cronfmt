# cronfmt

A multi-platform CLI tool that parses Vixie `cron` expressions and prints their extended format in stdout. It emits a pretty formatted table with the field name or section and the times as a space-separated list following it.

```bash
$ cronfmt "*/15" 0 1,15 "*" 1-5 "/usr/bin/find / -type f .terraform"
+---------------+------------------------------------+
| FIELD SECTION | EXTENDED FORMAT                    |
+---------------+------------------------------------+
| minute        | 0 15 30 45                         |
| hour          | 0                                  |
| day of month  | 1 15                               |
| month         | 1 2 3 4 5 6 7 8 9 10 11 12         |
| day of week   | 1 2 3 4 5                          |
| command       | /usr/bin/find / -type f .terraform |
+---------------+------------------------------------+
```

It supports especial-field value to further specify the execution time:

| Special field value |                         Description                         |                       Example                       |
|:-------------------:|:-----------------------------------------------------------:|:---------------------------------------------------:|
| Asterisk            | An asterisk represents every allowed value (first to last). | * (run every hour, month, etc.)                     |
| Range               | A range consists of two numbers separated by a hyphen.      | 0-5 (run from 0th to 5th hour, month, etc.)         |
| List                | A list is a set of numbers or ranges separated by commas.   | 0,1,2,3,4,5 (run from 0th to 5th hour, month, etc.) |
| Step                | A step is used in conjunction with asterisks.               | */2 (run every second hour, month, etc.)            |

## Usage

All dependencies are statically self-contained when compiled, and the resultant binary file weights a little more than 4.5MB.

Currently, just the `root` and the `version` commands are implemented leveraging the powerful [Cobra CLI framework](https://github.com/spf13/cobra)

## Development

Golang version 1.19+ is required to build the project. Although, a `Makefile` is provided to ease the build process. This process has been tested in Linux and macOS only, however, binaries for Windows x64 are also released.

```bash
make test
make binaries
```

After the successful execution of the above commands, the binary files will be available in the `binaries` directory. For example using a MacBook M1:

  ```bash
  $ binaries/osx_arm64/cronfmt --version
  cronfmt 0.0.1
  $ binaries/osx_arm64/cronfmt "*/5" 1,2,3,4,5 1 "*/4" 1-5 "/bin/cat /tmp/myHiddenFile"
  +-----------------+-----------------------------------+
  | CRON  | EXTENDED FORMAT                   |
  +-----------------+-----------------------------------+
  | minute          | 0 5 10 15 20 25 30 35 40 45 50 55 |
  | hour            | 1 2 3 4 5                         |
  | day of month    | 1                                 |
  | month           | 4 8 12                            |
  | day of week     | 1 2 3 4 5                         |
  | command         | /bin/cat /tmp/myHiddenFile        |
  +-----------------+-----------------------------------+
  ```

## Limitations

- The special `?` character is not supported.
- Named months and weekdays are not supported.
- Step values are only supported with asterisks and not ranges.
- The `@yearly` `@monthly` `@weekly` `@daily` `@hourly` and `@reboot` expressions are not supported.
