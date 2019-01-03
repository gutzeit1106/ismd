package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "html/template"
    "fmt"
    "os"
	"io/ioutil"
	"path/filepath"
	"bufio"
    "strings"
)

type MetaData struct {
    TimeStamp string
    ResponceCode string
    Message string
}

func fromFile(filePath string) []string {
	// ファイルを開く
	f, err := os.Open(filePath)
	if err != nil {
	  fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filePath, err)
	  os.Exit(1)
	}
  
	// 関数return時に閉じる
	defer f.Close()
  
	// Scannerで読み込む
	// lines := []string{}
	lines := make([]string, 0, 100)  // ある程度行数が事前に見積もれるようであれば、makeで初期capacityを指定して予めメモリを確保しておくことが望ましい
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
	  // appendで追加
	  lines = append(lines, scanner.Text())
	}
	if serr := scanner.Err(); serr != nil {
	  fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
	}
  
	return lines
}

func listFiles(rootPath, searchPath string) ([]string, []string ){
	fullpathlist := []string{}
	filelist := []string{}
	//old := time.Now().Add(time.Duration(-modifiedSpan) * time.Hour)
	fis, err := ioutil.ReadDir(searchPath)
	
    if err != nil {
        panic(err)
    }

    for _, fi := range fis {

		//if ( old.Before(fi.ModTime()) ){break}
		
		fullPath := filepath.Join(searchPath, fi.Name())
		fmt.Println(fullPath)
        if fi.IsDir() {
			subf,sub := listFiles(rootPath, fullPath)
			fullpathlist = append(fullpathlist, subf...)
			filelist = append(filelist, sub...)

		} else {
			fmt.Printf("Time:%v\n", fi.ModTime())
			fullpathlist = append(fullpathlist, fullPath)
			rel, err := filepath.Rel(rootPath, fullPath)

            if err != nil {
                panic(err)
			}
			filelist = append(filelist, rel)
		}

	}
	return fullpathlist, filelist
}

func ListMetaDataResult(path string)([] MetaData){

	//csv読み込んでMetaData配列でreturnする
	Metadatas := []MetaData{}
	fullpathlist, _ := listFiles(path, path)
	
	// 読み込み
	for fi := range fullpathlist {
		lines := fromFile(fullpathlist[fi])
		fmt.Printf("lines: %v\n", lines)
		
		for _, line := range lines {
			data := strings.Split(line, ";")
			Metadatas = append([] MetaData{ {data[0],data[1],data[2]} }, Metadatas...)
		}
	}
    
    return Metadatas
} 

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("views/*")
    router.Static("/assets", "./assets")

    router.GET("/", func(ctx *gin.Context){
        //ctx.HTML(http.StatusOK, "index.html", gin.H{})
        html := template.Must(template.ParseFiles("views/base.tmpl", "views/index.tmpl"))
        router.SetHTMLTemplate(html)
        ctx.HTML(http.StatusOK, "base.tmpl", gin.H{
            "Metadata": ListMetaDataResult("./logs"),
        })
    })

    router.Run(":8088")
}