//blissWebApp.go, A Sample Web Application

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
)

type page struct {
	Title string
	Body  []byte
}

const (
	cstStrApp = "BlissWebApp"
	cstStrVer = "v1.0.0"
	cstIPPort = "8081"
	cstStrUTC = "Build LIHUI 2021-02-01 22:20:09.1919361 +0000 UTC"
)

type stMainPage struct {
	App string
	Ver string
	Prt string
	UTC string
	Dir string
}

var (
	mainPage stMainPage
)

var tmplMainPage = template.Must(template.New("tmplMainPage").Parse(`
<h1>{{.App}} </h1>
<body>
<p>Version: {{.Ver}}</p>
<p>Port No: {{.Prt}}</p>
<p>{{.UTC}}</p>
<p>Path: {{.Dir}}</p>
<table>
<tr style='text-align: left'>
<th>No.</th>
<th>Item</th>
<th>Description</th>
</tr>
<tr>
<th>1</th>
<th>Modbus</th>
<th>Modbus Test Tool</th>
</tr>
</table>
<body>
`))

func init() {
	mainPage.App = cstStrApp
	mainPage.Ver = cstStrVer
	mainPage.UTC = cstStrUTC
	mainPage.Prt = cstIPPort
}

func webServerRoutine() {
	http.HandleFunc("/", webHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":"+cstIPPort, nil)
}

func printAppInfo(w http.ResponseWriter, r *http.Request) {
	mainPage.Dir = fmt.Sprintf("%s", r.URL.Path)
	tmplMainPage.Execute(w, mainPage)
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	printAppInfo(w, r)
}

func startBrower() {
	cmd := exec.Command("explorer", "http://127.0.0.1:"+cstIPPort)
	cmd.Run()
}

func main() {
	go startBrower()
	webServerRoutine()
}

func (p *page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			printAppInfo(w, r)
			//http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
