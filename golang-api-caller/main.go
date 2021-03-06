package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
   	"log"
	"time"
	appd "appdynamics"
)

type Response struct {
	Name         string `json:"name"`
        Manufacturer string `json:"manufacturer"`
        ID           string `json:"id"`
        InPark       string `json:"inPark"`
        Height       int    `json:"height"`
}

func main() {

	cfg := appd.Config{}
	cfg.AppName = "GoLangExample"
	cfg.TierName = "APICaller"
	cfg.NodeName = "ApiCaller1"
	cfg.Controller.Host = "*.saas.appdynamics.com"
	cfg.Controller.Port = 443
	cfg.Controller.UseSSL = true
	cfg.Controller.Account = "*"
	cfg.Controller.AccessKey = ""
	cfg.InitTimeoutMs = 1000  // Wait up to 1s for initialization to finish

	if err := appd.InitSDK(&cfg); err != nil {
    		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
    		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}
        backendName := "Api1"//this is really the name
        backendType := "API"
        backendProperties := map[string]string {
                "API": "API",
        }
        resolveBackend := true

        appd.AddBackend(backendName, backendType, backendProperties, resolveBackend)

//        time.Sleep(20 * time.Second)
	for {
//		btHandle := appd.StartBT("Setting Coasters - main", "")
//		addcoaster()
		addcoaster()
		for i := 0; i < 10; i++ {
			getcoasters()
			time.Sleep(2*time.Second)
		}
		time.Sleep(20 * time.Second)
		// end the transaction
//		appd.EndBT(btHandle)
	}

	appd.TerminateSDK()
}

func getcoasters() {
        backendName := "Api1"
	btHandle2 := appd.StartBT("getting Coasters", "")

	time.Sleep(4 * time.Second)
	fmt.Println("Calling API...")

	client := &http.Client{}

	//start the exit call
	ecHandle := appd.StartExitcall(btHandle2, backendName)
	ecHeader := appd.GetExitcallCorrelationHeader(ecHandle)
	fmt.Println("echandle" + ecHeader)

	req, err := http.NewRequest("GET", "http://192.168.86.246:9090/coasters", nil)

	req.Header.Set(appd.APPD_CORRELATION_HEADER_NAME, ecHeader)


	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	hdr := req.Header.Get(appd.APPD_CORRELATION_HEADER_NAME)
       	fmt.Print("get header:" +hdr )

	// end the exit call
	appd.EndExitcall(ecHandle)

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject Response
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)
 	appd.EndBT(btHandle2)
}

func addcoaster() {
        btHandle2 := appd.StartBT("Adding Coasters", "")
	//Encode the data
   	postBody, _ := json.Marshal(map[string]string{
      		"id":  "Toby",
      		"name": "Toby@example.com",
		"manufacturer": "jaosn",
		"inPark": "yes it is",
   	})
   	responseBody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
   	resp, err := http.Post("http://192.168.86.246:9090/coasters", "application/json", responseBody)

	//Handle Error
	if err != nil {
      		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	//Read the response body
   	body, err := ioutil.ReadAll(resp.Body)
   	if err != nil {
      		log.Fatalln(err)
   	}
   	sb := string(body)
   	log.Printf(sb)
        appd.EndBT(btHandle2)
}
