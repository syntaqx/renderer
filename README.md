# Renderer

[![Build Status](https://travis-ci.org/syntaqx/renderer.svg?branch=master)](https://travis-ci.org/syntaqx/renderer)
[![Coverage Status](https://coveralls.io/repos/github/syntaqx/renderer/badge.svg?branch=master)](https://coveralls.io/github/syntaqx/renderer?branch=master)
[![GoDoc](https://godoc.org/github.com/syntaqx/renderer?status.svg)](https://godoc.org/github.com/syntaqx/renderer)
![license](https://img.shields.io/github/license/syntaqx/renderer.svg)

[unrolled-render]: https://github.com/unrolled/render
[labstack-echo]: https://github.com/labstack/echo
[example-app]: ./example
[license]: ./LICENSE

Renderer provides a wrapper for [unrolled/render][unrolled-render]'s `Render`
instance that implements the [labstack/echo][labstack-echo] `Renderer` interface
for template rendering.

## Install

With a properly configured Go toolchain:

```sh
go get -u github.com/syntaqx/renderer
```

## Example

The minimum amount needed to set render as your apps template renderer is the
following:

```go
package main

import (
    "github.com/labstack/echo"
    "github.com/syntaqx/renderer"
    "github.com/unrolled/render"
)

func main() {
    e := echo.New()
    r := render.New()

    e.Renderer = renderer.Wrap(r)

    // ... the rest of your application

    e.Logger.Fatal(e.Start(":8080"))
}
```

However, that's not a particularly useful example, given there's no routes or
templates to reference. For a bit more elaborate example, check out our
[example][example-app].

## Contributing

Please feel free to submit issues and send pull requests!

## License

Renderer is licensed under the MIT license. See [LICENSE][] file for details.
