package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)



type DataChaincode struct {}




type StructuredData struct {
	ID 				string 	`json:"id"`
	FirstField 		string 	`json:"firstField"`
	SecondField 	string 	`json:"secondField"`
}



func main() {
	err := shim.Start(new(DataChaincode))
	if err != nil {
		fmt.Printf("Error starting Basic chaincode: %s", err)
	}
}



func (t *DataChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initialized our Chaincode")
	return shim.Success(nil)
}




func (t *DataChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoking our new Chaincode")

	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	
	if function == "getSimpleData" {
		result, err = t.getSimpleData(stub, args)
	} else if function == "getComplexData" {
		result, err = t.getComplexData(stub, args)
	} else if function == "createSimpleData" {
		result, err = t.createSimpleData(stub, args)
	} else if function == "createComplexData" {
		result, err = t.createComplexData(stub, args)
	} else if function == "searchComplexData" {
		var complexSearchResult []byte
		complexSearchResult, err = t.searchComplexData(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(complexSearchResult)
	} else {
		return shim.Error("Invalid invoke function name. ")
	}

	if err != nil {
			return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}







func (t *DataChaincode) getSimpleData(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
			return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
			return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}


func (t *DataChaincode) createSimpleData(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
			return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}






















func (t *DataChaincode) getComplexData(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	keys := make([]string, 1)
	keys[0] = args[0]
	ObjectKey,_ := stub.CreateCompositeKey("StructuredData", keys)

	value, err := stub.GetState(ObjectKey)

	if err != nil {
			return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
			return "", fmt.Errorf("Asset not found: %s", args[0])
	}

	var sd StructuredData
	json.Unmarshal(value, &sd) 

	if sd.FirstField == "1" {
		fmt.Println("We retrieved complex data from the Blockchain, and the FirstField attribute has a value of 1")
	}

	return string(value), nil
}

func (t *DataChaincode) createComplexData(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and 2 values")
	}

	keys := make([]string, 1)
	keys[0] = args[0]
	ObjectKey,_ := stub.CreateCompositeKey("StructuredData", keys)
	 
	var sdObject = StructuredData{ID: args[0], FirstField: args[1], SecondField: args[2]}
	sdBytes, _ := json.Marshal(sdObject)
	stub.PutState(ObjectKey, sdBytes)

	return ObjectKey, nil
}
















func (t *DataChaincode) searchComplexData(stub shim.ChaincodeStubInterface, args []string) ([] byte, error)  {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect arguments. Expecting a query")
	}

	resultsIterator, err := stub.GetQueryResult(args[0])
    defer resultsIterator.Close()
    if err != nil {
		fmt.Println("Failed to search Complex Data")
        return nil, err
	}
	

    // buffer is a JSON array containing QueryRecords
    var buffer bytes.Buffer
    buffer.WriteString("[")
    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
			fmt.Println("Failed while looping through the ResultSet of Complex Data")
            return nil, err
		}
		
        // Add a comma before array members, suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
        buffer.WriteString("\"")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString("\"")
        buffer.WriteString(", \"Record\":")
        // Record is a JSON object, so we write as-is
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")
	return buffer.Bytes(), nil

}

