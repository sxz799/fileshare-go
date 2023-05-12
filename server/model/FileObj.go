package model

import (
	"fileshare-server/gobalConfig"
	"fileshare-server/util"
	"log"
	"os"
	"time"
)

type FileObj struct {
	Id           int    `json:"id" gorm:"autoIncrement"`
	FileName     string `json:"fileName"`
	PathName     string `json:"pathName"`
	FileSize     int64  `json:"fileSize"`
	FileMd5      string `json:"fileMd5"`
	UploadDate   string `json:"uploadDate"`
	ShareCode    string `json:"shareCode"`
	FileLocation string `json:"fileLocation"`
	Finger       string `json:"finger"`
}

func listFilesByMd5(fileMd5 string) []FileObj {
	var fileObjs []FileObj
	util.DB.Where("file_md5=?", fileMd5).Find(&fileObjs)
	return fileObjs
}

func FileExist(fileName, fileFinger, fileMd5 string) (md5Exist, sameFinger, sameName bool, pathName, shareCode string) {

	files := listFilesByMd5(fileMd5)
	if len(files) > 0 {
		md5Exist = true
		pathName = files[0].PathName
	} else {
		return
	}
	for _, file := range files {
		if file.Finger == fileFinger {
			sameFinger = true
			break
		}
	}
	for _, file := range files {
		if file.Finger == fileFinger && file.FileName == fileName {
			sameName = true
			shareCode = file.ShareCode
			break
		}
	}

	return
}

func CodeExist(code string) bool {
	var fileObj FileObj
	return util.DB.Where("share_code=?", code).First(&fileObj) == nil
}

func CreateFile(fileName, fileCode, fileMd5, pathName, finger string, fileSize int64) {
	fileObj := FileObj{
		FileName:     fileName,
		PathName:     pathName,
		FileSize:     fileSize,
		FileMd5:      fileMd5,
		UploadDate:   time.Now().Format("2006-01-02 15:04:05"),
		ShareCode:    fileCode,
		FileLocation: "file/" + pathName,
		Finger:       finger,
	}
	util.DB.Create(&fileObj)
}

func GetFile(fileCode string) (FileObj, error) {
	var fileObj FileObj
	err := util.DB.Where("share_code=?", fileCode).First(&fileObj).Error
	return fileObj, err
}

func AutoDelFile() {
	var files []FileObj
	util.DB.Where("upload_date < ?", time.Now().Add(-time.Hour*time.Duration(gobalConfig.LimitFileLife)).Format("2006-01-02 15:04:05")).Find(&files)
	for _, file := range files {
		AddSystemLog("删除了文件："+file.PathName, "deleteFile")
		num := getNumsByPathName(file.PathName)
		if num == 1 {
			err := os.Remove("files/" + file.PathName)
			if err != nil {
				log.Println("自动文件删除失败：" + err.Error())
			}
		}
		util.DB.Delete(&file)
	}
}
func DeAllFile() {
	var files []FileObj
	util.DB.Find(&files)
	for _, file := range files {
		util.DB.Delete(&file)
		err := os.Remove("files/" + file.PathName)
		if err != nil {
			log.Println("全部文件删除失败：" + err.Error())
		}
	}
}

func ListFiles(finger string) (files []FileObj) {
	util.DB.Where("finger=?", finger).Find(&files)
	return
}

func getNumsByPathName(name string) (num int64) {
	util.DB.Debug().Model(FileObj{}).Where("path_name=?", name).Count(&num)
	return
}

func DelFile(shareCode string) error {
	var fileObj FileObj
	err := util.DB.Where("share_code=?", shareCode).First(&fileObj).Error
	if err != nil {
		return err
	}
	num := getNumsByPathName(fileObj.PathName)

	if num == 1 {
		err := os.Remove("files/" + fileObj.PathName)
		if err != nil {
			log.Println("文件删除失败：" + err.Error())
		}
	}
	util.DB.Delete(FileObj{}, "share_code=? ", shareCode)
	return nil
}
