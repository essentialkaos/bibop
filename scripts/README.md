### `bibop` scripts

- `bibop-dep` — utility for installing/uninstalling recipe dependecnies
- `bibop-docker` — `bibop` docker/podman wrapper
- `bibop-libtest-gen` — utility for generating compilation tests for libraries
- `bibop-libtest-gen` — utility listing linked shared libraries
- `bibop-massive` — utility for mass package testing
- `bibop-multi-check` — utility for checking different versions of package
- `bibop-so-exported` — utility for generating exported symbols tests

#### `bibop-dep`

```
Usage: bibop-dep {options} {action} <recipe>

Actions

  install, i      Install packages
  reinstall, r    Reinstall packages
  uninstall, u    Uninstall packages

Options

  --enablerepo, -ER repo     Enable repository
  --disablerepo, -DR repo    Disable repository
  --yes, -y                  Automatically answer yes for all questions
  --no-color, -nc            Disable colors in output
  --help, -h                 Show this help message
  --version, -v              Show information about version

Examples

  bibop-dep install -ER kaos-testing myapp.recipe
  Install packages for myapp recipe with enabled kaos-testing repository

  bibop-dep install -ER kaos-testing,epel,cbr -y myapp.recipe
  Install packages for myapp recipe with enabled repositories

  bibop-dep uninstall
  Uninstall all packages installed by previous transaction
```

#### `bibop-libtest-gen`

```
Usage: bibop-libtest-gen {options} devel-package

Options

  --list-libs, -L      List all libs in package
  --output, -o name    Output source file (default: test.c)
  --lib, -l name       Lib name
  --no-color, -nc      Disable colors in output
  --help, -h           Show this help message
  --version, -v        Show information about version

Examples

  bibop-libtest-gen dirac-devel-1.0.2-15.el7.x86_64.rpm
  Generate test.c with all required headers for RPM package

  bibop-libtest-gen dirac-devel
  Generate test.c with all required headers for installed package
```

#### `bibop-linked`

```
Usage: bibop-linked {options} binary-file

Options

  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show information about version

Examples

  bibop-linked /usr/bin/curl
  List required shared libraries for binary file

  bibop-linked /usr/lib64/libcurl.so.4
  List required shared libraries for other library
```

#### `bibop-massive`

```
Usage: bibop-massive {options} recipe…

Options

  --validate, -V             Just validate recipes
  --recheck, -R              Run only failed checks
  --fresh, -F                Clean all caches before run
  --interrupt, -X            Interrupt checks after first error
  --barcode, -B              Print unique barcode for every test
  --enablerepo, -ER repo     Enable repository
  --disablerepo, -DR repo    Disable repository
  --error-dir, -e dir        Path to directory with tests errors
  --log, -l file             Path to log file
  --no-color, -nc            Disable colors in output
  --help, -h                 Show this help message
  --version, -v              Show information about version

Examples

  bibop-massive ~/tests/
  Run all tests in given directory

  bibop-massive -ER kaos-testing ~/tests/package1.recipe ~/tests/package2.recipe
  Run 2 tests with enabled repository 'kaos-testing' for installing packages

  bibop-massive -ER kaos-testing,epel,cbr ~/tests/package1.recipe
  Run verbose test with enabled repositories for installing packages
```

#### `bibop-multi-check`

```
Usage: bibop-multi-check {options} recipe package-list

Options

  --enablerepo, -ER repo     Enable repository
  --disablerepo, -DR repo    Disable repository
  --error-dir, -e dir        Path to directory with tests errors
  --log, -l file             Path to log file
  --no-color, -nc            Disable colors in output
  --help, -h                 Show this help message
  --version, -v              Show information about version

Examples

  bibop-multi-check app.recipe package.list
  Run tests for every package in list

  bibop-multi-check -ER kaos-testing,epel,cbr ~/tests/package1.recipe app.recipe package.list
  Run tests with enabled repositories for installing packages
```

#### `bibop-so-exported`

```
Usage: bibop-so-exported {options} package-name

Options

  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show information about version

Examples

  bibop-so-exported zlib
  Create tests for exported symbols for shared libraries in package zlib
```
