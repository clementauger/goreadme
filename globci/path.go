// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package filepath implements utility routines for manipulating filename paths
// in a way compatible with the target operating system-defined file paths.

// The globci package uses either forward slashes or backslashes,
// depending on the operating system. To process paths such as URLs
// that always use forward slashes regardless of the operating
// system, see the path package.
package globci

import (
	"os"
	"path/filepath"
)

// re import
const (
	Separator     = os.PathSeparator
	ListSeparator = os.PathListSeparator
)

// re import
var (
	SkipDir    = filepath.SkipDir
	Clean      = filepath.Clean
	ToSlash    = filepath.ToSlash
	FromSlash  = filepath.FromSlash
	SplitList  = filepath.SplitList
	Split      = filepath.Split
	Join       = filepath.Join
	Ext        = filepath.Ext
	Abs        = filepath.Abs
	Rel        = filepath.Rel
	VolumeName = filepath.VolumeName
)
