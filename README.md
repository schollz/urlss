# urlss

*A URL Shortening Service.*

Install and run with

    go get -u -v github.com/schollz/urlss
    urlss -p 8009

All URLs are saved into a Gzipped JSON backend, `urls.json.gz`.

## Development

Make sure you have `go-bindata` installed so that templates are updated:

    go get -u github.com/jteeuwen/go-bindata/...


Then use the following to build a new version of the server (with builtin templates):


    go-bindata.exe templates/... && go build && ./urlss