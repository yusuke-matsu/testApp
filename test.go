package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Json struct {
	data interface{}
}

type Issue struct {
	PersonName string  `json:"personName`
	Amount     float64 `json:"amount"`
	IssueTime  string  `json:"issueTime"`
}

type IssueSet struct {
	Issues []Issue `json:"issues"`
}

var myLogger = shim.NewLogger("myChainCode")

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Entering into init ")

	// Nothing to do here, just return
	fmt.Println("Returning from init")
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Entering into Invoke : " + function)
	//called by issue func

	if function == "issue" {
		//to get (neme amount)
		err := t.addIssue(stub, args)

		if err != nil {
			return nil, err
		}
		return nil, nil
	} else if function == "payCrypto" {
		err := t.payCrypto(stub, args)
		if err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		return nil, errors.New("undifine such func")
	}
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println(function)

	if function == "getIssue" {

		/*if len(args) != 1 {
			fmt.Println("####Incorrect args number####")
			return nil, errors.New("#### Incorrect number of args###")
		}*/

		person_name := args[0]
		return t.getIssue(stub, person_name)

	} else if function == "getAllIssues" {
		if len(args) != 0 {
			fmt.Println("####Incorrect args number ####")
			return nil, errors.New("####Incorrect args number ####")
		}
		return t.getAllIssue(stub)
	}

	fmt.Println("query did not find function")
	return nil, errors.New("####query did not find function####")
}

func (t *SimpleChaincode) getIssue(stub shim.ChaincodeStubInterface, personName string) ([]byte, error) {

	var err error
	var record_issue Issue

	//get issue
	key := "issue/" + personName
	issueBytes, err := stub.GetState(key)
	fmt.Println(issueBytes)
	myLogger.Info(issueBytes)
	if err != nil {
		return nil, errors.New("#### failed to get state of" + key)
	}
	err = json.Unmarshal(issueBytes, &record_issue)
	myLogger.Info(record_issue)
	if err != nil {
		return nil, errors.New("#### failed to unmarshall state")
	}

	bytes, err := json.Marshal(record_issue)
	myLogger.Info(bytes)
	if err != nil {
		return nil, errors.New("####error creating data ####")
	}
	return []byte(bytes), nil
}

func (t *SimpleChaincode) getAllIssue(stub shim.ChaincodeStubInterface) ([]byte, error) {

	var err error
	var issue_record Issue
	var issue_set IssueSet

	//get issue
	iter, err := stub.RangeQueryState("issue/", "issue/~")
	fmt.Println(iter)
	if err != nil {
		return nil, errors.New("#### failed to get state ####")
	}
	defer iter.Close()
	for iter.HasNext() {
		_, issue_asbytes, iterErr := iter.Next()
		if iterErr != nil {
			return nil, errors.New("error")
		}
		err = json.Unmarshal(issue_asbytes, &issue_record)
		if err != nil {
			return nil, errors.New("####Error unmarmashalling data" + string(issue_asbytes) + "####")
		}
		issue_set.Issues = append(issue_set.Issues, issue_record)
	}

	bytes, err := json.Marshal(issue_set.Issues)

	if err != nil {
		return nil, errors.New("#### error creating return record ####")
	}
	return []byte(bytes), nil
}

func (t *SimpleChaincode) addIssue(stub shim.ChaincodeStubInterface, args []string) error {

	var err error
	var person_Name string
	var issue_amount float64
	var record_issue Issue

	myLogger.Info(args)
	slliceLength := len(args)
	myLogger.Info(slliceLength)

	if slliceLength == 3 {
		person_Name = args[2]
	} else {
		person_Name = args[0]
	}

	key := "issue/" + person_Name
	issue_amount, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return errors.New("can't parse float")
	}

	queryMethod := "getIssue"
	currentBytes, err := t.Query(stub, queryMethod, args)
	myLogger.Info(currentBytes)

	if err == nil && currentBytes != nil {
		err = json.Unmarshal(currentBytes, &record_issue)
		myLogger.Info(record_issue)
		newAmount := issue_amount + record_issue.Amount
		newTime := time.Now()
		timeString := ""
		timeString = newTime.String()

		newRecordIssue := Issue{
			PersonName: record_issue.PersonName,
			Amount:     newAmount,
			IssueTime:  timeString,
		}
		myLogger.Info(newRecordIssue)

		newBytes, err := json.Marshal(newRecordIssue)
		myLogger.Info(newBytes)
		err = stub.PutState(key, []byte(newBytes))
		if err != nil {
			return errors.New("#####  faild to update data #####")
		}

	} else {

		err = nil
		//Get current date and time
		newTime := time.Now()
		timeString := ""
		timeString = newTime.String()
		record_issue = Issue{
			PersonName: person_Name,
			Amount:     issue_amount,
			IssueTime:  timeString,
		}
		bytes, err := json.Marshal(record_issue)
		if err != nil {
			return errors.New("#####  Failed to convert json #####")
		}
		err = stub.PutState(key, []byte(bytes))
		if err != nil {
			return errors.New("unable to put the state")
		}

	}
	return nil
}

func (t *SimpleChaincode) payCrypto(stub shim.ChaincodeStubInterface, args []string) error {

	var err error

	/*if err != nil {
		fmt.Printf("Error getting transaction timestamp: %s", err)
	}*/

	err = t.paySubstruction(stub, args)
	if err != nil {
		return errors.New("error to substruct")
	}
	err = t.addIssue(stub, args)
	if err != nil {
		return errors.New("error to add")
	}
	return nil
}

func NewJson(body []byte) (*Json, error) {
	j := new(Json)
	err := j.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (j *Json) UnmarshalJSON(p []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(p))
	dec.UseNumber()
	return dec.Decode(&j.data)
}

func (j *Json) Get(key string) *Json {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}
		}
	}
	return &Json{nil}
}

func (j *Json) Map() (map[string]interface{}, error) {
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

func (j *Json) Set(key string, val interface{}) {
	m, err := j.Map()
	if err != nil {
		return
	}
	m[key] = val
}

func (j *Json) Float64() (float64, error) {
	switch j.data.(type) {
	case json.Number:
		return j.data.(json.Number).Float64()
	case float32, float64:
		return reflect.ValueOf(j.data).Float(), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(j.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(j.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")

}

func (t *SimpleChaincode) paySubstruction(stub shim.ChaincodeStubInterface, args []string) error {

	var record_issue Issue
	person_Name := args[0]
	key := "issue/" + person_Name
	queryMethod := "getIssue"
	currentBytes, err := t.Query(stub, queryMethod, args)
	myLogger.Info(currentBytes)
	err = json.Unmarshal(currentBytes, &record_issue)
	subAmount, err := strconv.ParseFloat(args[1], 64)
	newAmount := record_issue.Amount - subAmount
	newTime := time.Now()
	timeString := ""
	timeString = newTime.String()
	newRecordIssue := Issue{
		PersonName: record_issue.PersonName,
		Amount:     newAmount,
		IssueTime:  timeString,
	}
	myLogger.Info(newRecordIssue)
	newBytes, err := json.Marshal(newRecordIssue)
	myLogger.Info(newBytes)
	err = stub.PutState(key, []byte(newBytes))
	if err != nil {
		return errors.New("#####  faild to update data #####")
	}

	return nil
}

/*func (t *SimpleChaincode) payAdd(stub shim.ChaincodeStubInterface, args []string) error {
	var err error
	var person_Name string
	var issue_amount float64
	var record_issue Issue

	person_Name = args[0]
	key := "issue/" + person_Name
	issue_amount, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return errors.New("can't parse float")
	}
	queryMethod := "getIssue"
	currentBytes, err := t.Query(stub, queryMethod, args)
	myLogger.Info(currentBytes)

	if err == nil && currentBytes != nil {
		err = json.Unmarshal(currentBytes, &record_issue)
		myLogger.Info(record_issue)
		newAmount := issue_amount + record_issue.Amount
		newTime := time.Now()
		timeString := ""
		timeString = newTime.String()
		newRecordIssue := Issue{
			PersonName: record_issue.PersonName,
			Amount:     newAmount,
			IssueTime:  timeString,
		}
		myLogger.Info(newRecordIssue)
		newBytes, err := json.Marshal(newRecordIssue)
		myLogger.Info(newBytes)
		err = stub.PutState(key, []byte(newBytes))
		if err != nil {
			return errors.New("#####  faild to update data #####")
		}
		return nil
	}

	if err != nil {
		return errors.New("#####  faild to update data #####")
	}
	return nil
}*/

func main() {
	myLogger.SetLevel(shim.LogInfo)
	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	shim.SetLoggingLevel(logLevel)
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
