// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package specifications

// A Specification describes a boolean condition.
//
// Leafs are basic conditions that can be combined with Or, And and Not.
//
// The specification only describes the structure of the condition.
// To actually do something with it you also need an Evaluator.
//
// The zero value of a specification is not useful.
// Functions in this package will panic when passed a zero specification.
type Specification[Leaf any] struct {
	ops   []_OpCode
	leafs []Leaf
}

// An Evaluator describes how to convert a Specification[Leaf] to some Output type.
type Evaluator[Leaf, Output any] interface {

	// EvaluateLeaf turns a leaf specification into the output type.
	EvaluateLeaf(Leaf) Output

	// EvaluateNot negates the output for a specification that was wrapped with Not.
	EvaluateNot(Output) Output

	// EvaluateOr combines the outputs of multiple specifications that were combined with Or.
	EvaluateOr([]Output) Output

	// EvaluateAnd combines the outputs of multiple specifications that were combined with And.
	EvaluateAnd([]Output) Output
}

type _OpCode uint16

const (
	_OpCodeLeaf _OpCode = iota + 1
	_OpCodeNot
	_OpCodeAnd
	_OpCodeOr
)

// Leaf creates a leaf specification.
// Some specifications contain only leafs.
func Leaf[Leaf any](l Leaf) Specification[Leaf] {
	return Specification[Leaf]{
		leafs: []Leaf{l},
		ops:   []_OpCode{_OpCodeLeaf, _OpCode(0)},
	}
}

// Not creates a specification that negates s.
// It is true when s was false and false when s was true.
//
// Not panics if s is a zero specification.
func Not[Leaf any](s Specification[Leaf]) Specification[Leaf] {
	if s.Zero() {
		panic("Not: cannot use zero spec")
	}

	return Specification[Leaf]{
		leafs: s.leafs,
		ops:   append(append([]_OpCode{}, s.ops...), _OpCodeNot),
	}
}

// And creates a specification that is true when all the specs are true.
//
// And panics if specs is empty or contains a zero specification.
func And[Leaf any](specs ...Specification[Leaf]) Specification[Leaf] {
	verifyParts(
		specs,
		"And: cannot combine 0 specifications",
		"And: cannot combine the zero specification")

	return combineSpecs(_OpCodeAnd, specs)
}

// Or creates a specification that is true when any of the specs is true.
//
// Or panics if specs is empty or contains a zero specification.
func Or[Leaf any](specs ...Specification[Leaf]) Specification[Leaf] {
	verifyParts(
		specs,
		"Or: cannot combine 0 specifications",
		"Or: cannot combine the zero specification")

	return combineSpecs(_OpCodeOr, specs)
}

func verifyParts[Leaf any](
	specs []Specification[Leaf],
	whenEmpty, whenZeroFound string,
) {
	if len(specs) == 0 {
		panic(whenEmpty)
	}
	for _, s := range specs {
		if s.Zero() {
			panic(whenZeroFound)
		}
	}
}

func combineSpecs[Leaf any](op _OpCode, specs []Specification[Leaf]) Specification[Leaf] {
	out := preallocateCombined(specs)

	for _, s := range specs {
		out.ops = reencode(out.ops, s.ops, len(out.leafs))
		out.leafs = append(out.leafs, s.leafs...)
	}

	out.ops = append(out.ops, op)
	out.ops = append(out.ops, _OpCode(len(specs)))

	return out
}

func preallocateCombined[Leaf any](specs []Specification[Leaf]) Specification[Leaf] {
	const combEncLen = 2

	opCount, leafCount := 0, 0
	for _, s := range specs {
		opCount += len(s.ops)
		leafCount += len(s.leafs)
	}

	return Specification[Leaf]{
		ops:   make([]_OpCode, 0, opCount+combEncLen),
		leafs: make([]Leaf, 0, leafCount),
	}
}

func reencode(dst, src []_OpCode, offset int) []_OpCode {
	for len(src) > 0 {
		switch src[0] {

		case _OpCodeLeaf:
			dst = append(dst, _OpCodeLeaf)
			dst = append(dst, src[1]+_OpCode(offset))
			src = src[2:]

		case _OpCodeNot:
			dst = append(dst, _OpCodeNot)
			src = src[1:]

		case _OpCodeAnd, _OpCodeOr:
			dst = append(dst, src[0])
			dst = append(dst, src[1])
			src = src[2:]
		}
	}
	return dst
}

// Zero checks if s is a zero-value specification.
//
// If this is true, s is unusable.
// Functions in this package will panic when passed in zero-value specifications.
func (s Specification[Leaf]) Zero() bool {
	return s.ops == nil && s.leafs == nil
}

// Evaluate uses an Evaluator to convert a Specification an Output type.
func Evaluate[Leaf, Output any](
	spec Specification[Leaf],
	ev Evaluator[Leaf, Output],
) Output {

	opsLeft, stack := spec.ops, make([]Output, 0, 20)

	if len(opsLeft) == 0 {
		panic("Evaluate: cannot evaluate the zero specification")
	}

	for len(opsLeft) > 0 {

		switch opsLeft[0] {

		case _OpCodeLeaf:
			leafIx := int(opsLeft[1])
			leaf := spec.leafs[leafIx]
			out := ev.EvaluateLeaf(leaf)
			stack = append(stack, out)

			opsLeft = opsLeft[2:]

		case _OpCodeNot:
			top, rest := stack[0], stack[:len(stack)-1]
			out := ev.EvaluateNot(top)
			stack = append(rest, out)

			opsLeft = opsLeft[1:]

		case _OpCodeAnd:
			argCount := int(opsLeft[1])
			top, rest := pickTop(stack, argCount)
			out := ev.EvaluateAnd(top)
			stack = append(rest, out)

			opsLeft = opsLeft[2:]

		case _OpCodeOr:
			argCount := int(opsLeft[1])
			top, rest := pickTop(stack, argCount)
			out := ev.EvaluateOr(top)
			stack = append(rest, out)

			opsLeft = opsLeft[2:]
		}
	}

	return stack[0]
}

func pickTop[A any](stack []A, n int) (top, rest []A) {
	restLen := len(stack) - n
	rest = stack[:restLen]
	top = stack[restLen:]
	return top, rest
}
