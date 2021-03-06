package gogo

import (
	"net/http"
	"net/http/httputil"

	"github.com/dolab/gogo/pkgs/interceptors"
	"github.com/dolab/httpdispatch"
)

// A Middleware represents request filters or resource handler
//
// NOTE: It is the filter's responsibility to invoke ctx.Next() for chainning.
type Middleware func(ctx *Context)

// A Configer represents config interface
type Configer interface {
	RunMode() RunMode
	RunName() string
	Filename() string
	SetMode(mode RunMode)
	Section() *SectionConfig
	UnmarshalJSON(v interface{}) error
	UnmarshalYAML(v interface{}) error

	// for middlewares
	Interceptors() interceptors.Configer
	LoadInterceptors() error
}

// A Grouper represents router interface
type Grouper interface {
	NewGroup(prefix string, filters ...Middleware) Grouper
	SetHandler(handler Handler)
	Use(filters ...Middleware)
	OPTIONS(uri string, filter Middleware)
	HEAD(uri string, filter Middleware)
	POST(uri string, filter Middleware)
	GET(uri string, filter Middleware)
	PUT(uri string, filter Middleware)
	PATCH(uri string, filter Middleware)
	DELETE(uri string, filter Middleware)
	Any(uri string, filter Middleware)
	Static(uri, root string)
	Resource(uri string, resource interface{}) Grouper
	Proxy(method, uri string, proxy *httputil.ReverseProxy)
	HandlerFunc(method, uri string, fn http.HandlerFunc)
	Handler(method, uri string, handler http.Handler)
	Handle(method, uri string, filter Middleware)
	MountRPC(method string, rpc RPCServicer)
	MockHandle(method, uri string, recorder http.ResponseWriter, filter Middleware)
}

// A Servicer represents application interface
type Servicer interface {
	Init(config Configer, group Grouper)
	Middlewares()
	Resources()
}

// RPCServicer is the interface for rpc serve. It wraps HTTP handlers with
// additional methods for accessing metadata about the service.
type RPCServicer interface {
	// ProtocGenGogoVersion is the semantic version string of the version of
	// protoc-gen-gogo used to generate service.
	ProtocGenGogoVersion() string

	// ServiceRegistry returns a rpc method name to handlers map.
	ServiceRegistry(prefix string) map[string]Middleware

	// ServiceNames returns all rpc services registered.
	ServiceNames() []string

	// ServiceDescriptor returns gzipped bytes describing the .proto file that
	// this service was generated from. Once unzipped, the bytes can be
	// unmarshalled as a
	// github.com/golang/protobuf/protoc-gen-go/descriptor.FileDescriptorProto.
	//
	// The returned integer is the index of this particular service within that
	// FileDescriptorProto's 'Service' slice of ServiceDescriptorProtos. This is a
	// low-level field, expected to be used for reflection.
	ServiceDescriptor() ([]byte, int)
}

// A Handler represents handler interface
type Handler interface {
	http.Handler

	Handle(string, string, httpdispatch.Handler)
	ServeFiles(string, http.FileSystem)
}

// A Responser represents HTTP response interface
type Responser interface {
	http.ResponseWriter
	http.Flusher

	HeaderFlushed() bool        // whether response header has been sent?
	FlushHeader()               // send response header only if it has not sent
	Status() int                // response status code
	Size() int                  // return the size of response body
	Hijack(http.ResponseWriter) // hijack response with new http.ResponseWriter
}

// A StatusCoder represents HTTP response status code interface.
// it is useful for custom response data with response status code
type StatusCoder interface {
	StatusCode() int
}

// A Logger represents log interface
type Logger interface {
	New(requestID string) Logger
	Reuse(l Logger)
	RequestID() string
	SetLevelByName(level string) error
	SetColor(color bool)

	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}
