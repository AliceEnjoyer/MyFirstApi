package response

/*
эти данные будут в каждой response структуре,
поэтому мы выносим их в эту структуру.
*/
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}
