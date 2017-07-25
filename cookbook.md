# `bibop` Cookbook

* [Recipe Syntax](#recipe-syntax)
  * [Comments](#comments)
  * [Global](#global)
    * [`dir`](#dir)
    * [`unsafe-paths`](#unsafe-paths)
    * [`command`](#command)
  * [Actions](#actions)
    * [`exit`](#exit)
    * [`expect`](#expect)
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