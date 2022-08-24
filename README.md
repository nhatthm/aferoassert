# Assert with spf13/afero

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/aferoassert)](https://github.com/nhatthm/aferoassert/releases/latest)
[![Build Status](https://github.com/nhatthm/aferoassert/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/aferoassert/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/aferoassert/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/aferoassert)
[![Go Report Card](https://goreportcard.com/badge/go.nhat.io/aferoassert)](https://goreportcard.com/report/go.nhat.io/aferoassert)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/aferoassert)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

The logic is shamelessly copy from [stretchr/testify/assert](https://github.com/stretchr/testify/tree/master/assert)
with some salt and pepper.

## Prerequisites

- `Go >= 1.17`

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

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
