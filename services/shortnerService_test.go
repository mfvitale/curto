package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UrlRepositoryMock struct{
	mock.Mock
}

func (m *UrlRepositoryMock) Store(key string, value string) error {

	args := m.Called(key, value)
	return args.Error(0)

}

func (m *UrlRepositoryMock) Get(key string) (string, error) {

	args := m.Called(key)
	return args.String(0), args.Error(1)

}

type IdGeneratorMock struct{
	mock.Mock
}

func (m *IdGeneratorMock) NextID() (uint64, error) {

	args := m.Called()
	return uint64(args.Int(0)), args.Error(1)
}

func TestUrlShortner(t *testing.T) {

    assert := assert.New(t)

	urlRepositoryMock := new(UrlRepositoryMock)
	urlRepositoryMock.On("Store", "2TX", "https://www.mfvitale.me").Return(nil)

	identifierGenerator := new(IdGeneratorMock)
	identifierGenerator.On("NextID").Return(11157, nil)

    shortnertService := NewShortnerService(urlRepositoryMock, identifierGenerator)

	hashValue, _ := shortnertService.Shorten("https://www.mfvitale.me")

	urlRepositoryMock.AssertExpectations(t)
	urlRepositoryMock.AssertNumberOfCalls(t, "Store", 1)

	identifierGenerator.AssertExpectations(t)
	identifierGenerator.AssertNumberOfCalls(t, "NextID", 1)

	assert.Equal("2TX", hashValue)
}

func TestUrlResolve(t *testing.T) {

    assert := assert.New(t)

	urlRepositoryMock := new(UrlRepositoryMock)
	urlRepositoryMock.On("Get", "2TX").Return("https://www.mfvitale.me", nil)

	identifierGenerator := new(IdGeneratorMock)

    shortnertService := NewShortnerService(urlRepositoryMock, identifierGenerator)

	originalUrl, _ := shortnertService.Resolve("2TX")

	urlRepositoryMock.AssertExpectations(t)
	urlRepositoryMock.AssertNumberOfCalls(t, "Get", 1)

	assert.Equal("https://www.mfvitale.me", originalUrl)
}

func TestResolveNotShortnedUrl(t *testing.T) {

    assert := assert.New(t)

	urlRepositoryMock := new(UrlRepositoryMock)
	urlRepositoryMock.On("Get", "2TX").Return("https://www.mfvitale.me", nil)

	identifierGenerator := new(IdGeneratorMock)

    shortnertService := NewShortnerService(urlRepositoryMock, identifierGenerator)

	originalUrl, _ := shortnertService.Resolve("2TX")

	urlRepositoryMock.AssertExpectations(t)
	urlRepositoryMock.AssertNumberOfCalls(t, "Get", 1)

	assert.Equal("https://www.mfvitale.me", originalUrl)
}

func TestIdentifierError(t *testing.T) {

    assert := assert.New(t)

	urlRepositoryMock := new(UrlRepositoryMock)
	identifierGenerator := new(IdGeneratorMock)

	identifierGenerator.On("NextID").Return(0, errors.New("error"))

	shortnertService := NewShortnerService(urlRepositoryMock, identifierGenerator)

	_, err := shortnertService.Shorten("https://www.mfvitale.me")

	assert.ErrorContains(err, "Error shortening url")
}

func TestRepositoryError(t *testing.T) {

    assert := assert.New(t)

	urlRepositoryMock := new(UrlRepositoryMock)
	identifierGenerator := new(IdGeneratorMock)

	identifierGenerator.On("NextID").Return(0, nil)
	urlRepositoryMock.On("Store", "0", "https://www.mfvitale.me").Return(errors.New("connect"))

	shortnertService := NewShortnerService(urlRepositoryMock, identifierGenerator)

	_, err := shortnertService.Shorten("https://www.mfvitale.me")

	assert.ErrorContains(err, "Error shortening url")
}