package http_server

import "github.com/gorilla/mux"

type HTTPServerConfig struct {
	Port int `yaml:"port"`
}

func (s *HTTPServerConfig) Create(impl *Impl) (*Impl, error) {
	impl.router = mux.NewRouter()
	impl.config = s
	return impl, nil
}
