# uploader

A small server to receive files over HTTP

### Install

**Binaries**

See [the latest release](https://github.com/jpillora/uploader/releases/latest) or download it with this one-liner: `curl i.jpillora.com/uploader | bash`

**Source**

``` sh
$ go get -v github.com/jpillora/uploader
```

### Usage

```
$ uploader
2015/06/24 00:46:08 listening on 3000...

...

$ curl -F file=@my-file.txt localhost:3000

...

2015/06/24 01:10:59 #0001 receiving my-file.txt
2015/06/24 01:10:59 #0001 received 53B
```

### Programmatic Usage

[![GoDoc](https://godoc.org/github.com/jpillora/uploader/lib?status.svg)](https://godoc.org/github.com/jpillora/uploader/lib)

#### MIT License

Copyright © 2015 Jaime Pillora &lt;dev@jpillora.com&gt;

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