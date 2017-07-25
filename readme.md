<p align="center"><img src="https://gh.kaos.io/bibop.svg"/></p>

`bibop` is a utility for testing command-line tools.

* [Installation](#installation)
* [Usage](#usage)
* [Build Status](#build-status)
* [License](#license)

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

### Usage

```
Usage: bibop {options} recipe

Options

  --quiet, -q        Quiet mode
  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

```

### Build Status

| Branch | Status |
|------------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/bibop.svg?branch=master)](https://travis-ci.org/essentialkaos/bibop) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/bibop.svg?branch=develop)](https://travis-ci.org/essentialkaos/bibop) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.io/ekgh.svg"/></a></p>
