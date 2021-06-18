package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/shopspring/decimal"
)

type SmartContract struct {
	contractapi.Contract
}

//wallet structure
type Wallet struct {
	// MOBILENO      string `json:"mobileno"`      //owner no
	ID             string `json:"id"`            //mobileno
	CREATEDFROM    string `json:"createdfrom"`   //created from which app
	TYPE           string `json:"wallettype"`    //user/sosh/activity
	AMOUNT         string `json:"walletamount"`  //wallet balance
	COUPONS        int    `json:"coupons"`       //coupon codes
	UNSOLDTICKETS  int    `json:"unsoldtickets"` //unsoldtickets
	SOLDTICKETS    int    `json:"soldtickets"`   //soldtickets
	HOURS          int    `json:"hours"`         //hours traced
	GIFTCARD       int    `json:"giftcard"`      //Sponsorship
	ITEMS          int    `json:"items"`         //items
	VOLUNTEERHOURS int    `json: volunteerhours`
	STATUS         string `json:"status"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	// Key    string `json:"Key"`
	Record *Wallet
}

//Creating a wallet/Inserting data into wallet
func (s *SmartContract) CreateWallet(ctx contractapi.TransactionContextInterface, arrId []int, createdfrom string, wallettype string) (string, error) {
	//check if user exist
	for _, ids := range arrId {
		id := strconv.Itoa(ids)
		walletAsBytes, err := ctx.GetStub().GetState(id)

		if err != nil {
			return "", fmt.Errorf("Failed to read from world state")
		}
		if walletAsBytes == nil {
			wallet := Wallet{
				//MOBILENO:      mobileno,     //mobile no
				ID:             id,          //unique id
				CREATEDFROM:    createdfrom, //from which app
				TYPE:           wallettype,  // user/sosh/activity
				AMOUNT:         "0.0",       //wallet balance
				COUPONS:        0,           //coupon codes
				UNSOLDTICKETS:  0,           //unsold tickets
				SOLDTICKETS:    0,           // sold tickets
				HOURS:          0,           //hours
				GIFTCARD:       0,           //sponsorship
				ITEMS:          0,           //items donated/shared
				VOLUNTEERHOURS: 0,
				STATUS:         "ACTIVE",
			}
			walletAsBytes, _ := json.Marshal(wallet)
			err := ctx.GetStub().PutState(id, walletAsBytes)
			if err != nil {
				return "", fmt.Errorf("Failed to PutState into Network with error: " + err.Error())
			}
			// return ("Wallet with ID: " + id + " is created"), nil
		} else {
			wallet, err := s.QueryWalletData(ctx, id)
			if err != nil {
				return "", fmt.Errorf(" Failed to get wallet details with error: " + err.Error())
			} else if id == wallet.ID {
				return "", fmt.Errorf("User with ID: " + id + " already exist")
			} else {
				wallet := Wallet{
					//MOBILENO:      mobileno,     //mobile no
					ID:             id,          //unique id
					CREATEDFROM:    createdfrom, //from which app
					TYPE:           wallettype,  // user/sosh/activity
					AMOUNT:         "0.0",       //wallet balance
					COUPONS:        0,           //coupon codes
					UNSOLDTICKETS:  0,           //unsold tickets
					SOLDTICKETS:    0,           // sold tickets
					HOURS:          0,           //hours
					GIFTCARD:       0,           //sponsorship
					ITEMS:          0,           //items donated/shared
					VOLUNTEERHOURS: 0,
					STATUS:         "ACTIVE",
				}
				walletAsBytes, _ := json.Marshal(wallet)
				err := ctx.GetStub().PutState(id, walletAsBytes)

				if err != nil {
					return "", fmt.Errorf("Failed to PutState into Network with error: " + err.Error())
				}
				// return ("Wallet with ID: " + id + " is created"), nil
			}
		}
	}
	return "Wallet created", nil
}

//Query Wallet Data
func (s *SmartContract) QueryWalletData(ctx contractapi.TransactionContextInterface, id string) (*Wallet, error) {
	walletAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if walletAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	wallet := new(Wallet)
	_ = json.Unmarshal(walletAsBytes, wallet)

	return wallet, nil
}

func (s *SmartContract) QueryWallets(ctx contractapi.TransactionContextInterface, ids string) ([]QueryResult, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"id\":{\"$in\":[%s]}}}", ids)
	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Enter valid wallet id's")
	}
	defer resultIterator.Close()

	results := []QueryResult{}

	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()

		if err != nil {
			return nil, fmt.Errorf("Enter valid wallet id's")
		}

		wallet := new(Wallet)
		_ = json.Unmarshal(queryResponse.Value, wallet)

		queryResult := QueryResult{Record: wallet}
		results = append(results, queryResult)
	}
	return results, nil
}

// Update Amount in wallet
func (s *SmartContract) UpdateWalletDetails(ctx contractapi.TransactionContextInterface, id string, fieldname string, value string) error {

	//reflect is a package used by go for updating struct field
	//more details : https://samwize.com/2015/03/20/how-to-use-reflect-to-set-a-struct-field/

	wallet, err := s.QueryWalletData(ctx, id)
	// var str string
	// json.Unmarshal([]byte(wallet), &str)

	v := reflect.ValueOf(wallet).Elem().FieldByName(fieldname)
	if v.IsValid() {
		v.SetString(value)
	}

	// for i := 0; i < len(fieldname); i++ {
	// 	field := strings.ToUpper(fieldname[i])
	// 	v := reflect.ValueOf(wallet).Elem().FieldByName(field)
	// 	if v.IsValid() {
	// 		// v.SetString(value)
	// 		v.SetInt(value[i])
	// 	}
	// }

	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}

	walletAsBytes, _ := json.Marshal(wallet)

	return ctx.GetStub().PutState(id, walletAsBytes)
}

//Adding Funds to wallet
func (s *SmartContract) AddFunds(ctx contractapi.TransactionContextInterface, id string, amount string) (map[string]string, error) {
	response := make(map[string]string)
	if s.isWalletDeleted(ctx, id) {
		response[id] = "Wallet was deleted already"
		return response, nil
	}
	wallet, err := s.QueryWalletData(ctx, id)
	amt, err := decimal.NewFromString(amount)
	walletAmt, err := decimal.NewFromString(wallet.AMOUNT)
	newBal := amt.Add(walletAmt)

	if err != nil {
		return nil, fmt.Errorf("Error before PutState with error: " + err.Error())
	}

	wallet.AMOUNT = newBal.String()
	walletAsBytes, _ := json.Marshal(wallet)

	err1 := ctx.GetStub().PutState(id, walletAsBytes)
	if err1 != nil {
		return nil, fmt.Errorf("Failed to PutState into Network with error: " + err1.Error())
	}

	response[id] = newBal.String()
	return response, nil
}

//withdraw Funds from wallet
func (s *SmartContract) WithdrawFunds(ctx contractapi.TransactionContextInterface, id string, amount string) (map[string]string, error) {
	response := make(map[string]string)
	if s.isWalletDeleted(ctx, id) {
		response[id] = "Wallet was deleted already"
		return response, nil
	}
	wallet, err := s.QueryWalletData(ctx, id)
	amt, err := decimal.NewFromString(amount)
	walletAmt, err := decimal.NewFromString(wallet.AMOUNT)

	if err != nil {
		return nil, fmt.Errorf("Error while querying wallet data")
	}

	if amt.GreaterThan(walletAmt) == true || walletAmt.IsZero() == true || walletAmt.IsNegative() == true {
		return nil, fmt.Errorf("Insufficient balance of user " + id + " Available Balance: " + walletAmt.String())
	}
	newBal := walletAmt.Sub(amt)
	wallet.AMOUNT = newBal.String()
	walletAsBytes, _ := json.Marshal(wallet)

	err1 := ctx.GetStub().PutState(id, walletAsBytes)
	if err1 != nil {
		return nil, fmt.Errorf("Failed to PutState with error: " + err1.Error())
	}
	response[id] = newBal.String()
	return response, nil
}

//Transfering funds Wallet to Wallet
func (s *SmartContract) TransferFunds(ctx contractapi.TransactionContextInterface, fromId string, toData map[string]string) (map[string]string, error) {
	response := make(map[string]string)
	if s.isWalletDeleted(ctx, fromId) {
		response[fromId] = "Wallet was deleted already"
		return response, nil
	}

	var totalAmt decimal.Decimal
	for _, v := range toData {
		val, _ := decimal.NewFromString(v)
		totalAmt = totalAmt.Add(val)
	}

	fmt.Println("totalAmt of Amount %s", totalAmt.String())
	wallet1, err := s.QueryWalletData(ctx, fromId)

	if err != nil {
		return nil, fmt.Errorf("Failed To Get WalletDetails of Wallet: " + fromId + " with errors: " + err.Error())
	}

	wallet1Amt, err := decimal.NewFromString(wallet1.AMOUNT)

	if wallet1Amt.IsNegative() == true || wallet1Amt.IsZero() == true || totalAmt.GreaterThan(wallet1Amt) == true {
		return nil, fmt.Errorf("Insufficient balance of Wallet id: " + fromId + " Please Add Funds to the wallet, Available Balance is: " + wallet1Amt.String() + " & Required Amount is: " + totalAmt.String())
	}

	for toId, amount := range toData {
		fmt.Printf("key[%s] value[%s]\n", toId, amount)
		wallet2, err1 := s.QueryWalletData(ctx, toId)

		if err1 != nil {
			return nil, fmt.Errorf("Failed to Get Wallet details of wallet: " + toId + " with error: " + err1.Error())
		}

		newAmt, _ := decimal.NewFromString(amount)
		wallet2Amt, _ := decimal.NewFromString(wallet2.AMOUNT)
		fmt.Println("Adding amount: ", newAmt.String(), " having balance wallet2Amt: ", wallet2Amt.String())
		newWallet2Amt := wallet2Amt.Add(newAmt)
		fmt.Println("Added amount in Wallet ", newWallet2Amt.String())

		wallet2.AMOUNT = newWallet2Amt.String()
		wallet2AsBytes, _ := json.Marshal(wallet2)
		fmt.Println("PutState of Wallet: ", toId, " with amount: ", amount)
		err2 := ctx.GetStub().PutState(toId, wallet2AsBytes)
		if err2 != nil {
			return nil, fmt.Errorf("Failed to PutState of wallet: " + toId + " with error: " + err2.Error())
		}
		response[toId] = newWallet2Amt.String()
	}
	newWallet1Amt := wallet1Amt.Sub(totalAmt)
	wallet1.AMOUNT = newWallet1Amt.String()
	wallet1AsBytes, _ := json.Marshal(wallet1)
	err3 := ctx.GetStub().PutState(fromId, wallet1AsBytes)
	if err3 != nil {
		return nil, fmt.Errorf("Failed to PutState of wallet: " + fromId + "with error: " + err3.Error())
	}
	response[fromId] = newWallet1Amt.String()
	return response, nil
}

//Delete wallet
func (s *SmartContract) DeleteWallet(ctx contractapi.TransactionContextInterface, id []int) (string, error) {
	for _, item := range id {
		id := strconv.Itoa(item)
		valAsBytes, err := ctx.GetStub().GetState(id)

		if err != nil {
			return "", fmt.Errorf("Failed to GetState with error: " + err.Error())
		} else if valAsBytes == nil {
			return "", fmt.Errorf("valAsbytes are not nil ")
		}
		wallet, err := s.QueryWalletData(ctx, id)

		if err != nil {
			return "", fmt.Errorf("Failed to Query wallet with error:" + err.Error())
		}
		// Delete the key from the state in ledger
		walletAmt, err := decimal.NewFromString(wallet.AMOUNT)
		if walletAmt.IsZero() == true {
			s.UpdateWalletDetails(ctx, id, "STATUS", "DELETED")
			if err != nil {
				return "", fmt.Errorf("Failed to DelState with error:" + err.Error())
			}
			fmt.Sprintf("Wallet with ID: "+id+" is deleted", nil)
		} else {
			return "", fmt.Errorf("To delete wallet balance must be zero")
		}
	}
	return "Wallet deleted", nil

}

//AddTime to wallet
func (s *SmartContract) AddTime(ctx contractapi.TransactionContextInterface, id string, hours int) error {
	wallet, err := s.QueryWalletData(ctx, id)

	newBal := hours + wallet.HOURS
	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}

	wallet.HOURS = newBal
	walletAsBytes, _ := json.Marshal(wallet)

	return ctx.GetStub().PutState(id, walletAsBytes)
}

//AddTime to wallet
func (s *SmartContract) AddMultipleTimes(ctx contractapi.TransactionContextInterface, id1 string, id2 string, id3 string, hours int) error {
	wallet1, err := s.QueryWalletData(ctx, id1)
	wallet2, err := s.QueryWalletData(ctx, id2)
	wallet3, err := s.QueryWalletData(ctx, id3)

	newBal1 := hours + wallet1.HOURS
	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}
	newBal2 := hours + wallet2.HOURS
	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}
	newBal3 := hours + wallet3.HOURS
	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}
	wallet1.HOURS = newBal1
	walletAsBytes1, _ := json.Marshal(wallet1)
	wallet2.HOURS = newBal2
	walletAsBytes2, _ := json.Marshal(wallet2)
	wallet3.HOURS = newBal3
	walletAsBytes3, _ := json.Marshal(wallet3)

	err1 := ctx.GetStub().PutState(id1, walletAsBytes1)
	if err1 != nil {
		return fmt.Errorf("Error while updating wallet data")
	}
	err2 := ctx.GetStub().PutState(id2, walletAsBytes2)
	if err2 != nil {
		return fmt.Errorf("Error while updating wallet data")
	}
	err3 := ctx.GetStub().PutState(id3, walletAsBytes3)
	if err3 != nil {
		return fmt.Errorf("Error while updating wallet data")
	}
	return fmt.Errorf("Error while querying wallet data")
}

//
func (s *SmartContract) IsWalletEmpty(ctx contractapi.TransactionContextInterface, id string) error {
	wallet, err := s.QueryWalletData(ctx, id)

	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}

	walletamount, err := decimal.NewFromString(wallet.AMOUNT)
	if walletamount.IsZero() == true {
		return fmt.Errorf("Insufficient balance of user " + id + " Available Balance: " + walletamount.String())
	}

	return fmt.Errorf("Error while converting decinal to string")
}

// add volunteer hours in two wallets
func (s *SmartContract) AddVolunteerHours(ctx contractapi.TransactionContextInterface, id1 string, id2 string, hours int) error {
	wallet1, err := s.QueryWalletData(ctx, id1)
	wallet2, err := s.QueryWalletData(ctx, id2)

	newBal1 := hours + wallet1.VOLUNTEERHOURS
	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}
	newBal2 := hours + wallet2.VOLUNTEERHOURS
	if err != nil {
		return fmt.Errorf("Error while querying wallet data")
	}

	wallet1.HOURS = newBal1
	walletAsBytes1, _ := json.Marshal(wallet1)
	wallet2.HOURS = newBal2
	walletAsBytes2, _ := json.Marshal(wallet2)

	err1 := ctx.GetStub().PutState(id1, walletAsBytes1)
	if err1 != nil {
		return fmt.Errorf("Error while updating wallet data")
	}
	err2 := ctx.GetStub().PutState(id2, walletAsBytes2)
	if err2 != nil {
		return fmt.Errorf("Error while updating wallet data")
	}

	return fmt.Errorf("Error while querying wallet data")
}

// QueryAllWallets give all wallet details
func (s *SmartContract) QueryAllWallets(ctx contractapi.TransactionContextInterface) ([]Wallet, error) {
	queryString := fmt.Sprintf("{\"selector\":{}}}")
	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("Enter valid details")
	}
	defer resultIterator.Close()
	results := []Wallet{}
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Enter valid details")
		}

		w := new(Wallet)
		_ = json.Unmarshal(queryResponse.Value, w)
		results = append(results, *w)
	}
	return results, nil
}

// GetAvgBalance totalAmt of amount in all wallet / total no of wallets
func (s *SmartContract) GetAvgBalance(ctx contractapi.TransactionContextInterface) (string, error) {
	queryString := fmt.Sprintf("{\"selector\":{}}}")
	resultIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "0", fmt.Errorf("Enter valid details")
	}
	defer resultIterator.Close()
	results := []Wallet{}
	var totalBalance decimal.Decimal
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return "0", fmt.Errorf("Enter valid details")
		}

		w := new(Wallet)
		_ = json.Unmarshal(queryResponse.Value, w)
		walletAmt, err := decimal.NewFromString(w.AMOUNT)
		totalBalance = totalBalance.Add(walletAmt)
		results = append(results, *w)
	}
	fmt.Println(results)
	x := decimal.NewFromFloat32(float32(len(results)))
	return (totalBalance.Div(x).String()), nil
}

func (s *SmartContract) isWalletDeleted(ctx contractapi.TransactionContextInterface, walletId string) bool {
	wallet, _ := s.QueryWalletData(ctx, walletId)
	if wallet.STATUS == "DELETED" {
		return true
	} else {
		return false
	}
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
