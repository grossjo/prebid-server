package jsonutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDropElement(t *testing.T) {

	tests := []struct {
		description     string
		input           []byte
		elementToRemove string
		output          []byte
		errorExpected   bool
		errorContains   string
	}{
		{
			description:     "Drop Single Element After Another Element",
			input:           []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"test": 1,"consented_providers": [1608,765,492]}}`),
			elementToRemove: "consented_providers",
			output:          []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"test": 1}}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Single Element Before Another Element",
			input:           []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"consented_providers": [1608,765,492],"test": 1}}`),
			elementToRemove: "consented_providers",
			output:          []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"test": 1}}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Single Element",
			input:           []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"consented_providers": [1545,2563,1411]}}`),
			elementToRemove: "consented_providers",
			output:          []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {}}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Single Element string",
			input:           []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"consented_providers": "test"}}`),
			elementToRemove: "consented_providers",
			output:          []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {}}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Parent Element Between Two Elements",
			input:           []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"consented_providers": [1608,765,492], "test": 1},"test": 123}`),
			elementToRemove: "consented_providers_settings",
			output:          []byte(`{"consent": "TESTCONSENT","test": 123}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Parent Element Before Element",
			input:           []byte(`{"consented_providers_settings": {"consented_providers": [1608,765,492], "test": 1},"test": 123}`),
			elementToRemove: "consented_providers_settings",
			output:          []byte(`{"test": 123}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Parent Element After Element",
			input:           []byte(`{"consent": "TESTCONSENT","consented_providers_settings": {"consented_providers": [1608,765,492], "test": 1}}`),
			elementToRemove: "consented_providers_settings",
			output:          []byte(`{"consent": "TESTCONSENT"}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Parent Element Only",
			input:           []byte(`{"consented_providers_settings": {"consented_providers": [1608,765,492], "test": 1}}`),
			elementToRemove: "consented_providers_settings",
			output:          []byte(`{}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Parent Element List",
			input:           []byte(`{"consented_providers_settings":{"consented_providers":[1608,765,492],"test":1},"data": [{"test1":5},{"test2": [1,2,3]}]}`),
			elementToRemove: "data",
			output:          []byte(`{"consented_providers_settings":{"consented_providers":[1608,765,492],"test":1}}`),
			errorExpected:   false,
			errorContains:   "",
		},
		{
			description:     "Drop Element That Doesn't Exist",
			input:           []byte(`{"consented_providers_settings": {"consented_providers": [1608,765,492], "test": 1}}`),
			elementToRemove: "test2",
			output:          []byte(`{"consented_providers_settings": {"consented_providers": [1608,765,492], "test": 1}}`),
			errorExpected:   false,
			errorContains:   "",
		},
		//Errors
		{
			description:     "Error Decode",
			input:           []byte(`{"consented_providers_settings": {"consented_providers": ["123",1,,1365,5678,1545,2563,1411], "test": 1}}`),
			elementToRemove: "consented_providers",
			output:          []byte(``),
			errorExpected:   true,
			errorContains:   "looking for beginning of value",
		},
		{
			description:     "Error Malformed",
			input:           []byte(`{consented_providers_settings: {"consented_providers": [1365,5678,1545,2563,1411], "test": 1}}`),
			elementToRemove: "consented_providers",
			output:          []byte(``),
			errorExpected:   true,
			errorContains:   "invalid character",
		},
	}

	for _, tt := range tests {
		res, err := DropElement(tt.input, tt.elementToRemove)

		if tt.errorExpected {
			assert.Error(t, err, "Error should not be nil")
			assert.True(t, strings.Contains(err.Error(), tt.errorContains))
		} else {
			assert.NoError(t, err, "Error should be nil")
			assert.Equal(t, tt.output, res, "Result is incorrect")
		}

	}
}

func TestDropNestedElementSingleOccurrence(t *testing.T) {
	input := []byte(`{"consented_providers_settings":{"consented_providers":[1608,765,492],"test":1},"data": [{"test1":5},{"test2": [1,2,3]}]}`)
	output := []byte(`{"consented_providers_settings":{"consented_providers":[1608,765,492]},"data": [{"test1":5},{"test2": [1,2,3]}]}`)

	res, err := DropElement(input, "consented_providers_settings", "test")

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}

func TestDropNestedElementMultipleOccurrence(t *testing.T) {
	input := []byte(`{"consented_providers_settings":{"consented_providers":[1608,765,492],"test":1},"data": [{"test":5},{"test": [1,2,3]}]}`)
	output := []byte(`{"consented_providers_settings":{"consented_providers":[1608,765,492]},"data": [{"test":5},{"test": [1,2,3]}]}`)

	res, err := DropElement(input, "consented_providers_settings", "test")

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}

func TestDropNestedStructureSingleOccurrence(t *testing.T) {
	input := []byte(`{"consented_providers":{"providers":[1608,765,492],"test":{"nested":true}},"data": [{"test":5},{"test": [1,2,3]}]}`)
	output := []byte(`{"consented_providers":{"providers":[1608,765,492]},"data": [{"test":5},{"test": [1,2,3]}]}`)

	res, err := DropElement(input, "consented_providers", "test")

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}

func TestDropNestedStructureSingleOccurrenceDeepNested(t *testing.T) {
	input := []byte(`{"consented_providers":{"providers":[1608,765,492],"test":{"nested":true, "nested2": {"test6": 123}}},"data": [{"test":5},{"test": [1,2,3]}]}`)
	output := []byte(`{"consented_providers":{"providers":[1608,765,492],"test":{"nested":true, "nested2": {}}},"data": [{"test":5},{"test": [1,2,3]}]}`)

	res, err := DropElement(input, "consented_providers", "test6")

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}

func TestDropNestedStructureSingleOccurrenceDeepNestedFullPath(t *testing.T) {
	input := []byte(`{"consented_providers":{"providers":[1608,765,492],"test":{"nested":true,"nested2": {"test6": 123}}},"data": [{"test":5},{"test": [1,2,3]}]}`)
	output := []byte(`{"consented_providers":{"providers":[1608,765,492],"test":{"nested2": {"test6": 123}}},"data": [{"test":5},{"test": [1,2,3]}]}`)

	res, err := DropElement(input, "consented_providers", "test", "nested")

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}

func TestDropNestedStructureDoesntExist(t *testing.T) {
	input := []byte(`{"consented_providers":{"providers":[1608,765,492]},"test":{"nested":true}},"data": [{"test":5},{"test": [1,2,3]}]}`)
	output := []byte(`{"consented_providers":{"providers":[1608,765,492]},"test":{"nested":true}},"data": [{"test":5},{"test": [1,2,3]}]}`)

	res, err := DropElement(input, "consented_providers", "test2")

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}

func TestSetElement(t *testing.T) {
	input := []byte(`{"consented_providers":{"providers":[1608,765,492], "nested":{"data": 123}}}`)
	setValue := []byte(`"test":{"nested":true}`)
	output := []byte(`{"consented_providers":{"providers":[1608,765,492], "nested":{"data": 123,"test":{"nested":true}}}}`)

	res, err := SetElement(input, setValue, "consented_providers", "nested")

	fmt.Println(string(output))
	fmt.Println(string(res))

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}

func TestSetElementNestedExists(t *testing.T) {
	input := []byte(`{"data":{"sitedata": "mysitedata"}}`)
	setValue := []byte(`{"somefpd": "fpdDataTest"}`)
	output := []byte(`{"data":{"sitedata": "mysitedata", "somefpd": "fpdDataTest}}`)

	res, err := SetElement(input, setValue, "data")

	fmt.Println(string(output))
	fmt.Println(string(res))

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, output, res, "Result is incorrect")
}
