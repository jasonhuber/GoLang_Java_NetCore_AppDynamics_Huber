package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
        appd "appdynamics"
)

type Response struct {
 ID     string `json:"id"`
 Joke   string `json:"joke"`
 Status int    `json:"status"`
}

type Coaster struct {
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID           string `json:"id"`
	InPark       string `json:"inPark"`
}

type coasterHandlers struct {
	sync.Mutex
	store map[string]Coaster
}


func (h *coasterHandlers) coasters(w http.ResponseWriter, r *http.Request) {
	my_bt_guid := "1234533"
	switch r.Method {
	case "GET":
		hdr := r.Header.Get(appd.APPD_CORRELATION_HEADER_NAME)
		fmt.Println("123:"+ hdr)

//		if reqHeadersBytes, err := json.Marshal(r.Header); err != nil {
//		    fmt.Println("Could not Marshal Req Headers")
//		} else {
//		    fmt.Println(string(reqHeadersBytes))
//		}

		fmt.Println("GET header:" +hdr )
                btHandle := appd.StartBT("Someone is Getting Coasters", hdr)
		// Optionally store the handle in the global registry
		appd.StoreBT(btHandle, my_bt_guid)
		othercalls(my_bt_guid)
                time.Sleep(7 * time.Second)
		h.get(w, r)
                // end the transaction
                appd.EndBT(btHandle)
		return
	case "POST":
		hdr := r.Header.Get(appd.APPD_CORRELATION_HEADER_NAME)
//                fmt.Print(r.Header)
//                fmt.Print("POST header:" +hdr )
                btHandle := appd.StartBT("Someone is Posting Coasters", hdr)
                // Optionally store the handle in the global registry
		appd.StoreBT(btHandle, my_bt_guid)
		time.Sleep(9 * time.Second)
		h.post(w, r)
                // end the transaction
                appd.EndBT(btHandle)
		return
	default:
		hdr := r.Header.Get(appd.APPD_CORRELATION_HEADER_NAME)
                btHandle := appd.StartBT("Someone is invalidating Coasters",hdr)
		// Optionally store the handle in the global registry
		appd.StoreBT(btHandle, my_bt_guid)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
        	// end the transaction
	        appd.EndBT(btHandle)
		return
	}

	client := &http.Client{}
 	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
 	if err != nil {
  		fmt.Print(err.Error())
 	}

	req.Header.Add("Accept", "application/json")
 	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
  		fmt.Print(err.Error())
 	}

	defer resp.Body.Close()
 	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
  		fmt.Print(err.Error())
 	}

	var responseObject Response
 	json.Unmarshal(bodyBytes, &responseObject)
 	fmt.Printf("API Response as struct %+v\n", responseObject)

	appd.TerminateSDK()
}

func (h *coasterHandlers) get(w http.ResponseWriter, r *http.Request) {
	coasters := make([]Coaster, len(h.store))

	h.Lock()
	i := 0
	for _, coaster := range h.store {
		coasters[i] = coaster
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(coasters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *coasterHandlers) getRandomCoaster(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.store))
	h.Lock()
	i := 0
	for id := range h.store {
		ids[i] = id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	w.Header().Add("location", fmt.Sprintf("/coasters/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *coasterHandlers) getCoaster(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomCoaster(w, r)
		return
	}

	h.Lock()
	coaster, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(coaster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *coasterHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var coaster Coaster
	err = json.Unmarshal(bodyBytes, &coaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	coaster.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[coaster.ID] = coaster
	defer h.Unlock()
}

func newCoasterHandlers() *coasterHandlers {
	return &coasterHandlers{
		store: map[string]Coaster{},
	}

}

type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := "lakers" //os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}

	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - unauthorized"))
		return
	}

	w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}

func othercalls(my_bt_guid string) {
	backendName := "Third Party API"
	backendType := "API"
	backendProperties := map[string]string {
    		"API": "icanhazjoke",
	}
	resolveBackend := false

	appd.AddBackend(backendName, backendType, backendProperties, resolveBackend)
	// Retrieve a stored handle from the global registry
	btHandle := appd.GetBT(my_bt_guid)
	ecHandle := appd.StartExitcall(btHandle, backendName)

	client := &http.Client{}
        req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/", nil)
        if err != nil {
                fmt.Print(err.Error())
        }

        req.Header.Add("Accept", "application/json")
        req.Header.Add("Content-Type", "application/json")
        resp, err := client.Do(req)

        if err != nil {
                fmt.Print(err.Error())
        }

        defer resp.Body.Close()
        bodyBytes, err := ioutil.ReadAll(resp.Body)

        if err != nil {
                fmt.Print(err.Error())
        }

        var responseObject Response
        json.Unmarshal(bodyBytes, &responseObject)
        fmt.Printf("API Response as struct %+v\n", responseObject)
        // end the exit call
	appd.EndExitcall(ecHandle)

}

//func printheaders() {

   // Loop over all values for the name.
//    for _, value := range values {
   //     fmt.Println(name, value)
   // }
//}
//}

func main() {
        cfg := appd.Config{}
        cfg.AppName = "GoLangExample"
        cfg.TierName = "API"
        cfg.NodeName = "Api1"
        cfg.Controller.Host = "*.saas.appdynamics.com"
        cfg.Controller.Port = 443
        cfg.Controller.UseSSL = true
        cfg.Controller.Account = "*"
        cfg.Controller.AccessKey = ""
        cfg.InitTimeoutMs = 1000  // Wait up to 1s for initialization to fini>

        if err := appd.InitSDK(&cfg); err != nil {
                fmt.Printf("Error initializing the AppDynamics SDK\n")
        } else {
                fmt.Printf("Initialized AppDynamics SDK successfully\n")
        }

	admin := newAdminPortal()
	coasterHandlers := newCoasterHandlers()
	http.HandleFunc("/coasters", coasterHandlers.coasters)
	http.HandleFunc("/coasters/", coasterHandlers.getCoaster)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(err)
	}
}
