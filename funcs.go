package renderer

// text/template functions to add "essential" functionality to projects without
// overdoing it. The majority of these functions are extracted from spf13/hugo.
// Since the template.FuncMap wasn't explorted, we weren't able to include them
// without adding an unnecessary amount of overhead to the project.
// https://github.com/spf13/hugo/blob/master/tpl/template_funcs.go

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/url"
	"time"

	"github.com/spf13/cast"
)

// dateFormat formats a textual representation of a datetime string into the
// specified layout. If nil is provided as the textual datetime it will be
// replaced with time.Now.
//
// https://golang.org/pkg/time/#pkg-constants
func dateFormat(layout string, v interface{}) (string, error) {
	var t time.Time
	var err error

	if v == nil {
		t = time.Now()
	} else {
		t, err = cast.ToTimeE(v)
		if err != nil {
			return "", err
		}
	}

	return t.Format(layout), nil
}

// htmlEscape returns the given string with critical reserved HTML codes
// escaped, such that `&` becomes `&amp;` and so on. Only the `<`, `>`, `&`,
// `_`, `'`, and `"` characters are escaped.
//
// Keep in mind that, unless content is passed through `safeHTML`, output
// strings are escaped in default settings by the processor anyway.
func htmlEscape(in interface{}) (string, error) {
	conv, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}
	return html.EscapeString(conv), nil
}

// htmlUnescape returns the given string with HTML escape codes un-escaped. This
// un-escapes more codes than `htmlEscape` escapes, including `#` codes and
// pre-UTF8 escapes for accented characters. It defers completely to the native
// `html.UnescapeString` function, so it's functionally consistent with it.
//
// Remember to pass the output of this to `safeHTML` if fully unescaped
// characters are desired, otherwise the output will be escaped again as normal.
//
// https://golang.org/pkg/html/#EscapeString
func htmlUnescape(in interface{}) (string, error) {
	conv, err := cast.ToStringE(in)
	if err != nil {
		return "", err
	}
	return html.UnescapeString(conv), nil
}

// dictionary creates a map[string]interface{} from the given parameters by
// walking the parameters and treating them as key-value pairs. The number of
// parameters must be even.
func dictionary(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// querify encodes a set of key-value pairs into a "URL encoded" query string
// that can be appended to a URL after the `?` character.
func querify(params ...interface{}) (string, error) {
	qs := url.Values{}
	vals, err := dictionary(params...)
	if err != nil {
		return "", errors.New("querify keys must be strings")
	}

	for name, value := range vals {
		qs.Add(name, fmt.Sprintf("%v", value))
	}

	return qs.Encode(), nil
}

// safeHTML returns a given string as a html/template known-safe HTML document
// fragment, instructing template parsers to output its content verbatim.
//
// https://golang.org/pkg/html/template/#HTML
func safeHTML(a interface{}) (template.HTML, error) {
	s, err := cast.ToStringE(a)
	return template.HTML(s), err
}

// safeURL returns a given string as a html/template known-safe URL or URL
// substring, instructing template parsers to output its content verbatim.
//
// https://golang.org/pkg/html/template/#URL
func safeURL(a interface{}) (template.URL, error) {
	s, err := cast.ToStringE(a)
	return template.URL(s), err
}
