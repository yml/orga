# orga
Utility to ease the development process of go wasm applications

## Usage

```
yml@kodi$ orga -h
orga is a cli to work with golang webassembly compilation target.

Usage:
  orga generate (<directory> [--force|--file=<name>])
  orga serve (<directory> [--addr=<host:port>])
  orga -h | --help
  orga -v | --version

Options:
  -h --help            Show this screen.
  --version            Show version.
  --force              Overwrite existing file.
  --file=<name>        Name of the wasm file [default: main.wasm].
  --addr=<host:port>   Name of the wasm file [default: :3000].

```

### Generate

```
orga generate /tmp/example/
```

`generate` sub command creates an `index.html`, `wasm_exec.js` and `go_js_wasm_exec`. The first 2
files are used to bootstrap your webassembly application in the browser.
By default `main.wasm` is expected, this file can be created as follow if you have a `main.go`

```
cd /tmp/example/
GOOS=js GOARCH=wasm ~/gowasm/bin/go build -o /tmp/example/main.wasm .
```

*note:* I often run the command above on change of any `.go` file

### Serve

`serve` starts simple http server.

```
orga serve /tmp/example/
```



