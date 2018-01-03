<p align="center">
  <img alt="ConfWatch Logo" src="https://raw.githubusercontent.com/ConfWatch/confwatchd/master/static/img/avatar.png" height="140" />
  <h3 align="center"><a href="https://confwatch.ninja/" target="_blank">https://confwatch.ninja/</a></h3>
  <p align="center">Discover hacking conferences around the world.</p>
  <p align="center">
    <a href="https://github.com/ConfWatch/confwatchd/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/ConfWatch/confwatchd.svg?style=flat-square"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/ConfWatch/confwatchd"><img alt="Travis" src="https://img.shields.io/travis/ConfWatch/confwatchd/master.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/ConfWatch/confwatchd"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/ConfWatch/confwatchd?style=flat-square"></a>
  </p>
</p>

---

This repository contains the server source code for the ConfWatch project.

**VERY WORK IN PROGRESS, MUCH INCOMPLETE, DON'T USE**

Building
===

Make sure to have a working Go environment, that your `$GOPATH` is set correctly, and then:

    git clone https://github.com/ConfWatch/confwatchd $GOPATH/src/github.com/ConfWatch/confwatchd
    cd $GOPATH/src/github.com/ConfWatch/confwatchd
    make deps
    make

To run in a dev environment:

    ./confwatchd -config dev-config.json

Seeding the Database
===

    ./confwatchd -config config-file.json -seed /path/to/confwatch-data
