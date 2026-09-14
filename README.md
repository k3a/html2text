[![Go Reference](https://pkg.go.dev/badge/github.com/k3a/html2text.svg)](https://pkg.go.dev/github.com/k3a/html2text)
[![test](https://github.com/k3a/html2text/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/k3a/html2text/actions/workflows/test.yml)
[![coverage](https://raw.githubusercontent.com/k3a/html2text/badges/.badges/master/coverage.svg)](https://github.com/k3a/html2text/tree/badges)

# html2text

A simple Golang package to convert HTML to plain text, with no external dependencies.
It processes input character-by-character in a single pass, with occasional forward lookups (e.g., for HTML entities).

The conversion:

- Strips out the <head> section, along with <style> and <script> tags.
- Converts HTML tags to plain text and parses HTML entities into the characters they represent.
- Converts links into their href attribute, or to the `inner text <link>` format when the `WithLinksInnerText` option is enabled.

The primary use case is converting HTML emails into plain text.

The package includes a comprehensive test suite in html2text_test.go.
It follows semantic versioning, and no breaking changes are planned.

Feel free to submit a pull request with suggestions for improvement. However, please note that the library is now considered feature-complete and API‑stable. If you require more advanced functionality or need to regularly process malformed HTML input, you should consider using an alternative package listed below.

## Install
```bash
go get github.com/k3a/html2text
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/k3a/html2text"
)

func main() {
	html := `<html><head><title>Good</title></head><body><strong>clean</strong> text</body>`
	
	plain := html2text.HTML2Text(html)
			  
	fmt.Println(plain)
}

/*	Outputs:

	clean text
*/

```

Note: The output is plain text, do not embed the result back into HTML without escaping.

To see all features in action, take a look at [html2text_test.go](html2text_test.go).

## Alternatives
- https://github.com/jaytaylor/html2text (heavier, with more features, uses [golang.org/x/net/html parser](https://pkg.go.dev/golang.org/x/net/html))

## License

MIT

