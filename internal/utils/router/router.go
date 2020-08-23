package router

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/julienschmidt/httprouter"
	"github.com/nichojovi/tauria-test/internal/utils/response"
	"github.com/opentracing/opentracing-go"
	tlog "github.com/opentracing/opentracing-go/log"
)

type MyRouter struct {
	Httprouter     *httprouter.Router
	WrappedHandler http.Handler
	Options        *Options
}

var (
	defaultExcludeHeader = []string{
		"Accounts-Authorization",
		"Cookie",
	}
)

type Options struct {
	Prefix  string
	Timeout int
}

type WrittenResponseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *WrittenResponseWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *WrittenResponseWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

func (w *WrittenResponseWriter) Written() bool {
	return w.written
}

var HttpRouter *httprouter.Router

func init() {
	HttpRouter = httprouter.New()
}

func WrapperHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writtenResponseWriter := &WrittenResponseWriter{
			ResponseWriter: w,
			written:        false,
		}
		w = writtenResponseWriter
		_ = httpsnoop.CaptureMetrics(HttpRouter, w, r)
	})
}

func New(o *Options) *MyRouter {
	myrouter := &MyRouter{
		Options:    o,
		Httprouter: HttpRouter,
	}

	return myrouter
}

type Handle func(http.ResponseWriter, *http.Request) *response.JSONResponse

func (mr *MyRouter) GET(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	mr.Httprouter.GET(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) POST(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	mr.Httprouter.POST(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) PUT(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	mr.Httprouter.PUT(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) DELETE(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	mr.Httprouter.DELETE(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) handleNow(fullPath string, handle Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t := time.Now()
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(mr.Options.Timeout))
		defer cancel()
		ctx = context.WithValue(ctx, "HTTPParams", ps)

		span, ctx := opentracing.StartSpanFromContext(ctx, r.RequestURI)
		defer span.Finish()

		r.Header.Set("routePath", fullPath)
		r = r.WithContext(ctx)

		respChan := make(chan *response.JSONResponse)
		go func() {
			defer panicRecover(r, fullPath)
			resp := handle(w, r)
			respChan <- resp
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte("timeout"))
			}
		case resp := <-respChan:
			if resp != nil {
				span.LogFields(tlog.Object("log", resp.Log))
				span.SetTag("httpCode", resp.StatusCode)
				resp.SetLatency(time.Since(t).Seconds() * 1000)
				if resp.StatusCode > 499 {
					m := map[string]interface{}{}
					excludeHeaderDump(r, defaultExcludeHeader)
					httpDump, _ := httputil.DumpRequest(r, true)
					m["ERROR:"] = resp.RealError
					m["DUMP:"] = string(httpDump)
					log.Printf("%+v", m)
				}

				resp.Send(w)
			} else {
				if w, ok := w.(*WrittenResponseWriter); ok && !w.Written() {
					log.Println("Error nil response from the handler")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(""))
				}
			}
		}
		return
	}
}

func excludeHeaderDump(r *http.Request, excludeHeader []string) {
	for _, v := range excludeHeader {
		r.Header.Del(v)
	}
}

func GetHttpParam(ctx context.Context, name string) string {
	ps := ctx.Value("HTTPParams").(httprouter.Params)
	return ps.ByName(name)
}

func panicRecover(r *http.Request, path string) {
	if err := recover(); err != nil {
		log.Println("Error nil response from the handler")
	}
}
