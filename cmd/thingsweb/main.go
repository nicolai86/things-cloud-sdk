package main

//go:generate statik -src=./build/default

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	thingscloud "github.com/nicolai86/things-cloud-sdk"
	_ "github.com/nicolai86/things-cloud-sdk/cmd/thingsweb/statik"
	"github.com/nicolai86/things-cloud-sdk/state/memory"
	"github.com/rakyll/statik/fs"
)

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func sameHost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		handler.ServeHTTP(w, r)
	})
}

func proxyRequest(target string) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatal(err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	combined := logger(sameHost(proxy))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		combined.ServeHTTP(w, r)
	})
}

type projectAPI struct {
	state   *state
	project *thingscloud.Task
}

func (api *projectAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	listOpts := memory.ListOption{
		ExcludeCompleted: false,
		ExcludeInTrash:   false,
	}
	tasks := api.state.Subtasks(api.project, listOpts)
	clis := []*thingscloud.CheckListItem{}
	for _, task := range tasks {
		clis = append(clis, api.state.CheckListItemsByTask(task, listOpts)...)
	}
	json.NewEncoder(w).Encode(&taskResponse{
		Title:          fmt.Sprintf("Project %q", api.project.Title),
		Tasks:          tasks,
		ChecklistItems: clis,
	})
}

type areaAPI struct {
	state *state
	area  *thingscloud.Area
}

type taskResponse struct {
	Title          string                       `json:"title"`
	Tasks          []*thingscloud.Task          `json:"tasks"`
	ChecklistItems []*thingscloud.CheckListItem `json:"checklistItems"`
}

func (api *areaAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	listOpts := memory.ListOption{
		ExcludeCompleted: false,
		ExcludeInTrash:   false,
	}
	tasks := api.state.TasksByArea(api.area, listOpts)
	for _, task := range tasks {
		if !task.IsProject {
			continue
		}
		tasks = append(tasks, api.state.Subtasks(task, listOpts)...)
	}
	clis := []*thingscloud.CheckListItem{}
	for _, task := range tasks {
		clis = append(clis, api.state.CheckListItemsByTask(task, listOpts)...)
	}

	json.NewEncoder(w).Encode(&taskResponse{
		Title:          fmt.Sprintf("Area %q", api.area.Title),
		Tasks:          tasks,
		ChecklistItems: clis,
	})
}

func main() {
	historyID := flag.String("history", "", "things history id (optional)")
	projectName := flag.String("project", "", "things project to expose")
	areaName := flag.String("area", "", "things area to expose")
	username := flag.String("username", "", "things cloud username")
	password := flag.String("password", "", "things cloud password")
	store := flag.String("store", "", "persisted state")
	development := flag.Bool("development", false, "development mode")
	flag.Parse()

	if username == nil || *username == "" || password == nil || *password == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if (projectName == nil || *projectName == "") && (areaName == nil || *areaName == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	c := thingscloud.New(thingscloud.APIEndpoint, *username, *password)
	_, err := c.Verify()
	if err != nil {
		log.Fatalf("Login failed: %q\nPlease check your credentials.", err.Error())
	}

	hs, err := c.Histories()
	if err != nil {
		log.Fatalf("Failed to lookup histories: %q\n", err.Error())
	}
	var history *thingscloud.History
	if historyID == nil || *historyID == "" {
		history = hs[0]
	} else {
		for _, h := range hs {
			if h.ID == *historyID {
				history = h
			}
		}
	}

	s, err := load(*store)
	if err != nil {
		s = state{
			History: history,
			State:   memory.NewState(),
		}
	}
	s.History.Client = c
	fmt.Printf("using history %q since %d\n", s.History.ID, s.History.LatestServerIndex)
	items, _, err := s.History.Items(thingscloud.ItemsOptions{StartIndex: s.History.LatestServerIndex})
	if err != nil {
		log.Fatalf("Failed to lookup items: %q\n", err.Error())
	}
	log.Printf("Updating state with %d new items\n", len(items))
	if err := s.Update(items...); err != nil {
		log.Printf("Failed aggregating state: %q", err.Error())
	}
	save(*store, s)

	if *development {
		log.Println("Starting in development mode w/ proxy to polymer…")

		polymerServer := "http://localhost:8081/"
		http.Handle("/public/", http.StripPrefix("/public/", proxyRequest(polymerServer)))
		http.Handle("/bower_components/", proxyRequest(polymerServer))
		http.Handle("/src/", proxyRequest(polymerServer))
	} else {
		statikFS, err := fs.New()
		if err != nil {
			log.Fatalf(err.Error())
		}
		http.Handle("/public/", logger(http.StripPrefix("/public/", http.FileServer(statikFS))))
		http.Handle("/bower_components/", logger(http.FileServer(statikFS)))
		http.Handle("/src/", logger(http.FileServer(statikFS)))
	}

	if *projectName != "" {
		project := s.ProjectByName(*projectName)
		if project == nil {
			log.Fatalf("%q is no known project name\n", *projectName)
		}

		http.Handle("/api/", &projectAPI{&s, project})
	}

	if *areaName != "" {
		area := s.AreaByName(*areaName)
		if area == nil {
			log.Fatalf("%q is no known area name\n", *areaName)
		}

		http.Handle("/api/", &areaAPI{&s, area})
	}
	http.ListenAndServe(":8080", nil)

}
