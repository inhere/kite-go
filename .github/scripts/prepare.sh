#!/usr/bin/env bash

echo "Preparing for build"
pwd
env

ls -al ../

echo "Checkout deps repo:"
git clone --depth=1  https://github.com/gookit/goutil
git clone --depth=1  https://github.com/gookit/gcli
git clone --depth=1  https://github.com/gookit/greq
git clone --depth=1  https://github.com/gookit/gitw

ls -al ../
