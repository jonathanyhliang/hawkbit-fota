package main

import (
	"errors"
	"sync"
)

var (
	ErrDeploymentNotFound = errors.New("Deployment: not found")
	ErrDeploymentExist    = errors.New("Deployment: existed")
)

type distribution struct {
	Name     string `json:"name"`
	ACID     int    `json:"acid,omitempty"`
	Artifact struct {
		Name    string `json:"name"`
		File    string `json:"file"`
		Version string `json:"version"`
		Sha256  string `json:"sha256"`
		Size    int    `json:"size"`
	} `json:"artifact"`
}

type status struct {
	Execution string `json:"execution"`
	Result    struct {
		Finished string `json:"finished"`
	} `json:"result"`
}

type Deployment struct {
	BID          string       `json:"bid"`
	Distribution distribution `json:"distribution,omitempty"`
	Status       status       `json:"status,omitempty"`
}

type hawkbitDeployment struct {
	mtx         sync.Mutex
	deployments map[string]Deployment
}

func (h *hawkbitDeployment) GetDeployment(bid string) (Deployment, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if d, ok := h.deployments[bid]; ok {
		return d, nil
	}
	return Deployment{}, ErrDeploymentNotFound
}

func (h *hawkbitDeployment) PostDeployment(bid string, n Deployment) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if _, ok := h.deployments[bid]; !ok {
		h.deployments[bid] = n
		return nil
	}
	return ErrDeploymentExist
}

func (h *hawkbitDeployment) PutDeploymentDistribution(bid string, n Deployment) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if d, ok := h.deployments[bid]; ok {
		d.Distribution = n.Distribution
		h.deployments[bid] = d
		return nil
	}
	return ErrDeploymentNotFound
}

func (h *hawkbitDeployment) PutDeploymentStatus(bid string, n Deployment) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	if d, ok := h.deployments[bid]; ok {
		d.Status = n.Status
		h.deployments[bid] = d
		return nil
	}
	return ErrDeploymentNotFound
}

var dp = &hawkbitDeployment{
	deployments: map[string]Deployment{},
}
