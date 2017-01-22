package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/op/go-logging"
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
	//	IssueYear  uint16  `json:"issue_year"`
	//	IssueMonth uint8   `json:"issue_month"`
	//	IssueDay   uint8   `json:"issue_day"`
}

type IssueSet struct {
	Issues []Issue `json:"issues"`
}

//var myLogger = logging.MustGetLogger("asset_mgm")

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Entering into init ")

	// Nothing to do here, just return
	fmt.Println("Returning from init")
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	fmt.Println("Entering into Invoke : " + function)
	//user, err := t.get_username(stub)
	//if err != nil{
	//	return nil, errors.New("#####  Failed to get username for function: " + function + " #####")
	//}
	//fmt.Println(" function called by: "+ user)

	//called by issue func
	if function == "issue" {
		//to get (neme amount)
		if len(args) != 2 {
			return nil, errors.New("##### Incorect number of arg#####")
		}

		var person_Name string
		var issue_amount float64

		//need to check person has already been registered or not

		//get and set values
		person_Name = args[0]
		key := "issue/" + person_Name
		issue_amount, err = strconv.ParseFloat(args[1], 64)
		if err != nil {
			return nil, errors.New("can't parse float")
		}

		//myLogger.Info("start get issue")
		//currentBytes, err := t.getIssue(stub, key)
		currentBytes, err := stub.GetState(key)
		fmt.Println(currentBytes)
		fmt.Println(err)
		//myLogger.Info(currentBytes)
		//err := json.Unmarshal(currentBytes, &record_issue)

		if err == nil && currentBytes != nil {
			//myLogger.Info("enter into reged person")

			regData, err := NewJson(currentBytes)
			currentAmount, err := regData.Get("amount").Float64()
			newAmount := currentAmount + issue_amount
			regData.Set("amount", newAmount)
			t := time.Now()
			timeString := ""
			timeString = t.String()
			regData.Set("issueTime", timeString)
			newBytes, err := json.Marshal(regData)
			err = stub.PutState(key, newBytes)
			if err != nil {
				return nil, errors.New("#####  faild to update data #####")
			}

			return nil, nil

		}
		err = nil
		//myLogger.Info("enter into new person")
		//Get current date and time
		t := time.Now()

		/*time record
		var year uint16
		var month uint8
		var day uint8

		year = uint16(t.Year())
		month = uint8(t.Month())
		day = uint8(t.Day())

		var record_issue Issue
		record_issue = Issue{
			PersonName: person_Name,
			Amount:     issue_amount,
			IssueYear:  year,
			IssueMonth: month,
			IssueDay:   day,
		}*/
		timeString := ""
		timeString = t.String()

		var record_issue Issue

		record_issue = Issue{
			PersonName: person_Name,
			Amount:     issue_amount,
			IssueTime:  timeString,
		}

		bytes, err := json.Marshal(record_issue)

		if err != nil {
			return nil, errors.New("#####  Failed to convert json #####")
		}
		err = stub.PutState(key, []byte(bytes))

		if err != nil {
			return nil, errors.New("unable to put the state")
		}

		fmt.Println("complete register name and amount")
		return nil, nil
	} else {
		return nil, errors.New("undifine such func")
	}
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println(function)

	if function == "getIssue" {

		if len(args) != 1 {
			fmt.Println("####Incorrect args number####")
			return nil, errors.New("#### Incorrect number of args###")
		}

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

	if err != nil {
		return nil, errors.New("#### failed to get state of" + key)
	}
	err = json.Unmarshal(issueBytes, &record_issue)

	if err != nil {
		return nil, errors.New("#### failed to unmarshall state")
	}

	bytes, err := json.Marshal(record_issue)
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

/*func (t *SimpleChaincode) payCrypto( payPerson string, getPerson strimg,payAmount float64) ([]byte, error) {

	var err err

	ts, err = stub.GetTxTimestamp()

	if err != nil {
		fmt.Printf("Error getting transaction timestamp: %s", err)
	}

  err := paySubstruction(payPerson,payAmount)

  if t.getIssue(getPerson) = nil, _ {
	// add crypto to new person

}else{
  currentBytes, _ := t.getIssu(getPerson)
  err  := payAdd(currentBytes,payAmount)

}
	if err != nil {
		return nil, errors.New("#### failed to get state of" + payPerson)
	}
}*/

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

/*func paySubstruction(payPerson string, amount float64)(error) {

	payPersonBytes, err := t.getIssue(payPerson)
  fmt.Println(payPersonBytes)
	payPersonRecords = NewJson(payPersonBytes)
  payPersonAmout,err := payPersonRecords.Get("amount").Float64()
  afterPayAmount := payPersonAmout - amount
	payPersonRecords.Set("amount",afterPayAmount)
	t := time.Now()
	timeString := ""
	timeString = t.String()
	payPersonRecords.Set("issueTime",timeString)
  afterPayBytes, err = json.Marshal(payPersonRecords)
  err = stub.PutState(payPerson,afterPayBytes)
	if err != nil{
    return errors.New("miss updat")
}
	return nil
}

 func payAdd(currentBytes []byte(), addAmount float64)(error){

		getParsonJson, err := NewJson(currentBytes)
		currentAmount, err := getParsonJson.Get("amount").Float64()
		newAmount := currentAmount + addAmount
		getParsonJson.Set("amount", newAmount)
		t := time.Now()
		timeString := ""
		timeString = t.String()
		payPersonRecords.Set("issueTime",timeString)
		newBytes, err := json.Marshal(getParsonJson)
		err = stub.PutState(key, newBytes)
		if err != nil {
			return errors.New("#####  faild to update data #####")
		}

}
*/

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
