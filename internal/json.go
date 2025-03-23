/*
 * MIT License
 *
 * Copyright (c) 2025 Peter Vrba
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package internal

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
)

// ExtractJsonPath extracts the json path from the given raw json
func ExtractJsonPath(_ require.TestingT, jsonPath string, raw json.RawMessage) (json.RawMessage, error) {
	jsonPath = strings.TrimSpace(jsonPath)

	var (
		err error
		ok  bool
	)

main:
	for _, part := range strings.Split(jsonPath, ".") {
		if part == "" {
			continue main
		}
		// try to parse the part as a number so we know we need to unmarshal array
		if number, numberError := strconv.ParseUint(part, 10, 64); numberError == nil {
			var arr []json.RawMessage

			// try to unmarshal array
			if err = json.Unmarshal(raw, &arr); err != nil {
				err = fmt.Errorf("cannot unmarshal object into slice")
				break main
			}

			// check if the number is in bounds
			if int(number) >= len(arr) {
				err = fmt.Errorf("%w: out of bounds", ErrNotPresent)
				break main
			}

			raw = arr[number]
			continue main
		}
		// parse object
		obj := make(map[string]json.RawMessage)

		// try to unmarshal object, if it fails there is a problem
		if err = json.Unmarshal(raw, &obj); err != nil {
			err = fmt.Errorf("cannot unmarshal action")
			break main
		}

		// check if object contains the part
		if raw, ok = obj[part]; !ok {
			err = ErrNotPresent
			break main
		}
	}

	if err != nil {
		return nil, err
	}

	return raw, nil
}
