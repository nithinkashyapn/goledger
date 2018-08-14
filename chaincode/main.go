package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}

func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	fmt.Println("----INIT----")

	function, _ := stub.GetFunctionAndParameters()

	if function != "init" {
		return shim.Error("Wrong function call for init")
	}

	err := stub.PutState("hello", []byte("world!"))

	if err != nil {
		return shim.Error("Cannot put state in Init() error is " + err.Error())
	}

	return shim.Success(nil)
}

func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fmt.Println("----INVOKE----")

	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Wrong function call for invoke")
	}

	if len(args) < 1 {
		return shim.Error("Number of arguments to invoke is insufficient")
	}

	if args[0] == "query" {
		return cc.query(stub, args)
	}

	if args[0] == "add" {
		return cc.add(stub, args)
	}

	return shim.Error("Wrong action passes in first argument")
}

func (cc *Chaincode) query(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("----QUERY----")

	if len(args) < 2 {
		return shim.Error("Number of arguments to query is wrong")
	}

	if args[1] == "hello" {
		state, err := stub.GetState("hello")

		if err != nil {
			return shim.Error("Could not get state in query")
		}

		return shim.Success(state)
	}

	return shim.Error("Wrong argument passed for key")
}

func (cc *Chaincode) add(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("----ADD----")

	if len(args) < 2 {
		return shim.Error("Number of arguments to add is insufficient")
	}

	if args[1] == "hello" && len(args) == 3 {

		err := stub.PutState("hello", []byte(args[2]))
		if err != nil {
			return shim.Error("Failed to put state in add " + err.Error())
		}

		err = stub.SetEvent("added", []byte{})
		if err != nil {
			return shim.Error("Failed to set event 'added' in add " + err.Error())
		}
		return shim.Success(nil)
	}

	return shim.Error("Unknown command check the second argument")
}

func main() {
	fmt.Println("----MAIN----")

	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting server cause of %v\n", err)
	}
}
