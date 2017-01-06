package renderer

import (
	"bytes"
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateFormat(t *testing.T) {
	for i, this := range []struct {
		layout string
		value  interface{}
		expect interface{}
	}{
		{"Monday, Jan 2, 2006", nil, time.Now().Format("Monday, Jan 2, 2006")},
		{"Monday, Jan 2, 2006", "2015-01-21", "Wednesday, Jan 21, 2015"},
		{"Monday, Jan 2, 2006", time.Date(2015, time.January, 21, 0, 0, 0, 0, time.UTC), "Wednesday, Jan 21, 2015"},
		{"This isn't a date layout string", "2015-01-21", "This isn't a date layout string"},
		// The following test case gives either "Tuesday, Jan 20, 2015" or "Monday, Jan 19, 2015" depending on the local time zone
		{"Monday, Jan 2, 2006", 1421733600, time.Unix(1421733600, 0).Format("Monday, Jan 2, 2006")},
		{"Monday, Jan 2, 2006", 1421733600.123, false},
		{time.RFC3339, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "2016-03-03T04:05:00Z"},
		{time.RFC1123, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "Thu, 03 Mar 2016 04:05:00 UTC"},
		{time.RFC3339, "Thu, 03 Mar 2016 04:05:00 UTC", "2016-03-03T04:05:00Z"},
		{time.RFC1123, "2016-03-03T04:05:00Z", "Thu, 03 Mar 2016 04:05:00 UTC"},
	} {
		result, err := dateFormat(this.layout, this.value)
		if b, ok := this.expect.(bool); ok && !b {
			assert.Error(t, err, "[%d] DateFormat didn't return an expected error", i)
		} else {
			if assert.NoError(t, err, "[%d] DateFormat failed", i) {
				assert.Equal(t, this.expect, result, "[%d] DateFormat failed", i)
			}
		}
	}
}

func TestDictionary(t *testing.T) {
	for i, this := range []struct {
		v1            []interface{}
		expecterr     bool
		expectedValue map[string]interface{}
	}{
		{[]interface{}{"a", "b"}, false, map[string]interface{}{"a": "b"}},
		{[]interface{}{5, "b"}, true, nil},
		{[]interface{}{"a", 12, "b", []int{4}}, false, map[string]interface{}{"a": 12, "b": []int{4}}},
		{[]interface{}{"a", "b", "c"}, true, nil},
	} {
		r, err := dictionary(this.v1...)

		assert.False(t, (this.expecterr && err == nil) || (!this.expecterr && err != nil), "[%d] got an unexpected error", i)
		if !this.expecterr {
			assert.Equal(t, this.expectedValue, r, "[%d] got an unexpected value", i)
		}
	}
}

func TestSafeHTML(t *testing.T) {
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`<div></div>`, `{{ . }}`, `&lt;div&gt;&lt;/div&gt;`, `<div></div>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		assert.NoError(t, err, "[%d| unable to create new html template %q", i, this.tmplStr)

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		assert.NoError(t, err, "[%d] execute template with a raw string value returns unexpected error", i)
		assert.Equal(t, this.expectWithoutEscape, buf.String(), "[%d] execute template with a raw string value", i)

		buf.Reset()
		v, err := safeHTML(this.str)
		assert.NoError(t, err, "[%d] unexpected error in safeHTML", i)

		err = tmpl.Execute(buf, v)
		assert.NoError(t, err, "[%d] execute template with an escaped string value by safeHTML returns unexpected error", i)
		assert.Equal(t, this.expectWithEscape, buf.String(), "[%d] execute template with an escaped string value by safeHTML", i)
	}
}

func TestSafeURL(t *testing.T) {
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`irc://irc.freenode.net/#golang`, `<a href="{{ . }}">IRC</a>`, `<a href="#ZgotmplZ">IRC</a>`, `<a href="irc://irc.freenode.net/#golang">IRC</a>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		assert.NoError(t, err, "[%d] unable to create new html template %q", i, this.tmplStr)

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if assert.NoError(t, err, "[%d] execute template with a raw string value returns unexpected error", i) {
			assert.Equal(t, this.expectWithoutEscape, buf.String(), "[%d] execute template with a raw string value", i)
		}

		buf.Reset()
		v, err := safeURL(this.str)
		assert.NoError(t, err, "[%d] unexpected error in safeURL", i)

		err = tmpl.Execute(buf, v)
		if assert.NoError(t, err, "[%d] execute template with an escaped string value by safeURL", i) {
			assert.Equal(t, this.expectWithEscape, buf.String(), "[%d] execute template with an escaped string value by safeURL", i)
		}
	}
}
