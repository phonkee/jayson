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

package jayson

import "sync"

// newRegistry creates a new registry
func newRegistry[T comparable]() *registry[T] {
	return &registry[T]{
		items: make(map[T]*registryItem[T]),
	}
}

// registry holds ext for given types
type registry[T comparable] struct {
	shared []Extension
	items  map[T]*registryItem[T]
	mutex  sync.RWMutex
}

// AddShared adds shared ext
func (r *registry[T]) AddShared(ext ...Extension) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.shared = append(r.shared, ext...)
}

// WithShared prepends shared ext to given ext
func (r *registry[T]) WithShared(ext ...Extension) []Extension {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return append(r.shared, ext...)
}

// Exists checks if ext for given type Exists
func (r *registry[T]) Exists(typ T) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.exists(typ)
}

// exists checks if ext for given type Exists (without lock)
func (r *registry[T]) exists(typ T) bool {
	_, ok := r.items[typ]
	return ok
}

// Get return ext for given type if Exists
func (r *registry[T]) Get(typ T) ([]Extension, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Get ext for given type
	if item, ok := r.items[typ]; ok {
		return item.ext, true
	}
	return nil, false
}

// Register registers ext for given type
func (r *registry[T]) Register(typ T, ext []Extension) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	exists := r.exists(typ)

	r.items[typ] = &registryItem[T]{
		typ: typ,
		ext: ext,
	}

	// warn if already registered
	if exists {
		return WarnAlreadyRegistered
	}

	return nil
}

// registryItem holds ext for given type
type registryItem[T comparable] struct {
	typ T
	ext []Extension
}
