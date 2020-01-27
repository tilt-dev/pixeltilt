package api

type WriteRequest struct {
	Name string
	Body []byte
}

type ReadRequest struct {
	Name string
}

type ReadResponse struct {
	Body []byte
}

type ListResponse struct {
	Names []string
}
