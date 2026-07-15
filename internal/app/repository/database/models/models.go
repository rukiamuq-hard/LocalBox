package models

type UploadedFiles struct {
	ID            int
	Original_name string
	Stored_name   string
	Upload_date   string
	Size          int
	Uploader_id   string
}
