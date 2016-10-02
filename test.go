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

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}



type Issue struct{
	PersonName     string   `json:"personName`
	Amount         float64  `json:"amount"`
	IssueYear      uint16   `json:"issue_year"`
	IssueMonth     uint8    `json:"issue_month"`
	IssueDay       uint8    `json:"issue_day"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Entering into init ")

	// Nothing to do here, just return
	fmt.Println("Returning from init")
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

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

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Entering into query ")

	// Nothing to do here, just return
	fmt.Println("Returning from query")
	return nil, nil

}
