package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract sc struct
type SmartContract struct {
	contractapi.Contract
}

//Transactions structure
type Transactions struct {
	ObjectType         string    `json:"docType"`
	TxID               string    `json:"txId"`      //generate using uuid
	MultiTxID          string    `json:"multiTxID"` // to identify tx which are related to single multi wallet transfer
	ReferenceID        string    `json:"referenceId"`
	FromName           string    `json:"fromName"`   //user, activity, sosh names
	SoshID             string    `json:"soshId"`     //if transactions are happening in sosh
	ActivityID         string    `json:"activityId"` //if transactions are happening in activity
	FromID             string    `json:"fromId"`     // from id from mongo of user/sosh/activity
	FromWalletID       string    `json:"fromWalletId"`
	FromType           string    `json:"fromType"`        //SOSH, ACTIVITY, INDIVIDUAL
	FromAccountType    string    `json:"fromAccountType"` //Individual, Business, Charity
	ToName             string    `json:"toName"`          // user, activity, sosh name
	ToID               string    `json:"toId"`            //to id from mongo of user/sosh/activity
	ToWalletID         string    `json:"toWalletId"`
	ToAccountType      string    `json:"toAccountType"`
	ToType             string    `json:"toType"`
	Amount             string    `json:"txamount"`
	Date               time.Time `json:"date"`
	TransactionType    string    `json:"txType"`             // Transferred, Withdrawn, Add Funds
	TransactionComment string    `json:"txComment"`          //Custom Message for transactions
	PaymentMethod      string    `json:"paymentMethod"`      // SW, CC, DC, B/A
	PaymentGatewayTxID string    `json:"paymentGatewayTxId"` //Only for add funds
	TransactionCharges string    `json:"txCharges"`
	PaymentStatus      string    `json:"paymentStatus"` //Pending, Failed, Success
	EntryID            string    `json:"entryId"`
	EntryName          string    `json:"entryName"`
	EntryType          string    `json:"entryType"`
	PlatformFees       string    `json:"platformFees"`
	FromEmail          string    `json:"fromEmail"`
	FromMobile         string    `json:"fromMobile"`
}

// Page for pagination
type Page struct {
	Data     []Transactions `json:"data"`
	Count    int32          `json:"count"`
	Bookmark string         `json:"bookmark"`
}

// QueryResult structure used for handling result of query
// simple user to user, user to sosh, user to activity
type QueryResult struct {
	Record *Transactions
}

// AddTransactions Creating a Transaction
func (s *SmartContract) AddTransactions(ctx contractapi.TransactionContextInterface, txid string, date string, referenceId string, multiTxID string, fromName string, fromId string, fromWalletId string, fromType string, fromAccountType string, toName string, toId string, toWalletId string, toType string, toAccountType string, txamount string, txType string, txComment string, paymentMethod string, paymentGatewayTxId string, txCharges string, paymentStatus string, soshId string, activityId string, entryId string, entryName string, entryType string, platformFees string, fromEmail string, fromMobile string) error {

	txAsBytes, err := ctx.GetStub().GetState(txid)
	if err != nil {
		return fmt.Errorf("Failed to get transaction")
	} else if txAsBytes != nil {
		fmt.Println("Transaction id already exist " + txid)
		return fmt.Errorf("Transaction id already exists: " + txid)
	} else {
		myDate, err := time.Parse(time.RFC3339, date)
		if err != nil {
			panic(err)
		}

		tx := Transactions{
			ObjectType:         "tx",
			TxID:               txid,
			ReferenceID:        referenceId,
			MultiTxID:          multiTxID,
			FromName:           fromName,
			FromID:             fromId,
			FromWalletID:       fromWalletId,
			FromType:           fromType,
			FromAccountType:    fromAccountType,
			ToName:             toName,
			ToID:               toId,
			ToWalletID:         toWalletId,
			ToType:             toType,
			ToAccountType:      toAccountType,
			Amount:             txamount,
			Date:               myDate,
			TransactionType:    txType,
			TransactionComment: txComment,
			PaymentMethod:      paymentMethod,
			PaymentGatewayTxID: paymentGatewayTxId,
			TransactionCharges: txCharges,
			PaymentStatus:      paymentStatus,
			SoshID:             soshId,
			ActivityID:         activityId,
			EntryID:            entryId,
			EntryName:          entryName,
			EntryType:          entryType,
			PlatformFees:       platformFees,
			FromEmail:          fromEmail,
			FromMobile:         fromMobile,
		}
		txAsBytes, _ := json.Marshal(tx)

		return ctx.GetStub().PutState(txid, txAsBytes)
	}
	// return fmt.Errorf("Something went wrong")
}

// QuerySingleTransaction Query Tx Data
func (s *SmartContract) QuerySingleTransaction(ctx contractapi.TransactionContextInterface, id string) (*Transactions, error) {
	txAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state")
	}

	if txAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	tx := new(Transactions)
	_ = json.Unmarshal(txAsBytes, tx)

	return tx, nil
}

// QueryAllTransactions query all tx from cocuch
func (s *SmartContract) QueryAllTransactions(ctx contractapi.TransactionContextInterface, richQuery string) ([]QueryResult, error) {
	//
	queryString := fmt.Sprintf(richQuery)
	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Enter valid Query")
	}
	defer resultIterator.Close()
	//
	results := []QueryResult{}
	//
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		//
		if err != nil {
			return nil, fmt.Errorf("Enter valid details")
		}

		tx := new(Transactions)
		_ = json.Unmarshal(queryResponse.Value, tx)

		queryResult := QueryResult{Record: tx}
		results = append(results, queryResult)
	}
	return results, nil
}

// GetTxDetailsViaIds query wallet data using the soshId
func (s *SmartContract) GetTxDetailsViaIds(ctx contractapi.TransactionContextInterface, userid string, soshid string, activityid string) ([]QueryResult, error) {

	userId := strings.ToLower(userid)
	soshId := strings.ToLower(soshid)
	activityId := strings.ToLower(activityid)

	queryString := fmt.Sprintf("{\"selector\":{\"fromId\":\"%s\",\"soshId\":\"%s\",\"activityId\":\"%s\"}}", userId, soshId, activityId)
	queryString2 := fmt.Sprintf("{\"selector\":{\"toId\":\"%s\",\"soshId\":\"%s\",\"activityId\":\"%s\"}}", userId, soshId, activityId)

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Enter valid wallet sosh id")
	}
	defer resultsIterator.Close()

	resultsIterator2, err := ctx.GetStub().GetQueryResult(queryString2)
	if err != nil {
		return nil, fmt.Errorf("Enter valid wallet sosh id")
	}
	defer resultsIterator2.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, fmt.Errorf("Enter valid id's")
		}

		tx := new(Transactions)
		_ = json.Unmarshal(queryResponse.Value, tx)

		queryResult := QueryResult{Record: tx}
		results = append(results, queryResult)
	}

	for resultsIterator2.HasNext() {
		queryResponse, err := resultsIterator2.Next()

		if err != nil {
			return nil, fmt.Errorf("Enter valid id's")
		}

		tx := new(Transactions)
		_ = json.Unmarshal(queryResponse.Value, tx)

		queryResult := QueryResult{Record: tx}
		results = append(results, queryResult)
	}

	return results, nil

}

// pagination
func pagination(ctx contractapi.TransactionContextInterface, queryString string, recordcount string, bookmark string) (*Page, error) {
	pageSize, err := strconv.ParseInt(recordcount, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error while converting pageSize to int")
	}
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, int32(pageSize), bookmark)
	if err != nil {
		return nil, fmt.Errorf("Error while getting iterator")
	}
	defer resultsIterator.Close()

	results := []Transactions{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, fmt.Errorf("Error while iterating records")
		}

		t := new(Transactions)
		_ = json.Unmarshal(queryResponse.Value, t)
		queryResult := Transactions{
			ObjectType:         t.ObjectType,
			TxID:               t.TxID,
			ReferenceID:        t.ReferenceID,
			MultiTxID:          t.MultiTxID,
			FromName:           t.FromName,
			FromID:             t.FromID,
			FromWalletID:       t.FromWalletID,
			FromType:           t.FromType,
			FromAccountType:    t.FromAccountType,
			ToName:             t.ToName,
			ToID:               t.ToID,
			ToWalletID:         t.ToWalletID,
			ToType:             t.ToType,
			ToAccountType:      t.ToAccountType,
			Amount:             t.Amount,
			Date:               t.Date,
			TransactionType:    t.TransactionType,
			TransactionComment: t.TransactionComment,
			PaymentMethod:      t.PaymentMethod,
			PaymentGatewayTxID: t.PaymentGatewayTxID,
			TransactionCharges: t.TransactionCharges,
			PaymentStatus:      t.PaymentStatus,
			SoshID:             t.SoshID,
			ActivityID:         t.ActivityID,
			EntryID:            t.EntryID,
			EntryName:          t.EntryName,
			EntryType:          t.EntryType,
			PlatformFees:       t.PlatformFees,
			FromEmail:          t.FromEmail,
			FromMobile:         t.FromMobile,
		}
		results = append(results, queryResult)

	}
	p := Page{Data: results, Count: responseMetadata.FetchedRecordsCount, Bookmark: responseMetadata.Bookmark}
	return &p, nil
}

// GetTxWithPage get tx
func (s *SmartContract) GetTxWithPage(ctx contractapi.TransactionContextInterface, recordcount string, bookmark string) (*Page, error) {
	queryString := `{
		"selector": {
			"docType":"tx",
			"date": {"$gt":null}
		},
		"sort": [
			{"date":"asc"}
		]
	}`
	p, err := pagination(ctx, queryString, recordcount, bookmark)
	return p, err
}

// GetTxBetweenDates get txs between two dates
func (s *SmartContract) GetTxBetweenDates(ctx contractapi.TransactionContextInterface, fromDate string, toDate string, userid string, soshid string, activityid string, recordcount string, bookmark string) (*Page, error) {
	queryString := fmt.Sprintf(`{
		"selector": {
		   "docType": "tx",
		   "$or": [
			   {
				   "fromId": "%s"
				},
				{
					"toId": "%s"
				}
			],
		   "soshId" : "%s",
		   "activityId" : "%s",
		   "date": {
			  "$gte": "%s",
			  "$lte": "%s"
		   }
		},
		"sort": [
		   {
			  "date": "asc"
		   }
		]
	}`, userid, userid, soshid, activityid, fromDate, toDate)
	p, err := pagination(ctx, queryString, recordcount, bookmark)
	return p, err
}

// FilterTx filter tx for user, activity and sosh
func (s SmartContract) FilterTx(ctx contractapi.TransactionContextInterface, filterBy string, fromDate string, toDate, soshid string, activityid string, recordcount string, bookmark string) (*Page, error) {

	if filterBy == "Individual" || filterBy == "Business" || filterBy == "Charity" {
		if fromDate == "" || toDate == "" {
			queryString := fmt.Sprintf(`{
				"selector": {
				   "docType": "tx",
				   "SoshID" : "%s",
				   "ActivityId" : "%s",
				   "$and": [
						{
						   "fromAccountType": "%s"
						},
						{
							"toAccountType": "%s"
						}
					]
				},
				"sort": [
				   {
					  "date": "asc"
				   }
				]
			}`, soshid, activityid, filterBy, filterBy)
			p, err := pagination(ctx, queryString, recordcount, bookmark)
			return p, err
		} else {
			queryString := fmt.Sprintf(`{
				"selector": {
				   "docType": "tx",
				   "SoshID" : "%s",
				   "ActivityId" : "%s",
				   "date": {
					"$gte": "%s",
					"$lte": "%s"
				   },
				   "$and": [
						{
						   "fromAccountType": "%s"
						},
						{
							"toAccountType": "%s"
						}
					]
				},
				"sort": [
				   {
					  "date": "asc"
				   }
				]
			}`, soshid, activityid, fromDate, toDate, filterBy, filterBy)
			p, err := pagination(ctx, queryString, recordcount, bookmark)
			return p, err
		}
	} else {
		if fromDate == "" || toDate == "" {
			queryString := fmt.Sprintf(`{
				"selector": {
				   "docType": "tx",
				   "SoshID" : "%s",
				   "ActivityId" : "%s",
				   "$or": [
						{
						   "fromAccountType": "%s"
						},
						{
							"toAccountType": "%s"
						}
					]
				},
				"sort": [
				   {
					  "date": "asc"
				   }
				]
			}`, soshid, activityid, filterBy, filterBy)
			p, err := pagination(ctx, queryString, recordcount, bookmark)
			return p, err
		} else {
			queryString := fmt.Sprintf(`{
				"selector": {
				   "docType": "tx",
				   "SoshID" : "%s",
				   "ActivityId" : "%s",
				   "date": {
					"$gte": "%s",
					"$lte": "%s"
				   },
				   "$or": [
						{
						   "fromAccountType": "%s"
						},
						{
							"toAccountType": "%s"
						}
					]
				},
				"sort": [
				   {
					  "date": "asc"
				   }
				]
			}`, soshid, activityid, fromDate, toDate, filterBy, filterBy)
			p, err := pagination(ctx, queryString, recordcount, bookmark)
			return p, err
		}
	}
}

// GetTransactions get txs with custom query
func (s *SmartContract) GetTransactions(ctx contractapi.TransactionContextInterface, richQuery string, recordcount string, bookmark string) (*Page, error) {
	queryString := fmt.Sprintf(richQuery)
	p, err := pagination(ctx, queryString, recordcount, bookmark)
	return p, err
}

// DynamicPagination pagination for admin panel screen to jump on any page
func (s *SmartContract) DynamicPagination(ctx contractapi.TransactionContextInterface, richQuery string, currentPage string, jumpPage string, recordcount string, bookmark string) (*Page, error) {
	queryString := fmt.Sprintf(richQuery)
	jumpPageInt, _ := strconv.Atoi(jumpPage)
	currentPageInt, _ := strconv.Atoi(currentPage)
	recordCountInt, _ := strconv.Atoi(recordcount)
	tempRecords := ((jumpPageInt - currentPageInt) - 1) * recordCountInt
	if tempRecords == 0 {
		p, err := pagination(ctx, queryString, recordcount, bookmark)
		return p, err
	}
	tempRecordsString := strconv.Itoa(tempRecords)
	tempPage, err := pagination(ctx, queryString, tempRecordsString, bookmark)
	p, err := pagination(ctx, queryString, recordcount, tempPage.Bookmark)
	return p, err
}

// DeleteTx delete tx from world state
func (s *SmartContract) DeleteTx(ctx contractapi.TransactionContextInterface, id string) error {

	A := id
	valAsbytes, err := ctx.GetStub().GetState(A) //get the transaction from chaincode state
	if err != nil {
		return fmt.Errorf("Error while getting tx")
	} else if valAsbytes == nil {
		return fmt.Errorf("Error while getting tx")
	}
	// // Delete the key from the state in ledger
	err1 := ctx.GetStub().DelState(A)
	if err1 != nil {
		return fmt.Errorf("Error while deleting transaction")
	}

	return fmt.Errorf("Something went wrong")
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
