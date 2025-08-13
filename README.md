<div align="center">
    <img src="icon.svg">
    <h1>Grawler</h1>
    <blockquote>
        <p dir="auto">Web crawler that discovers and visits relative links on a website.</p>
    </blockquote>
</div>

Grawler is a web crawler written in go. It scrapes the website of the given url and finds all relative 
links and visit these urls. Initially this application was developed to build up the cache of a page and to 
check the availability of existing pages.

## Install
                                         
### Binary

Download and use a binary suitable for your system from the [prebuild releases](https://github.com/robole-dev/grawler/releases). 

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

All features can be read via the help flag

```bash
grawler -h
```
       
More examples below.


## Examples

### Crawl website

```bash
grawler grawl https://books.toscrape.com
```

- crawles the given url
- search for anchor tags href elements (`<a href="...">`) and crawls these urls too

### Save result to a CSV-file

```bash
grawler grawl https://books.toscrape.com -o out.csv 
```

### Allow parallel requests
          
Set to 8 requests in parallel

```bash
grawler grawl https://books.toscrape.com -l 8 
```
                                    
### Limit the search depth 

Limit to a search recursion depth to 2 

```bash
grawler grawl https://books.toscrape.com --max-depth 2 
```

### Set a delay for each request
               
Set a delay of 500 milliseconds

```bash
grawler grawl https://books.toscrape.com --delay 500 
```

### Request a page with http basic auth

To rrequest a website that uses/requires a http basic auth you can set the username and password as flags 

```bash
grawler grawl https://books.toscrape.com --username user_xy --password mypassword 
```   

Optionally you can ommit the password. Then you will be asked to enter the password when you start grawling

```bash
grawler grawl https://books.toscrape.com --username user_xy
No config file found.
Grawling https://books.toscrape.com
✔ Password: █
```         

### Add allowed domains

By default, only the domain of the start url is allowed to be crawled. All other urls from other domains are being skipped.
You can allow more domains with the `-a` flag

```bash
grawler grawl https://quotes.toscrape.com -a example.com  
```   

You can also add multiple domains

```bash
grawler grawl https://quotes.toscrape.com -a example.com -a google.de  
```   
                                                                     
### Skip/Disallow urls

You can define one or multiple regular expression to skip urls when they match this/these expressions.

Here we skip all urls starting with `https://books.toscrape.com/catalogue/category/books/` with a max depth of 2: 

```bash
grawler grawl https://books.toscrape.com --disallowed-url-filters "^https://books.toscrape.com/catalogue/category/books/.*" --max-depth 2
```   

Here we skip all urls which contain the word `category` and the word `art`:

```bash
grawler grawl https://books.toscrape.com --disallowed-url-filters "category" --disallowed-url-filters "art"
```   


## Configuration

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
