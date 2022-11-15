package core

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSingleNumber(t *testing.T) {

    assert := assert.New(t)

    hashValue := Base62hash(0)

    assert.Equal("0", hashValue)
}

func TestMultipleNumber(t *testing.T) {

    assert := assert.New(t)

    hashValue := Base62hash(11157)

    assert.Equal("2TX", hashValue)
}