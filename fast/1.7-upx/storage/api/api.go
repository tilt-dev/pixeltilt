package api

type WriteRequest struct {
	Name string
	Body []byte
}

type WriteResponse struct {
	Name string
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
