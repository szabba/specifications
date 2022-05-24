// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package specifications_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/specifications"
)

func TestConstructors(t *testing.T) {
	t.Run("Not/panicsWhenGivenAZeroSpec", func(t *testing.T) {
		// given
		var zero specifications.Specification[bool]

		// when
		p := catchPanic(func() { not(zero) })

		// then
		msg := "Not: cannot use zero spec"
		assert.That(p == msg, t.Errorf, "got %#v panic, not %q", p, msg)
	})

	t.Run("And/panicsWhenGivenANoSpecs", func(t *testing.T) {
		// given
		// when
		p := catchPanic(func() { and() })

		// then
		msg := "And: cannot combine 0 specifications"
		assert.That(p == msg, t.Errorf, "got %#v panic, not %q", p, msg)
	})

	t.Run("And/panicsWhenGivenAZeroSpec", func(t *testing.T) {
		// given
		var zero specifications.Specification[bool]

		// when
		p := catchPanic(func() { and(zero) })

		// then
		msg := "And: cannot combine the zero specification"
		assert.That(p == msg, t.Errorf, "got %#v panic, not %q", p, msg)
	})

	t.Run("Or/panicsWhenGivenANoSpecs", func(t *testing.T) {
		// given
		// when
		p := catchPanic(func() { or() })

		// then
		msg := "Or: cannot combine 0 specifications"
		assert.That(p == msg, t.Errorf, "got %#v panic, not %q", p, msg)
	})

	t.Run("And/panicsWhenGivenAZeroSpec", func(t *testing.T) {
		// given
		var zero specifications.Specification[bool]

		// when
		p := catchPanic(func() { or(zero) })

		// then
		msg := "Or: cannot combine the zero specification"
		assert.That(p == msg, t.Errorf, "got %#v panic, not %q", p, msg)
	})
}

func catchPanic(f func()) (p any) {
	defer func() { p = recover() }()
	f()
	return p
}

func TestEvaluatePanicsWhenGivenTheZeroSpec(t *testing.T) {
	// given
	var zero specifications.Specification[bool]

	// when
	p := catchPanic(func() {
		specifications.Evaluate[bool, bool](zero, Evaluator{})
	})

	// then
	msg := "Evaluate: cannot evaluate the zero specification"
	assert.That(p == msg, t.Errorf, "got %#v panic, not %q", p, msg)
}

func TestEvaluate(t *testing.T) {
	tts := map[string]struct {
		Spec specifications.Specification[bool]
		Out  bool
	}{
		"false": {leaf(false), false},
		"true":  {leaf(true), true},

		"Not(false)": {not(leaf(false)), true},
		"Not(true)":  {not(leaf(true)), false},

		"And(false)": {and(leaf(false)), false},
		"And(true)":  {and(leaf(true)), true},

		"And(false,false)": {and(leaf(false), leaf(false)), false},
		"And(false,true)":  {and(leaf(false), leaf(true)), false},
		"And(true,false)":  {and(leaf(true), leaf(false)), false},
		"And(true,true)":   {and(leaf(true), leaf(true)), true},

		"Or(false,false)": {or(leaf(false), leaf(false)), false},
		"Or(false,true)":  {or(leaf(false), leaf(true)), true},
		"Or(true,false)":  {or(leaf(true), leaf(false)), true},
		"Or(true,true)":   {or(leaf(true), leaf(true)), true},

		// These cases are there to provide coverage for reencode in the library's internals.
		"And(Not(true)}": {and(not(leaf(true))), false},
		"And(And(true))": {and(and(leaf(true))), true},
	}

	for name, tt := range tts {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			out := specifications.Evaluate[bool, bool](tt.Spec, Evaluator{})

			// then
			assert.That(
				out == tt.Out,
				t.Errorf, "spec evaluated to %#v, not %#v", out, tt.Out)
		})
	}

}

var (
	leaf = specifications.Leaf[bool]
	not  = specifications.Not[bool]
	and  = specifications.And[bool]
	or   = specifications.Or[bool]
)

type Evaluator struct{}

var _ specifications.Evaluator[bool, bool] = Evaluator{}

func (Evaluator) EvaluateLeaf(v bool) bool { return v }

func (Evaluator) EvaluateNot(v bool) bool { return !v }

func (Evaluator) EvaluateAnd(vs []bool) bool {
	out := true
	for _, v := range vs {
		out = out && v
	}
	return out
}

func (Evaluator) EvaluateOr(vs []bool) bool {
	out := false
	for _, v := range vs {
		out = out || v
	}
	return out
}
