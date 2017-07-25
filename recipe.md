* [Syntax](#syntax)
  * [Comments](#comments)
  * [Global](#global)
    * [`dir`](#dir)
    * [`unsafe-paths`](#unsafe-paths)
    * [`command`](#command)
* [Example](#example)

---

### Syntax

#### Comments

In bibop recipe all comments must have `#` prefix. 

Example:

```
# Set working directory to home dir
dir "/home/john"
```

#### Global

##### `dir`

Set working directory to given path

Syntax: `dir <path>`

*Arguments:*

* `path` - Absolute path to working directory

*Example:*

```
dir "/var/tmp"
```

<br/>

##### `unsafe-paths`

Allow to create and remove files and directories outside working directory.

*Syntax:* `unsafe-paths true`

*Example:*
```
unsafe-paths true
```

<br/>

##### `command`

Execute command.

*Syntax:* `command <cmd-line> [description]`

*Arguments:*

* `cmd-line` - Full command with all arguments
* `descriprion` - Command description (Optional)

*Example:*

```
command "echo 'ABCD' 'Simple echo command'"
  expect "ABCD" 
  exit 0
```

<br/>

### Example

