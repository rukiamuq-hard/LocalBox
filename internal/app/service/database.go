package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (s *Service) GetIdFromDB(login string) (int, error) {
	id, err := s.db.GetIdFromLogin(login)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) StoreFileToDB(r io.Reader, fileName string, storeFileName string, size int64, uploader_id string) error {
	Ext := filepath.Ext(fileName)
	ServCreateFile, err := os.Create("./uploads/" + storeFileName + Ext)
	if err != nil {
		return err
	}
	defer ServCreateFile.Close()

	if _, err = io.Copy(ServCreateFile, r); err != nil {
		return err
	}
	fmt.Println("db")
	TimeIs := time.Now().String()
	err = s.db.StoreFile(fileName, storeFileName+Ext, TimeIs, size, uploader_id)
	if err != nil {
		return err
	}
	fmt.Println("File uploaded: ", fileName)
	return nil
}

func (s *Service) GetFilesFromDB() ([]map[string]any, error) {
	files, err := s.db.GetFile()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, 0, len(files))
	for _, f := range files {
		result = append(result, map[string]any{
			"id":   f.ID,
			"name": f.Original_name,
			"type": "file",
			"date": f.Upload_date,
			"size": f.Size,
		})
	}
	return result, nil
}
