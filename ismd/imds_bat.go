package main

import (
    "fmt"
    "os"
    "log"
    "io/ioutil"
    "net/http"
    "time"
    "encoding/json"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
    "strconv"
)
type scheduledevents struct {
    DocumentIncarnation int `json:"DocumentIncarnation"`
    Events []Event `json:"Events"`
}

type Event struct {
    EventType int    `json:"EventType"`
    ResourceType string `json:"ResourceType"`
    Resources []string `json:"Resources"`
    EventStatus string `json:"EventStatus"`
    NotBefore string `json:"NotBefore"`
}   

func sendAlert(event string) {
	from := mail.NewEmail("xxx xxxx", "xxx@gmail.com")
	subject := "[Alert] Monitoring Instance Metadata Scheduled Events"
	to := mail.NewEmail("xxxx", "xxx@gmail.com")
    plainTextContent := "Instance Metadata に Scheduled Eventsが登録されました。"
    htmlContent := "<strong>Instance Metadata に Scheduled Eventsが登録されました。</strong>"
    htmlContent += "<strong>"+event+"</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	//client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	client := sendgrid.NewSendClient("xxxxxxx")
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
    }
}


func main() {

    var layout = "2006-01-02 15:04:05"
    var layout2 = "20060102"

    client := &http.Client{}

    req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/scheduledevents", nil)
    req.Header.Add("Metadata", "True")

    q := req.URL.Query()
    q.Add("format", "json")
    q.Add("api-version", "2017-08-01")
    req.URL.RawQuery = q.Encode()

    t := time.Now()
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Errored when sending request to the server")
        return
    }

    defer resp.Body.Close()
    resp_body, _ := ioutil.ReadAll(resp.Body)

    fmt.Println(resp.Status)
    fmt.Println(string(resp_body))
    fmt.Println(t)
    str := t.Format(layout)    
    line := str + ";" + resp.Status + ";" + string(resp_body)
    str2 := t.Format(layout2) 
    filepath := "/home/azureuser/go/src/ismd/logs/" + str2 + ".log"

    f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    fmt.Fprintln(f, line) //書き込み

    e := new(scheduledevents)
    if err := json.Unmarshal(resp_body, e);err != nil {
        log.Fatal(err)
    }

    fmt.Println(strconv.Itoa(e.DocumentIncarnation))
    if(len(e.Events)> 0){
        sendAlert(line) 
    }
}
