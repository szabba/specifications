// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package specifications provides a type-parametric implementation of the specification pattern.
//
// It lets you write a boolean condition once and reuse it in different contexts.
// You could use a single specification to:
//     - generate queries for different databases,
//     - check if a value meets a condition,
//     - generate a human-readable description of the condition,
//     - explain why a value does not meet a condition.
//
// The library does not provide all that functionality out of the box.
// What it gives you is the plumbing common in all these use cases.
package specifications
