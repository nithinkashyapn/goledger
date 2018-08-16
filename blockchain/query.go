package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (setup *FabricSetup) QueryHello() (string, error) {

	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "hello")

	response, err := setup.client.Query(channel.Request{
		ChaincodeID: setup.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1]), []byte(args[2])}})
	if err != nil {
		return "", fmt.Errorf("Failed to query %v\n", err)
	}

	return string(response.Payload), nil
}
