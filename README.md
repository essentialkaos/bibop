<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/bibop"><img src="https://kaos.sh/r/bibop.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/y/bibop"><img src="https://kaos.sh/y/892b9659ae8e4a5495bd9f8971bb31a2.svg" alt="Codacy badge" /></a>
  <a href="https://kaos.sh/w/bibop/ci-push"><img src="https://kaos.sh/w/bibop/ci-push.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/bibop/codeql"><img src="https://kaos.sh/w/bibop/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="COOKBOOK.md"><img src=".github/images/cookbook.svg"></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#upgrading">Upgrading</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`bibop` is a utility for testing command-line tools, packages and daemons. Initially, this utility was created for testing packages from [ESSENTIAL KAOS Public Repository](https://kaos.sh/kaos-repo).

Information about bibop recipe syntax you can find in our [cookbook](COOKBOOK.md).

### Usage demo

https://github.com/essentialkaos/bibop/assets/182020/c63dc147-fa44-40df-92e2-12f530c411af

### Installation

#### From source

To build the `bibop` from scratch, make sure you have a working Go [1.23+](https://github.com/essentialkaos/.github/blob/master/GO-VERSION-SUPPORT.md) workspace ([instructions](https://go.dev/doc/install)), then:

```
go install github.com/essentialkaos/bibop@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/bibop/latest).

To install the latest prebuilt version of bibop, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) bibop
```

### Upgrading

Since version `8.2.0` you can update `bibop` to the latest release using [self-update feature](https://github.com/essentialkaos/.github/blob/master/APPS-UPDATE.md):

```bash
bibop --update
```

This command will runs a self-update in interactive mode. If you want to run a quiet update (_no output_), use the following command:

```bash
bibop --update=quiet
```

### Docker support

Official webkaos images available on [GitHub Container Registry](https://kaos.sh/p/bibop) and [Docker Hub](http://kaos.sh/d/bibop). Install the latest version of Docker, then:

```bash
curl -fL# -o bibop-docker https://kaos.sh/bibop/bibop-docker
chmod +x bibop-docker
sudo mv bibop-docker /usr/bin/

bibop-docker your.recipe your-package.rpm
# or
bibop-docker --image ghcr.io/essentialkaos/bibop:ol9 your.recipe your-package.rpm
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

<img src=".github/images/usage.svg" />

### CI Status

| Branch | Status |
|------------|--------|
| `master` | [![CI](https://kaos.sh/w/bibop/ci-push.svg?branch=master)](https://kaos.sh/w/bibop/ci-push?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/bibop/ci-push.svg?branch=develop)](https://kaos.sh/w/bibop/ci-push?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/.github/blob/master/CONTRIBUTING.md).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://kaos.dev"><img src="https://raw.githubusercontent.com/essentialkaos/.github/refs/heads/master/images/ekgh.svg"/></a></p>
