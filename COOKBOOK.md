<p align="center"><img src="https://gh.kaos.st/bibop-cookbook.svg"/></p>

* [Recipe Syntax](#recipe-syntax)
  * [Comments](#comments)
  * [Global](#global)
    * [`unsafe-actions`](#unsafe-actions)
    * [`require-root`](#require-root)
    * [`fast-finish`](#fast-finish)
    * [`lock-workdir`](#lock-workdir)
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
      * [`perms`](#perms)
      * [`owner`](#owner)
      * [`exist`](#exist)
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
      * [`backup`](#backup)
      * [`backup-restore`](#backup-restore)
    * [System](#system)
      * [`process-works`](#process-works)
      * [`wait-pid`](#wait-pid)
      * [`wait-fs`](#wait-fs)
      * [`connect`](#connect)
      * [`app`](#app)
      * [`env`](#env)
    * [Users/Groups](#users-groups)
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
    * [Libraries](#libraries)
      * [`lib-loaded`](#lib-loaded)
      * [`lib-header`](#lib-header)
      * [`lib-config`](#lib-config)
* [Examples](#examples)

## Recipe Syntax

### Comments

In bibop recipe all comments must have `#` prefix. 

**Example:**

```yang
# Set working directory to home dir
dir "/home/john"

```

<br/>

### Global

#### `unsafe-actions`

Allow doing unsafe actions (_like removing files outside of working directory_).

**Syntax:** `unsafe-actions <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`no` by default]

**Example:**

```yang
unsafe-actions yes

```

<br/>

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

<br/>

#### `fast-finish`

If set to Yes, the test will be finished after the first failure.

**Syntax:** `fast-finish <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`no` by default]

**Example:**

```yang
fast-finish yes

```

<br/>

#### `lock-workdir`

Changes current directory to working dir for every command.

**Syntax:** `lock-workdir <flag>`

**Arguments:**

* `flag` - Flag (_Boolean_) [`yes` by default]

**Example:**

```yang
lock-workdir no

```

<br/>

#### `command`

Execute command. If you want to do some actions and checks without executing any command or binary, you can use "-" (_minus_) as a command name.

You can execute the command as another user. For using this feature, you should define user name at the start of the command, e.g. `nobody:echo 'ABCD'`. This feature requires that bibop utility was executed with super user privileges (e.g. `root`).

You can define tag and execute the command with a tag on demand (using `-t` /` --tag` option of CLI). By default, all commands with tags are ignored.

**Syntax:** `command:tag <cmd-line> [description]`

**Arguments:**

* `cmd-line` - Full command with all arguments
* `descriprion` - Command description [Optional]

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  expect "ABCD" 
  exit 0

```

```yang
command "postgres:echo 'ABCD'" "Simple echo command as postgres user"
  expect "ABCD" 
  exit 0

```

```yang
command "-" "Check configuration files"
  exist "/etc/myapp.conf"
  owner "/etc/myapp.conf" "root"
  perms "/etc/myapp.conf" 644

```

```yang
command:init "my app initdb" "Init database"
  exist "/var/db/myapp.db"

```

<br/>

### Variables

You can define global variables using keyword `var` and than use them in actions and commands.
Variables defined with `var` keyword is read-only and cannot be set by `*-read` actions.

Also, there are some run-time variables:

* `WORKDIR` - Path to working directory
* `TIMESTAMP` - Unix timestamp
* `DATE` - Current date
* `HOSTNAME` - Hostname
* `IP` - Host IP
* `PYTHON_SITELIB`, `PYTHON2_SITELIB` - Path to Python 2 platform-independent library installation
* `PYTHON_SITEARCH`, `PYTHON2_SITEARCH` - Path to Python 2 platform-dependent library installation
* `PYTHON3_SITELIB` - Path to Python 3 platform-independent library installation
* `PYTHON3_SITEARCH` - Path to Python 3 platform-dependent library installation

**Example:**

```yang
dir "/tmp"

var service      nginx
var service_user nginx

command "service start {service}" "Starting service"
  service-works {service}

```

<br/>

### Actions

Action do or check something after executing command.

All action must have prefix (two spaces or horizontal tab) and follow command token.

#### Common

##### `exit`

Waits till command will be finished and then checks exit code.

**Syntax:** `exit <code> [max-wait]`

**Arguments:**

* `code` - Exit code (_Integer_)
* `timeout` - Max wait time in seconds (_Float_) [Optional | 60 seconds]

**Negative form:** Yes

**Example:**

```yang
command "git clone git@github.com:user/repo.git" "Repository clone"
  exit 0 60

```

<br/>

##### `wait`

Pause before next action.

**Syntax:** `wait <duration>`

**Arguments:**

* `duration` - Duration in seconds (_Float_)

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  wait 3.5

```

<br/>

#### Input/Output

Be aware that the output store limited to 2 Mb of data for each stream (stdout and stderr). So if command generates lots of output data, it better to use `expect` action to working with the output.

Also, `expect` checks output store every 25 milliseconds. It means that `expect` action can handle 80 Mb/s output stream without losing data. But if command generates such amount of output data it is not the right decision to test it with bibop.

##### `expect`

Expect some substring in command output.

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

<br/>

##### `print`

Print some data to `stdin`.

**Syntax:** `print <data>`

**Arguments:**

* `data` - Some text (_String_)

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  print "abcd"

```

<br/>

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

Check output with some Regexp.

**Syntax:** `output-match <regexp>`

**Arguments:**

* `regexp` - Regexp pattern (_String_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-match "[A-Z]{4}"

```

<br/>

##### `output-contains`

Check if output contains some substring.

**Syntax:** `output-contains <substr>`

**Arguments:**

* `substr` - Substring for search (_String_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-contains "BC"

```

<br/>

##### `output-trim`

Trim output (remove output data from store).

**Syntax:** `output-trim`

**Negative form:** No

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-trim

```

<br/>

#### Filesystem

##### `chdir`

Changes current directory to given path.

Be aware that if you don't set `lock-workdir` to `no` for every command we will set current dir to working dir defined in the recipe or through cli options.

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

<br/>

##### `perms`

Check file or directory permissions.

**Syntax:** `perms <path> <mode>`

**Arguments:**

* `path` - Path to file or directory (_String_)
* `mode` - Mode (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  perms "/home/user/file.log" 644

```

<br/>

##### `owner`

Check file or directory owner.

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

<br/>

##### `exist`

Check if file or directory exist.

**Syntax:** `exist <path>`

**Arguments:**

* `path` - Path to file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  exist "/home/john/file.log"

```

<br/>

##### `readable`

Check if file or directory is readable for some user.

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

<br/>

##### `writable`

Check if file or directory is writable for some user.

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

<br/>

##### `executable`

Check if file or directory is executable for some user.

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

<br/>

##### `dir`

Check if given target is directory.

**Syntax:** `dir <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  dir "/home/john/abcd"

```

<br/>

##### `empty`

Check if file is empty.

**Syntax:** `empty <path>`

**Arguments:**

* `path` - Path to file (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  empty "/home/john/file.log"

```

<br/>

##### `empty-dir`

Check if directory is empty.

**Syntax:** `empty-dir <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  empty-dir "/home/john/file.log"

```

<br/>

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

<br/>

##### `checksum-read`

Checks file SHA256 checksum and write it into the variable.

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

<br/>

##### `file-contains`

Check if file contains some substring.

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

<br/>

##### `copy`

Make copy of file or directory.

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

<br/>

##### `move`

Move file or directory.

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

<br/>

##### `touch`

Change file timestamps.

**Syntax:** `touch <path>`

**Arguments:**

* `path` - Path to file (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  touch "/home/john/file.log"

```

<br/>

##### `mkdir`

Create directory.

**Syntax:** `mkdir <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  mkdir "/home/john/abcd"

```

<br/>

##### `remove`

Remove file or directory.

**Syntax:** `remove <target>`

**Arguments:**

* `target` - Path to file or directory (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Check environment"
  remove "/home/john/abcd"

```

<br/>

##### `chmod`

Change file mode bits.

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

<br/>

##### `backup`

Create backup for the file.

**Syntax:** `backup <path>`

**Arguments:**

* `path` - Path to file (_String_)

**Negative form:** No

**Example:**

```yang
command "-" "Configure environment"
  backup /etc/myapp.conf

```

<br/>

##### `backup-restore`

Restore file from backup. `backup-restore` can be executed multiple times with different commands.

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

<br/>

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

<br/>

##### `wait-pid`

Waits for PID file.

**Syntax:** `wait-pid <pid-file> <timeout>`

**Arguments:**

* `pid-file` - Path to PID file (_String_)
* `timeout` - Timeout in seconds (_Float_) [Optional | 60 seconds]

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  wait-pid /var/run/service.pid 90

```

<br/>

##### `wait-fs`

Waits for file/directory.

**Syntax:** `wait-fs <target> <timeout>`

**Arguments:**

* `target` - Path to file or directory (_String_)
* `timeout` - Timeout in seconds (_Float_) [Optional | 60 seconds]

**Negative form:** Yes

**Example:**

```yang
command "service myapp start" "Starting MyApp"
  wait-fs /var/log/myapp.log 180

```

<br/>

##### `connect`

Checks if it possible to connect to some address.

**Syntax:** `connect <network> <address>`

**Arguments:**

* `network` - Network name (`udp`, `tcp`, `ip`) (_String_)
* `address` - Network address (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  connect tcp 127.0.0.1:6379
  connect tcp 192.0.2.1:http
  connect udp [fe80::1%lo0]:53

```

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

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

<br/>

#### HTTP

##### `http-status`

Makes HTTP request and checks status code.

**Syntax:** `http-status <method> <url> <code>`

**Arguments:**

* `method` - Method (_String_)
* `url` - URL (_String_)
* `code` - Status code (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  http-status GET "http://127.0.0.1:19999" 200

```

<br/>

##### `http-header`

Makes HTTP request and checks response header value.

**Syntax:** `http-header <method> <url> <code> <header-name> <header-value>`

**Arguments:**

* `method` - Method (_String_)
* `url` - URL (_String_)
* `header-name` - Header name (_String_)
* `header-value` - Header value (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  http-header GET "http://127.0.0.1:19999" strict-transport-security "max-age=32140800"

```

<br/>

##### `http-contains`

Makes HTTP request and checks response data for some substring.

**Syntax:** `http-contains <method> <url> <substr>`

**Arguments:**

* `method` - Method (_String_)
* `url` - URL (_String_)
* `substr` - Substring for search (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  http-contains GET "http://127.0.0.1:19999/info" "version: 1"

```

<br/>

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

<br/>

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

<br/>

##### `lib-config`

Checks if the library has a configuration file for pkg-config.

**Syntax:** `lib-config <lib>`

**Arguments:**

* `lib` - Library name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "-" "Check environment"
  lib-config expat

```

<br/>

## Examples

```yang
# Simple recipe for mkcryptpasswd utility

command "mkcryptpasswd" "Generate basic hash for password"
  expect "Please enter password"
  print "MyPassword1234"
  expect "Hash: "
  exit 0

command "mkcryptpasswd -sa SALT1234" "Generate hash for password with predefined salt"
  expect "Please enter password"
  print "MyPassword1234"
  wait 1
  output-contains "$6$SALT1234$lTxNu4.6r/j81sirgJ.s9ai8AA3tJdp67XBWLFiE10tIharVYtzRJ9eJ9YEtQsiLzVtg94GrXAYjf40pWEEg7/"
  exit 0

command "mkcryptpasswd -sa SALT1234 -1" "Generate MD5 hash for password with predefined salt"
  expect "Please enter password"
  print "MyPassword1234"
  wait 1
  output-contains "$1$SALT1234$zIPLJYODoLlesdP3bf95S1"
  exit 0

command "mkcryptpasswd -sa SALT1234 -5" "Generate SHA256 hash for password with predefined salt"
  expect "Please enter password"
  print "MyPassword1234"
  wait 1
  output-contains "$5$SALT1234$HOV.9Dkp4HSDzcfizNDG7x5ST4e74zcezvCJ8BWHuK8"
  exit 0

command "mkcryptpasswd -S" "Return error if password is too weak"
  expect "Please enter password"
  print "password"
  expect "Password is too weak: it is based on a dictionary word"
  print "password"
  wait 0.5
  print "password"
  wait 0.5
  exit 1

command "mkcryptpasswd --abcd" "Return error about unsupported argument"
  expect "Error! You used unsupported argument --abcd. Please check command syntax."
  exit 1


```

```yang
# Bibop recipe for webkaos

require-root yes
unsafe-actions yes

var service_name webkaos
var user_name webkaos
var config /etc/webkaos/webkaos.conf
var pid_file /var/run/webkaos.pid
var log_dir /var/log/webkaos

command "-" "System environment validation"
  user-exist {user_name}
  group-exist {user_name}
  service-present {service_name}
  service-enabled {service_name}
  exist {config}
  exist {log_dir}

command "systemctl start {service_name}" "Starting service"
  wait-pid {pid_file} 180
  service-works {service_name}
  http-status GET "http://127.0.0.1:80" 200
  http-header GET "http://127.0.0.1:80" server webkaos
  !empty {log_dir}/access.log
  checksum-read {pid_file} pid_sha

command "service {service_name} upgrade" "Upgrading binary"
  wait 3
  exist {pid_file}
  service-works {service_name}
  http-status GET "http://127.0.0.1:80" 200
  !checksum {pid_file} {pid_sha}

command "-" "Updating config to broken one"
  copy webkaos-broken.conf {config}

command "service {service_name} check" "Checking broken config"
  !exit 0
  !empty {log_dir}/error.log

command "service {service_name} reload" "Reloading broken config"
  !exit 0

command "service {service_name} restart" "Restarting with broken config"
  !exit 0

command "-" "Updating config to working one"
  copy webkaos-ok.conf {config}

command "service {service_name} reload" "Reloading working config"
  exit 0

command "systemctl stop {service_name}" "Stopping service"
  !wait-pid {pid_file} 5
  !service-works {service_name}
  !connect tcp ":http"
  !exist {pid_file}

```

More working examples you can find in [our repository](https://github.com/essentialkaos/kaos-repo/tree/master/tests) with recipes for our rpm packages.
