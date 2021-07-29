package rdtp

import "encoding/json"

// Request is the json model of a request for the rdtp service
type Request struct {
	Type       *RequestType `json:"type"`
	LocalAddr  *Addr        `json:"local_addr"`
	RemoteAddr *Addr        `json:"remote_addr"`
}

// Response is the json model of a response from the rdtp service
// TODO: responses currently not being used
type Response struct {
	Type       *ResponseType `json:"type"`
	LocalAddr  *Addr         `json:"local_addr"`
	RemoteAddr *Addr         `json:"remote_addr"`
	Error      *ErrorType    `json:"error"`
}

type RequestType string
type ResponseType string
type ErrorType string

const (
	RequestTypeAccept = RequestType("ACCEPT")
	RequestTypeDial   = RequestType("DIAL")
	RequestTypeListen = RequestType("LISTEN")

	ResponseTypeOK    = ResponseType("OK")
	ResponseTypeError = ResponseType("ERROR")

	ErrorTypeInvalidAddress = ErrorType("INVALID_ADDRESS")
)

func NewRequest(requestType RequestType, laddr, raddr *Addr) ([]byte, error) {
	return json.Marshal(Request{Type: &requestType, LocalAddr: laddr, RemoteAddr: raddr})
}

func NewResponse(responseType ResponseType, laddr, raddr *Addr, errorType ErrorType) ([]byte, error) {
	return json.Marshal(Response{Type: &responseType, LocalAddr: laddr, RemoteAddr: raddr, Error: &errorType})
}
