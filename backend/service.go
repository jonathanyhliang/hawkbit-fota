package backend

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jonathanyhliang/hawkbit-fota/deployment"
)

var (
	ErrBackendDownload   = errors.New("Backend: image download failed")
	ErrBackendBadRequest = errors.New("Backend: bad request")
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

type BackendService interface {
	GetController(ctx context.Context, bid string) (Controller, error)
	PostCancelActionFeedback(ctx context.Context, bid string, fb CancelActionFeedback) error
	PostConfigData(ctx context.Context, bid string, cfg ConfigData) error
	GetDeplymentBase(ctx context.Context, bid string, acid string) (DeploymentBase, error)
	PostDeploymentBaseFeedback(ctx context.Context, bid string, fb DeploymentBaseFeedback) error
	GetDownloadHttp(ctx context.Context, bid string, img string) ([]byte, error)
}

type hawkbitBackendService struct{}

func NewHawkbitBackendService() BackendService {
	return &hawkbitBackendService{}
}

func (h *hawkbitBackendService) GetController(ctx context.Context, bid string) (Controller, error) {
	var c Controller
	if d, err := deployment.GetDeployment(bid); err == nil {
		href := fmt.Sprintf("/default/controller/v1/%s/deploymentBase/%d", d.Target, d.ActionId)
		c.Links.DeploymentBase.Href = href
	}
	c.Config.Polling.Sleep = "00:05:00"
	c.Links.ConfigData.Href = "/default/controller/v1/" + bid + "/configData"
	return c, nil
}

func (h *hawkbitBackendService) PostCancelActionFeedback(ctx context.Context, bid string,
	fb CancelActionFeedback) error {
	actionId, err := strconv.Atoi(fb.ID)
	if err != nil {
		return ErrBackendBadRequest
	}
	if err = deployment.UpdateStatus(bid, actionId, fb.Status); err != nil {
		return err
	}
	return nil

}

func (h *hawkbitBackendService) PostConfigData(ctx context.Context, bid string, cfg ConfigData) error {
	return nil
}

func (h *hawkbitBackendService) GetDeplymentBase(ctx context.Context, bid string,
	acid string) (DeploymentBase, error) {
	actionId, err := strconv.Atoi(acid)
	if err != nil {
		return DeploymentBase{}, ErrBackendBadRequest
	}
	d, err := deployment.GetDeployment(bid)
	if err != nil {
		return DeploymentBase{}, err
	}
	if actionId != d.ActionId {
		return DeploymentBase{}, ErrBackendBadRequest
	}

	var db DeploymentBase
	db.ID = d.Target
	db.Deployment.Chunks[0].Part = "bApp"
	db.Deployment.Chunks[0].Version = d.Artifact.Version
	db.Deployment.Chunks[0].Artifacts[0].Hashes.SHA256 = d.Artifact.Upload.Sha256
	db.Deployment.Chunks[0].Artifacts[0].Size = d.Artifact.Upload.Size
	href := "/DEFAULT/controller/v1/" + bid + "/softwareModules/" + d.Artifact.Upload.Version
	db.Deployment.Chunks[0].Artifacts[0].Links.DownloadHttp.Href = href
	return db, nil
}

func (h *hawkbitBackendService) PostDeploymentBaseFeedback(ctx context.Context, bid string,
	fb DeploymentBaseFeedback) error {
	acid, err := strconv.Atoi(fb.ID)
	if err != nil {
		return ErrBackendBadRequest
	}
	if err = deployment.UpdateStatus(bid, acid, fb.Status); err != nil {
		return err
	}
	return nil
}

func (h *hawkbitBackendService) GetDownloadHttp(ctx context.Context, bid string, ver string) ([]byte, error) {
	d, err := deployment.GetDeployment(bid)
	if err != nil {
		return nil, err
	}
	if ver != d.Artifact.Upload.Version {
		return nil, err
	}
	f, err := os.ReadFile(d.Artifact.Upload.File)
	if err != nil {
		return nil, ErrBackendDownload
	}
	return f, nil
}
