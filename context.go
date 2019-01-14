package gogo

import (
	"crypto"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Context defines context of a request
type Context struct {
	Response Responser
	Request  *http.Request
	Params   *AppParams
	Logger   Logger

	controller     string
	action         string
	settings       map[string]interface{}
	frozenSettings map[string]interface{}
	middlewares    []Middleware
	issuedAt       time.Time

	mux    sync.RWMutex
	cursor int8
}

// NewContext returns a *Context with *AppServer
func NewContext() *Context {
	return &Context{
		Response: NewResponse(nil),
		cursor:   -1,
	}
}

// Controller returns controller name of routed handler
func (c *Context) Controller() string {
	return c.controller
}

// Action returns action name of routed handler
func (c *Context) Action() string {
	return c.action
}

// Set binds a new value with key for the context
func (c *Context) Set(key string, value interface{}) {
	c.mux.Lock()

	if c.settings == nil {
		c.settings = make(map[string]interface{})
	}

	c.settings[key] = value

	c.mux.Unlock()
}

// Get returns a value of the key
func (c *Context) Get(key string) (v interface{}, ok bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.settings == nil {
		return
	}

	v, ok = c.settings[key]
	return
}

// MustGet returns a value of key or panic when key doesn't exist
func (c *Context) MustGet(key string) interface{} {
	v, ok := c.Get(key)
	if !ok {
		c.Logger.Panicf("Key %s doesn't exist", key)
	}

	return v
}

// SetFinal binds a value with key for the context and freezes it
func (c *Context) SetFinal(key string, value interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.frozenSettings == nil {
		c.frozenSettings = make(map[string]interface{})
	}

	if _, ok := c.frozenSettings[key]; ok {
		return ErrSettingsKey
	}

	c.frozenSettings[key] = value

	return nil
}

// GetFinal returns a frozen value of the key
func (c *Context) GetFinal(key string) (v interface{}, ok bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.frozenSettings == nil {
		return
	}

	v, ok = c.frozenSettings[key]
	return
}

// MustSetFinal likes SetFinal, but it panics if key is duplicated.
func (c *Context) MustSetFinal(key string, value interface{}) {
	err := c.SetFinal(key, value)
	if err != nil {
		c.Logger.Panicf("Freese key %s duplicated!", key)
	}
}

// MustGetFinal returns a frozen value of key or panic when it doesn't exist
func (c *Context) MustGetFinal(key string) interface{} {
	v, ok := c.GetFinal(key)
	if !ok {
		c.Logger.Panicf("Freeze key %s doesn't exist", key)
	}

	return v
}

// RequestID returns request id of the Context
func (c *Context) RequestID() string {
	if c.Logger == nil {
		return ""
	}

	return c.Logger.RequestID()
}

// RequestURI returns request raw uri of http.Request
func (c *Context) RequestURI() string {
	if c.Request.RequestURI != "" {
		return c.Request.RequestURI
	}

	uri, _ := url.QueryUnescape(c.Request.URL.EscapedPath())
	return uri
}

// HasRawHeader returns true if request sets its header with specified key
func (c *Context) HasRawHeader(key string) bool {
	for yek := range c.Request.Header {
		if key == yek {
			return true
		}
	}

	return false
}

// RawHeader returns request header value of specified key
func (c *Context) RawHeader(key string) string {
	for yek, val := range c.Request.Header {
		if key == yek {
			return strings.Join(val, ",")
		}
	}

	return ""
}

// HasHeader returns true if request sets its header for canonicaled specified key
func (c *Context) HasHeader(key string) bool {
	_, ok := c.Request.Header[http.CanonicalHeaderKey(key)]

	return ok
}

// Header returns request header value of canonicaled specified key
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// AddHeader adds response header with key/value pair
func (c *Context) AddHeader(key, value string) {
	c.Response.Header().Add(key, value)
}

// SetHeader sets response header with key/value pair
func (c *Context) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}

// SetStatus sets response status code
func (c *Context) SetStatus(code int) {
	c.Response.WriteHeader(code)
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(location string) {
	// always abort
	c.Abort()

	// adjust status code, default to 302
	status := c.Response.Status()
	switch status {
	case http.StatusMovedPermanently, http.StatusSeeOther, http.StatusTemporaryRedirect:
		// skip

	default:
		status = http.StatusFound
	}

	c.SetHeader("Location", location)

	http.Redirect(c.Response, c.Request, location, status)
}

// Return returns response with auto detected Content-Type
func (c *Context) Return(body ...interface{}) error {
	var render Render

	// auto detect response content-type from request header of accept
	if len(c.Response.Header().Get("Content-Type")) == 0 {
		accept := c.Request.Header.Get("Accept")
		for _, enc := range strings.Split(accept, ",") {
			switch {
			case strings.Contains(enc, RenderJsonContentType), strings.Contains(enc, RenderJsonPContentType):
				render = NewJsonRender(c.Response)

			case strings.Contains(enc, RednerXmlContentType):
				render = NewXmlRender(c.Response)

			}
			if render != nil {
				break
			}
		}
	}

	// third, use default render
	if render == nil {
		render = NewDefaultRender(c.Response)
	}

	if len(body) > 0 {
		return c.Render(render, body[0])
	}

	return c.Render(render, nil)
}

// HashedReturn returns response with ETag header calculated hash of response.Body dynamically
func (c *Context) HashedReturn(hasher crypto.Hash, body ...interface{}) error {
	if len(body) > 0 {
		return c.Render(NewHashRender(c.Response, hasher), body[0])
	}

	return c.Render(NewHashRender(c.Response, hasher), "")
}

// Text returns response with Content-Type: text/plain header
func (c *Context) Text(data interface{}) error {
	return c.Render(NewTextRender(c.Response), data)
}

// Json returns response with json codec and Content-Type: application/json header
func (c *Context) Json(data interface{}) error {
	return c.Render(NewJsonRender(c.Response), data)
}

// JsonP returns response with json codec and Content-Type: application/javascript header
func (c *Context) JsonP(callback string, data interface{}) error {
	return c.Render(NewJsonpRender(c.Response, callback), data)
}

// Xml returns response with xml codec and Content-Type: text/xml header
func (c *Context) Xml(data interface{}) error {
	return c.Render(NewXmlRender(c.Response), data)
}

// Render responses client with data rendered by Render
func (c *Context) Render(w Render, data interface{}) error {
	// always abort
	c.Abort()

	// currect response status code
	if coder, ok := data.(StatusCoder); ok {
		c.SetStatus(coder.StatusCode())
	}

	// currect response header of content-type
	c.SetHeader("Content-Type", w.ContentType())

	// flush header
	c.Response.FlushHeader()

	// shortcut for nil
	if data == nil {
		return nil
	}

	err := w.Render(data)
	if err != nil {
		c.Logger.Errorf("%T.Render(?): %v", w, err)
	}

	return err
}

// Next executes the remain middlewares in the chain.
// NOTE: It ONLY used in the middlewares!
func (c *Context) Next() {
	c.cursor++

	if c.cursor >= 0 && c.cursor < int8(len(c.middlewares)) {
		c.middlewares[c.cursor](c)
	} else {
		c.Logger.Warn("No more executer in the chain.")
	}
}

// Abort forces to stop call chain.
func (c *Context) Abort() {
	c.cursor = math.MaxInt8
}

// run starting request chan with new envs.
func (c *Context) run(w http.ResponseWriter, handler http.Handler, middlewares []Middleware) {
	// hijack Responser with correct http.ResponseWriter
	c.Response.Hijack(w)

	// reset internal
	c.settings = nil
	c.frozenSettings = nil
	c.middlewares = middlewares
	c.issuedAt = time.Now()
	c.cursor = -1

	// start chains
	c.Next()

	// invoke http.Handler if chained
	if c.cursor >= 0 && c.cursor < math.MaxInt8 {
		c.Abort()

		if handler != nil {
			handler.ServeHTTP(c.Response, c.Request)
		} else {
			// ghost, response status code only
			if !c.Response.HeaderFlushed() {
				w.WriteHeader(c.Response.Status())
			}
		}
	}
}
