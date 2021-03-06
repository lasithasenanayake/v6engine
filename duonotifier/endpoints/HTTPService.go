package endpoints

import (
	"duov6.com/duonotifier/client"
	"duov6.com/duonotifier/messaging"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"io/ioutil"
	"net/http"
)

type HTTPService struct {
}

func (h *HTTPService) Start() {
	fmt.Println("DuoNotifier Listening on Port : 7000")
	m := martini.Classic()
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"securityToken", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	m.Get("/", versionHandler)
	m.Post("/:namespace", handleRequest)
	m.RunOnAddr(":7000")
}

func versionHandler(params martini.Params, w http.ResponseWriter, r *http.Request) {
	versionData := "{\"name\": \"DuoNotifier\",\"version\": \"1.0.0-a\",\"Change Log\":\"Added Logs!\",\"author\": {\"name\": \"Duo Software\",\"url\": \"http://www.duosoftware.com/\"},\"repository\": {\"type\": \"git\",\"url\": \"https://github.com/DuoSoftware/v6engine/\"}}"
	fmt.Fprintf(w, versionData)
}

func handleRequest(params martini.Params, w http.ResponseWriter, r *http.Request) {
	namespace := params["namespace"]
	requestBody, _ := ioutil.ReadAll(r.Body)

	switch namespace {
	case "GetTemplate":
		request := getTemplateRequest(requestBody)
		response := client.GetTemplate(request)
		temp, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", string(temp))
		break
	default:
		fmt.Fprintf(w, "%s", "Method Not Found!")
		break
	}

}

func getTemplateRequest(body []byte) messaging.TemplateRequest {
	var templateRequest messaging.TemplateRequest

	fmt.Println("--------------------------------------------")
	fmt.Println("Request in String : ")
	fmt.Println(string(body))
	fmt.Println("Request in Map : ")
	err := json.Unmarshal(body, &templateRequest)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(templateRequest)
	}
	fmt.Println("--------------------------------------------")
	return templateRequest
}
