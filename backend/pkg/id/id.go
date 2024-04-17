package id

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/bwmarrin/snowflake"
)

const (
	// maximum number of bits for a to divide between node and step
	maxSnowflakeBits = 22
	// the number of netmask bits in an ip4
	netmaskBits = 16
	// the number of bits in an ip4
	ipBits = 32
)

func init() {
	// get the amount of host bits in an ip and make that the node bits
	snowflake.NodeBits = ipBits - netmaskBits
	snowflake.StepBits = maxSnowflakeBits - snowflake.NodeBits
}

// Service is the interface for generating unique ids
// Using twitter snowflake to generate unique id with timestamp, node id and step
// Is a 64 bit id compared to 128 bit id from UUID; inclusion of timestamp allows for sorting
type Service interface {
	// GenerateID generates a unique auth id
	GenerateID() string
}

// ServiceImpl is the implementation of the IDService interface
type ServiceImpl struct {
	node *snowflake.Node
}

// GenerateID generates a unique auth ID
func (i *ServiceImpl) GenerateID() string {
	return fmt.Sprintf("%020d", i.node.Generate())
}

func NewServiceFromIP(iPv4 string) (*ServiceImpl, error) {
	nodeId, err := nodeIDFromIP(iPv4)
	if err != nil {
		return nil, err
	}
	return NewService(nodeId)
}

func nodeIDFromIP(iPv4 string) (int64, error) {
	ip := net.ParseIP(iPv4)
	if ip == nil {
		return 0, fmt.Errorf("invalid ip address: %v", iPv4)
	}
	if ip.To4() == nil {
		return 0, fmt.Errorf("not an ipv4 address: %v", iPv4)
	}
	nodeId := int64(binary.BigEndian.Uint32(ip.To4()) & 0x0000FFFF)
	return nodeId, nil
}

var service *ServiceImpl

func NewService(nodeId int64) (*ServiceImpl, error) {
	if service != nil {
		return service, nil
	}

	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		return nil, err
	}

	service = &ServiceImpl{node: node}
	return service, nil
}
