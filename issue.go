package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Cert struct {
	ID				string `json:"id"`			// ID is unique for a person
	Name			string `json:"sname"`		// Name of person
	CertName		string `json:"cname"`		// Name of certificate eg:- CCNS, 
	CertDetails		string `json:"cdets"`
	Organisation	string `json:"org"`
	TxHash			string `json:"txhash"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCert" {
		return s.queryCert(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createCert" {
		return s.createCert(APIstub, args)
	} else if function == "queryAllCerts" {
		return s.queryAllCerts(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCert(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	certAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(certAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	i1 := APIstub.GetTxID()
	i2 := APIstub.GetTxID()
	i3 := APIstub.GetTxID()
	certs := []Cert{
		Cert{ID: "STUD0", Name: "Arnav", CertName: "ABCD", Organisation: "BlockVidhya", TxHash: i1},
		Cert{ID: "STUD1", Name: "Rachit", CertName: "SSMS", Organisation: "BlockVidhya", TxHash: i2},
		Cert{ID: "STUD2", Name: "Puneet", CertName: "XJZZ", Organisation: "BlockVidhya", TxHash: i3},
	}

	i := 0
	for i < len(certs) {
		fmt.Println("i is ", i)
		certAsBytes, _ := json.Marshal(certs[i])
		APIstub.PutState("CERT"+strconv.Itoa(i), certAsBytes)
		fmt.Println("Added", certs[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createCert(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	i := APIstub.GetTxID()

	var cert = Cert{ID: args[1], Name: args[2], CertName: args[3], Organisation: args[4], TxHash: i}	// Change

	certAsBytes, _ := json.Marshal(cert)
	APIstub.PutState(args[0], certAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllCerts(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CERT0"
	endKey := "CERT999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	fmt.Printf("- queryAllCerts:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}