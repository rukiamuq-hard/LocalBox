package endpoint

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (e *EndPoint) GetFiles(ctx *echo.Context) error {
	files, err := e.s.GetFilesFromDB()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, files)
}

func (e *EndPoint) UploadFile(ctx *echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "file not found"})
	}

	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()
	cookie, err := ctx.Cookie("account-id")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "can`t get cookie"})
	}

	err = e.s.StoreFileToDB(srcFile, file.Filename, e.s.MakeUUID(), file.Size, cookie.Value)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "file not stored"})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"message":  "File loaded succesfully",
		"filename": file.Filename,
	})
}

func (e *EndPoint) DownloadFile(ctx *echo.Context) error {
	id := ctx.Param("id")

	file, err := e.s.DownloadFile(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "file not found"})
	}

	return ctx.Attachment("./uploads/"+file.Stored_name, file.Original_name)
}
