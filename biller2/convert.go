package main

import (
	"encoding/xml"
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

func toXML(iso string) {
	something := Spec{}
	nice := iso8583.NewISOStruct("../spec1987.yml", false)
	e := something.readFromFile("../spec1987.yml")

	if e != nil {
		fmt.Println(e.Error())
	}
	if len(iso) < 4 {
		fmt.Println("message seems incorrect")
		return
	}
	lnt, err := strconv.Atoi(iso[:4])

	if len(iso) != lnt+4 || err != nil {
		logWriter("New request ISO:8583 to JSON")
		logWriter("Incorrect format")
		logWriter(fmt.Sprintf("request for %s", iso))
		return
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
	payment := PaymentResponse{}
	payment.TransactionData.Pan = elm[2]
	payment.TransactionData.ProcessingCode = elm[3]
	payment.TransactionData.TotalAmount = amountTotal
	payment.TransactionData.TransmissionDateTime = elm[7]
	payment.TransactionData.LocalTransactionTime = elm[12]
	payment.TransactionData.LocalTransactionDate = elm[13]
	payment.TransactionData.CaptureDate = elm[17]
	payment.TransactionData.AdditionalData = elm[48]
	payment.TransactionData.Stan = elm[11]
	payment.TransactionData.Refnum = elm[37]
	payment.TransactionData.Currency = elm[49]
	payment.TransactionData.CategoryCode = elm[18]
	payment.TransactionData.SettlementAmount = elm[5]
	payment.TransactionData.CardholderBillingAmount = elm[6]
	payment.TransactionData.SettlementConversionrate = elm[9]
	payment.TransactionData.CardHolderBillingConvRate = elm[10]
	payment.TransactionData.PointOfServiceEntryMode = elm[22]
	payment.TransactionData.SettlementCurrencyCode = elm[50]
	payment.TransactionData.CardHolderBillingCurrencyCode = elm[51]
	payment.TransactionData.AdditionalDataNational = elm[57]
	payment.TransactionData.CardAcceptorData.CardAcceptorTerminalId = elm[41]
	if elm[43] != "" {
		payment.TransactionData.CardAcceptorData.CardAcceptorName = elm[43][:24]
		payment.TransactionData.CardAcceptorData.CardAcceptorCity = elm[43][25:38]
		payment.TransactionData.CardAcceptorData.CardAcceptorCountryCode = elm[43][38:40]
	}
	payment.ResponseStatus.ResponseCode = 200
	payment.ResponseStatus.ResponseDescription = "success"
	//fmt.Print(payment)
	resXML, err := xml.MarshalIndent(payment, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resXML))
}
