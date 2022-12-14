package ginp

import (
	"io/ioutil"
	"net/http"
	"sync"
)

type engine struct {
	ctxPool sync.Pool
}

func (e *engine) Post(url string, handlerFunc ...HandlerFunc) {
	PushToChain(url)
	PushToChain(url, handlerFunc...)
}

func (e *engine) Use(handlerFunc ...HandlerFunc) {
	PushMiddlewareToChain("", handlerFunc...)
}

func (e *engine) UseByRegex(reg string, handlerFunc ...HandlerFunc) {
	PushMiddlewareToChain(reg, handlerFunc...)
}

// ServeHTTP conforms to the http.Handler interface.
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//log.Println(w)
	//log.Println(req)

	c := e.ctxPool.Get().(*Context)
	c.Path = req.URL.Path
	c.Url = req.URL
	c.Host = req.Host
	c.HanderIndex = 0
	c.Writer = &w
	c.Request = req

	if req.Body != nil {
		body, _ := ioutil.ReadAll(req.Body)
		c.Body = body
	} else {
		c.Body = []byte{}
	}

	r, isOk := router[c.Path]
	if !isOk {
		// TODO: 这儿直接返回，不是很合适
		return
	}

	c.Handers = r

	e.handleHTTPRequest(c)

	e.ctxPool.Put(c)
}

// 这儿是处理的业务逻辑
func (e *engine) handleHTTPRequest(c *Context) {
	//log.Println(c.Path)

	for int(c.HanderIndex) < len(c.Handers) {
		h := c.Handers[c.HanderIndex]

		h(c)
		c.HanderIndex += 1
	}

	// 找到handlers

	// Find root of the tree for the given HTTP method
}

func (e *engine) Run(host string) (err error) {
	return http.ListenAndServe(host, e)
}
