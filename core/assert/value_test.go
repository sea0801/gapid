// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package assert_test

import (
	"fmt"
	"os"

	"github.com/google/gapid/core/assert"
)

// An example that shows the simplest of value equality tests with a message
func ExampleAssertMessage() {
	ctx := assert.Context(nil)
	assert.For(ctx, "A message").That(false).Equals(true)
	fmt.Fprintf(os.Stdout, "Test complete")
	// Output:
	// Error:A message
	//     Got       false
	//     Expect == true
	// Test complete
}

// An example that shows a critical error
func ExampleAssertCritical() {
	defer func() { recover() }() // Consume the critical level panic
	ctx := assert.Context(nil)
	assert.For(ctx, "A message").Critical().That(false).Equals(true)
	fmt.Fprintf(os.Stdout, "Test complete")
	// Output:
	// Critical:A message
	//     Got       false
	//     Expect == true
}

// An example of testing untyped nil values
func ExampleNil() {
	ctx := assert.Context(nil)
	assert.For(ctx, "nil equals nil").That(nil).Equals(nil)
	assert.For(ctx, "nil does not equal nil").That(nil).NotEquals(nil)
	assert.For(ctx, "nil deep equals nil").That(nil).DeepEquals(nil)
	assert.For(ctx, "nil is nil").That(nil).IsNil()
	assert.For(ctx, "nil is not nil").That(nil).IsNotNil()
	// Output:
	// Error:nil does not equal nil
	//     Got       <nil>
	//     Expect != <nil>
	// Error:nil is not nil
	//     Got       <nil>
	//     Expect != `nil`
}

// An example of testing typed nil values
func ExampleTypedNil() {
	var typedNil *int
	ctx := assert.Context(nil)
	assert.For(ctx, "typed_nil equals nil").That(typedNil).Equals(nil)
	assert.For(ctx, "typed_nil does not equal nil").That(typedNil).NotEquals(nil)
	assert.For(ctx, "typed_nil deep equals nil").That(typedNil).DeepEquals(nil)
	assert.For(ctx, "typed_nil is nil").That(typedNil).IsNil()
	assert.For(ctx, "typed_nil is not nil").That(typedNil).IsNotNil()
	// Output:
	// Error:typed_nil equals nil
	//     Got       <nil>
	//     Expect == <nil>
	// Error:typed_nil deep equals nil
	//     nil ⟦<nil>⟧ != ⟦<nil>⟧
	// Error:typed_nil is not nil
	//     Got       <nil>
	//     Expect != `nil`
}

// An example of testing non nil values
func ExampleNotNil() {
	notNil := &struct{ s string }{"not_nil"}
	ctx := assert.Context(nil)
	assert.For(ctx, "not_nil equals nil").That(notNil).Equals(nil)
	assert.For(ctx, "not_nil does not equal nil").That(notNil).NotEquals(nil)
	assert.For(ctx, "not_nil deep equals nil").That(notNil).DeepEquals(nil)
	assert.For(ctx, "not_nil is nil").That(notNil).IsNil()
	assert.For(ctx, "not_nil is not nil").That(notNil).IsNotNil()
	// Output:
	// Error:not_nil equals nil
	//     Got       &{not_nil}
	//     Expect == <nil>
	// Error:not_nil deep equals nil
	//     nil ⟦&{not_nil}⟧ != ⟦<nil>⟧
	// Error:not_nil is nil
	//     Got       &{not_nil}
	//     Expect == `nil`
}

// An example of using value Equals
func ExampleValueEquals() {
	ctx := assert.Context(nil)
	assert.With(ctx).That(1).Equals(1)
	assert.With(ctx).That(2).Equals(3)
	// Output:
	// Error:
	//     Got       2
	//     Expect == 3
}

// An example of using value NotEquals
func ExampleValueNotEquals() {
	ctx := assert.Context(nil)
	assert.With(ctx).That(1).NotEquals(1)
	assert.With(ctx).That(2).NotEquals(3)
	// Output:
	// Error:
	//     Got       1
	//     Expect != 1
}

// An example of using value DeepEquals
func ExampleValueDeepEquals() {
	a := []string{"1", "2"}
	values := []struct{ V []string }{{a}, {a}, {[]string{"1", "2"}}, {[]string{"1", "3"}}, {[]string{"2", "4"}}}
	ctx := assert.Context(nil)
	assert.For(ctx, "deep equals same object").That(values[0]).DeepEquals(values[1])
	assert.For(ctx, "deep equals same value").That(values[0]).DeepEquals(values[2])
	assert.For(ctx, "deep equals different value").That(values[0]).DeepEquals(values[3])
	assert.For(ctx, "deep equals all different").That(values[0]).DeepEquals(values[4])
	// Output:
	// Error:deep equals different value
	//     ⟦2⟧ != ⟦3⟧ for v.V[1]
	// Error:deep equals all different
	//     ⟦1⟧ != ⟦2⟧ for v.V[0]
	//     ⟦2⟧ != ⟦4⟧ for v.V[1]
}

// An example of using value NotEquals
func ExampleValueDeepNotEquals() {
	a := []string{"1", "2"}
	values := []struct{ V []string }{{a}, {a}, {[]string{"1", "2"}}, {[]string{"1", "3"}}, {[]string{"2", "4"}}}
	ctx := assert.Context(nil)
	assert.For(ctx, "deep equals same object").That(values[0]).DeepNotEquals(values[1])
	assert.For(ctx, "deep equals same value").That(values[0]).DeepNotEquals(values[2])
	assert.For(ctx, "deep equals different value").That(values[0]).DeepNotEquals(values[3])
	assert.For(ctx, "deep equals all different").That(values[0]).DeepNotEquals(values[4])
	// Output:
	// Error:deep equals same object
	//     Got            {[1 2]}
	//     Expect deep != {[1 2]}
	// Error:deep equals same value
	//     Got            {[1 2]}
	//     Expect deep != {[1 2]}
}
