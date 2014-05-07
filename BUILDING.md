Building Cloud Foundry CLI
==========================

For developing on unix systems:

1. Run `./bin/build`
1. The binary will be built into the `./out` directory.

Optionally, you can use `bin/run` to compile and run the executable in one step.

For developing on windows with powershell.exe:
1. $Env:GODEP_PATH=C:\path\to\go-path\src\github.com\cloudfoundry\cli\Godeps\_workspace;
1. $Env:GOPATH = $Env:GODEP_PATH + ";" + "C:\path\to\go-path\"

Building Installers and Cross Compiling On Unix Systems
=======================================================
1. [Configure your go installation for cross compilation](https://stackoverflow.com/questions/12168873/cross-compile-go-on-osx)
1. Run `bin/build-all.sh`
1. Run `ci/scripts/build-installers`
1. Installers will all be in the `release` dir

How We Test Build and Release The CLI
=====================================

