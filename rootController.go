package myhouse

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// use a binifier go-bindata later for now load f rom ./templates
type asset struct {
	bytes []byte
	info  os.FileInfo
}

var homepageTpl *template.Template

//NewRouter produces the router with configured handlers
func MyNewRouter() *mux.Router {
	r := mux.NewRouter()
	// handle home
	r.HandleFunc("/", HomeHandler)
	// handle api routes
	r.HandleFunc("/api/{owner}/{device}/{state}", ActionHandler)
	// handle static routes
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	// Load templates
	homepageTpl = template.Must(template.New("homepage_view").Parse(string(homepageHTML.bytes)))
	return r
}

// Render a template, or server error.
func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}
func flagValue(flag bool, tval string, fval string) string {
	if flag {
		return tval
	}
	return fval
}

//HomeHandler handles the homepage and anything matching "/"
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//push(w, "/static/style.css")
	//push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	c, err := NewClient()
	if err != nil {
		log.Printf("Failed to get new client %v\n", err)
		http.Error(w, "Error fetching states!", http.StatusInternalServerError)
		return
	}

	cjLapLocked, cjPadLocked, sjLapLocked, sjPadLocked, err := c.Scanstates()
	if err != nil {
		log.Printf("Failed to get firewall status %v\n", err)
		http.Error(w, "Error fetching states!", http.StatusInternalServerError)
		return
	}
	fullData := map[string]interface{}{
		"time":          time.Now().Format(time.UnixDate),
		"sjLockClass":   flagValue(sjLapLocked && sjPadLocked, "btn-secondary", "btn-primary"),
		"sjUnlockClass": flagValue(!sjLapLocked && !sjPadLocked, "btn-secondary", "btn-primary"),
		"sjLapClass":    flagValue(sjLapLocked, "btn-danger", "btn-success"),
		"sjPadClass":    flagValue(sjPadLocked, "btn-danger", "btn-success"),
		"cjLockClass":   flagValue(cjLapLocked && cjPadLocked, "btn-secondary", "btn-primary"),
		"cjUnlockClass": flagValue(!cjLapLocked && !cjPadLocked, "btn-secondary", "btn-primary"),
		"cjLapClass":    flagValue(cjLapLocked, "btn-danger", "btn-success"),
		"cjPadClass":    flagValue(cjPadLocked, "btn-danger", "btn-success"),
		"cjLapAction":   flagValue(cjLapLocked, "unlock", "lock"),
		"cjPadAction":   flagValue(cjPadLocked, "unlock", "lock"),
		"sjLapAction":   flagValue(sjLapLocked, "unlock", "lock"),
		"sjPadAction":   flagValue(sjPadLocked, "unlock", "lock"),
	}

	render(w, r, homepageTpl, "homepage_view", fullData)
}

//ActionHandler handles the homepage and anything matching "/api/{owner}/{device}/{state}"
// Owner = cj or sj
// Device = ipad, laptop, or all
// State = lock or unlock
func ActionHandler(w http.ResponseWriter, r *http.Request) {
	//push(w, "/static/style.css")
	//push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	c, err := NewClient()
	if err != nil {
		log.Printf("Failed to get new client %v\n", err)
		http.Error(w, "Error fetching states!", http.StatusInternalServerError)
		return
	}
	pathVariables := mux.Vars(r)
	owner := pathVariables["owner"]
	device := pathVariables["device"]
	state := pathVariables["state"]
	enable := "0"
	log.Println("STATE = " + state)
	if state == "lock" {
		enable = ""
	}
	ruleNames := []string{}
	if device == "laptop" || device == "all" {
		//lock/unlock laptop
		ruleNames = append(ruleNames, "reject-"+owner+"-laptop-out")
		ruleNames = append(ruleNames, "reject-"+owner+"-laptop-in")
	}
	if device == "ipad" || device == "all" {
		//lock ipad
		ruleNames = append(ruleNames, "reject-"+owner+"-ipad-out")
		ruleNames = append(ruleNames, "reject-"+owner+"-ipad-in")
	}
	err = c.SetRulesEnabled(ruleNames, enable)
	if err != nil {
		log.Printf("Failed to get new client %v\n", err)
		http.Error(w, "Error fetching states!", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

var homepageHTML = asset{
	bytes: []byte(`<!DOCTYPE html>

	<head>
	  <meta charset="utf-8" />
	  <meta name="viewport" content="width=device-width, initial-scale=1" />
	  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/css/bootstrap.min.css" integrity="sha384-/Y6pD6FV/Vv2HJnA6t+vslU6fwYXjCFtcEpHbNJ0lyAFsXTsjBbfaDjzALeQsN6M" crossorigin="anonymous">
	  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
	  <meta http-equiv="refresh" content="120" />
	  <title>Get Control!</title>
	</head>
	
	<body>
	  <div class="container-fluid">
		<div class="jumbotron text-center">
		  <h1>Firewall</h1>
		  <p>At {{.time}}</p>
		</div>
		<div class="row">
		  <div class="col-1">&nbsp;</div>
		  
		  <div class="col-4">
			<div class="row">
			  <a  href="/api/savannah/all/lock" role="button" class="btn {{.sjLockClass}} btn-block">
				Savannah Lock
			  </a>
			</div>
			<div class="row">
			  <a  href="/api/savannah/laptop/{{.sjLapAction}}" role="button" class="btn {{.sjLapClass}} btn-block">
			  {{.sjLapAction}} Laptop
			  </a>
			</div>
			<div class="row">
			  <a  href="/api/savannah/ipad/{{.sjPadAction}}" role="button" class="btn {{.sjPadClass}} btn-block">
			  {{.sjPadAction}} Ipad
			  </a>
			</div>
			<div class="row">
			  <a  href="/api/savannah/all/unlock" role="button" class="btn {{.sjUnlockClass}} btn-block">
				Savannah Unlock
			  </a>
			</div>
		  </div>
		  <div class="col-1">&nbsp;</div>
		  <div class="col-4">
			<div class="row">
			  <a  href="/api/charlie/all/lock" role="button" class="btn {{.cjLockClass}} btn-block">
				Charlie Lock
			  </a>
			</div>
			<div class="row">
			  <a  href="/api/charlie/laptop/{{.cjLapAction}}" role="button" class="btn {{.cjLapClass}} btn-block">
			  {{.cjLapAction}} Laptop
			  </a>
			</div>
			<div class="row">
			  <a  href="/api/charlie/ipad/{{.cjPadAction}}" role="button" class="btn {{.cjPadClass}} btn-block">
			  {{.cjPadAction}} Ipad
			  </a>
			</div>
			<div class="row">
			  <a  href="/api/charlie/all/unlock" role="button" class="btn {{.cjUnlockClass}} btn-block">
				Charlie Unlock
			  </a>
			</div>
		  </div>
		  <div class="col-1">&nbsp;</div>
		  
		</div>
	  </div>
	
	  <script src="static/js/bootstrap.min.js"></script>
	</body>	
	`),
}
