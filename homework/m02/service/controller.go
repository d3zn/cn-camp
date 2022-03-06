package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Controller struct {
	s Service
}

func NewController(service Service) *Controller {
	return &Controller{s: service}
}

func (c *Controller) Header(w http.ResponseWriter, r *http.Request) {
	c.s.Header()
	w.Header().Add("name", r.Header.Get("Name"))
	w.WriteHeader(http.StatusOK)
	return
}

func (c *Controller) Env(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(buf) > 0 {
		req := make(map[string]interface{})
		err := json.Unmarshal(buf, &req)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			_, _ = fmt.Fprintf(w, "get body failed, err = %v", err)
			return
		}
		key := req["env"].(string)
		env, err := c.s.Env(key)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, "get env failed, err = %v", err)
			return
		}
		w.Header().Add(key, env)
		_, _ = fmt.Fprintf(w, "%s=%s", key, env)
	}
	return
}

func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	if err := c.s.Health(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, "health")
	return
}
