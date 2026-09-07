# Assert with spf13/afero

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/aferoassert)](https://github.com/nhatthm/aferoassert/releases/latest)
[![Build Status](https://github.com/nhatthm/aferoassert/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/aferoassert/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/aferoassert/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/aferoassert)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/aferoassert)
[![Donate](https://img.shields.io/badge/%20-Donate-%20?style=flat&logo=githubsponsors&color=E5E4E2)](http://donate.nhat.me)

The logic is shamelessly copy from [stretchr/testify/assert](https://github.com/stretchr/testify/tree/master/assert)
with some salt and pepper.

## Prerequisites

- `Go >= 1.25`

## Install

```bash
go get go.nhat.io/aferoassert
```

## Usage

```go
package mypackage_test

import (
	"testing"

	"github.com/spf13/afero"
	"go.nhat.io/aferoassert"
)

func TestTreeEqual_Success(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows:
    - golangci-lint.yaml
    - test.yaml 'perm:"0644"'
`

	aferoassert.DirExists(t, osFs, ".github")
	aferoassert.YAMLTreeEqual(t, osFs, tree, ".github")
}
```

## Donation

If this project saved you some development time, buy me a cup of coffee :)

[![donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](http://donate.nhat.me)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://github.com/nhatthm/donate.nhat.me/blob/master/images/qr_sponsor.png" width="147px" />

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)
