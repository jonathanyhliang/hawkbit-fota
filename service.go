package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

var (
	ErrDownloadImage = errors.New("download image failed")
)

type Controller struct {
	Config struct {
		Polling struct {
			Sleep string `json:"sleep"`
		} `json:"polling"`
	} `json:"config"`
	Links struct {
		DeploymentBase struct {
			Href string `json:"href,omitempty"`
		} `json:"deploymentBase"`
		CancelAction struct {
			Href string `json:"href,omitempty"`
		} `json:"cancelAction"`
		ConfigData struct {
			Href string `json:"href"`
		} `json:"configData"`
	} `json:"_links"`
}

type CancelActionFeedback struct {
	ID     string `json:"id"`
	Time   string `json:"time"`
	Status struct {
		Execution string `json:"execution"`
		Result    struct {
			Finished string `json:"finished"`
		} `json:"result"`
	} `json:"status"`
}

type ConfigData struct {
	Mode string `json:"mode"`
	Data struct {
		VIN        string `json:"vin"`
		HwRevision string `json:"hwRevision"`
	} `json:"data"`
	ID     string `json:"id"`
	Time   string `json:"time"`
	Status struct {
		Execution string `json:"execution"`
		Result    struct {
			Finished string `json:"finished"`
		} `json:"result"`
	} `json:"status"`
}

type artifacts struct {
	Filename string `json:"filename"`
	Hashes   struct {
		SHA1   string `json:"sha1"`
		MD5    string `json:"md5"`
		SHA256 string `json:"sha256"`
	} `json:"hashes"`
	Size  int `json:"size"`
	Links struct {
		DownloadHttp struct {
			Href string `json:"href"`
		} `json:"download-http"`
		MD5SumHttp struct {
			Href string `json:"href"`
		} `json:"md5sum-http"`
	} `json:"_links"`
}

type chunks struct {
	Part      string       `json:"part"`
	Name      string       `json:"name"`
	Version   string       `json:"version"`
	Artifacts [1]artifacts `json:"artifacts"`
}

type DeploymentBase struct {
	ID         string `json:"id"`
	Deployment struct {
		Download string    `json:"download"`
		Update   string    `json:"update"`
		Chunks   [1]chunks `json:"chunks"`
	} `json:"deployment"`
}

type DeploymentBaseFeedback struct {
	ID     string `json:"id"`
	Status struct {
		Execution string `json:"execution"`
		Result    struct {
			Finished string `json:"finished"`
		} `json:"result"`
	} `json:"status"`
}

type Service interface {
	GetController(ctx context.Context, bid string) (Controller, error)
	PostCancelActionFeedback(ctx context.Context, bid string, acid string, fb CancelActionFeedback) error
	PostConfigData(ctx context.Context, bid string, cfg ConfigData) error
	GetDeplymentBase(ctx context.Context, bid string) (DeploymentBase, error)
	PostDeploymentBaseFeedback(ctx context.Context, bid string, acid string, fb DeploymentBaseFeedback) error
	GetDownloadHttp(ctx context.Context, bid string, img string) ([]byte, error)
}

type hawkbitService struct{}

func NewHawkbitService() Service {
	return &hawkbitService{}
}

func (h *hawkbitService) GetController(ctx context.Context, bid string) (Controller, error) {
	var c Controller
	if d, err := dp.GetDeployment(bid); err != nil {
		d.BID = bid
		dp.PostDeployment(bid, d)
	} else {
		href := fmt.Sprintf("/default/controller/v1/%s/deploymentBase/%d", d.BID, d.Distribution.ACID)
		c.Links.DeploymentBase.Href = href
	}
	c.Config.Polling.Sleep = "00:05:00"
	c.Links.ConfigData.Href = "/default/controller/v1/" + bid + "/configData"
	return c, nil
}

func (h *hawkbitService) PostCancelActionFeedback(ctx context.Context, bid string, acid string,
	fb CancelActionFeedback) error {
	var d Deployment
	d.Status = fb.Status
	if err := dp.PutDeploymentStatus(bid, d); err != nil {
		return err
	}
	return nil

}

func (h *hawkbitService) PostConfigData(ctx context.Context, bid string, cfg ConfigData) error {
	return nil
}

func (h *hawkbitService) GetDeplymentBase(ctx context.Context, bid string) (DeploymentBase, error) {
	var db DeploymentBase
	d, err := dp.GetDeployment(bid)
	if err != nil {
		return DeploymentBase{}, err
	}
	db.ID = d.BID
	db.Deployment.Chunks[0].Part = "bApp"
	db.Deployment.Chunks[0].Version = d.Distribution.Artifact.Version
	db.Deployment.Chunks[0].Artifacts[0].Hashes.SHA256 = d.Distribution.Artifact.Sha256
	db.Deployment.Chunks[0].Artifacts[0].Size = d.Distribution.Artifact.Size
	href := "/DEFAULT/controller/v1/" + d.BID + "/softwareModules/" + d.Distribution.Artifact.Version
	db.Deployment.Chunks[0].Artifacts[0].Links.DownloadHttp.Href = href
	return db, nil
}

func (h *hawkbitService) PostDeploymentBaseFeedback(ctx context.Context, bid string, acid string,
	fb DeploymentBaseFeedback) error {
	var d Deployment
	d.Status = fb.Status
	if err := dp.PutDeploymentStatus(bid, d); err != nil {
		return err
	}
	return nil
}

func (h *hawkbitService) GetDownloadHttp(ctx context.Context, bid string, ver string) ([]byte, error) {
	d, err := dp.GetDeployment(bid)
	if err != nil {
		return nil, err
	}
	if ver != d.Distribution.Artifact.Version {
		return nil, err
	}
	f, err := os.ReadFile(d.Distribution.Artifact.File)
	if err != nil {
		return nil, ErrDownloadImage
	}
	return f, nil
}
