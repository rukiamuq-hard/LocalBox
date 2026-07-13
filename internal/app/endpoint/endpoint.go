package endpoint

import (
	"Umbrella/internal/app/repository/SQLite/models"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

const CookieLiveTime = 20

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

func (e *EndPoint) Register(ctx *echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")

	fmt.Println("User try to register: ", login)
	if len(password) > 14 || len(password) < 4 {
		fmt.Println("Password size more than 12 length")
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	err := e.s.RegisterUser(login, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		return errors.New("error with register user, try again later")
	}
	return ctx.Redirect(http.StatusMovedPermanently, "/login.html")
}

func (e *EndPoint) Login(ctx *echo.Context) error {
	login := ctx.FormValue("login")
	password := ctx.FormValue("password")

	fmt.Println("User try to login: ", login)
	fmt.Println("Password for Login: ", password)

	err := e.s.LoginUser(login, password)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {

			fmt.Println("Error, account alreade exist: ", err.Error())
			return ctx.Redirect(http.StatusSeeOther, "/")
		}
		fmt.Println("Error with login: ", err)
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	cookie := new(http.Cookie)
	cookie.Name = "loggin_token"
	cookie.Value = e.s.MakeUUID()
	cookie.Path = "/"

	id, err := e.s.GetIdFromDB(login)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}
	err = e.s.RedisSetKeyValue(cookie.Value, id)
	if err != nil {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}
	cookie.Expires = time.Now().Add(20 * time.Minute)
	cookie.HttpOnly = true

	cookieID := &http.Cookie{
		Name:     "account-id",
		Value:    strconv.Itoa(id),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(20 * time.Minute),
	}

	ctx.SetCookie(cookie)
	ctx.SetCookie(cookieID)

	return ctx.Redirect(http.StatusSeeOther, "/dashboard.html")
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

func (e *EndPoint) GetFiles(ctx *echo.Context) error {
	files, err := e.s.GetFilesFromDB()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, files)
}

func (e *EndPoint) DownloadFile(ctx *echo.Context) error {
	id := ctx.Param("id")

	file, err := e.s.DownloadFile(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "file not found"})
	}

	return ctx.Attachment("./uploads/"+file.Stored_name, file.Original_name)
}
