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
grawler grawl <url>
```

Example

```bash
grawler grawl https://www.google.de
```

                 
## Features

All features can be read via the help flag

```bash
grawler -h
```

### Command `grawl`

Search and find all URLs that exist on a given page (`grawler grawl <url>`)

Options:
 
- Save informations of each request to an CSV file (`--output-filepath <path>` flag)
- Make calls in parallel (`--parallel <num>` flag)
- Limit the recursion depth (`--max-depth <num>` flag)
- Set a delay on each request (`--delay <num>` flag)
- Http Basic Auth (`--username` and `--password` flags. If you omit the password-flag you will get prompted.)

### Configuration

Precedence for configuration is first given to the flags set on the command-line, then to what's set in your configuration file.
                           
Grawler looks first for the command-line flag `--config` (path to the config file), then to the file `grawler.yaml`
in the current working directory and at least to the path `$HOME/.config/grawler/conf.yaml`.

You can **generate a config file** with default values with the `init` command.

A sample config files can be found here: [sample-conf.yaml](./sample-conf.yaml).
                                                                                                  
## Need to know

Currently we have some trouble to track the redirect http status codes.

More infos about that:

- <https://github.com/gocolly/colly/issues/298>
- <https://github.com/gocolly/colly/issues/212>
