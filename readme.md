<p align="center"><img src="https://gh.kaos.io/bibop.svg"/></p>
<p align="center">
<a href="https://travis-ci.org/essentialkaos/bibop"><img src="https://travis-ci.org/essentialkaos/bibop.svg?branch=master" /></a> 
<a href="https://goreportcard.com/report/github.com/essentialkaos/bibop"><img src="https://goreportcard.com/badge/github.com/essentialkaos/bibop" /></a> 
<a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.io/ekol.svg" /></a>
</p>
<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#from-source">From source</a> • <a href="#prebuilt-binaries">Prebuilt binaries</a> • <a href="#using-getsh">Using `get.sh`</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>


`bibop` is a utility for testing command-line tools. Information about bibop recipe syntax you can find in our [cookbook](cookbook.md).

_**Note, that this is beta software, so it's entirely possible that there will be some significant bugs. Please report bugs so that we are aware of the issues.**_

* [Usage demo](#usage-demo)
* [Installation](#installation)
* [Usage](#usage)
* [Build Status](#build-status)
* [License](#license)

### Usage demo

[![demo](https://gh.kaos.io/bibop-001.gif)](#usage-demo)

### Installation

#### From source

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `bibop` from scratch, make sure you have a working Go 1.5+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/bibop
```

If you want to update `bibop` to latest stable release, do:

```
go get -u github.com/essentialkaos/bibop
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.io/bibop/latest).

#### Using `get.sh`

If you want to use `bibop` in your CI environment, you can simply install latest prebuilt `bibop` binary using `get.sh` script:

```
bash <(curl -Ls https://raw.githubusercontent.com/essentialkaos/bibop/master/get.sh) && export "PATH=$PATH:$(pwd)"
```

Also you can configure install path:

```
bash <(curl -Ls https://raw.githubusercontent.com/essentialkaos/bibop/master/get.sh) /usr/local/bin/bibop
```

### Usage

```
Usage: bibop {options} recipe

Options

  --log, -l          Path to log file for verbose info about errors
  --quiet, -q        Quiet mode
  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

Examples

  bibop application.recipe
  Run tests from application.recipe

  bibop application.recipe --quiet --log errors.log 
  Run tests from application.recipe in quiet mode and log errors to errors.log

```

### Build Status

| Branch | Status |
|------------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/bibop.svg?branch=master)](https://travis-ci.org/essentialkaos/bibop) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/bibop.svg?branch=develop)](https://travis-ci.org/essentialkaos/bibop) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.io/ekgh.svg"/></a></p>
