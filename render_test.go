package render

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-baa/baa"
	. "github.com/smartystreets/goconvey/convey"
)

var b = baa.New()

func TestRender1(t *testing.T) {
	Convey("render", t, func() {
		b.SetDI("render", New(Options{
			Baa:        b,
			Root:       "_fixture/templates/",
			Extensions: []string{".html"},
			FuncMap:    template.FuncMap{},
		}))

		Convey("normal render", func() {
			b.Get("/", func(c *baa.Context) {
				c.HTML(200, "index")
			})
			w := request("GET", "/")
			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("embed render", func() {
			b.Get("/i2", func(c *baa.Context) {
				body, err := c.Fetch("index2")
				So(err, ShouldBeNil)
				So(strings.Contains(string(body), "header"), ShouldBeTrue)
			})
			w := request("GET", "/i2")
			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("change file", func() {
			file := "_fixture/templates/index.html"
			body, err := ioutil.ReadFile(file)
			So(err, ShouldBeNil)
			err = ioutil.WriteFile(file, body, 0664)
			So(err, ShouldBeNil)
		})
	})
}

func request(method, uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, uri, nil)
	w := httptest.NewRecorder()
	b.ServeHTTP(w, req)
	return w
}
