<p align="center"><img src="https://gh.kaos.st/bibop-cookbook.svg"/></p>

* [Recipe Syntax](#recipe-syntax)
  * [Comments](#comments)
  * [Data types](#data-types)
  * [Global keywords](#global-keywords)
    * [`pkg`](#pkg)
    * [`unsafe-actions`](#unsafe-actions)
    * [`require-root`](#require-root)
    * [`fast-finish`](#fast-finish)
    * [`lock-workdir`](#lock-workdir)
    * [`unbuffer`](#unbuffer)
    * [`https-skip-verify`](#https-skip-verify)
    * [`delay`](#delay)
    * [`command`](#command)
  * [Variables](#variables)
  * [Actions](#actions)
    * [Common](#common)
      * [`exit`](#exit)
      * [`wait`](#wait)
    * [Input/Output](#inputoutput)
      * [`expect`](#expect)
      * [`print`](#print)
      * [`wait-output`](#wait-output)
      * [`output-match`](#output-match)
      * [`output-contains`](#output-contains)
      * [`output-trim`](#output-trim)
    * [Filesystem](#filesystem)
      * [`chdir`](#chdir)
      * [`mode`](#mode)
      * [`owner`](#owner)
      * [`exist`](#exist)
      * [`link`](#link)
      * [`readable`](#readable)
      * [`writable`](#writable)
      * [`executable`](#executable)
      * [`dir`](#dir)
      * [`empty`](#empty)
      * [`empty-dir`](#empty-dir)
      * [`checksum`](#checksum)
      * [`checksum-read`](#checksum-read)
      * [`file-contains`](#file-contains)
      * [`copy`](#copy)
      * [`move`](#move)
      * [`touch`](#touch)
      * [`mkdir`](#mkdir)
      * [`remove`](#remove)
      * [`chmod`](#chmod)
      * [`truncate`](#truncate)
      * [`cleanup`](#cleanup)
      * [`backup`](#backup)
      * [`backup-restore`](#backup-restore)
    * [System](#system)
      * [`process-works`](#process-works)
      * [`wait-pid`](#wait-pid)
      * [`wait-fs`](#wait-fs)
      * [`wait-connect`](#wait-connect)
      * [`connect`](#connect)
      * [`app`](#app)
      * [`signal`](#signal)
      * [`env`](#env)
      * [`env-set`](#env-set)
    * [Users/Groups](#usersgroups)
      * [`user-exist`](#user-exist)
      * [`user-id`](#user-id)
      * [`user-gid`](#user-gid)
      * [`user-shell`](#user-shell)
      * [`user-home`](#user-home)
      * [`group-exist`](#group-exist)
      * [`group-id`](#group-id)
    * [Services](#services)
      * [`service-present`](#service-present)
      * [`service-enabled`](#service-enabled)
      * [`service-works`](#service-works)
    * [HTTP](#http)
      * [`http-status`](#http-status)
      * [`http-header`](#http-header)
      * [`http-contains`](#http-contains)
      * [`http-json`](#http-json)
      * [`http-set-auth`](#http-set-auth)
      * [`http-set-header`](#http-set-header)
    * [Libraries](#libraries)
      * [`lib-loaded`](#lib-loaded)
      * [`lib-header`](#lib-header)
      * [`lib-config`](#lib-config)
      * [`lib-exist`](#lib-exist)
      * [`lib-linked`](#lib-linked)
      * [`lib-rpath`](#lib-rpath)
      * [`lib-soname`](#lib-soname)
      * [`lib-exported`](#lib-exported)
    * [Python](#python)
      * [`python2-package`](#python2-package)
      * [`python3-package`](#python3-package)
* [Examples](#examples)

## Recipe Syntax

### Comments

In `bibop` recipe all comments must have `#` prefix.

**Example:**

```yang
# Logs directory must be empty before tests
command "-" "Check environment"
  empty-dir /var/log/my-app
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

### Data types

Every action or variable can have values with next types:

* String
* Number (_integer or floating point_)
* Boolean (`true`/`false` _or_ `yes`/`no`)
* File mode (_integer with leading zero_)

▲ _To avoid splitting strings with whitespaces, value can be wrapped by double quotes:_

```yang
command "echo 'My message'" "Simple echo command"
  expect "My message"
  exit 0
```

▲ _If value contains double quotes, it must be wrapped by singular quotes:_

```yang
command "myapp john" "Check user"
  expect 'Unknown user "john"'
  exit 1
```

### Global keywords

#### `pkg`

One or more required packages for tests.

**Syntax:** `pkg <package-name>…`

**Arguments:**

* `package-name` - Package name (_String_)

**Example:**

```yang
pkg php nginx libhttp2 libhttp2-devel
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `unsafe-actions`

Allows doing unsafe actions (_like removing files outside of working directory_).

**Syntax:** `unsafe-actions <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`no` by default]

**Example:**

```yang
unsafe-actions yes
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `require-root`

Requires root privileges for the recipe.

If you use command syntax for executing the command as another user, this requirement will be enabled by automatically.

**Syntax:** `require-root <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`no` by default]

**Example:**

```yang
require-root yes
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `fast-finish`

If set to Yes, the test will be finished after the first failure.

**Syntax:** `fast-finish <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`no` by default]

**Example:**

```yang
fast-finish yes
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `lock-workdir`

If set to Yes, the current directory will be changed to working dir for every command.

**Syntax:** `lock-workdir <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`yes` by default]

**Example:**

```yang
lock-workdir no
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `unbuffer`

Disables I/O stream buffering.

**Syntax:** `unbuffer <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`no` by default]

**Example:**

```yang
unbuffer yes
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `https-skip-verify`

Disables TLS/SSL certificates verification.

**Syntax:** `https-skip-verify <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`no` by default]

**Example:**

```yang
https-skip-verify yes
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `delay`

Delay between commands.

**Syntax:** `delay <seconds>`

**Arguments:**

* `delay` - Delay in seconds (_Float_)

**Example:**

```yang
delay 1.5
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### `command`

Executes command. If you want to do some actions and checks without executing any binary (_"hollow" command_), you can use "-" (_minus_) as a command name.

You can execute the command as another user. For using this feature, you should define user name at the start of the command, e.g. `nobody:echo 'ABCD'`. This feature requires that `bibop` utility was executed with super user privileges (e.g. `root`).

Commands could be combined into groups. By default, every command has its own group. If you want to add a command to the group, use `+` as a prefix (e.g., `+command`). See the example below. If any command from the group fails, all the following commands in the group will be skipped.

You can define tag and execute the command with a tag on demand (using `-t` /` --tag` option of CLI). By default, all commands with tags are ignored.

Also, there is a special tag — `teardown`. If a command has this tag, this command will be executed even if `fast-finish` is set to true.

**Syntax:** `command:tag <cmd-line> [description]`

**Arguments:**

* `cmd-line` - Full command with all arguments
* `descriprion` - Command description [Optional]

**Examples:**

```yang
command "echo 'ABCD'" "Simple echo command"
  expect "ABCD"
  exit 0
```

```yang
command "USER=john ID=123 echo 'ABCD'" "Simple echo command with environment variables"
  expect "ABCD"
  exit 0
```

```yang
command "postgres:echo 'ABCD'" "Simple echo command as postgres user"
  expect "ABCD"
  exit 0
```

```yang
command "-" "Check configuration files (hollow command)"
  exist "/etc/myapp.conf"
  owner "/etc/myapp.conf" "root"
  mode "/etc/myapp.conf" 644
```

```yang
command:init "myapp initdb" "Init database"
  exist "/var/db/myapp.db"
```

```yang
command "-" "Replace configuration file"
  backup {redis_config}
  copy redis.conf {redis_config}

command "systemctl start {service_name}" "Start Redis service"
  wait {delay}
  service-works {service_name}
  connect tcp :6379

+command "systemctl status {service_name}" "Check status of Redis service"
  expect "active (running)"

+command "systemctl restart {service_name}" "Restart Redis service"
  wait {delay}
  service-works {service_name}
  connect tcp :6379

+command "redis-cli CONFIG GET logfile" "Check Redis Client"
  exit 0
  output-contains "/var/log/redis/redis.log"

+command "systemctl stop {service_name}" "Stop Redis service"
  wait {delay}
  !service-works {service_name}
  !connect tcp :6379

command "-" "Configuration file restore"
  backup-restore {redis_config}
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

### Variables

You can define global variables using keyword `var` and than use them in actions and commands.
Variables defined with `var` keyword is read-only and cannot be set by `*-read` actions.

Variables can contain other variables defined earlier.

Also, there are some run-time variables:

| Name | Description |
|------|-------------|
| `ENV:*` | Environment variable (_see example below_) |
| `DATE:*` | Current date with given [format](https://pkg.go.dev/github.com/essentialkaos/ek/v12/timeutil#Format) (_see example below_) |
| `WORKDIR` | Path to working directory |
| `TIMESTAMP` | Unix timestamp |
| `HOSTNAME` | Hostname |
| `IP` | Host IP |
| `OS` | OS name (_linux/darwin/freebsd…_) |
| `ARCH` | System architecture (_i386/i686/x86_64/arm…_) |
| `ARCH_NAME` | System architecture name (_386/686/amd64/arm…_) |
| `ARCH_BITS` | System architecture (_32/64_) |
| `LIBDIR` | Path to directory with shared libraries |
| `LIBDIR_LOCAL` | Path to local directory with shared libraries |
| `PYTHON2_VERSION` | Python 2.x version |
| `PYTHON2_SITELIB` | Path to directory where pure Python 2 modules are installed (`/usr/lib/python2.X/site-packages`) |
| `PYTHON2_SITEARCH` | Path where Python 2 extension modules (_e.g. C compiled_) are installed (`/usr/local/lib64/python2.X/site-packages`) |
| `PYTHON3_VERSION` | Python 3.x version |
| `PYTHON3_SITELIB` | Path to directory where pure Python 3 modules are installed (`/usr/lib/python3.X/site-packages`) |
| `PYTHON3_SITEARCH` | Path to directory where Python 3 extension modules (_e.g. C compiled_) are installed (`/usr/lib64/python3.X/site-packages`) |
| `PYTHON3_BINDING_SUFFIX` | Suffix for Python 3.x bindings |
| `ERLANG_BIN_DIR` | Path to directory with Erlang executables |

You can view and check all recipe variables using `-V`/`--variables` option:

```bash
bibop my-app.recipe --variables
```

**Examples:**

```yang
var service      nginx
var service_user nginx
var data_dir     /var/cache/{service}

command "service start {service}" "Starting service"
  service-works {service}
  exist {data_dir}
```

```yang
command "-" "Check shared library"
  exist {LIBDIR}/mylib.so
  mode {LIBDIR}/mylib.so 755
```

```yang
var app_name mysuppaapp

command "go build {app_name}.go" "Build application"
  exist {ENV:GOPATH}/bin/{app_name}_{DATE:%Y%m%d}
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

### Actions

Action do or check something after executing command.

▲ _All actions must have prefix (two spaces or horizontal tab) and follow command token._

#### Common

##### `exit`

Waits till command will be finished and then checks exit code.

**Syntax:** `exit <code> [max-wait]`

**Arguments:**

* `code` - Exit code (_Integer_)
* `timeout` - Max wait time in seconds (_Float_) [Optional | 60 seconds]

**Negative form:** Yes

**Examples:**

```yang
command "git clone git@github.com:user/repo.git" "Repository clone"
  exit 0
```

```yang
command "git clone git@github.com:user/repo.git" "Repository clone"
  exit 0 60
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `wait`

Makes pause before the next action.

**Syntax:** `wait <duration>`

**Arguments:**

* `duration` - Duration in seconds (_Float_)

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  wait 3.5
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### Input/Output

Be aware that the output store limited to 2 Mb of data for each stream (`stdout` _and_ `stderr`). So if command generates lots of output data, it better to use `expect` action to working with the output.

Also, `expect` checks output store every 25 milliseconds. It means that `expect` action can handle 80 Mb/s output stream without losing data. But if the command generates such an amount of output data it is not the right decision to test it with `bibop`.

##### `expect`

Expects some substring in command output.

**Syntax:** `expect <substr> [max-wait]`

**Arguments:**

* `substr` - Substring for search (_String_)
* `max-wait` - Max wait time in seconds (_Float_) [Optional | 5 seconds]

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  expect "ABCD"
```

```yang
command "echo 'ABCD'" "Simple echo command with 1 seconds timeout"
  expect "ABCD" 1
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `print`

Prints some data to `stdin`.

**Syntax:** `print <data>`

**Arguments:**

* `data` - Some text (_String_)

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  print "abcd"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `wait-output`

Waits till command prints any data.

**Syntax:** `wait-output <timeout>`

**Arguments:**

* `timeout` - Max wait time in seconds (_Float_)

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  wait-output 10.0
```

##### `output-match`

Checks output with given [regular expression](https://en.wikipedia.org/wiki/Regular_expression).

**Syntax:** `output-match <regexp>`

**Arguments:**

* `regexp` - Regexp pattern (_String_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-match "[A-Z]{4}"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `output-contains`

Checks if the output contains given substring.

**Syntax:** `output-contains <substr>`

**Arguments:**

* `substr` - Substring for search (_String_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-contains "BC"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `output-trim`

Trims output (_remove output data from store_).

**Syntax:** `output-trim`

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-trim
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### Filesystem

##### `chdir`

Changes current directory to given path.

▲ _Be aware that if you don't set `lock-workdir` to `no` for every command we will set current directory to working directory defined through CLI option._

**Syntax:** `chdir <path>`

**Arguments:**

* `path` - Path to file or directory (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  chdir /var/log
  exist secure.log
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `mode`

Checks file or directory mode bits.

**Syntax:** `mode <path> <mode>`

**Arguments:**

* `path` - Path to file or directory (_String_)
* `mode` - Mode (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  mode "/home/user/file.log" 644
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `owner`

Checks file or directory owner.

**Syntax:** `owner <path> <user>:<group>`

**Arguments:**

* `path` - Path to file or directory (_String_)
* `user` - User name (_String_)
* `group` - Group name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  owner "/home/john/file.log" "john"
  owner "/home/john/file1.log" ":sysadmins"
  owner "/home/john/file1.log" "john:sysadmins"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `exist`

Checks if file or directory exist.

**Syntax:** `exist <path>`

**Arguments:**

* `path` - Path to file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  exist "/home/john/file.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `link`

Checks if link points on given file or directory. Action follows all links until it finds a non-link object.

**Syntax:** `link <link> <target>`

**Arguments:**

* `link` - Path to link (_String_)
* `target` - Path to target file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  link "/etc/myapp.conf" "/srv/myapp/common/myapp.conf"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `readable`

Checks if file or directory is readable for some user.

**Syntax:** `readable <username> <path>`

**Arguments:**

* `username` - User name (_String_)
* `path` - Path to file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  readable john "/home/john/file.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `writable`

Checks if file or directory is writable for some user.

**Syntax:** `writable <username> <path>`

**Arguments:**

* `username` - User name (_String_)
* `path` - Path to file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  writable john "/home/john/file.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `executable`

Checks if file or directory is executable for some user.

**Syntax:** `executable <username> <path>`

**Arguments:**

* `username` - User name (_String_)
* `path` - Path to file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  executable john "/usr/bin/myapp"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `dir`

Checks if given target is directory.

**Syntax:** `dir <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  dir "/home/john/abcd"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `empty`

Checks if file is empty.

**Syntax:** `empty <path>`

**Arguments:**

* `path` - Path to file (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  empty "/home/john/file.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `empty-dir`

Checks if directory is empty.

**Syntax:** `empty-dir <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  empty-dir /var/log/my-app
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `checksum`

Checks file SHA256 checksum.

**Syntax:** `checksum <path> <hash>`

**Arguments:**

* `path` - Path to file (_String_)
* `hash` - SHA256 checksum (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  checksum "/home/john/file.log" "88D4266FD4E6338D13B845FCF289579D209C897823B9217DA3E161936F031589"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `checksum-read`

Checks file SHA256 checksum and writes it into the variable.

**Syntax:** `checksum-read <path> <variable>`

**Arguments:**

* `path` - Path to file (_String_)
* `variable` - Variable name (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  checksum-read "/home/john/file.log" log_checksum
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `file-contains`

Checks if file contains some substring.

**Syntax:** `file-contains <path> <substr>`

**Arguments:**

* `path` - Path to file (_String_)
* `substr` - Substring for search (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  file-contains "/home/john/file.log" "abcd"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `copy`

Makes copy of file or directory.

**Syntax:** `copy <source> <dest>`

**Arguments:**

* `source` - Path to source file or directory (_String_)
* `dest` - Path to destination (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  copy "/home/john/file.log" "/home/john/file2.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `move`

Moves file or directory.

**Syntax:** `move <source> <dest>`

**Arguments:**

* `source` - Path to source file or directory (_String_)
* `dest` - New destination (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  move "/home/john/file.log" "/home/john/file2.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `touch`

Changes file timestamps.

**Syntax:** `touch <path>`

**Arguments:**

* `path` - Path to file (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  touch "/home/john/file.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `mkdir`

Creates a directory.

**Syntax:** `mkdir <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  mkdir "/home/john/abcd"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `remove`

Removes file or directory.

▲ _Deleting files created due to the test in working dir is not required. `bibop` automatically deletes all files created due test process._

**Syntax:** `remove <target>`

**Arguments:**

* `target` - Path to file or directory (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  remove "/home/john/abcd"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `chmod`

Changes file mode bits.

**Syntax:** `chmod <target> <mode>`

**Arguments:**

* `target` - Path to file or directory (_String_)
* `mode` - Mode (_Integer_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  chmod "/home/john/abcd" 755
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `truncate`

Changes the size of the file to zero.

**Syntax:** `truncate <target>`

**Arguments:**

* `target` - Path to file (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Clear log file"
  truncate "/var/log/my-app/app.log"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `cleanup`

Removes all directories and files in the given directory.

**Syntax:** `cleanup <dir>`

**Arguments:**

* `dir` - Path to directory with data (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Remove app data"
  cleanup "/srv/myapp/data"
  cleanup "/srv/myapp/backups"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `backup`

Creates backup for the file.

**Syntax:** `backup <path>`

**Arguments:**

* `path` - Path to file (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Configure environment"
  backup /etc/myapp.conf
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `backup-restore`

Restores file from backup. `backup-restore` can be executed multiple times with different commands.

**Syntax:** `backup-restore <path>`

**Arguments:**

* `path` - Path to file (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Configure environment"
  backup /etc/myapp.conf
  backup-restore /etc/myapp.conf
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### System

##### `process-works`

Checks if process is works.

**Syntax:** `process-works <pid-file>`

**Arguments:**

* `pid-file` - Path to PID file (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  process-works /var/run/service.pid
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `wait-pid`

Waits for PID file.

**Syntax:** `wait-pid <pid-file> [timeout]`

**Arguments:**

* `pid-file` - Path to PID file (_String_)
* `timeout` - Timeout in seconds (_Float_) [Optional | 60 seconds]

**Negative form:** Yes

**Examples:**

```yang
command "-" "Check environment"
  wait-pid /var/run/service.pid
```

```yang
command "-" "Check environment"
  wait-pid /var/run/service.pid 90
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `wait-fs`

Waits for file/directory.

**Syntax:** `wait-fs <target> [timeout]`

**Arguments:**

* `target` - Path to file or directory (_String_)
* `timeout` - Timeout in seconds (_Float_) [Optional | 60 seconds]

**Negative form:** Yes

**Examples:**

```yang
command "service myapp start" "Starting MyApp"
  wait-fs /var/log/myapp.log
```

```yang
command "service myapp start" "Starting MyApp"
  wait-fs /var/log/myapp.log 180
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `wait-connect`

Waits for connection.

**Syntax:** `wait-connect <network> <address> [timeout]`

**Arguments:**

* `network` - Network name (`udp`, `tcp`, `ip`) (_String_)
* `address` - Network address (_String_)
* `timeout` - Timeout in seconds (_Float_) [Optional | 60 seconds]

**Negative form:** Yes

**Examples:**

```yang
command "service myapp start" "Starting MyApp server"
  wait-connect tcp :80
```

```yang
command "service myapp start" "Starting MyApp server"
  wait-connect tcp 127.0.0.1:80 15
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `connect`

Checks if it possible to connect to some address.

**Syntax:** `connect <network> <address>`

**Arguments:**

* `network` - Network name (`udp`, `tcp`, `ip`) (_String_)
* `address` - Network address (_String_)
* `timeout` - Timeout in seconds (_Float_) [Optional | 1 second]

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  connect tcp :6379
  connect tcp 192.0.2.1:http
  connect tcp 192.0.2.1:http 60
  connect udp [fe80::1%lo0]:53
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `app`

Checks if application binary is present in PATH.

**Syntax:** `app <name>`

**Arguments:**

* `name` - Application name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  app wget
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `signal`

Sends signal to process.

If `pid-file` not defined signal will be sent to current process.

**Syntax:** `signal <sig> [pid-file]`

**Arguments:**

* `sig` - Signal name or code (_String_ or _Number_)
* `pid-file` - Path to PID file (_String_) [Optional]

**Negative form:** No

**Examples:**

```yang
command "myapp --daemon" "Check my app"
  signal HUP
```

```yang
command "myapp --daemon" "Check my app"
  signal HUP /var/run/myapp.pid
```

```yang
command "myapp --daemon" "Check my app"
  signal 16
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `env`

Checks environment variable value.

**Syntax:** `env <name> <value>`

**Arguments:**

* `name` - Environment variable name (_String_)
* `value` - Environment variable value (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  env LANG en_US.UTF-8
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `env-set`

Sets environment variable.

**Syntax:** `env-set <name> <value>`

**Arguments:**

* `name` - Environment variable name (_String_)
* `value` - Environment variable value (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Prepare environment"
  env-set HTTP_PROXY "http://127.0.0.1:3300"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### Users/Groups

##### `user-exist`

Checks if user is exist on system.

**Syntax:** `user-exist <username>`

**Arguments:**

* `username` - User name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  user-exist nginx
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `user-id`

Checks if user has some ID (UID).

**Syntax:** `user-id <username> <id>`

**Arguments:**

* `username` - User name (_String_)
* `id` - UID (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  user-id nginx 345
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `user-gid`

Checks if user has some group ID (GID).

**Syntax:** `user-gid <username> <id>`

**Arguments:**

* `username` - User name (_String_)
* `id` - GID (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  user-gid nginx 994
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `user-group`

Checks if user is a part of some group.

**Syntax:** `user-group <username> <groupname>`

**Arguments:**

* `username` - User name (_String_)
* `groupname` - Group name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  user-group nginx nobody
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `user-shell`

Checks if user has some shell.

**Syntax:** `user-shell <username> <shell>`

**Arguments:**

* `username` - User name (_String_)
* `shell` - Shell binary (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  user-shell nginx /usr/sbin/nologin
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `user-home`

Checks if user has some home directory.

**Syntax:** `user-shell <username> <home-dir>`

**Arguments:**

* `username` - User name (_String_)
* `home-dir` - Directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  user-home nginx /usr/share/nginx
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `group-exist`

Checks if group is exist on system.

**Syntax:** `group-exist <groupname>`

**Arguments:**

* `groupname` - Group name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  group-exist nginx
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `group-id`

Checks if group has some ID (GID).

**Syntax:** `group-id <groupname> <id>`

**Arguments:**

* `groupname` - Group name (_String_)
* `id` - GID (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  group-id nginx 994
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### Services

##### `service-present`

Checks if service is present on the system.

**Syntax:** `service-present <service-name>`

**Arguments:**

* `service-name` - Service name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  service-present nginx
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `service-enabled`

Checks if auto start is enabled for service.

**Syntax:** `service-enabled <service-name>`

**Arguments:**

* `service-name` - Service name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  service-enabled nginx
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `service-works`

Checks if service is works.

**Syntax:** `service-works <service-name>`

**Arguments:**

* `service-name` - Service name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  service-works nginx
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### HTTP

##### `http-status`

Makes HTTP request and checks status code.

**Syntax:** `http-status <method> <url> <code> [payload]`

**Arguments:**

* `method` - Method (_String_)
* `url` - URL (_String_)
* `code` - Status code (_Integer_)
* `payload` - Request payload (_String_) [Optional]

**Negative form:** Yes

**Examples:**

```yang
command "-" "Make HTTP request"
  http-status GET "http://127.0.0.1:19999" 200
```

```yang
command "-" "Make HTTP request"
  http-status PUT "http://127.0.0.1:19999" 200 '{"id":103}'
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `http-header`

Makes HTTP request and checks response header value.

**Syntax:** `http-header <method> <url> <header-name> <header-value> [payload]`

**Arguments:**

* `method` - Method (_String_)
* `url` - URL (_String_)
* `header-name` - Header name (_String_)
* `header-value` - Header value (_String_)
* `payload` - Request payload (_String_) [Optional]

**Negative form:** Yes

**Examples:**

```yang
command "-" "Make HTTP request"
  http-header GET "http://127.0.0.1:19999" strict-transport-security "max-age=32140800"
```

```yang
command "-" "Make HTTP request"
  http-header PUT "http://127.0.0.1:19999" x-request-status "OK" '{"id":103}'
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `http-contains`

Makes HTTP request and checks response data for some substring.

**Syntax:** `http-contains <method> <url> <substr> [payload]`

**Arguments:**

* `method` - Method (_String_)
* `url` - URL (_String_)
* `substr` - Substring for search (_String_)
* `payload` - Request payload (_String_) [Optional]

**Negative form:** Yes

**Example:**

```yang
command "-" "Make HTTP request"
  http-contains GET "http://127.0.0.1:19999/info" "version: 1"
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `http-json`

Makes HTTP request and allows to check JSON value.

**Syntax:** `http-json <method> <url> <query> <value>`

**Arguments:**

* `method` - Method (_String_)
* `url` - URL (_String_)
* `query` - Query (_String_)
* `value` - Value for check (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Make HTTP request and check domain info"
  http-json GET https://dns.google/resolve?name=andy.one Question[0].name andy.one.
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `http-set-auth`

Sets username and password for Basic Auth.

_Notice that auth data will be set only for current command scope._

**Syntax:** `http-set-auth <username> <password>`

**Arguments:**

* `username` - User name (_String_)
* `password` - Password (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Make HTTP request with auth"
  http-set-auth admin test1234
  http-status GET "http://127.0.0.1:19999" 200

command "-" "Make HTTP request without auth"
  http-status GET "http://127.0.0.1:19999" 403
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `http-set-header`

Sets request header.

_Notice that header will be set only for current command scope._

**Syntax:** `http-set-header <header-name> <header-value>`

**Arguments:**

* `header-name` - Header name (_String_)
* `header-value` - Header value (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Make HTTP request"
  http-set-header Accept application/vnd.myapp.v3+json
  http-status GET "http://127.0.0.1:19999" 200
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### Libraries

##### `lib-loaded`

Checks if the library is loaded to the dynamic linker.

**Syntax:** `lib-loaded <glob>`

**Arguments:**

* `glob` - Shared library file name glob (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  lib-loaded libreadline.so.*
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `lib-header`

Checks if library header files are present on the system.

**Syntax:** `lib-header <lib>`

**Arguments:**

* `lib` - Library name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  lib-header expat
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `lib-config`

Checks if the library has a valid configuration file for pkg-config.

**Syntax:** `lib-config <lib>`

**Arguments:**

* `lib` - Library name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  lib-config expat
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `lib-exist`

Checks if library file is exist in libraries directory.

**Syntax:** `lib-exist <filename>`

**Arguments:**

* `filename` - Library file name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  lib-exist libc.so.1
  lib-exist libc_nonshared.a
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `lib-linked`

Checks if binary file has link with given library.

**Syntax:** `lib-linked <binary> <glob>`

**Arguments:**

* `binary` - Path to binary file (_String_)
* `glob` - Shared library file name glob (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check linking"
  lib-linked /usr/bin/myapp libc.so.*
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `lib-rpath`

Checks if binary file has [rpath](https://en.wikipedia.org/wiki/Rpath) field with given path.

**Syntax:** `lib-rpath <binary> <path>`

**Arguments:**

* `binary` - Path to binary file (_String_)
* `path` - Path to directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check rpath"
  lib-rpath /usr/bin/myapp /usr/share/myapp/lib64
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `lib-soname`

Checks if shared library file has [soname](https://en.wikipedia.org/wiki/Soname) field with given value.

**Syntax:** `lib-soname <lib> <glob>`

**Arguments:**

* `lib` - Path to shared library file (_String_)
* `glob` - Shared library soname glob (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check library soname"
  lib-soname /usr/lib64/libz.so.1.2.11 libz.so.1
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `lib-exported`

Checks if shared library exported a [symbol](https://www.gnu.org/software/gnulib/manual/html_node/Exported-Symbols-of-Shared-Libraries.html).

▲ _You can use script [`bibop-so-exported`](scripts/bibop-so-exported) for generating these tests._

**Syntax:** `lib-exported <lib> <symbol>`

**Arguments:**

* `lib` - Name or path to shared library file (_String_)
* `symbol` - Exported symbol (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check symbols exported by libcurl.so.4"
  lib-exported libcurl.so.4 curl_url_dup
  lib-exported libcurl.so.4 curl_url_get
  lib-exported libcurl.so.4 curl_url_set
  lib-exported libcurl.so.4 curl_version
  lib-exported libcurl.so.4 curl_version_info
```

```yang
command "-" "Check symbols exported by mylib.so"
  lib-exported /srv/myapp/libs/myapp-lib.so suppa_duppa_method
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

#### Python

##### `python2-package`

Checks if a given Python 2.x package could be loaded.

**Syntax:** `python-package <name>`

**Arguments:**

* `name` - Module name (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check Python package loading"
  python-package certifi
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

##### `python3-package`

Checks if a given Python 3.x package could be loaded.

**Syntax:** `python3-package <name>`

**Arguments:**

* `name` - Module name (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check Python 3 package loading"
  python3-package certifi
```

<a href="#"><img src="https://gh.kaos.st/separator.svg"/></a>

## Examples

```yang
# Bibop recipe for MkCryptPasswd

pkg mkcryptpasswd

fast-finish yes

var password MyPassword1234
var salt SALT1234
var salt_length 9

command "mkcryptpasswd -s" "Generate basic hash for password"
  expect "Please enter password:"
  print "{password}"
  expect "Hash: "
  exit 0

command "mkcryptpasswd -s -sa {salt}" "Generate hash for password with predefined salt"
  expect "Please enter password"
  print "{password}"
  expect "$6${salt}$lTxNu4.6r/j81sirgJ.s9ai8AA3tJdp67XBWLFiE10tIharVYtzRJ9eJ9YEtQsiLzVtg94GrXAYjf40pWEEg7/"
  exit 0

command "mkcryptpasswd -s -sa {salt} -1" "Generate MD5 hash for password with predefined salt"
  expect "Please enter password"
  print "{password}"
  expect "$1${salt}$zIPLJYODoLlesdP3bf95S1"
  exit 0

command "mkcryptpasswd -s -sa {salt} -5" "Generate SHA256 hash for password with predefined salt"
  expect "Please enter password"
  print "{password}"
  expect "$5${salt}$HOV.9Dkp4HSDzcfizNDG7x5ST4e74zcezvCJ8BWHuK8"
  exit 0

command "mkcryptpasswd -s -S" "Return error if password is too weak"
  expect "Please enter password"
  print "password"
  expect "Password is too weak: it is based on a dictionary word"
  print "password"
  expect "Password is too weak: it is based on a dictionary word"
  print "password"
  expect "Password is too weak: it is based on a dictionary word"
  !exit 0

command "mkcryptpasswd --abcd" "Return error about unsupported argument"
  expect "Error! You used unsupported argument --abcd. Please check command syntax."
  !exit 0
```

```yang
# Bibop recipe for webkaos (CentOS 7+)

pkg webkaos webkaos-debug webkaos-nginx webkaos-module-brotli webkaos-module-naxsi

require-root yes
unsafe-actions yes
https-skip-verify yes

var service_name webkaos
var user_name webkaos
var prefix_dir /etc/webkaos
var config {prefix_dir}/webkaos.conf
var binary /usr/sbin/webkaos
var modules_config {prefix_dir}/modules.conf
var modules_dir /usr/share/webkaos/modules
var pid_file /var/run/webkaos.pid
var ssl_dir {prefix_dir}/ssl
var dh_param {ssl_dir}/dhparam.pem
var log_dir /var/log/webkaos

var lua_ver 2.1.0-beta3
var lua_dir /usr/share/webkaos/luajit/share/luajit-{lua_ver}

command "-" "System environment validation"
  user-exist {user_name}
  group-exist {user_name}
  service-present {service_name}
  service-enabled {service_name}
  exist {config}
  exist {log_dir}

command "-" "Debug version"
  exist {binary}.debug
  service-present webkaos-debug

command "-" "Check linking with LuaJIT"
  lib-rpath {binary} /usr/share/webkaos/luajit/lib
  lib-linked {binary} "libluajit-5.1.so.*"

command "-" "Check Resty core and lrucache"
  dir {lua_dir}/ngx
  dir {lua_dir}/ngx/ssl
  dir {lua_dir}/resty
  dir {lua_dir}/resty/core
  dir {lua_dir}/resty/lrucache

  exist {lua_dir}/ngx/balancer.lua
  exist {lua_dir}/ngx/base64.lua
  exist {lua_dir}/ngx/errlog.lua
  exist {lua_dir}/ngx/ocsp.lua
  exist {lua_dir}/ngx/pipe.lua
  exist {lua_dir}/ngx/process.lua
  exist {lua_dir}/ngx/re.lua
  exist {lua_dir}/ngx/req.lua
  exist {lua_dir}/ngx/resp.lua
  exist {lua_dir}/ngx/semaphore.lua
  exist {lua_dir}/ngx/ssl.lua
  exist {lua_dir}/ngx/ssl/session.lua
  exist {lua_dir}/resty/core.lua
  exist {lua_dir}/resty/lrucache.lua
  exist {lua_dir}/resty/core/base.lua
  exist {lua_dir}/resty/core/base64.lua
  exist {lua_dir}/resty/core/ctx.lua
  exist {lua_dir}/resty/core/exit.lua
  exist {lua_dir}/resty/core/hash.lua
  exist {lua_dir}/resty/core/misc.lua
  exist {lua_dir}/resty/core/ndk.lua
  exist {lua_dir}/resty/core/phase.lua
  exist {lua_dir}/resty/core/regex.lua
  exist {lua_dir}/resty/core/request.lua
  exist {lua_dir}/resty/core/response.lua
  exist {lua_dir}/resty/core/shdict.lua
  exist {lua_dir}/resty/core/socket.lua
  exist {lua_dir}/resty/core/time.lua
  exist {lua_dir}/resty/core/uri.lua
  exist {lua_dir}/resty/core/utils.lua
  exist {lua_dir}/resty/core/var.lua
  exist {lua_dir}/resty/core/worker.lua
  exist {lua_dir}/resty/lrucache/pureffi.lua

  mode {lua_dir}/ngx/balancer.lua 644
  mode {lua_dir}/ngx/base64.lua 644
  mode {lua_dir}/ngx/errlog.lua 644
  mode {lua_dir}/ngx/ocsp.lua 644
  mode {lua_dir}/ngx/pipe.lua 644
  mode {lua_dir}/ngx/process.lua 644
  mode {lua_dir}/ngx/re.lua 644
  mode {lua_dir}/ngx/req.lua 644
  mode {lua_dir}/ngx/resp.lua 644
  mode {lua_dir}/ngx/semaphore.lua 644
  mode {lua_dir}/ngx/ssl.lua 644
  mode {lua_dir}/ngx/ssl/session.lua 644
  mode {lua_dir}/resty/core.lua 644
  mode {lua_dir}/resty/lrucache.lua 644
  mode {lua_dir}/resty/core/base.lua 644
  mode {lua_dir}/resty/core/base64.lua 644
  mode {lua_dir}/resty/core/ctx.lua 644
  mode {lua_dir}/resty/core/exit.lua 644
  mode {lua_dir}/resty/core/hash.lua 644
  mode {lua_dir}/resty/core/misc.lua 644
  mode {lua_dir}/resty/core/ndk.lua 644
  mode {lua_dir}/resty/core/phase.lua 644
  mode {lua_dir}/resty/core/regex.lua 644
  mode {lua_dir}/resty/core/request.lua 644
  mode {lua_dir}/resty/core/response.lua 644
  mode {lua_dir}/resty/core/shdict.lua 644
  mode {lua_dir}/resty/core/socket.lua 644
  mode {lua_dir}/resty/core/time.lua 644
  mode {lua_dir}/resty/core/uri.lua 644
  mode {lua_dir}/resty/core/utils.lua 644
  mode {lua_dir}/resty/core/var.lua 644
  mode {lua_dir}/resty/core/worker.lua 644
  mode {lua_dir}/resty/lrucache/pureffi.lua 644

command "-" "Nginx compatibility package"
  exist /etc/nginx
  exist /var/log/nginx
  exist /etc/nginx/nginx.conf
  exist /usr/sbin/nginx
  service-present nginx
  service-present nginx-debug

command "-" "Original configuration backup"
  backup {config}
  backup {modules_config}

command "-" "Add modules configuration"
  copy modules.conf {modules_config}

command "-" "Replace original configuration"
  copy webkaos.conf {config}

command "-" "Add test DH params file"
  copy dhparam.pem {dh_param}
  chmod {dh_param} 600

command "-" "Add self-signed certificate"
  copy ssl.key {ssl_dir}/ssl.key
  copy ssl.crt {ssl_dir}/ssl.crt
  chmod {ssl_dir}/ssl.key 600
  chmod {ssl_dir}/ssl.crt 600

command "-" "Clear old log files"
  touch {log_dir}/access.log
  touch {log_dir}/error.log
  truncate {log_dir}/access.log
  truncate {log_dir}/error.log

command "-" "Check brotli module"
  exist {prefix_dir}/xtra/brotli.conf
  exist {modules_dir}/ngx_http_brotli_filter_module.so
  exist {modules_dir}/ngx_http_brotli_static_module.so
  mode {prefix_dir}/xtra/brotli.conf 644
  mode {modules_dir}/ngx_http_brotli_filter_module.so 755
  mode {modules_dir}/ngx_http_brotli_static_module.so 755

command "-" "Check NAXSI module"
  exist {prefix_dir}/naxsi_core.rules
  exist {modules_dir}/ngx_http_naxsi_module.so
  mode {prefix_dir}/naxsi_core.rules 644
  mode {modules_dir}/ngx_http_naxsi_module.so 755

command "systemctl start {service_name}" "Start service"
  wait-pid {pid_file} 5
  service-works {service_name}

command "-" "Make HTTP requests"
  http-status GET "http://127.0.0.1" 200
  http-header GET "http://127.0.0.1" server webkaos
  http-contains GET "http://127.0.0.1/lua" "LUA MODULE WORKS"
  !empty {log_dir}/access.log
  truncate {log_dir}/access.log

command "-" "Make HTTPS requests"
  http-status GET "https://127.0.0.1" 200
  http-header GET "https://127.0.0.1" server webkaos
  http-contains GET "https://127.0.0.1/lua" "LUA MODULE WORKS"
  !empty {log_dir}/access.log
  truncate {log_dir}/access.log

command "-" "Save PID file checksum"
  checksum-read {pid_file} pid_sha

command "service {service_name} upgrade" "Binary upgrade"
  wait 3
  exist {pid_file}
  service-works {service_name}
  http-status GET "http://127.0.0.1" 200
  !checksum {pid_file} {pid_sha}

command "-" "Update configuration to broken one"
  copy broken.conf {config}

command "service {service_name} check" "Broken config check"
  !exit 0
  !empty {log_dir}/error.log

command "service {service_name} reload" "Broken config reload"
  !exit 0

command "service {service_name} restart" "Restart with broken config"
  !exit 0

command "-" "Restore working configuration"
  copy webkaos.conf {config}

command "service {service_name} reload" "Reload with original config"
  exit 0

command "systemctl stop {service_name}" "Stop service"
  !wait-pid {pid_file} 5
  !service-works {service_name}
  !connect tcp ":http"
  !exist {pid_file}

command "-" "Clear old log files"
  truncate {log_dir}/access.log
  truncate {log_dir}/error.log

command "systemctl start {service_name}-debug" "Start debug version of service"
  wait-pid {pid_file} 5
  service-works {service_name}-debug

+command "-" "Make HTTP requests"
  http-status GET "http://127.0.0.1" 200
  http-header GET "http://127.0.0.1" server webkaos
  http-contains GET "http://127.0.0.1/lua" "LUA MODULE WORKS"
  !empty {log_dir}/access.log
  truncate {log_dir}/access.log

+command "-" "Make HTTPS requests"
  http-status GET "https://127.0.0.1" 200
  http-header GET "https://127.0.0.1" server webkaos
  http-contains GET "https://127.0.0.1/lua" "LUA MODULE WORKS"
  !empty {log_dir}/access.log
  truncate {log_dir}/access.log

+command "systemctl stop {service_name}-debug" "Stop debug version of service"
  !wait-pid {pid_file} 5
  !service-works {service_name}-debug
  !connect tcp ":http"
  !exist {pid_file}

command:teardown "-" "Configuration restore"
  backup-restore {config}
  backup-restore {modules_config}

command:teardown "-" "DH param cleanup"
  remove {dh_param}

command:teardown "-" "Self-signed certificate cleanup"
  remove {ssl_dir}/ssl.key
  remove {ssl_dir}/ssl.crt
```

More working examples you can find in [our repository](https://github.com/essentialkaos/kaos-repo/tree/master/tests) with recipes for our RPM packages.
