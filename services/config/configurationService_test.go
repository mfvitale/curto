package config

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCorrectLoading(t *testing.T) {

    assert := assert.New(t)

    appConfigurationService := NewAppConfigurationService("../../config.yml", logrus.New())

    assert.Equal("8080", appConfigurationService.GetConfig().App.Port)
}

func TestLoadFail(t *testing.T) {

    assert := assert.New(t)

    appConfigurationService := NewAppConfigurationService("notexistingfile", logrus.New())

    assert.Equal("", appConfigurationService.GetConfig().App.Port)
}