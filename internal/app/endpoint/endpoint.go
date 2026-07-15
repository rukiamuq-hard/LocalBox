package endpoint

import (
	"Umbrella/internal/app/repository/database/models"
	"fmt"
	"io"

	"github.com/labstack/echo/v5"
)

type Service interface {
	RegisterUser(login string, password string) error
	LoginUser(login string, password string) error
	GetIdFromDB(login string) (int, error)
	RedisSetKeyValue(key string, value int) error
	MakeUUID() string
	StoreFileToDB(r io.Reader, fileName string, storedFileName string, size int64, uploader_id string) error
	GetFilesFromDB() ([]map[string]any, error)
	DownloadFile(id string) (models.UploadedFiles, error)
}

type EndPoint struct {
	s Service
}

func New(svc Service) *EndPoint {
	return &EndPoint{
		s: svc,
	}
}

func (e *EndPoint) LoadMainHTML(ctx *echo.Context) error {
	err := ctx.File("website/LocalCloudMain.html")
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (e *EndPoint) LoadDashboard(ctx *echo.Context) error {
	err := ctx.File("website/dashboard.html")
	if err != nil {
		fmt.Println(err)
	}
	return err
}
