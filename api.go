package rdtp

import "encoding/json"

// ClientMessage is the json model of a request for the rdtp service
type ClientMessage struct {
	Type       ClientMessageType `json:"type"`
	LocalAddr  Addr              `json:"local_addr"`
	RemoteAddr Addr              `json:"remote_addr"`
}

// ServiceMessage is the json model of a message/response from the rdtp service
type ServiceMessage struct {
	Type       ServiceMessageType `json:"type"`
	LocalAddr  Addr               `json:"local_addr"`
	RemoteAddr Addr               `json:"remote_addr"`
	Error      ServiceErrorType   `json:"error,omitempty"`
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

	// ServiceMessageTypeNotify is the message type sent from rdtp-service to clients
	// to notify them that there is a new remote client for the client's listener
	ServiceMessageTypeNotify = ServiceMessageType("NOTIFY")

	// ServiceMessageTypeError is the message type sent from rdtp-service to
	// clients to acklowledge their request and indicate that there was an error
	ServiceMessageTypeError = ServiceMessageType("ERROR")

	// ServiceErrorTypeConnClosedByClient is the error type for errors caused by
	// the rdtp client closing the client -> rdtp-service connection
	ServiceErrorTypeConnClosedByClient = ServiceErrorType("CONN_CLOSED_BY_CLIENT")
)

// NewClientMessage returns a serialized client message
func NewClientMessage(clientMessageType ClientMessageType, laddr, raddr *Addr) ([]byte, error) {
	laddrd, raddrd := getAddressesDereferenced(laddr, raddr)
	return json.Marshal(ClientMessage{Type: clientMessageType, LocalAddr: laddrd, RemoteAddr: raddrd})
}

// NewServiceMessage returns a serialized service message
func NewServiceMessage(serviceMessageType ServiceMessageType, laddr, raddr *Addr, errorType *ServiceErrorType) ([]byte, error) {
	laddrd, raddrd := getAddressesDereferenced(laddr, raddr)
	if serviceMessageType == ServiceMessageTypeError && errorType != nil {
		return json.Marshal(ServiceMessage{Type: serviceMessageType, LocalAddr: laddrd, RemoteAddr: raddrd, Error: *errorType})

	}
	return json.Marshal(ServiceMessage{Type: serviceMessageType, LocalAddr: laddrd, RemoteAddr: raddrd})
}

func getAddressesDereferenced(laddr, raddr *Addr) (local Addr, remote Addr) {
	if laddr != nil {
		local.Host = laddr.Host
		local.Port = laddr.Port
	}
	if raddr != nil {
		remote.Host = raddr.Host
		remote.Port = raddr.Port
	}
	return
}
