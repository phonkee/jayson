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

package tester

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"testing"
	"time"
)

// WithHttpServer runs a new HTTP server with the provided handler and calls the provided function with the server address.
func WithHttpServer(t *testing.T, ctx context.Context, handler http.Handler, fn func(t *testing.T, ctx context.Context, address string)) {
	if handler == nil {
		panic("handler is nil")
	}
	if fn == nil {
		panic("closure is nil")
	}
	if ctx == nil {
		panic("context is nil")
	}

	// add cancel to context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// listen to the first available port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	// close listener on exit
	defer func() {
		_ = lis.Close()
	}()

	// prepare http server
	srv := &http.Server{
		// Handler to serve
		Handler: handler,
		// Base context
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	// run server in goroutine
	go func() {
		_ = srv.Serve(lis)
	}()

	// prepare message
	message := fmt.Sprintf("test http server on %v", lis.Addr().String())

	// run test
	t.Run(message, func(t *testing.T) {
		// this should not panic
		assert.NotPanicsf(t, func() {
			// call function with address
			fn(t, ctx, lis.Addr().String())
		}, "closure should not panic")
	})

	// fast shutdown
	ctx, cancelTimeout := context.WithTimeout(context.Background(), time.Second)
	defer cancelTimeout()

	// shutdown server with timeout
	shutdownErr := srv.Shutdown(ctx)

	// shutdown server immediately
	if shutdownErr != nil {
		// bye bye
		_ = srv.Close()
	}
}
