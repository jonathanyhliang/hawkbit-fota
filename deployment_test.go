package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentGetNotFound(t *testing.T) {
	d, err := dp.GetDeployment("123")
	assert.Empty(t, d)
	assert.Equal(t, ErrDeploymentNotFound, err)
}

func TestDeploymentPutNotFound(t *testing.T) {
	var d Deployment
	err = dp.PutDeploymentDistribution("123", d)
	assert.Equal(t, ErrDeploymentNotFound, err)

	err = dp.PutDeploymentStatus("123", d)
	assert.Equal(t, ErrDeploymentNotFound, err)
}

func TestDeploymentPost(t *testing.T) {
	var d Deployment
	d.BID = "123"
	err = dp.PostDeployment("123", d)
	assert.Equal(t, nil, err)

	err = dp.PostDeployment("123", d)
	assert.Equal(t, ErrDeploymentExist, err)
}

func TestDeploymentGet(t *testing.T) {
	d, err := dp.GetDeployment("123")
	assert.Equal(t, "123", d.BID)
	assert.Equal(t, nil, err)
}

func TestDeploymentPut(t *testing.T) {
	var d, n Deployment
	d, err := dp.GetDeployment("123")
	assert.Equal(t, nil, err)

	d.Distribution.Name = "hawkbit"
	err = dp.PutDeploymentDistribution("123", d)
	n, err = dp.GetDeployment("123")
	assert.Equal(t, n, d)
	assert.Equal(t, nil, err)

	d.Status.Execution = "done"
	err = dp.PutDeploymentStatus("123", d)
	n, err = dp.GetDeployment("123")
	assert.Equal(t, n, d)
	assert.Equal(t, nil, err)
}
