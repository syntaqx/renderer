package renderer

import (
	"bytes"
	"html/template"
	"testing"
	"time"
)

func TestDateFormat(t *testing.T) {
	for i, this := range []struct {
		layout string
		value  interface{}
		expect interface{}
	}{
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
			if err == nil {
				t.Errorf("[%d] DateFormat didn't return an expected error, got %v", i, result)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] DateFormat failed: %s", i, err)
				continue
			}
			if result != this.expect {
				t.Errorf("[%d] DateFormat got %v but expected %v", i, result, this.expect)
			}
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
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		v, err := safeHTML(this.str)
		if err != nil {
			t.Fatalf("[%d] unexpected error in safeHTML: %s", i, err)
		}

		err = tmpl.Execute(buf, v)
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by safeHTML returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by safeHTML, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
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
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		v, err := safeURL(this.str)
		if err != nil {
			t.Fatalf("[%d] unexpected error in safeURL: %s", i, err)
		}

		err = tmpl.Execute(buf, v)
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by safeURL returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by safeURL, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}
