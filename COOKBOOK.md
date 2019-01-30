# Bibop Cookbook

* [Recipe Syntax](#recipe-syntax)
  * [Comments](#comments)
  * [Global](#global)
    * [`dir`](#dir)
    * [`unsafe-paths`](#unsafe-paths)
    * [`command`](#command)
  * [Actions](#actions)
    * [Common](#common)
      * [`exit`](#exit)
      * [`wait`](#wait)
    * [Input/Output](#input-output)
      * [`expect`](#expect)
      * [`print`](#print)
      * [`output-match`](#output-match)
      * [`output-prefix`](#output-prefix)
      * [`output-suffix`](#output-suffix)
      * [`output-length`](#output-length)
      * [`output-contains`](#output-contains)
      * [`output-equal`](#output-equal)
      * [`output-trim`](#output-trim)
    * [Filesystem](#filesystem)
      * [`perms`](#perms)
      * [`owner`](#owner)
      * [`exist`](#exist)
      * [`readable`](#readable)
      * [`writable`](#writable)
      * [`directory`](#directory)
      * [`empty`](#empty)
      * [`empty-directory`](#empty-directory)
      * [`checksum`](#checksum)
      * [`file-contains`](#file-contains)
      * [`copy`](#copy)
      * [`move`](#move)
      * [`touch`](#touch)
      * [`mkdir`](#mkdir)
      * [`remove`](#remove)
      * [`chmod`](#chmod)
    * [Processes](#processes)
      * [`process-works`](#process-works)
    * [Users/Groups](#users-groups)
      * [`user-exist`](#user-exist)
      * [`user-id`](#user-id)
      * [`user-gid`](#user-gid)
      * [`user-shell`](#user-shell)
      * [`user-home`](#user-home)
      * [`group-exist`](#group-exist)
      * [`group-id`](#group-id)
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

#### `dir`

Set working directory to given path.

**Syntax:** `dir <path>`

**Arguments:**

* `path` - Absolute path to working directory

**Example:**

```yang
dir "/var/tmp"
```

<br/>

#### `unsafe-paths`

Allow to create/remove files and directories outside working directory.

**Syntax:** `unsafe-paths true`

**Example:**

```yang
unsafe-paths true
```

<br/>

#### `command`

Execute command. If you want to do some actions and checks without executing any command or binary, you can use "-" (_minus_) as a command name.

**Syntax:** `command <cmd-line> [description]`

**Arguments:**

* `cmd-line` - Full command with all arguments
* `descriprion` - Command description (Optional)

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  expect "ABCD" 
  exit 0
```

```yang
command "-" "Check configuration files"
  exist "/etc/myapp.conf"
  owner "/etc/myapp.conf" "root"
  perms "/etc/myapp.conf" 644
```

<br/>

### Actions

Action do or check something after executing command.

All action must have prefix (two spaces or horizontal tab) and follow command token.

#### Common

##### `exit`

Check command exit code.

**Syntax:** `exit <code> [max-wait]`

**Arguments:**

* `code` - Exit code (_Integer_)
* `max-wait` - Max wait time in seconds (Optional, 60 by default) (_Float_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  exit 1 30
```

<br/>

##### `wait`

Waits before next action.

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

##### `expect`

Expect some substring in command output.

**Syntax:** `expect <substr> [max-wait]`

**Arguments:**

* `substr` - Substring for search (_String_)
* `max-wait` - Max wait time in seconds (Optional, 5 by default) (_Float_)

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

##### `output-prefix`

Check output prefix.

**Syntax:** `output-prefix <substr>`

**Arguments:**

* `substr` - Substring for search (_String_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-prefix "AB"
```

<br/>

##### `output-suffix`

Check output suffix.

**Syntax:** `output-suffix <substr>`

**Arguments:**

* `substr` - Substring for search (_String_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-suffix "CD"
```

<br/>

##### `output-length`

Check output length.

**Syntax:** `output-length <length>`

**Arguments:**

* `length` - Output length (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-length 4
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

##### `output-equal`

Check if output is equal to given value.

**Syntax:** `output-equal <substr>`

**Arguments:**

* `substr` - Substring for search (_String_)

**Negative form:** Yes

**Example:**

```yang
command "echo 'ABCD'" "Simple echo command"
  output-equal "ABCD"
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

##### `perms`

Check file or directory permissions.

**Syntax:** `perms <path> <mode>`

**Arguments:**

* `path` - Path to file or directory (_String_)
* `mode` - Mode (_Integer_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  perms "/home/user/file.log" 644
```

<br/>

##### `owner`

Check file or directory owner.

**Syntax:** `owner <path> <owner-name>`

**Arguments:**

* `path` - Path to file or directory (_String_)
* `owner-name` - Owner name (_String_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  owner "/home/john/file.log" "john"
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
command "" "Check environment"
  exist "/home/john/file.log"
```

<br/>

##### `readable`

Check if file or directory is readable.

**Syntax:** `readable <path>`

**Arguments:**

* `path` - Path to file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  readable "/home/john/file.log"
```

<br/>

##### `writable`

Check if file or directory is writable.

**Syntax:** `writable <path>`

**Arguments:**

* `path` - Path to file or directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  writable "/home/john/file.log"
```

<br/>

##### `directory`

Check if given target is directory.

**Syntax:** `directory <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  directory "/home/john/abcd"
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
command "" "Check environment"
  empty "/home/john/file.log"
```

<br/>

##### `empty-directory`

Check if directory is empty.

**Syntax:** `empty-directory <path>`

**Arguments:**

* `path` - Path to directory (_String_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  empty-directory "/home/john/file.log"
```

<br/>

##### `checksum`

Check file SHA256 checksum.

**Syntax:** `checksum <path> <hash>`

**Arguments:**

* `path` - Path to file (_String_)
* `hash` - SHA256 checksum (_String_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  checksum "/home/john/file.log" "88D4266FD4E6338D13B845FCF289579D209C897823B9217DA3E161936F031589"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
  chmod "/home/john/abcd" 755
```

<br/>

#### Processes

##### `process-works`

Checks if process is works.

**Syntax:** `process-works <pid-file>`

**Arguments:**

* `pid-file` - Path to PID file (_String_)

**Negative form:** Yes

**Example:**

```yang
command "" "Check environment"
  process-works "/var/run/service.pid"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
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
command "" "Check environment"
  group-id nginx 994
```

<br/>

## Examples

```yang
dir "/tmp"

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
