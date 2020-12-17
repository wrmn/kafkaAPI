package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/mofax/iso8583"
)

func toJSON(data Transaction) (string, error) {
	logWriter("New request Json to iso:8583")
	logWriter("original : " + fmt.Sprint(data))
	iso := iso8583.NewISOStruct("spec1987.yml", false)
	var cardAcceptor string
	if data.CardAcceptorData.CardAcceptorCity != "" ||
		data.CardAcceptorData.CardAcceptorCountryCode != "" ||
		data.CardAcceptorData.CardAcceptorName != "" {
		for len(data.CardAcceptorData.CardAcceptorCity) < 13 {
			data.CardAcceptorData.CardAcceptorCity += " "
		}
		for len(data.CardAcceptorData.CardAcceptorName) < 25 {
			data.CardAcceptorData.CardAcceptorName += " "
		}

		cardAcceptor = data.CardAcceptorData.CardAcceptorName +
			data.CardAcceptorData.CardAcceptorCity +
			data.CardAcceptorData.CardAcceptorCountryCode
	}
	amount := strconv.Itoa(data.TotalAmount)
	something := Spec{}
	e := something.readFromFile("spec1987.yml")
	if e != nil {
		fmt.Println(e.Error())
	}
	val := map[int]string{2: data.Pan,
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
	iso.AddMTI("0200")

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
			logWriter(missing)
			logWriter("request aborted")
			return "", errors.New(missing)
		}

		if val[id] != "" {
			iso.AddField(int64(id), val[id])
		}

	}

	result, _ := iso.ToString()
	lnth := strconv.Itoa(len(result))
	for len(lnth) < 4 {
		lnth = "0" + lnth
	}

	mti := result[:4]
	res := result[4:20]
	ele := result[20:]
	bitmap, _ := iso8583.HexToBitmapArray(res)
	logWriter("Full message	: " + lnth + result)
	logWriter("Length		: " + lnth)
	logWriter("Msg Only		: " + result)
	logWriter("MTI			: " + mti)
	logWriter("Hexmap		: " + res)
	logWriter("Bitmap		: " + fmt.Sprintf("%d", bitmap))
	logWriter("Element		: " + ele)
	return lnth + result, nil

}
