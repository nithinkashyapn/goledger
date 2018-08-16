package blockchain

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (setup *FabricSetup) InvokeHello(value string) (string, error) {

	var args []string
	args = append(args, "invoke")
	args = append(args, "add")
	args = append(args, "hello")
	args = append(args, value)

	eventID := "added"

	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChaincodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	response, err := setup.client.Execute(channel.Request{
		ChaincodeID:  setup.ChaincodeID,
		Fcn:          args[0],
		Args:         [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("Failed to insert value %v\n", err)
	}

	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received chaincode event %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("Did not receive chaincode event for %v\n", eventID)
	}

	return string(response.TransactionID), nil
}
