package models

// Protocol represents the network protocol type
type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

// String returns the string representation of the protocol
func (p Protocol) String() string {
	return string(p)
}
