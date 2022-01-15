// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package typeparams provides functions to work indirectly with type parameter
// data stored in go/ast and go/types objects, while these API are guarded by a
// build constraint.
//
// This package exists to make it easier for tools to work with generic code,
// while also compiling against older Go versions.
package typeparams

import (
	"go/ast"
	"go/token"
<<<<<<< HEAD:vendor/golang.org/x/tools/internal/typeparams/common.go
=======
	"go/types"
>>>>>>> b5e1417 (add operator-registry as dep and use gzip method to check bundle size):vendor/golang.org/x/tools/internal/typeparams/doc.go
)

// A IndexExprData holds data from both ast.IndexExpr and the new
// ast.MultiIndexExpr, which was introduced in Go 1.18.
type IndexExprData struct {
	X       ast.Expr   // expression
	Lbrack  token.Pos  // position of "["
	Indices []ast.Expr // index expressions
	Rbrack  token.Pos  // position of "]"
}
<<<<<<< HEAD:vendor/golang.org/x/tools/internal/typeparams/common.go
=======

// IsTypeParam reports whether t is a type parameter.
func IsTypeParam(t types.Type) bool {
	_, ok := t.(*TypeParam)
	return ok
}
>>>>>>> b5e1417 (add operator-registry as dep and use gzip method to check bundle size):vendor/golang.org/x/tools/internal/typeparams/doc.go
