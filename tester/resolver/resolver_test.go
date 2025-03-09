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

package resolver_test

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
)

// exampleResponse is a response struct for testing
type exampleResponse struct {
	Status string `json:"status"`
	Host   string `json:"host"`
}

// exampleHandler is a handler for testing
func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(exampleResponse{
		Status: "something",
		Host:   "localhost",
	}); err != nil {
		panic(err)
	}
}

// MatchByStringContains matches string by substring
func matchByStringContains(s string) func(in string) bool {
	return func(in string) bool {
		return strings.Contains(in, s)
	}
}

func newHealthRouter(t require.TestingT) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/health", exampleHandler).Methods(http.MethodGet).Name("api:v1:health")
	router.HandleFunc("/api/v1/health/{component}", exampleHandler).Methods(http.MethodGet).Name("api:v1:health:extra")
	return router
}
