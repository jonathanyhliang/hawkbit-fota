package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBodyEmpty(t *testing.T) {
	response, err := parseBody([]byte{})
	assert.Empty(t, response)
	assert.NotEqual(t, nil, err)
}

func TestParseBodyDeployNotFound(t *testing.T) {
	var d Deployment
	d.BID = "12345"
	request, err := json.Marshal(d)
	response, err := parseBody(request)
	assert.Empty(t, response)
	assert.Equal(t, ErrDeploymentNotFound, err)
}

func TestParseBody(t *testing.T) {
	var d Deployment
	d.BID = "12345"
	d.Status.Execution = "done"
	d.Status.Result.Finished = "none"
	err := dp.PostDeployment("12345", d)
	assert.Equal(t, nil, err)

	d.ACID = 1
	d.Distribution.Name = "cc3220"
	d.Distribution.Artifact.Name = "hawkbit"
	d.Distribution.Artifact.File = "/usr/bin"
	d.Distribution.Artifact.Version = "0.0.0"
	d.Distribution.Artifact.Sha256 = "abcde"
	d.Distribution.Artifact.Size = 2000
	request, err := json.Marshal(d)
	response, err := parseBody(request)

	var n Deployment
	err = json.Unmarshal(response, &n.Status)
	assert.Equal(t, n.Status, d.Status)
	assert.Equal(t, nil, err)
}
