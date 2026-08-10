package handler

import "net/http"

type Renderer interface {
	Render(w http.ResponseWriter)
	GetStatusCode() int
	GetBody() string
}
