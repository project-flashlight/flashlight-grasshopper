# Assert
Assert is a lite weight testing package which provides basic assertion operations.

## Installation
```bash
go get -u github.com/vwdilab/mango/assert
```

## Usage
**Example**:

```go
package ...

import "github.com/vwdilab/mango/assert"

func Test_SomethingToBeTrue(t *testing.T) {

	// given

	// when
	someCondtion := false

	// then
	assert.True(t, someCondtion, "expected to be true") // will fail

}

func Test_SomethingToBeEqual(t *testing.T) {

	// given
	exp := 42

	// when
	act := 3141

	// then
	assert.Equals(t, exp, act) // will fail

}

func Test_NoErrorOccured(t *testing.T) {

	// given

	// when
	err := errors.New("An unexpected error")

	// then
	assert.NoError(t, err) // will fail

}

```