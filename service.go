package rdtp

import "encoding/json"

// ClientMessage is the json model of a request for the rdtp service
type ClientMessage struct {
	Type       *ClientMessageType `json:"type"`
	LocalAddr  *Addr              `json:"local_addr"`
	RemoteAddr *Addr              `json:"remote_addr"`
}

// ServiceMessage is the json model of a message/response from the rdtp service
// TODO: responses currently not being used
type ServiceMessage struct {
	Type       *ServiceMessageType `json:"type"`
	LocalAddr  *Addr               `json:"local_addr"`
	RemoteAddr *Addr               `json:"remote_addr"`
	Error      *ServiceErrorType   `json:"error"`
}

// ClientMessageType is the go type for
// client -> rdtp-service messages (requests)
type ClientMessageType string

// ServiceMessageType is the go type for
// rdtp-service -> client messages (responses/notifications)
type ServiceMessageType string

// ServiceErrorType is the go type included within a
// ServiceMessageType to provide more detail on failures
type ServiceErrorType string

const (
	// ClientMessageTypeAccept is the message type sent from clients
	// to rdtp-service to acknowledge a notification and receive a full duplex
	// communication channel between the dialing caller and the rdtp port
	// being listened on by a client
	// Note: RemoteAddress **must** be defined
	ClientMessageTypeAccept = ClientMessageType("ACCEPT")

	// ClientMessageTypeDial is the message type sent from clients
	// to rdtp-service to "dial" a remote rdtp address
	// Note: RemoteAddr **must** be defined
	ClientMessageTypeDial = ClientMessageType("DIAL")

	// ClientMessageTypeListen is the message type sent from clients
	// to rdtp-service to "listen" for inbound connections on a local rdtp port
	// Note: LocalAddr **must** be defined
	ClientMessageTypeListen = ClientMessageType("LISTEN")

	// ServiceMessageTypeOK is the message type sent from rdtp-service to clients
	// to acknowledge their request and indicate that it was served successfully
	ServiceMessageTypeOK = ServiceMessageType("OK")

	// ServiceMessageTypeError is the message type sent from rdtp-service to
	// clients to acklowledge their request and indicate that there was an error
	ServiceMessageTypeError = ServiceMessageType("ERROR")

	// ErrorTypeInvalidAddress is the error type sent from rdtp-service to clients
	// whenever a client makes a request with an invalid address
	ErrorTypeInvalidAddress = ServiceErrorType("INVALID_ADDRESS")
)

func NewClientMessage(clientMessageType ClientMessageType, laddr, raddr *Addr) ([]byte, error) {
	return json.Marshal(ClientMessage{Type: &clientMessageType, LocalAddr: laddr, RemoteAddr: raddr})
}

func NewServiceMessage(serviceMessageType ServiceMessageType, laddr, raddr *Addr, errorType ServiceErrorType) ([]byte, error) {
	return json.Marshal(ServiceMessage{Type: &serviceMessageType, LocalAddr: laddr, RemoteAddr: raddr, Error: &errorType})
}
