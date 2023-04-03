package internalhttp

type Server struct { // TODO
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{}
}

func (s *Server) Start() error {
	// TODO
	return nil
}

func (s *Server) Stop() error {
	// TODO
	return nil
}

// TODO
