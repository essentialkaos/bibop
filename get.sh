#!/bin/bash

################################################################################

NORM=0
BOLD=1
UNLN=4
RED=31
GREEN=32
YELLOW=33
BLUE=34
MAG=35
CYAN=36
GREY=37
DARK=90

CL_NORM="\e[0m"
CL_BOLD="\e[0;${BOLD};49m"
CL_UNLN="\e[0;${UNLN};49m"
CL_RED="\e[0;${RED};49m"
CL_GREEN="\e[0;${GREEN};49m"
CL_YELLOW="\e[0;${YELLOW};49m"
CL_BLUE="\e[0;${BLUE};49m"
CL_MAG="\e[0;${MAG};49m"
CL_CYAN="\e[0;${CYAN};49m"
CL_GREY="\e[0;${GREY};49m"
CL_DARK="\e[0;${DARK};49m"
CL_BL_RED="\e[1;${RED};49m"
CL_BL_GREEN="\e[1;${GREEN};49m"
CL_BL_YELLOW="\e[1;${YELLOW};49m"
CL_BL_BLUE="\e[1;${BLUE};49m"
CL_BL_MAG="\e[1;${MAG};49m"
CL_BL_CYAN="\e[1;${CYAN};49m"
CL_BL_GREY="\e[1;${GREY};49m"

################################################################################

APP="bibop"
APPS_REPO="https://apps.kaos.io/bibop/latest"

################################################################################

main() {
  local output="${1:-$APP}"
  
  local os=$(getOSName)
  local arch=$(getArch)

  if [[ -z "$os" ]] ; then
    error "$os is not supported" $RED
    exit 1
  fi

  if [[ -z "$arch" ]] ; then
    error "$arch is not supported" $RED
    exit 1
  fi

  download "$output" "$os" "$arch"
}

download() {
  local output="$1"
  local os="$2"
  local arch="$3"

  show "Downloading ${CL_CYAN}${APP}${CL_NORM} for ${CL_CYAN}${os}/${arch}${CL_NORM}..."

  curl -# "$APPS_REPO/$os/$arch/$APP" -o "$output"

  if [[ $? -ne 0 ]] ; then
    error "Can't download ${APP} binary" $RED
    exit 1
  fi

  chmod +x "$output" &> /dev/null

  show "${APP} binary successfully downloaded!" $GREEN
}

getOSName() {
  local os=$(uname -s)

  case $os in
    "Linux")  echo "linux" ;;
    "Darwin") echo "macosx" ;;
    *)        echo "" ;;
  esac
}

getArch() {
  local arch=$(uname -m)

  case $arch in
    "i386")   echo "$arch" ;;
    "x86_64") echo "$arch" ;;
    *)        echo "" ;;
  esac
}

################################################################################

show() {
  if [[ -n "$2" && -z "$no_colors" ]] ; then
    echo -e "\e[${2}m${1}\e[0m"
  else
    echo -e "$*"
  fi
}

error() {
  show "$@" 1>&2
}

################################################################################

main "$@"
