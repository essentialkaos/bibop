# Bibop Cookbook

* [Recipe Syntax](#recipe-syntax)
  * [Comments](#comments)
  * [Global](#global)
    * [`dir`](#dir)
    * [`unsafe-paths`](#unsafe-paths)
    * [`command`](#command)
  * [Actions](#actions)
    * [`exit`](#exit)
    * [`expect`](#expect)
    * [`output-match`](#output-match)
    * [`output-prefix`](#output-prefix)
    * [`output-suffix`](#output-suffix)
    * [`output-length`](#output-length)
    * [`output-contains`](#output-contains)
    * [`output-equal`](#output-equal)
    * [`output-trim`](#output-trim)
    * [`print`](#print)
    * [`wait`](#wait)
    * [`perms`](#perms)
    * [`owner`](#owner)
    * [`exist`](#exist)
    * [`readable`](#readable)
    * [`writable`](#writable)
    * [`directory`](#directory)
    * [`empty`](#empty)
    * [`empty-directory`](#empty-directory)
    * [`not-exist`](#not-exist)
    * [`not-readable`](#not-readable)
    * [`not-writable`](#not-writable)
    * [`not-directory`](#not-directory)
    * [`not-empty`](#not-empty)
    * [`not-empty-directory`](#not-empty-directory)
    * [`checksum`](#checksum)
    * [`file-contains`](#file-contains)
    * [`copy`](#copy)
    * [`move`](#move)
    * [`touch`](#touch)
    * [`mkdir`](#mkdir)
    * [`remove`](#remove)
    * [`chmod`](#chmod)
    * [`process-works`](#process-works)
* [Example](#example)

## Recipe Syntax

### Comments

In bibop recipe all comments must have `#` prefix. 

**Example:**

```
# Set working directory to home dir
dir "/home/john"
```

### Global

#### `dir`

Set working directory to given path.

**Syntax:** `dir <path>`

**Arguments:**

* `path` - Absolute path to working directory

**Example:**

```
dir "/var/tmp"
```

<br/>

#### `unsafe-paths`

Allow to create and remove files and directories outside working directory.

**Syntax:** `unsafe-paths true`

**Example:**
```
unsafe-paths true
```

<br/>

#### `command`

Execute command.

**Syntax:** `command <cmd-line> [description]`

**Arguments:**

* `cmd-line` - Full command with all arguments
* `descriprion` - Command description (Optional)

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  expect "ABCD" 
  exit 0
```

<br/>

### Actions

Action do or check something after executing command.

All action must have prefix (two spaces or horizontal tab) and follow command token.

#### `exit`

Check command exit code.

**Syntax:** `exit <code> [max-wait]`

**Arguments:**

* `code` - Exit code
* `max-wait` - Max wait time in seconds (Optional, 60 by default)

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  exit 1 30
```

<br/>

#### `expect`

Expect some substring in command output.

**Syntax:** `expect <substr> [max-wait]`

**Arguments:**

* `substr` - Substring for search
* `max-wait` - Max wait time in seconds (Optional, 5 by default)

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  expect "ABCD"
```

<br/>

#### `output-match`

Check output with some Regexp.

**Syntax:** `output-match <regexp>`

**Arguments:**

* `regexp` - Regexp pattern

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  output-match "[A-Z]{4}"
```

<br/>

#### `output-prefix`

Check output prefix.

**Syntax:** `output-prefix <substr>`

**Arguments:**

* `substr` - Substring for search

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  output-prefix "AB"
```

<br/>

#### `output-suffix`

Check output suffix.

**Syntax:** `output-suffix <substr>`

**Arguments:**

* `substr` - Substring for search

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  output-suffix "CD"
```

<br/>

#### `output-length`

Check output length.

**Syntax:** `output-length <length>`

**Arguments:**

* `length` - Output length

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  output-length 4
```

<br/>

#### `output-contains`

Check if output contains some substring.

**Syntax:** `output-contains <substr>`

**Arguments:**

* `substr` - Substring for search

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  output-contains "BC"
```

<br/>

#### `output-equal`

Check if output is equal to given value.

**Syntax:** `output-equal <substr>`

**Arguments:**

* `substr` - Substring for search

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  output-equal "ABCD"
```

<br/>

#### `output-trim`

Trim output.

**Syntax:** `output-trim`

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  output-trim
```

<br/>

#### `print`

Print some data to `stdin`.

**Syntax:** `print <data>`

**Arguments:**

* `data` - Some text

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  print "abcd"
```

<br/>

#### `wait`

Waits before next action.

**Syntax:** `wait <duration>`

**Arguments:**

* `duration` - Duration in seconds

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  wait 3.5
```

<br/>

#### `perms`

Checks file or directory permissions.

**Syntax:** `perms <target> <mode>`

**Arguments:**

* `target` - Path to file or directory
* `mode` - Mode

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  perms "/home/user/file.log" 644
```

<br/>

#### `owner`

Checks file or directory owner.

**Syntax:** `perms <target> <owner-name>`

**Arguments:**

* `target` - Path to file or directory
* `owner-name` - Owner name

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  owner "/home/john/file.log" "john"
```

<br/>

#### `exist`

Checks if file or directory exist.

**Syntax:** `exist <target>`

**Arguments:**

* `target` - Path to file or directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  exist "/home/john/file.log"
```

<br/>

#### `readable`

Checks if file or directory is readable.

**Syntax:** `readable <target>`

**Arguments:**

* `target` - Path to file or directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  readable "/home/john/file.log"
```

<br/>

#### `writable`

Checks if file or directory is writable.

**Syntax:** `writable <target>`

**Arguments:**

* `target` - Path to file or directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  writable "/home/john/file.log"
```

<br/>

#### `directory`

Checks if given target is directory.

**Syntax:** `directory <target>`

**Arguments:**

* `target` - Path to directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  directory "/home/john/abcd"
```

<br/>

#### `empty`

Checks if file is empty.

**Syntax:** `empty <target>`

**Arguments:**

* `target` - Path to file

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  empty "/home/john/file.log"
```

<br/>

#### `empty-directory`

Checks if directory is empty.

**Syntax:** `empty-directory <target>`

**Arguments:**

* `target` - Path to directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  empty-directory "/home/john/file.log"
```

<br/>

#### `not-exist`

Checks if file or directory doesn't exist.

**Syntax:** `exist <target>`

**Arguments:**

* `target` - Path to file or directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  not-exist "/home/john/file.log"
```

<br/>

#### `not-readable`

Checks if file or directory is not readable.

**Syntax:** `readable <target>`

**Arguments:**

* `target` - Path to file or directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  not-readable "/home/john/file.log"
```

<br/>

#### `not-writable`

Checks if file or directory is not writable.

**Syntax:** `writable <target>`

**Arguments:**

* `target` - Path to file or directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  not-writable "/home/john/file.log"
```

<br/>

#### `not-directory`

Checks if given target is not a directory.

**Syntax:** `directory <target>`

**Arguments:**

* `target` - Path to directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  not-directory "/home/john/abcd"
```

<br/>

#### `not-empty`

Checks if file is not empty.

**Syntax:** `empty <target>`

**Arguments:**

* `target` - Path to file

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  not-empty "/home/john/file.log"
```

<br/>

#### `not-empty-directory`

Checks if directory is not empty.

**Syntax:** `empty-directory <target>`

**Arguments:**

* `target` - Path to directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  not-empty-directory "/home/john/file.log"
```

<br/>

#### `checksum`

Checks file SHA256 checksum.

**Syntax:** `checksum <target> <hash>`

**Arguments:**

* `target` - Path to file
* `hash` - SHA256 checksum

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  checksum "/home/john/file.log" "88D4266FD4E6338D13B845FCF289579D209C897823B9217DA3E161936F031589"
```

<br/>

#### `file-contains`

Checks if file contains some substring.

**Syntax:** `file-contains <target> <substr>`

**Arguments:**

* `target` - Path to file
* `substr` - Substring for search

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  file-contains "/home/john/file.log" "abcd"
```

<br/>

#### `copy`

Make copy of file or directory.

**Syntax:** `copy <from> <to>`

**Arguments:**

* `from` - Path to source file or directory
* `to` - Path to copy

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  copy "/home/john/file.log" "/home/john/file2.log"
```

<br/>

#### `move`

Move file or directory.

**Syntax:** `move <from> <to>`

**Arguments:**

* `from` - Path to source file or directory
* `to` - New path

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  move "/home/john/file.log" "/home/john/file2.log"
```

<br/>

#### `touch`

Change file timestamps.

**Syntax:** `touch <target>`

**Arguments:**

* `target` - Path to file

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  touch "/home/john/file.log"
```

<br/>

#### `mkdir`

Create directory.

**Syntax:** `mkdir <target>`

**Arguments:**

* `target` - Path to directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  mkdir "/home/john/abcd"
```

<br/>

#### `remove`

Remove file or directory.

**Syntax:** `remove <target>`

**Arguments:**

* `target` - Path to file or directory

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  remove "/home/john/abcd"
```

<br/>

#### `chmod`

Change file mode bits.

**Syntax:** `chmod <target> <mode>`

**Arguments:**

* `target` - Path to file or directory
* `mode` - Mode

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  chmod "/home/john/abcd" 755
```

<br/>

#### `process-works`

Checks if process is works.

**Syntax:** `process-works <pid-file>`

**Arguments:**

* `pid-file` - Path to PID file

**Example:**

```
command "echo 'ABCD'" "Simple echo command"
  process-works "/var/run/service.pid"
```

<br/>

## Example

```
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