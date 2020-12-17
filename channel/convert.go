package main

import (
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

func fromJSON(data Transaction) (string, error) {
	logWriter("New request Json to iso:8583")
	logWriter("original : " + fmt.Sprint(data))

	cardAcceptor := convCardAcc(data.CardAcceptorData)
	amount := strconv.Itoa(data.TotalAmount)

	val := map[int]string{
		2:  data.Pan,
		3:  data.ProcessingCode,
		4:  amount,
		5:  data.SettlementAmount,
		6:  data.CardholderBillingAmount,
		7:  data.TransmissionDateTime,
		9:  data.SettlementConversionrate,
		10: data.CardHolderBillingConvRate,
		11: data.Stan,
		12: data.LocalTransactionTime,
		13: data.LocalTransactionDate,
		17: data.CaptureDate,
		18: data.CategoryCode,
		22: data.PointOfServiceEntryMode,
		37: data.Refnum,
		41: data.CardAcceptorData.CardAcceptorTerminalId,
		43: cardAcceptor,
		48: data.AdditionalData,
		49: data.Currency,
		50: data.SettlementCurrencyCode,
		51: data.CardHolderBillingCurrencyCode,
		57: data.AdditionalDataNational,
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
