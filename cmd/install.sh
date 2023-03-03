#!/usr/bin/env sh
#
# This is an script for install inhere/kite-go
# More please see https://github.com/inhere/kite-go
#
# Local run:
#   bash ./cmd/install.sh
#
set -e

# check OS info
case $(uname | tr '[:upper:]' '[:lower:]') in
  linux*)
    CUR_OS_NAME=linux
    ;;
  darwin*)
    CUR_OS_NAME=darwin
    ;;
  msys*)
    CUR_OS_NAME=windows
    ;;
  *)
    CUR_OS_NAME=
    ;;
esac

if [ -z "$CUR_OS_NAME" ]; then
    echo "Unsupported OS";
    exit 0;
fi

# get arch name
arch_val=$(arch)
# shellcheck disable=SC2039
if [[ $arch_val =~ "x86_64" ]];then
    ARCH_NAME="amd64"
elif [[ $arch_val =~ "i386" ]];then
    ARCH_NAME="amd64"
elif [[ $arch_val =~ "aarch64" ]];then
    ARCH_NAME="arm"
else
    ARCH_NAME=
fi

if [ -z "$CUR_OS_NAME" ]; then
    echo "Unsupported OS";
    exit 0;
fi

# install bin name
#BIN_NAME=kitego
BIN_NAME=kite
# eg: kite-linux-amd64
BUILD_BIN_NAME="$BIN_NAME-$CUR_OS_NAME-$ARCH_NAME"

DOWNLOAD_URL="https://github.com/inhere/kite-go/releases/latest/download/$BUILD_BIN_NAME"
#echo $DOWNLOAD_URL;

INSTALL_DIR=/usr/local/bin
INSTALL_FILE=$INSTALL_DIR/$BIN_NAME

if [ -f "$INSTALL_FILE" ]; then
    echo "🙈 SKIP install, the kite bin file exists!"
    exit
fi

echo "🟢  Download kite bin file github release"
curl $DOWNLOAD_URL -L -o $INSTALL_FILE

echo "🟢  Initialize kite config"
# add exec perm
chmod a+x $INSTALL_FILE

set -x

# init user config
$BIN_NAME app init

echo "✅  Install kite successful"
$BIN_NAME --version

