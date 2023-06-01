package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUpload(t *testing.T) {
	var u Upload
	u.File = "./dummy.bin"
	u.Version = "1.2.3"
	u.Name = "test"
	err := SetUpload(u)
	assert.Equal(t, nil, err)
	u = dp.uploads["test"]
	assert.Equal(t, 27876, u.Size)
	assert.Equal(t, "dcb4462ed887656327579517f5dfc88daf7f6db5afdece7c945de5243d0fc305", u.Sha256)
}

func TestSetUploadInvalidFile(t *testing.T) {
	var u Upload
	u.File = "./dum.bin"
	u.Version = "1.2.3"
	u.Name = "test"
	err := SetUpload(u)
	assert.NotEqual(t, nil, err)
}

func TestSetUploadEmptyVersion(t *testing.T) {
	var u Upload
	u.File = "./dummy.bin"
	u.Version = ""
	u.Name = "test"
	err := SetUpload(u)
	assert.NotEqual(t, nil, err)
}

func TestSetUploadEmptyName(t *testing.T) {
	var u Upload
	u.File = "./dummy.bin"
	u.Version = "1.2.3"
	u.Name = ""
	err := SetUpload(u)
	assert.NotEqual(t, nil, err)
}
