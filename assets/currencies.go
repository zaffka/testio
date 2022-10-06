package assets

import (
	_ "embed"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

//go:embed dictionaries/currencies.json
var currenciesJSON []byte

//go:embed dictionaries/crypto.json
var cryptoJSON []byte

var (
	json             = jsoniter.ConfigCompatibleWithStandardLibrary
	Currencies       []Currency
	CryptoCurrencies []CryptoCurrency

	currenciesListedByCode   map[string]*Currency
	currenciesListedByName   map[string]*Currency
	currenciesListedByIsoNum map[uint16]*Currency

	cryptoCurrenciesListedByCode map[string]*CryptoCurrency
	cryptoCurrenciesListedByName map[string]*CryptoCurrency
)

type Currency struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	ISONum    uint16 `json:"iso_num"`
	Precision uint8  `json:"precision"`
}

type CryptoCurrency struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Precision uint8  `json:"precision"`
}

func init() {
	err := json.Unmarshal(currenciesJSON, &Currencies)
	if err != nil {
		panic(err.Error())
	}

	l := len(Currencies)
	currenciesListedByName = make(map[string]*Currency, l)
	currenciesListedByCode = make(map[string]*Currency, l)
	currenciesListedByIsoNum = make(map[uint16]*Currency, l)

	for i, c := range Currencies {
		currenciesListedByName[strings.ToUpper(c.Name)] = &Currencies[i]
		currenciesListedByCode[strings.ToUpper(c.Code)] = &Currencies[i]
		currenciesListedByIsoNum[c.ISONum] = &Currencies[i]
	}

	err = json.Unmarshal(cryptoJSON, &CryptoCurrencies)
	if err != nil {
		panic(err.Error())
	}

	l = len(CryptoCurrencies)
	cryptoCurrenciesListedByCode = make(map[string]*CryptoCurrency, l)
	cryptoCurrenciesListedByName = make(map[string]*CryptoCurrency, l)

	for i, c := range CryptoCurrencies {
		cryptoCurrenciesListedByCode[strings.ToUpper(c.Code)] = &CryptoCurrencies[i]
		cryptoCurrenciesListedByName[strings.ToUpper(c.Name)] = &CryptoCurrencies[i]
	}
}

func CurrenciesListLen() int {
	return len(Currencies)
}

func CurrencyByName(name string) (*Currency, bool) {
	c, exist := currenciesListedByName[strings.ToUpper(name)]

	return c, exist
}

func CurrencyByCode(code string) (*Currency, bool) {
	c, exist := currenciesListedByCode[strings.ToUpper(code)]

	return c, exist
}

func CurrencyByISONum(isoNum uint16) (*Currency, bool) {
	c, exist := currenciesListedByIsoNum[isoNum]

	return c, exist
}

func CryptoCurrenciesListLen() int {
	return len(CryptoCurrencies)
}

func CryptoCurrencyByName(name string) (*CryptoCurrency, bool) {
	c, exist := cryptoCurrenciesListedByName[strings.ToUpper(name)]

	return c, exist
}

func CryptoCurrencyByCode(code string) (*CryptoCurrency, bool) {
	c, exist := cryptoCurrenciesListedByCode[strings.ToUpper(code)]

	return c, exist
}
