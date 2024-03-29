<p align="center"><a href="#readme"><img src="https://gh.kaos.st/bibop.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/bibop"><img src="https://kaos.sh/r/bibop.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/l/bibop"><img src="https://kaos.sh/l/d306c6b094eae9f6ade7.svg" alt="Code Climate Maintainability" /></a>
  <a href="https://kaos.sh/b/bibop"><img src="https://kaos.sh/b/a03d5074-eea9-48a7-848c-dacbe7a9bf04.svg" alt="codebeat badge" /></a>
  <br/>
  <a href="COOKBOOK.md"><img src="https://gh.kaos.st/cookbook.svg"></a>
  <a href="https://kaos.sh/w/bibop/ci"><img src="https://kaos.sh/w/bibop/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/bibop/codeql"><img src="https://kaos.sh/w/bibop/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#license">License</a></p>

<br/>

`bibop` is a utility for testing command-line tools, packages and daemons. Initially, this utility was created for testing packages from [ESSENTIAL KAOS Public Repository](https://kaos.sh/kaos-repo).

Information about bibop recipe syntax you can find in our [cookbook](COOKBOOK.md).

### Usage demo

https://github.com/essentialkaos/bibop/assets/182020/c63dc147-fa44-40df-92e2-12f530c411af

### Installation

#### From source

To build the `bibop` from scratch, make sure you have a working Go 1.20+ workspace ([instructions](https://go.dev/doc/install)), then:

```
go install github.com/essentialkaos/bibop@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/bibop/latest).

To install the latest prebuilt version of bibop, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) bibop
```

### Docker support

Official webkaos images available on [GitHub Container Registry](https://kaos.sh/p/bibop) and [Docker Hub](http://kaos.sh/d/bibop). Install the latest version of Docker, then:

```bash
curl -fL# -o bibop-docker https://kaos.sh/bibop/bibop-docker
chmod +x bibop-docker
sudo mv bibop-docker /usr/bin/

bibop-docker your.recipe your-package.rpm
# or
bibop-docker --image ghcr.io/essentialkaos/bibop:centos7 your.recipe your-package.rpm
```

### Recipe syntax highlighting

* [Sublime Text 3/4](https://kaos.sh/blackhole-theme-sublime/bibop-recipe.sublime-syntax)
* [nano](https://kaos.sh/blackhole-theme-nano/bibop.nanorc)

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo bibop --completion=bash 1> /etc/bash_completion.d/bibop
```


ZSH:
```bash
sudo bibop --completion=zsh 1> /usr/share/zsh/site-functions/bibop
```


Fish:
```bash
sudo bibop --completion=fish 1> /usr/share/fish/vendor_completions.d/bibop.fish
```

### Man documentation

You can generate man page for bibop using next command:

```bash
bibop --generate-man | sudo gzip > /usr/share/man/man1/bibop.1.gz
```

### Usage

```
Usage: bibop {options} recipe

Options

  --dry-run, -D                Parse and validate recipe
  --extra, -X lines            Number of output lines for failed action (default: 10)
  --pause, -P duration         Pause between commands in seconds
  --list-packages, -L          List required packages
  --list-packages-flat, -L1    List required packages in one line (useful for scripts)
  --variables, -V              List recipe variables
  --barcode, -B                Show unique barcode for test (based on recipe and required packages)
  --time, -T                   Print execution time for every action
  --format, -f format          Output format (tap13|tap14|json|xml)
  --dir, -d dir                Path to working directory
  --path, -p path              Path to directory with binaries
  --error-dir, -e dir          Path to directory for errors data
  --tag, -t tag                One or more command tags to run
  --quiet, -q                  Quiet mode
  --ignore-packages, -ip       Do not check system for installed packages
  --no-cleanup, -nl            Disable deleting files created during tests
  --no-color, -nc              Disable colors in output
  --help, -h                   Show this help message
  --version, -v                Show version

Examples

  bibop app.recipe
  Run tests from app.recipe

  bibop app.recipe --quiet --error-dir bibop-errors
  Run tests from app.recipe in quiet mode and save errors data to bibop-errors directory

  bibop app.recipe --tag init,service
  Run tests from app.recipe and execute commands with tags init and service

  bibop app.recipe --extra
  Run tests from app.recipe and print the last 10 lines from command output if action was failed

  bibop app.recipe --extra=50
  Run tests from app.recipe and print the last 50 lines from command output if action was failed

  bibop app.recipe --format json 1> ~/results/app.json
  Run tests from app.recipe and save result in JSON format

  sudo dnf install $(bibop app.recipe -L1)
  Install all packages required for tests
```

### CI Status

| Branch | Status |
|------------|--------|
| `master` | [![CI](https://kaos.sh/w/bibop/ci.svg?branch=master)](https://kaos.sh/w/bibop/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/bibop/ci.svg?branch=develop)](https://kaos.sh/w/bibop/ci?query=branch:develop) |

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
