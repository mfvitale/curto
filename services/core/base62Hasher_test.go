package core

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSingleNumber(t *testing.T) {

    assert := assert.New(t)

    hashValue, err := Base62hash(0)

    assert.Nil(err)
    assert.Equal("0", hashValue)
}

func TestMultipleNumber(t *testing.T) {

    assert := assert.New(t)

    hashValue, err := Base62hash(11157)

    assert.Nil(err)
    assert.Equal("2TX", hashValue)
}