package assets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zaffka/testio/assets"
)

func TestCurrenciesList(t *testing.T) {
	assert.Equal(t, 3, assets.CurrenciesListLen())

	testCurrency := assets.Currency{
		Code:      "EUR",
		Name:      "Euro",
		ISONum:    978,
		Precision: 2,
	}

	c, ok := assets.CurrencyByName("Euro")
	assert.True(t, ok)
	assert.Equal(t, &testCurrency, c)

	c, ok = assets.CurrencyByCode("Eur")
	assert.True(t, ok)
	assert.Equal(t, &testCurrency, c)

	c, ok = assets.CurrencyByCode("EUR")
	assert.True(t, ok)
	assert.Equal(t, &testCurrency, c)

	c, ok = assets.CurrencyByISONum(978)
	assert.True(t, ok)
	assert.Equal(t, &testCurrency, c)

	c, ok = assets.CurrencyByName("XXXX")
	assert.False(t, ok)
	assert.Nil(t, c)

	c, ok = assets.CurrencyByCode("XXXX")
	assert.False(t, ok)
	assert.Nil(t, c)

	c, ok = assets.CurrencyByISONum(999)
	assert.False(t, ok)
	assert.Nil(t, c)
}
