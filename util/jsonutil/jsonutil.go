package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
)

var comma = []byte(",")[0]
var colon = []byte(":")[0]
var sqBracket = []byte("]")[0]
var openCurlyBracket = []byte("{")[0]
var closingCurlyBracket = []byte("}")[0]
var quote = []byte(`"`)[0]

func FindElement(extension []byte, elementNames ...string) (bool, int64, int64, error) {

	elementName := elementNames[0]

	buf := bytes.NewBuffer(extension)
	dec := json.NewDecoder(buf)
	found := false
	var startIndex, endIndex int64
	var i interface{}
	for {
		token, err := dec.Token()
		if err == io.EOF {
			// io.EOF is a successful end
			break
		}
		if err != nil {
			return false, -1, -1, err
		}

		if token == elementName {

			err := dec.Decode(&i)
			if err != nil {
				return false, -1, -1, err
			}
			endIndex = dec.InputOffset()

			if dec.More() {
				//if there were other elements before
				if extension[startIndex] == comma {
					startIndex++
				}

				for {
					//structure has more elements, need to find index of comma
					if extension[endIndex] == comma {
						endIndex++
						break
					}
					endIndex++
				}
			}
			found = true
			break
		} else {
			startIndex = dec.InputOffset()
		}

	}
	if found {
		if len(elementNames) == 1 {
			return found, startIndex, endIndex, nil
		} else if len(elementNames) > 1 {

			for {
				//find the beginning of nested element
				if extension[startIndex] == colon {
					startIndex++
					break
				}
				startIndex++
			}

			for {
				if endIndex == int64(len(extension)) {
					endIndex--
				}

				//if structure had more elements, need to find index of comma at the end
				if extension[endIndex] == sqBracket || extension[endIndex] == closingCurlyBracket {
					break
				}

				if extension[endIndex] == comma {
					endIndex--
					break
				} else {
					endIndex--
				}

			}

			if found {
				found, startInd, endInd, err := FindElement(extension[startIndex:endIndex], elementNames[1:]...)
				return found, startIndex + startInd, startIndex + endInd, err
			}
			return found, startIndex, startIndex, nil
		}

	}
	return found, startIndex, endIndex, nil
}

func DropElement(extension []byte, elementNames ...string) ([]byte, error) {
	//Doesnt support drop element from array
	found, startIndex, endIndex, err := FindElement(extension, elementNames...)
	if err != nil {
		return nil, err
	}

	if found {
		extension = append(extension[:startIndex], extension[endIndex:]...)
	}

	return extension, nil

}

func SetElement(data []byte, setValue []byte, keys ...string) ([]byte, error) {
	// Doesn't support set value in array
	// Doesn't support creation of nested elements
	// Can create one nested element if len(keys) == 1 and this element was not found
	// Can only insert value to existing element
	// Doesn't support type check. Specified element should be a non-empty object

	if len(keys) < 1 {
		return data, nil
	}

	if len(setValue) <= 2 { //empty object to set
		return data, nil
	}
	if len(data) <= 2 { //data is empty
		return buildNewObject(keys[len(keys)-1], setValue), nil
	}

	found, _, endIndex, err := FindElement(data, keys...)
	if err != nil {
		return nil, err
	}
	if found {
		if setValue[0] == openCurlyBracket && setValue[len(setValue)-1] == closingCurlyBracket {
			//delete open and closing brackets
			setValue = setValue[1 : len(setValue)-1]
		}

		for {
			//find the end of element
			if data[endIndex] == closingCurlyBracket {
				endIndex--
				break
			}
			endIndex--
		}
		result := insertValueToArray(data, setValue, int(endIndex))
		return result, nil
	} else if len(keys) == 1 {
		//if element not found, set new object to last key
		newObject := buildNewObject(keys[len(keys)-1], setValue)
		indexToInsert := len(data) - 1 //insert before the closing curly bracket
		result := insertValueToArray(data, newObject, indexToInsert)
		return result, nil
	}
	return data, nil
}

func insertValueToArray(data, setValue []byte, indexToInsert int) []byte {
	dataCopy := make([]byte, len(data))
	//Make a copy of data, otherwise it modifies original data array
	copy(dataCopy, data)
	res := append(dataCopy[:indexToInsert], comma)
	res = append(res, setValue...)
	res = append(res, data[indexToInsert:]...)
	return res
}

func buildNewObject(key string, setValue []byte) []byte {
	//build json like "key":setValue
	result := make([]byte, 0)
	result = append(result, quote)
	result = append(result, []byte(key)...)
	result = append(result, quote)
	result = append(result, colon)
	result = append(result, setValue...)
	return result
}

func SetElement2(originDataInput []byte, setValue []byte, key string) ([]byte, error) {

	originData := make(map[string]interface{})
	setValueData := make(map[string]interface{})

	err := json.Unmarshal(originDataInput, &originData)
	if err != nil {
		return originDataInput, err
	}
	err = json.Unmarshal(setValue, &setValueData)
	if err != nil {
		return originDataInput, err
	}

	if val, ok := originData[key]; ok {
		//element exists already - add new element(s) to it
		data := val.(map[string]interface{})
		for k, v := range setValueData {
			data[k] = v
		}
		originData[key] = data
	} else {
		//element doesn't exist - set value as is
		originData[key] = setValueData
	}
	res, err := json.Marshal(originData)
	return res, err
}
