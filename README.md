<p align="center"><img src="https://gh.kaos.st/bibop.svg"/></p>
<p align="center">
<a href="https://travis-ci.org/essentialkaos/bibop"><img src="https://travis-ci.org/essentialkaos/bibop.svg?branch=master" /></a> 
<a href="https://goreportcard.com/report/github.com/essentialkaos/bibop"><img src="https://goreportcard.com/badge/github.com/essentialkaos/bibop" /></a>
<a href="https://codebeat.co/projects/github-com-essentialkaos-bibop-master"><img alt="codebeat badge" src="https://codebeat.co/badges/a03d5074-eea9-48a7-848c-dacbe7a9bf04" /></a> 
<a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg" /></a>
</p>
<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>


`bibop` is a utility for testing command-line tools. Information about bibop recipe syntax you can find in our [cookbook](COOKBOOK.md).

_**Note, that this is beta software, so it's entirely possible that there will be some significant bugs. Please report bugs so that we are aware of the issues.**_

### Usage demo

[![demo](https://gh.kaos.st/bibop-001.gif)](#usage-demo)

### Installation

#### From source

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `bibop` from scratch, make sure you have a working Go 1.10+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/bibop
```

If you want to update `bibop` to latest stable release, do:

```
go get -u github.com/essentialkaos/bibop
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/bibop/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) bibop 0.0.1
```

### Docker support

You can use Docker containers for testing your packages. Install latest version of Docker, then:

```bash
curl -o bibop-docker https://kaos.sh/bibop/bibop-docker
chmod +x bibop-docker
[sudo] mv bibop-docker /usr/bin/
bibop-docker --image essentialkaos/bibop:centos6 your.recipe your-package.rpm
```

Official Docker images with bibop:

- `essentialkaos/bibop:centos6`
- `essentialkaos/bibop:centos7`

### Recipe syntax highlighting

* [Sublime Text 3](https://github.com/essentialkaos/blackhole-theme-sublime/blob/master/bibop-recipe.sublime-syntax)
* [nano](https://github.com/essentialkaos/blackhole-theme-nano/blob/master/bibop.nanorc)

### Usage

```
Usage: bibop {options} recipe

Options

  --dir, -d dir          Path to working directory
  --error-dir, -e dir    Path to directory for errors data
  --tag, -t tag          Command tag
  --quiet, -q            Quiet mode
  --dry-run, -D          Parse and validate recipe
  --no-color, -nc        Disable colors in output
  --help, -h             Show this help message
  --version, -v          Show version

Examples

  bibop app.recipe
  Run tests from app.recipe

  bibop app.recipe --quiet --error-dir bibop-errors
  Run tests from app.recipe in quiet mode and save errors data to bibop-errors directory

  bibop app.recipe --tag init,service
  Run tests from app.recipe and execute commands with tags init and service

```

### Build Status

| Branch | Status |
|------------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/bibop.svg?branch=master)](https://travis-ci.org/essentialkaos/bibop) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/bibop.svg?branch=develop)](https://travis-ci.org/essentialkaos/bibop) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
