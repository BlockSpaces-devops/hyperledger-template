package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)



type BasicChaincode struct {}



func main() {
	err := shim.Start(new(BasicChaincode))
	if err != nil {
		fmt.Printf("Error starting Basic chaincode: %s", err)
	}
}



func (t *BasicChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initialized our Chaincode")
	return shim.Success(nil)
}




func (t *BasicChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoking our new Chaincode")

	function, args := stub.GetFunctionAndParameters()

	if function == "hello" {
		return t.hello(stub, args)
	} else if function == "goodbye" {
		return t.goodbye(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"hello\" \"goodbye\"")
}






func (t *BasicChaincode) hello(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Println("Why did this not work")

	fmt.Printf("Hello, %s", args[0])

	return shim.Success(nil)
}


func (t *BasicChaincode) goodbye(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fmt.Printf("Goodbye, %s", args[0])

	return shim.Success(nil)
}

