package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// 下载图片
func DownImage(imgUrl string, savePath string) (filename string, err error) {
	//fmt.Println(beego.WorkPath)
	appPath, _ := os.Getwd()
	fileBaseName := path.Base(imgUrl)
	// 图片保存的相对路径
	imgPath := filepath.Join(savePath, "/")
	//imgPath := filepath.Join(savePath, time.Now().Format("200601"))
	filename = filepath.Join(imgPath, fileBaseName)
	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	// 自动创建文件夹
	if err = CheckDir(appPath + imgPath); err != nil {
		return
	}
	f, err := os.OpenFile(appPath+filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return
	} else {
		_, err = f.Write(b)
		if err != nil {
			return
		}
	}
	return
}

func CheckDir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else {
		err := os.MkdirAll(path, 0711)
		if err != nil {
			return err
		}
	}
	// check again
	_, err := os.Stat(path)
	return err
}
