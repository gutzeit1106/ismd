package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "html/template"
    "io/ioutil"
)

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("views/*")
    router.Static("/assets", "./assets")

    router.GET("/", func(ctx *gin.Context){
        //ctx.HTML(http.StatusOK, "index.html", gin.H{})
        html := template.Must(template.ParseFiles("views/base.tmpl", "views/index.tmpl"))
        router.SetHTMLTemplate(html)
        c.HTML(http.StatusOK, "base.tmpl", gin.H{})
    })

    router.Run(":8088")
}