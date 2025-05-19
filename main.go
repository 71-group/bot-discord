package main

import (
	"bot-discord/controller"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	files, err := loadTemplates("./website/html/")
	if err != nil {
		log.Println(err)
	}
	r.LoadHTMLFiles(files...)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/user-list")
	})

	r.GET("/user-list", controller.GetUserList)

	r.Static("/static", "./website/static/")
	err = r.Run(":80")
	if err != nil {
		fmt.Println(err)
	}
}

func loadTemplates(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if path != root {
				if _, err := loadTemplates(path); err != nil {
					fmt.Println("Erro ao carregar templates em", path, ":", err)
					return err
				}
			}
		} else {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
