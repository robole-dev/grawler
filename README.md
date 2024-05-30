# Grawler

This app scrapes the website of the given url and finds all relative links and visit these urls. It can be used to
force the cache building of a website or to test the availability of existing pages.

## Install

### Go

If you have `go` installed (Go version >= **1.22.3**) you can use `go install` to install the application on your system.

```
go install github.com/robole-dev/grawler
```

## Usage

```
grawler <url>
```

Example

```
grawler https://www.google.de
```

All options can be viewed via this command

```
grawler -h
```
