package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"sms-webapp/controllers"
)

// "log"
// "net/http"

// "github.com/gorilla/mux"
// "github.com/jubrilissa/sms-webapp/controllers"

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", controllers.DashboardHandler)

	router.HandleFunc("/add-student", controllers.AddStudentHandler)
	router.HandleFunc("/add-subject", controllers.AddSubjectHandler)
	router.HandleFunc("/add-teacher", controllers.AddTeacherHandler)
	router.HandleFunc("/add-class", controllers.AddClassHandler).Methods("GET", "POST")

	router.HandleFunc("/students", controllers.ViewAllStudentHandler)
	router.HandleFunc("/subjects", controllers.ViewAllSubjectHandler)
	router.HandleFunc("/teachers", controllers.ViewAllTeacherHandler)
	router.HandleFunc("/classes", controllers.ViewAllClassHandler)

	// router.HandleFunc("/student/{id}", controllers.ViewSingleStudentHandler).Methods("GET", "POST")
	router.HandleFunc("/student-profile/{id:[0-9]+}", controllers.ViewSingleStudentHandler).Methods("GET", "POST")

	// s := http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/")))
	// router.PathPrefix("/templates/").Handler(s)
	// ServeStatic(router, "/templates/")
	router.PathPrefix("/js").Handler(http.HandlerFunc(handleJs))
	router.PathPrefix("/css").Handler(http.HandlerFunc(handleJs))
	router.PathPrefix("/img").Handler(http.HandlerFunc(handleJs))
	router.PathPrefix("/fonts").Handler(http.HandlerFunc(handleJs))
	router.PathPrefix("/style.css").Handler(http.HandlerFunc(handleJs))
	// router.Handle("/js", http.StripPrefix("/", http.FileServer(http.Dir("templates"))))

	log.Printf("Serving site on port 8000")
	http.ListenAndServe(":8000", router)

	// var dir = flag.String("dir", "../sms-frontend", "directory to serve")     // using
	// var listen = flag.String("listen", "localhost:8000", "Port to listen on") // using
	// flag.Parse()                                                              // using

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello Go!"))
	// })
	// http.HandleFunc("/tester", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Testing Go!"))
	// })
	// http.ListenAndServe(":8000", nil)
	// r := mux.NewRouter()                                       // using
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir(*dir))) // using
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("../sms-frontend")))

	// fs := http.FileServer(http.Dir("sms-frontend"))
	// http.Handle("/", fs)
	// listen := ":8000"
	// log.Printf("Serving static sites on %s from directory %s", *listen, *dir) // using
	// // http.ListenAndServe(":8000", r)
	// http.ListenAndServe(*listen, r) // using
}

func ServeStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"css":   staticDirectory + "/css/",
		"fonts": staticDirectory + "/fonts/",
		"pdf":   staticDirectory + "/pdf/",
		"img":   staticDirectory + "/img/",
		"js":    staticDirectory + "/js/",
	}
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir(pathValue))))
	}
}

// router := mux.NewRouter()
// ServeStatic(router, "/static/")

func handleJs(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.RemoteAddr)
	// path := strings.TrimPrefix(r.URL.Path, "/js/")
	fn := "templates" + r.URL.Path
	fmt.Printf("The path = %s\n the file is %s", r.URL.Path, fn)

	http.ServeFile(w, r, fn)
}
