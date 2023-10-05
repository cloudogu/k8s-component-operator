package client

import (
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/action"
	"testing"
)

func Test_provider_newInstall(t *testing.T) {
	// given
	sut := &provider{
		Configuration: &action.Configuration{},
		plainHttp:     true,
	}

	// when
	result := sut.newInstall()

	// then
	assert.NotEmpty(t, result.raw())
	assert.True(t, result.raw().PlainHTTP)
}

func Test_provider_newUpgrade(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newUpgrade()

	// then
	assert.NotEmpty(t, result.raw())
	assert.False(t, result.raw().PlainHTTP)
}

func Test_provider_newUninstall(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newUninstall()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newLocateChart(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newLocateChart()

	// then
	assert.NotEmpty(t, result)
}

func Test_provider_newListReleases(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newListReleases()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newGetReleaseValues(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newGetReleaseValues()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newGetRelease(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newGetRelease()

	// then
	assert.NotEmpty(t, result.raw())
}

func Test_provider_newRollbackRelease(t *testing.T) {
	// given
	sut := &provider{Configuration: &action.Configuration{}}

	// when
	result := sut.newRollbackRelease()

	// then
	assert.NotEmpty(t, result.raw())
}
