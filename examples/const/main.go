package main

import (
	"fmt"

	"github.com/szabba/specifications"
)

func main() {
	for _, ex := range []specifications.Specification[bool]{
		specifications.Leaf(false),
		specifications.Leaf(true),
		specifications.Not(specifications.Leaf(false)),
		specifications.Not(specifications.Leaf(true)),
		specifications.And(specifications.Leaf(true)),
		specifications.Or(specifications.Leaf(false), specifications.Leaf(true)),
		specifications.And(specifications.Leaf(false), specifications.Leaf(true)),
	} {
		s := specifications.Evaluate[bool, string](ex, ToString{})
		b := specifications.Evaluate[bool, bool](ex, ToBool{})
		fmt.Printf("%s => %v\n", s, b)
	}
}

type ToBool struct{}

var _ specifications.Evaluator[bool, bool] = ToBool{}

func (ToBool) EvaluateLeaf(v bool) bool { return v }

func (ToBool) EvaluateNot(v bool) bool { return !v }

func (ToBool) EvaluateAnd(vs []bool) bool {
	out := true
	for _, v := range vs {
		out = out && v
	}
	return out
}

func (ToBool) EvaluateOr(vs []bool) bool {
	out := false
	for _, v := range vs {
		out = out || v
	}
	return out
}

type ToString struct{}

var _ specifications.Evaluator[bool, string] = ToString{}

func (ToString) EvaluateLeaf(v bool) string {
	return fmt.Sprint(v)
}

func (ToString) EvaluateNot(v string) string {
	return "!" + v
}

func (ts ToString) EvaluateAnd(vs []string) string {
	return ts.combine("&&", vs)
}

func (ts ToString) EvaluateOr(vs []string) string {
	return ts.combine("||", vs)
}

func (ToString) combine(op string, vs []string) string {
	if len(vs) == 1 {
		return vs[0]
	}

	out := "("
	for i, v := range vs {
		if i > 0 {
			out = out + " " + op + " "
		}
		out = out + v
	}
	out += ")"
	return out
}
