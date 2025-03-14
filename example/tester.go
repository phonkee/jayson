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

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/phonkee/jayson"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
)

const (
	// ErrorDetailKey is a key for error detail
	ErrorDetailKey = "errorDetail"

	// TypeKey is a key for error type
	TypeKey = "type"
)

var (
	ErrorNotFound       = errors.New("not found")
	ErrorClientNotFound = errors.New("client not found")
)

// ExtErrorDetail is an extFunc that adds an error detail to the response object.
func ExtErrorDetail(detail string) jayson.Extension {
	return jayson.ExtObjectKeyValue(ErrorDetailKey, detail)
}

func init() {
	// setup custom global settings
	jayson.ReplaceGlobal(
		jayson.New(jayson.Settings{
			DefaultErrorStatus:        http.StatusInternalServerError,
			DefaultErrorMessageKey:    "message",
			DefaultErrorStatusCodeKey: "code",
			DefaultErrorStatusTextKey: "status_text",
			DefaultResponseStatus:     http.StatusOK,
			DefaultUnwrapObjectKey:    "_object",
		}),
	)

	// register errors
	jayson.Must(
		jayson.G().RegisterError(jayson.Any, jayson.ExtObjectKeyValuef(TypeKey, "error")),
	)
	jayson.Must(
		jayson.G().RegisterError(ErrorNotFound, jayson.ExtStatus(http.StatusNotFound)),
	)
	jayson.Must(
		jayson.G().RegisterError(ErrorClientNotFound,
			jayson.ExtStatus(http.StatusBadRequest),
			ExtErrorDetail("client not found"),
			jayson.ExtOmitObjectKey(TypeKey),
		),
	)
}

func main() {
	fmt.Print("\nExample errors:\n\n")
	ExampleError(ErrorNotFound)
	ExampleError(ErrorNotFound, jayson.ExtOmitSettingsKey(func(settings jayson.Settings) []string {
		return []string{settings.DefaultErrorStatusCodeKey, settings.DefaultErrorStatusTextKey}
	}))
	ExampleError(ErrorClientNotFound)
	ExampleError(ErrorClientNotFound, ExtErrorDetail("client not found?"))
	ExampleError(ErrorClientNotFound, jayson.ExtOmitObjectKey(ErrorDetailKey))
	ExampleError(ErrorClientNotFound, jayson.ExtOmitSettingsKey(func(settings jayson.Settings) []string {
		return []string{settings.DefaultErrorStatusCodeKey}
	}))
	ExampleError(ErrorClientNotFound,
		jayson.ExtOmitSettingsKey(func(settings jayson.Settings) []string {
			return []string{settings.DefaultErrorStatusCodeKey, settings.DefaultErrorStatusTextKey}
		}),
		jayson.ExtOmitObjectKey(ErrorDetailKey),
		jayson.ExtStatus(http.StatusAccepted),
	)

	// Example of wrapping error
	ExampleError(jayson.WrapError(ErrorClientNotFound, jayson.ExtStatus(http.StatusTeapot)))

	fmt.Print("\nExample responses:\n\n")

	// Response examples
	ExampleResponse(
		&TestStruct{
			Name: "John",
		},
	)
	ExampleResponse(
		&TestStruct{
			Name: "John",
		},
		jayson.ExtStatus(http.StatusTeapot),
	)
	ExampleResponse(
		jayson.ExtObjectUnwrap(&TestStruct{
			Name: "John",
		}),
		jayson.ExtStatus(http.StatusTeapot),
		jayson.ExtObjectKeyValue("key", "value"),
	)
}

type TestStruct struct {
	Name string `json:"name"`
}

func ExampleError(err error, ext ...jayson.Extension) {
	rw := httptest.NewRecorder()
	jayson.G().Error(context.Background(), rw, err, ext...)

	// get caller info
	_, file, line, _ := runtime.Caller(1)
	file = file[strings.LastIndex(file, "/")+1:]

	fmt.Printf("File: %v, Line: %v, Error: `%v`, Status: `%d`, Body: `%s`\n", file, line, err.Error(), rw.Code, strings.TrimSpace(rw.Body.String()))
}

func ExampleResponse(what any, ext ...jayson.Extension) {
	rw := httptest.NewRecorder()
	jayson.G().Response(context.Background(), rw, what, ext...)

	// get caller info
	_, file, line, _ := runtime.Caller(1)
	file = file[strings.LastIndex(file, "/")+1:]

	fmt.Printf("File: %v, Line: %v, Status: `%d`, Body: `%s`\n", file, line, rw.Code, strings.TrimSpace(rw.Body.String()))
}
