/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package comp

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuperSharing(t *testing.T) {
	inheritanceDebug = t.Logf
	type A struct{ int }
	type B struct{ *A }
	type C struct{ *A }
	type D struct{ *A }
	type E struct {
		*D
	}
	type F struct {
		*B
		*C
		*D
		*E
	}
	// D has B and C and B and C share A
	newA := func(args Args, kwArgs KWArgs, options ...Option) (*A, error) {
		t.Log("building A")
		return &A{rand.Int()}, nil
	}
	newB := func(args Args, kwArgs KWArgs, options ...Option) (*B, error) {
		t.Log("building B")
		b := &B{}
		var err error
		b.A, err = CreateSharedSuperIfAbundant(options, newA, args, kwArgs, options...)
		require.NoError(t, err)
		return b, nil
	}
	newC := func(args Args, kwArgs KWArgs, options ...Option) (*C, error) {
		t.Log("building C")
		c := &C{}
		var err error
		c.A, err = CreateSharedSuperIfAbundant(options, newA, args, kwArgs, options...)
		require.NoError(t, err)
		return c, nil
	}
	newD := func(args Args, kwArgs KWArgs, options ...Option) (*D, error) {
		t.Log("building D")
		d := &D{}
		options = AddSharedSuperIfAbundant[A](options) // enrich
		var err error
		d.A, err = CreateSharedSuperIfAbundant(options, newA, args, kwArgs, options...)
		require.NoError(t, err)
		return d, nil
	}
	newE := func(args Args, kwArgs KWArgs, _ ...Option) (*E, error) {
		t.Log("building E")
		e := &E{}
		options := []Option{WithSharedSuper[A]()} // replace
		var err error
		e.D, err = newD(args, kwArgs, options...)
		require.NoError(t, err)
		return e, nil
	}
	newF := func(args Args, kwArgs KWArgs, options ...Option) (*F, error) {
		t.Log("building F")
		f := &F{}
		if len(options) == 0 {
			options = []Option{WithSharedSuper[A]()}
		}
		var err error
		f.B, err = newB(args, kwArgs, options...)
		require.NoError(t, err)
		f.C, err = newC(args, kwArgs, options...)
		require.NoError(t, err)
		f.D, err = newD(args, kwArgs, options...)
		require.NoError(t, err)
		f.E, err = newE(args, kwArgs, options...)
		require.NoError(t, err)
		return f, nil
	}
	t.Run("only set the option for F", func(t *testing.T) {
		inheritanceDebug = t.Logf
		f, err := newF(NoArgs, NoKWArgs())
		require.NoError(t, err)
		require.NotNil(t, f)
		ba := f.B.A
		ca := f.C.A
		da := f.D.A
		ea := f.E.A
		assert.Same(t, ba, ca)
		assert.Same(t, ba, da)
		assert.NotSame(t, ba, ea) // This is not the same
		assert.NotSame(t, da, ea) // This is not the same as e replaces all options on purpose
		assert.Same(t, ba, ca)
	})
	t.Run("use catch em all approach", func(t *testing.T) {
		f, err := newF(NoArgs, NoKWArgs(), WithMultiInheritanceSupport())
		require.NoError(t, err)
		require.NotNil(t, f)
		t.Logf("f\t\t%p", f)
		t.Logf("f.b\t%p", f.B)
		ba := f.B.A
		t.Logf("f.b.a\t%p", ba)
		t.Logf("f.c\t%p", f.C)
		ca := f.C.A
		t.Logf("f.c.a\t%p", ca)
		t.Logf("f.d\t%p", f.D)
		da := f.D.A
		t.Logf("f.d.a\t%p", da)
		t.Logf("f.e\t%p", f.E)
		ea := f.E.A
		t.Logf("f.e.a\t%p", ea)
		assert.Same(t, ba, ca)
		assert.Same(t, ba, da)
		assert.NotSame(t, ba, ea) // This is not the same
		assert.NotSame(t, da, ea) // This is not the same as e replaces all options on purpose
		assert.Same(t, ba, ca)
	})
	t.Run("use catch em all approach With addition", func(t *testing.T) {
		var options []Option
		options = AddMultiInheritanceSupportIfAbundant(options)
		f, err := newF(NoArgs, NoKWArgs(), options...)
		require.NoError(t, err)
		require.NotNil(t, f)
		ba := f.B.A
		ca := f.C.A
		da := f.D.A
		ea := f.E.A
		assert.Same(t, ba, ca)
		assert.Same(t, ba, da)
		assert.NotSame(t, ba, ea) // This is not the same
		assert.NotSame(t, da, ea) // This is not the same as e replaces all options on purpose
		assert.Same(t, ba, ca)
	})
}
