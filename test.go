package main

import(
	"errors"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}


type Issue struct{
	PersonName     string   `json:"personName`
	Amount         float64  `json:"amount"`
	IssueYear      uint16   `json:"issue_year"`
	IssueMonth     uint8    `json:"issue_month"`
	IssueDay       uint8    `json:"issue_day"`
}

type IssueSet struct{
	Issues []Issue `json:"issues"`
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	fmt.Println("Entering into init ")

	// Nothing to do here, just return
	fmt.Println("Returning from init")
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error
	fmt.Println("Entering into Invoke : "+ function)
	//user, err := t.get_username(stub)
	//if err != nil{
	//	return nil, errors.New("#####  Failed to get username for function: " + function + " #####")
	//}
	//fmt.Println(" function called by: "+ user)

	//called by issue func
	if function == "issue"{
		//to get (neme amount)
		if len(args) != 2{
			return nil,errors.New("##### Incorect number of arg#####")
		}

		var person_Name string
		var issue_amount  float64

		//need to check person has already been registered or not

		//get and set values
		person_Name = args[0]
		key := "issue/"+ person_Name
		issue_amount, err = strconv.ParseFloat(args[1],64)
		if err != nil{
			return nil, errors.New("can't parse float")
		}

		//Get current date and time
		t := time.Now()

		//time record
		var year uint16
		var month uint8
		var day uint8

		year = uint16(t.Year())
		month = uint8(t.Month())
		day = uint8(t.Day())

		var record_issue Issue
		record_issue = Issue{
			PersonName: person_Name,
			Amount: issue_amount,
			IssueYear: year,
			IssueMonth:month,
			IssueDay:day,
		}

		bytes, err := json.Marshal(record_issue)
		if err != nil{
			return nil, errors.New("#####  Failed to convert json #####")
		}
		err = stub.PutState(key,[]byte(bytes))

		if err != nil {
			return nil, errors.New("unable to put the state")
		}

		fmt.Println("complete register name and amount")
		return nil,nil
	}else {
		return nil, errors.New("undifine such func")
	}

}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println(function)

	if function == "getIssue"{

		if len(args) !=1{
			fmt.Println("####Incorrect args number####")
			return nil, errors.New("#### Incorrect number of args###")
		}
		person_name := args[0]
		return  t.getIssue(stub,person_name)

	}else if function == "getAllIssues"{
		if len(args) != nil{
			fmt.Println("####Incorrect args number ####")
			return nil, errors.New("####Incorrect args number ####")
		}
		return  t.getAllIssue(stub)
	}

	fmt.Println("query did not find function")
	return nil,errors.New("####query did not find function####")

}

func (t *SimpleChaincode) getIssue(stub *shim.ChaincodeStub, personName string) ([]byte, error) {

	var err error
        var record_issue Issue

	//get issue
	key :="issue/"+ personName
	issueBytes, err := stub.GetState(key)
	fmt.Println(issueBytes)

	if err != nil{
		return nil, errors.New("#### failed to get state of"+ key)
	}
	err = json.Unmarshal(issueBytes,&record_issue)

	if err != nil {
		return nil,errors.New("#### failed to unmarshall state"+)
	}

	bytes, err := json.Marshal(record_issue)
	if err != nil{
		return nil ,errors.New("####error creating data ####")
	}
	return []byte(bytes), nil
}

func (t *SimpleChaincode) getAllIssue(stub *shim.ChaincodeStub,) ([]byte, error) {

	var err error
	var issue_recrod Issue
	var issue_set IssueSet

	//get issue
	iter , err := stub.RangeQueryState("issue/", "issue/~")
	fmt.Println(iter)
	if err != nil {
		return nil, errors.New("#### failed to get state ####" )
	}
	defer iter.Close()
	for iter.HasNext(){
		_,issue_asbytes, iterErr := iter.Next()
		if iterErr !=nil{
			return nil,errors.New("error")
		}
		err = json.Unmarshal(issue_asbytes, &issue_recrod)
		if err != nil {
			return nil, errors.New("####Error unmarmashalling data"+string(issue_asbytes)+"####")
		}
		issue_set.Issues = append(issue_set.Issues,issue_recrod)
	}

	bytes, err := json.Marshal(issue_set.Issues)

	if err != nil {
		return nil, errors.New("#### error creating return record ####")
	}
	return []byte(bytes), nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}