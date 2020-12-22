package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/go-yaml/yaml"
	"github.com/mofax/iso8583"
)

func (s *Spec) readFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	yaml.Unmarshal(content, &s.fields)
	return nil
}

func toJson(iso string) string {
	something := Spec{}
	nice := iso8583.NewISOStruct("../spec1987.yml", false)
	e := something.readFromFile("../spec1987.yml")

	if e != nil {
		fmt.Println(e.Error())
	}
	if len(iso) < 4 {
		fmt.Println("message seems incorrect")
		return ""
	}
	lnt, err := strconv.Atoi(iso[:4])

	if len(iso) != lnt+4 || err != nil {
		logWriter("New request ISO:8583 to JSON")
		logWriter("Incorrect format")
		logWriter(fmt.Sprintf("request for %s", iso))
		return ""
	}

	mti := iso[4:8]
	res := iso[8:24]
	ele := iso[24:]
	bitmap, _ := iso8583.HexToBitmapArray(res)

	logWriter("New request ISO:8583 to Json")
	logWriter("Full message	: " + iso)
	logWriter("Length		: " + iso[:4])
	logWriter("Msg Only		: " + iso[4:])
	logWriter("MTI			: " + mti)
	logWriter("Hexmap		: " + res)
	logWriter("Bitmap		: " + fmt.Sprintf("%d", bitmap))
	logWriter("Element		: " + ele)

	tlen := len(ele)
	mark := 0

	nice.AddMTI(mti)
	nice.Bitmap = bitmap
	for idx := range bitmap {
		if bitmap[idx] == 1 {
			element := something.fields[idx+1]
			len := element.MaxLen
			if element.LenType == "llvar" {
				clen, _ := strconv.Atoi(ele[mark : mark+2])
				msg := fmt.Sprintf("[%d] length %d = %s", idx, clen, ele[mark+2:mark+clen+2])
				logWriter(msg)
				nice.AddField(int64(idx+1), ele[mark+2:mark+clen+2])
				tlen -= clen + 2
				mark += clen + 2
			} else if element.LenType == "lllvar" {
				clen, _ := strconv.Atoi(ele[mark : mark+3])
				msg := fmt.Sprintf("[%d] length %d =  %s", idx, clen, ele[mark+3:mark+clen+3])
				logWriter(msg)
				nice.AddField(int64(idx+1), ele[mark+3:mark+clen+3])
				tlen -= clen + 3
				mark += clen + 3
			} else {
				msg := fmt.Sprintf("[%d] length %d = %s", idx, len, ele[mark:mark+len])
				logWriter(msg)
				nice.AddField(int64(idx+1), ele[mark:mark+len])
				tlen -= len
				mark += len
			}
		}
	}
	elm := nice.Elements.GetElements()

	amountTotal, _ := strconv.Atoi(elm[4])
	payment := Transaction{}
	payment.Pan = elm[2]
	payment.ProcessingCode = elm[3]
	payment.TotalAmount = amountTotal
	payment.TransmissionDateTime = elm[7]
	payment.LocalTransactionTime = elm[12]
	payment.LocalTransactionDate = elm[13]
	payment.CaptureDate = elm[17]
	payment.AdditionalData = elm[48]
	payment.Stan = elm[11]
	payment.Refnum = elm[37]
	payment.Currency = elm[49]
	payment.CategoryCode = elm[18]
	payment.SettlementAmount = elm[5]
	payment.CardholderBillingAmount = elm[6]
	payment.SettlementConversionrate = elm[9]
	payment.CardHolderBillingConvRate = elm[10]
	payment.PointOfServiceEntryMode = elm[22]
	payment.SettlementCurrencyCode = elm[50]
	payment.CardHolderBillingCurrencyCode = elm[51]
	payment.AdditionalDataNational = elm[57]
	payment.CardAcceptorData.CardAcceptorTerminalId = elm[41]
	if elm[43] != "" {
		payment.CardAcceptorData.CardAcceptorName = elm[43][:24]
		payment.CardAcceptorData.CardAcceptorCity = elm[43][25:38]
		payment.CardAcceptorData.CardAcceptorCountryCode = elm[43][38:40]
	}
	//fmt.Print(payment)
	//json.NewEncoder(w).Encode(payment)
	resJson, err := json.MarshalIndent(payment, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(resJson)
}

func convCardAcc(cardAcceptorData CardAcceptorData) string {
	var cardAcceptor string
	if cardAcceptorData.CardAcceptorCity != "" ||
		cardAcceptorData.CardAcceptorCountryCode != "" ||
		cardAcceptorData.CardAcceptorName != "" {
		for len(cardAcceptorData.CardAcceptorCity) < 13 {
			cardAcceptorData.CardAcceptorCity += " "
		}
		for len(cardAcceptorData.CardAcceptorName) < 25 {
			cardAcceptorData.CardAcceptorName += " "
		}

		cardAcceptor = cardAcceptorData.CardAcceptorName +
			cardAcceptorData.CardAcceptorCity +
			cardAcceptorData.CardAcceptorCountryCode
	}
	return cardAcceptor
}

func resultLog(result string) {
	lnth := result[:4]
	mti := result[4:8]
	res := result[8:24]
	ele := result[24:]
	bitmap, _ := iso8583.HexToBitmapArray(res)
	logWriter("Full message	: " + result)
	logWriter("Length		: " + lnth)
	logWriter("Msg Only		: " + result[4:])
	logWriter("MTI			: " + mti)
	logWriter("Hexmap		: " + res)
	logWriter("Bitmap		: " + fmt.Sprintf("%d", bitmap))
	logWriter("Element		: " + ele)

}

func toISO(val map[int]string) (string, error) {
	iso := iso8583.NewISOStruct("../spec1987.yml", false)
	iso.AddMTI("0200")

	something := Spec{}

	e := something.readFromFile("../spec1987.yml")
	if e != nil {
		fmt.Println(e.Error())
	}

	for id := range something.fields {
		ele := something.fields[id]
		if ele.LenType == "fixed" && val[id] != "" {
			if id == 4 {
				for len(val[id]) < ele.MaxLen {
					val[id] = "0" + val[id]
				}
			} else {
				for len(val[id]) < ele.MaxLen {
					val[id] = val[id] + " "
				}
			}
			if len(val[id]) > ele.MaxLen {
				val[id] = val[id][:ele.MaxLen]
			}
			logWriter(fmt.Sprintf("[%d] length %d = %s", id, ele.MaxLen, val[id]))
		} else if val[id] != "" {
			logWriter(fmt.Sprintf("[%d] length %d = %s", id, len(val[id]), val[id]))
		}

		if ele.ContentType == "m" && val[id] == "" {
			missing := fmt.Sprintf("mandatory field required \n%s is empty", ele.Label)
			logWriter("request aborted")
			return "", errors.New(missing)
		}

		if val[id] != "" {
			iso.AddField(int64(id), val[id])
		}

	}

	result, _ := iso.ToString()

	return result, nil

}

func fromJSON(data PaymentResponse) (string, error) {
	logWriter("New request Json to iso:8583")
	logWriter("original : " + fmt.Sprint(data))

	cardAcceptor := convCardAcc(data.TransactionData.CardAcceptorData)
	amount := strconv.Itoa(data.TransactionData.TotalAmount)
	responseStatus := convResp(data.ResponseStatus)

	val := map[int]string{
		2:  data.TransactionData.Pan,
		3:  data.TransactionData.ProcessingCode,
		4:  amount,
		5:  data.TransactionData.SettlementAmount,
		6:  data.TransactionData.CardholderBillingAmount,
		7:  data.TransactionData.TransmissionDateTime,
		9:  data.TransactionData.SettlementConversionrate,
		10: data.TransactionData.CardHolderBillingConvRate,
		11: data.TransactionData.Stan,
		12: data.TransactionData.LocalTransactionTime,
		13: data.TransactionData.LocalTransactionDate,
		17: data.TransactionData.CaptureDate,
		18: data.TransactionData.CategoryCode,
		22: data.TransactionData.PointOfServiceEntryMode,
		37: data.TransactionData.Refnum,
		39: responseStatus,
		41: data.TransactionData.CardAcceptorData.CardAcceptorTerminalId,
		43: cardAcceptor,
		48: data.TransactionData.AdditionalData,
		49: data.TransactionData.Currency,
		50: data.TransactionData.SettlementCurrencyCode,
		51: data.TransactionData.CardHolderBillingCurrencyCode,
		57: data.TransactionData.AdditionalDataNational,
	}

	result, err := toISO(val)

	lnth := strconv.Itoa(len(result))
	for len(lnth) < 4 {
		lnth = "0" + lnth
	}
	finResult := lnth + result

	if err != nil {
		return finResult, err
	}

	resultLog(finResult)

	return finResult, nil
}

func convResp(resp Response) string {
	response := fmt.Sprintf("%d%d%s", resp.ResponseCode, resp.ReasonCode, resp.ResponseDescription)
	res := len(response)
	return fmt.Sprintf("%d%s", res, response)
}
