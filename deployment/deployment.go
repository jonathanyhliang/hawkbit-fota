package deployment

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	ErrDeploymentUpload         = errors.New("Deployment: upload image failed")
	ErrDeploymentUploadNotFound = errors.New("Deployment: upload not found")
	ErrDeploymentDist           = errors.New("Deployment: distribution set failed")
	ErrDeploymentDistNotFound   = errors.New("Deployment: distribution not found")
	ErrDeployment               = errors.New("Deployment: deployment set failed")
	ErrDeploymentNotFound       = errors.New("Deployment: deployment not found")
)

type Upload struct {
	Name    string `json:"name" example:"zephyr_cc3220sf_signed"`
	Version string `json:"version" example:"1.0.0+1"`
	File    string `json:"file" example:"/workdir/build/artifact.bin"`
	Sha256  string `json:"sha256,omitempty"`
	Size    int    `json:"size,omitempty"`
}

type Distribution struct {
	Name    string `json:"name" example:"hawkbit"`
	Version string `json:"version" example:"1.0.0+1"`
	Upload  Upload `json:"image,omitEmpty"`
}

type Status struct {
	Execution string `json:"execution"`
	Result    struct {
		Finished string `json:"finished"`
	} `json:"result"`
}

type Deployment struct {
	Target   string       `json:"target"`
	ActionId string       `json:"actionid"`
	Artifact Distribution `json:"artifact"`
	Status   Status       `json:"status"`
}

type hawkbitDeployment struct {
	mtx         sync.Mutex
	uploads     map[string]Upload
	artifacts   map[string]Distribution
	deployments map[string]Deployment
}

func SetUpload(u Upload) error {
	if u.Name == "" || u.Version == "" {
		return ErrDeploymentUpload
	}
	f, err := os.ReadFile(u.File)
	if err != nil {
		return ErrDeploymentUpload
	}
	u.Size = len(f)
	hash := sha256.Sum256(f)
	u.Sha256 = fmt.Sprintf("%x", hash)
	dp.mtx.Lock()
	defer dp.mtx.Unlock()
	dp.uploads[u.Name] = u

	return nil
}

func GetUpload(n string) (Upload, error) {
	dp.mtx.Lock()
	defer dp.mtx.Unlock()
	u, ok := dp.uploads[n]
	if !ok {
		return Upload{}, ErrDeploymentUploadNotFound
	}
	return u, nil
}

func SetDistribution(d Distribution, u string) error {
	if d.Name == "" || d.Version == "" {
		return ErrDeploymentDist
	}
	upl, ok := dp.uploads[u]
	if !ok {
		return ErrDeploymentDist
	}
	d.Upload = upl
	dp.mtx.Lock()
	defer dp.mtx.Unlock()
	dp.artifacts[d.Name] = d

	return nil
}

func GetDistribution(n string) (Distribution, error) {
	dp.mtx.Lock()
	defer dp.mtx.Unlock()
	d, ok := dp.artifacts[n]
	if !ok {
		return Distribution{}, ErrDeploymentDistNotFound
	}
	return d, nil
}

func SetDeployment(t string, d string) error {
	dp.mtx.Lock()
	defer dp.mtx.Unlock()
	a, ok := dp.artifacts[d]
	if !ok {
		return ErrDeployment
	}
	var n Deployment
	n.Target = t
	n.Artifact = a
	n.ActionId = a.Upload.Sha256[0:7]
	dp.deployments[t] = n

	return nil
}

func GetDeployment(t string) (Deployment, error) {
	dp.mtx.Lock()
	defer dp.mtx.Unlock()
	d, ok := dp.deployments[t]
	if !ok {
		return Deployment{}, ErrDeploymentNotFound
	}
	return d, nil
}

func UpdateStatus(t string, acid string, s Status) error {
	dp.mtx.Lock()
	defer dp.mtx.Unlock()
	d, ok := dp.deployments[t]
	if !ok {
		return ErrDeploymentNotFound
	}
	if acid == d.ActionId {
		d.Status = s
		dp.deployments[t] = d
	}
	return nil
}

var dp = &hawkbitDeployment{
	uploads:     map[string]Upload{},
	artifacts:   map[string]Distribution{},
	deployments: map[string]Deployment{},
}
