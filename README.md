# uploader

A small server to receive files over HTTP, with an embedded web UI:

<img width="300px" alt="screenshot" src="https://user-images.githubusercontent.com/633843/227771230-347164e2-61d6-4e00-a4a2-e0662a5d5dbf.png" />

### Features

* Simple
* Single binary; Tiny Docker image
* Classic write-only "drop box"; No file reads ever
* Renames files, not overwrite, by default

### Install

**Binaries**

See [the latest release](https://github.com/jpillora/uploader/releases/latest) or download it with this one-liner: `curl https://i.jpillora.com/uploader! | bash`

**Docker**

``` sh
$ docker run --rm -it -p 3000:3000 -v /tmp:/tmp ghcr.io/jpillora/uploader
```

**Source**

``` sh
$ go install -v github.com/jpillora/uploader@latest
```

### Usage

```sh
# server
$ uploader
2023/03/26 21:55:58 listening on 3000...
2023/03/26 21:55:58 saving files to: /tmp

...

# client (cli)
$ curl -F file=@my-file.txt localhost:3000
# client (browser)
# OPEN BROWSER, DRAG AND DROP FILE
...

# server
2015/06/24 01:10:59 #0001 receiving my-file.txt
2015/06/24 01:10:59 #0001 received 53B
```

### CLI Usage

```
  Usage: uploader [options]

  Options:
  --port, -p       listening port (default 3000)
  --no-log, -n     disable request logging
  --dir, -d        output directory (defaults to tmp)
  --overwrite, -o  duplicates are overwritten (auto-renames files by default)
  --auth, -a       require basic auth 'username:password'
  --version, -v    display version
  --help, -h       display help

  Version:
    0.0.0

  Read more:
    github.com/jpillora/uploader
```

### Programmatic Usage

[![GoDoc](https://godoc.org/github.com/jpillora/uploader/lib?status.svg)](https://godoc.org/github.com/jpillora/uploader/lib)

#### MIT License

Copyright © 2023 Jaime Pillora &lt;dev@jpillora.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
