package rdtp

import "encoding/json"

// Request is the json model of a request for the rdtp service
type Request struct {
	Type       RequestType `json:"type"`
	LocalAddr  *Addr       `json:"local_addr"`
	RemoteAddr *Addr       `json:"remote_addr"`
}

type RequestType string

const (
	RequestTypeAccept = RequestType("ACCEPT")
	RequestTypeDial   = RequestType("DIAL")
	RequestTypeListen = RequestType("LISTEN")
)

func NewRequest(requestType RequestType, laddr, raddr *Addr) ([]byte, error) {
	return json.Marshal(Request{Type: requestType, LocalAddr: laddr, RemoteAddr: raddr})
}
