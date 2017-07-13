<p align="center">
<img
    src="logo.png"
    width="260" height="80" border="0" alt="urlss">
<br>
<a href="https://travis-ci.org/schollz/urlss"><img src="https://img.shields.io/travis/schollz/urlss.svg?style=flat-square" alt="Build Status"></a>
<a href="http://gocover.io/github.com/schollz/urlss/lib"><img src="https://img.shields.io/badge/coverage-42%25-yellow.svg?style=flat-square" alt="Code Coverage"></a>
</p>

<p align="center">A URL shorterning service.</p>

This is a very simple URL shortening service. All URLs are saved into a Gzipped JSON backend, `urls.json.gz`. Try it out at [urls.schollz.com](https://urls.schollz.com).

Getting Started
===============

## Install

If you have Go installed, just do

    go get -u -v github.com/schollz/urlss

Otherwise, use the releases and [download urlss](https://github.com/schollz/urlss/releases/latest).


## Run

Once installed you can run the URL shortening service

    urlss -p 8009

and open a web browser to http://localhost:8009 to view it (or use a reverse proxy to attach it to a domain name).


## Development

Make sure you have `go-bindata` installed so that templates are updated:

    go get -u github.com/jteeuwen/go-bindata/...


Then use the following to build a new version of the server (with builtin templates):


    go-bindata.exe templates/... && go build && ./urlss


## License

MIT


