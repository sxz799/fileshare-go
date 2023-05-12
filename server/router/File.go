package router

import (
	"crypto/md5"
	"encoding/hex"
	"fileshare-server/gobalConfig"
	"fileshare-server/model"
	"fileshare-server/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/url"
	"time"
)

func upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	finger, _ := c.GetPostForm("finger")
	pFile, err := file.Open()
	if err != nil {
		c.JSON(200, model.Result{
			Status:  "error",
			Success: false,
			Message: fmt.Sprintf("校验文件失败：%s", err),
		})
		return
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)
	fileMd5 := hex.EncodeToString(md5h.Sum(nil))
	md5Exist, sameFinger, sameName, pathName, shareCode := model.FileExist(file.Filename, finger, fileMd5)
	if len(pathName) < 1 {
		pathName = time.Now().Format("20060102150405")
	}
	if md5Exist {
		if sameFinger {
			if sameName {
				c.JSON(200, model.Result{
					Status:  "success",
					Success: true,
					FileObj: model.FileObj{ShareCode: shareCode},
					Message: "文件已存在！",
				})
			} else {
				success, code := GenerateCode()
				if !success {
					c.JSON(200, model.Result{
						Status:  "error",
						Success: true,
						Message: "提取码生成失败,请重试！",
					})
					return
				}
				model.CreateFile(file.Filename, code, fileMd5, pathName, finger, file.Size)
				model.AddSystemLog(c.RemoteIP()+"...上传了文件："+file.Filename, "upload")
				c.JSON(200, model.Result{
					Status:  "success",
					Success: true,
					FileObj: model.FileObj{ShareCode: code},
					Message: "文件秒传传成功！",
				})
			}
		} else {
			success, code := GenerateCode()
			if !success {
				c.JSON(200, model.Result{
					Status:  "error",
					Success: true,
					Message: "提取码生成失败,请重试！",
				})
				return
			}
			model.CreateFile(file.Filename, code, fileMd5, pathName, finger, file.Size)
			model.AddSystemLog(c.RemoteIP()+"...上传了文件："+file.Filename, "upload")
			c.JSON(200, model.Result{
				Status:  "success",
				Success: true,
				FileObj: model.FileObj{ShareCode: code},
				Message: "文件秒传传成功！",
			})
		}

	} else {
		err2 := c.SaveUploadedFile(file, "./files/"+pathName)
		if err2 != nil {
			c.JSON(200, model.Result{
				Status:  "error",
				Success: false,
				Message: fmt.Sprintf("%s 保存失败!", file.Filename),
			})
			return
		}
		success, code := GenerateCode()
		if !success {
			c.JSON(200, model.Result{
				Status:  "error",
				Success: true,
				Message: "提取码生成失败,请重试！",
			})
			return
		}
		model.CreateFile(file.Filename, code, fileMd5, pathName, finger, file.Size)
		model.AddSystemLog(c.RemoteIP()+"...上传了文件："+file.Filename, "upload")
		c.JSON(200, model.Result{
			Status:  "success",
			Success: true,
			FileObj: model.FileObj{ShareCode: code},
			Message: "文件上传成功！",
		})
	}

}

func download(c *gin.Context) {
	code := c.Param("code")
	file, err := model.GetFile(code)
	if err != nil {
		c.String(200, "提取码不存在或文件已过期！")
	} else {
		model.AddSystemLog(c.RemoteIP()+"...下载了文件："+file.FileName, "download")
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+file.FileName)
		c.Header("filename", url.QueryEscape(file.FileName))
		c.Header("Content-Transfer-Encoding", "binary")
		c.File("./files/" + file.PathName)
	}
}

func exist(c *gin.Context) {
	code := c.Param("code")
	file, err := model.GetFile(code)
	if err != nil {
		c.JSON(200, model.Result{
			Status:  "error",
			Success: false,
			Message: "提取码不存在或文件已过期！",
		})
	} else {
		c.JSON(200, model.Result{
			Status:  "success",
			Success: true,
			FileObj: file,
		})
	}
}
func config(c *gin.Context) {
	c.JSON(200, gin.H{
		"limitFileLife": gobalConfig.LimitFileLife,
		"limitFileSize": gobalConfig.LimitFileSize,
	})
}
func list(c *gin.Context) {
	finger := c.Query("finger")
	files := model.ListFiles(finger)
	c.JSON(200, files)
}

func reset(c *gin.Context) {
	model.DeAllFile()
}
func del(c *gin.Context) {
	shareCode := c.Param("shareCode")
	err := model.DelFile(shareCode)
	if err != nil {
		c.JSON(200, model.Result{
			Status:  "success",
			Success: false,
			Message: "文件删除失败！" + err.Error(),
		})
	} else {
		c.JSON(200, model.Result{
			Status:  "success",
			Success: true,
			Message: "文件删除成功！",
		})
	}
}

func File(e *gin.Engine) {
	g := e.Group("/api")
	{
		g.POST("/upload", upload)
		g.GET("/exist/:code", exist)
		g.GET("/download/:code", download)
		g.GET("/config", config)
		g.GET("/list", list)
		g.GET("/del/:shareCode", del)
		g.GET("/reset", reset)
	}
}

func GenerateCode() (bool, string) {
	t := 0
	for {
		if t > 100 {
			return false, ""
		}
		codeLength := gobalConfig.ShareCodeLength
		var code string
		switch gobalConfig.ShareCodeType {
		case 1:
			code = util.RandStringAll(codeLength)
		case 2:
			code = util.RandStringLarge(codeLength)
		case 3:
			code = util.RandStringSmall(codeLength)
		case 4:
			code = util.RandStringNum(codeLength)
		case 5:
			code = util.RandStringLargeSmall(codeLength)
		case 6:
			code = util.RandStringLargeNum(codeLength)
		case 7:
			code = util.RandStringSmallNum(codeLength)
		default:
			code = util.RandStringAll(codeLength)
		}

		if model.CodeExist(code) {
			t++
			continue
		} else {
			return true, code
		}
	}
}
