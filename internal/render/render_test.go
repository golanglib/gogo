package render

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golib/assert"
)

func Test_DefaultRender(t *testing.T) {
	it := assert.New(t)
	recorder := httptest.NewRecorder()

	// render with normal string
	s := "Hello, world!"

	render := NewDefaultRender(recorder)
	it.Equal(ContentTypeDefault, render.ContentType())

	err := render.Render(s)
	if it.Nil(err) {
		it.Equal(http.StatusOK, recorder.Code)
		it.Equal(ContentTypeDefault, recorder.Header().Get("Content-Type"))
		it.Equal(s, recorder.Body.String())
	}
}

func Benchmark_DefaultRender(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	recorder := httptest.NewRecorder()

	s := "Hello, world!"

	render := NewDefaultRender(recorder)
	for i := 0; i < b.N; i++ {
		recorder.Body.Reset()

		render.Render(s)
	}
}

func Test_DefaultRenderWithReader(t *testing.T) {
	it := assert.New(t)
	recorder := httptest.NewRecorder()

	s := "Hello, world!"

	recorder.Header().Add("Content-Length", fmt.Sprintf("%d", len(s)))
	recorder.Header().Add("Content-Type", "text/plain")

	// render with normal string
	reader := strings.NewReader(s)

	render := NewDefaultRender(recorder)
	it.Equal("text/plain", render.ContentType())

	err := render.Render(reader)
	if it.Nil(err) {
		it.Equal(http.StatusOK, recorder.Code)
		it.Equal(s, recorder.Body.String())
	}
}

func Benchmark_DefaultRenderWithReader(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	recorder := httptest.NewRecorder()

	reader := strings.NewReader("Hello, world!")

	render := NewDefaultRender(recorder)
	for i := 0; i < b.N; i++ {
		recorder.Body.Reset()

		render.Render(reader)
	}
}

func Test_DefaultRenderWithJson(t *testing.T) {
	it := assert.New(t)
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "application/json")

	// render with complex data type
	data := struct {
		Name string
		Age  int
	}{"gogo", 5}

	render := NewDefaultRender(recorder)

	err := render.Render(data)
	if it.Nil(err) {
		it.Equal(`{"Name":"gogo","Age":5}`, strings.TrimSpace(recorder.Body.String()))
	}
}

func Benchmark_DefaultRenderWithJson(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "application/json")

	data := struct {
		Name string
		Age  int
	}{"gogo", 5}

	render := NewDefaultRender(recorder)
	for i := 0; i < b.N; i++ {
		recorder.Body.Reset()

		render.Render(data)
	}
}

func Test_DefaultRenderWithXml(t *testing.T) {
	it := assert.New(t)
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "text/xml")

	// render with complex data type
	data := struct {
		XMLName xml.Name `xml:"recorder"`
		Success bool     `xml:"Result>Success"`
		Content string   `xml:"Result>Content"`
	}{
		Success: true,
		Content: "Hello, world!",
	}

	render := NewDefaultRender(recorder)

	err := render.Render(data)
	if it.Nil(err) {
		it.Equal("<recorder><Result><Success>true</Success><Content>Hello, world!</Content></Result></recorder>", recorder.Body.String())
	}
}

func Benchmark_DefaultRenderWithXml(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	recorder := httptest.NewRecorder()
	recorder.Header().Set("Content-Type", "text/xml")

	data := struct {
		XMLName xml.Name `xml:"recorder"`
		Success bool     `xml:"Result>Success"`
		Content string   `xml:"Result>Content"`
	}{
		Success: true,
		Content: "Hello, world!",
	}

	render := NewDefaultRender(recorder)
	for i := 0; i < b.N; i++ {
		recorder.Body.Reset()

		render.Render(data)
	}
}

func Test_DefaultRenderWithStringify(t *testing.T) {
	it := assert.New(t)
	recorder := httptest.NewRecorder()

	// render with complex data type
	data := struct {
		Name string
		Age  int
	}{"gogo", 5}

	render := NewDefaultRender(recorder)

	err := render.Render(data)
	if it.Nil(err) {
		it.Equal(`{gogo 5}`, recorder.Body.String())
	}
}

func Benchmark_DefaultRenderWithStringify(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	recorder := httptest.NewRecorder()

	data := struct {
		Name string
		Age  int
	}{"gogo", 5}

	render := NewDefaultRender(recorder)
	for i := 0; i < b.N; i++ {
		recorder.Body.Reset()

		render.Render(data)
	}
}
