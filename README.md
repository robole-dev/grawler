# Grawler

This app scrapes the website of the given url and finds all relative links and visit these urls. It can be used to
force the cache building of a website or to test the availability of existing pages.

## Install

### Go

If you have `go` installed (Go version >= **1.22.3**) you can use `go install` to install the application on your system.

```bash
go install github.com/robole-dev/grawler@latest
```

## Usage

```bash
grawler <url>
```

Example

```bash
grawler https://www.google.de
```

                 
## Features

- Search and find all URLs that exist on a given page (`grawler <url>`)
- Save informations of each request to an CSV file (`--output-filepath <path>` flag)
- Make calls in parallel (`--parallel <num>` flag)
- Limit the recursion depth (`--max-depth <num>` flag)
- Set a delay on each request (`--delay <num>` flag) 

More informations can be viewed via the help flag

```bash
grawler -h
```
