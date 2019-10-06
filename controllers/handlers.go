package controllers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"sms-webapp/models"
)

const (
	templatesDir = "templates"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("templates/index.html"))

	files := []string{
		filepath.Join(templatesDir, "index.html"),
		filepath.Join(templatesDir, "base.html"),
	}
	// tmpl := template.Must(template.
	// 	ParseFiles(
	// 		filepath.Join(templatesDir, "base.html"),
	// 		filepath.Join(templatesDir, "index.html"),
	// 	))

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func ViewAllTeacherHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "all-teachers.html"),
		filepath.Join(templatesDir, "base.html"),
	}
	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func ViewAllStudentHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		data := models.GetAllStudents()
		// fmt.Printf("%+v\n", data[0])
		fmt.Println(data[0].Name)

		files := []string{
			filepath.Join(templatesDir, "all-students.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			panic(err.Error())
		}

		// tmpl := template.Must(template.
		// 	ParseFiles(files...))

		tmpl.Execute(w, &data)
		// tmpl.Execute(w, nil)

	} else {

		fmt.Println("I got to the else block")

		r.ParseMultipartForm(32 << 20)
		// r.ParseForm()

		// r.ParseForm()

		// r.ParseMultipartForm(32 << 20)
		// fmt.Println("I am testing the file name")
		// fmt.Println(r.FormFile("image"))

		// var Buf bytes.Buffer
		// // in your case file would be fileupload
		// file, header, err := r.FormFile("image")
		// if err != nil {
		// 	panic(err)
		// }
		// defer file.Close()
		// file_name := strings.Split(header.Filename, ".")
		// fmt.Printf("File name %s\n", file_name[0])
		// Copy the file data to my buffer
		// io.Copy(&Buf, file)
		// do something with the contents...
		// I normally have a struct defined and unmarshal into a struct, but this will
		// work as an example
		// contents := Buf.String()
		// fmt.Println(contents)
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		// Buf.Reset()
		// do something else
		// etc write header

		// logic part of log in
		name := r.FormValue("fullName")
		address := r.FormValue("address")
		mobileno := r.FormValue("phoneNo")
		religion := r.FormValue("religion")
		dateOfBirthString := r.FormValue("dateOfBirth")
		gender := r.FormValue("gender")
		class := r.FormValue("class")
		religionInterest := r.FormValue("religionInterest")

		fmt.Println("Firstname:", name)
		fmt.Println("address:", address)
		fmt.Println("mobileno:", mobileno)
		fmt.Println("gender:", religion)
		dateOfBirth, _ := time.Parse("2006-01-02", dateOfBirthString)

		student := &models.Student{}
		student.Name = name
		student.Address = address
		student.PhoneNo = mobileno
		student.Religion = religion
		student.Class = class
		student.Gender = gender
		student.DateOfBirth = &dateOfBirth
		// student.DateOfBirth = dateOfBirth

		// student.DateOfBirth = dateOfBirth

		fmt.Println("The date of birth is", dateOfBirth)
		fmt.Println("The religion interest is", religionInterest)
		// student.DateOfBirth = dateOfBirth

		file, handler, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
			return
		}

		student.Image = handler.Filename
		fmt.Println(student.ID)
		studentID := student.Create()

		fmt.Println("Time to redirect")
		redirectURL := fmt.Sprintf("/student-profile/%d", studentID)
		fmt.Println(redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusFound)

		defer file.Close()

		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./templates/img/student/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)

	}

}

func ViewAllSubjectHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "all-subjects.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func ViewAllClassHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "classes.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func ViewSingleStudentHandler(w http.ResponseWriter, r *http.Request) {
	// files := []string{
	// 	filepath.Join(templatesDir, "student-profile.html"),
	// 	filepath.Join(templatesDir, "base.html"),
	// }

	// tmpl := template.Must(template.
	// 	ParseFiles(files...))

	// tmpl.Execute(w, nil)

	// switch r.Method {
	// 	case "GET":
	// 		 http.ServeFile(w, r, "form.html")
	// 	case "POST":

	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		data := models.GetSingleStudentById(uint(id))

		files := []string{
			filepath.Join(templatesDir, "student-profile.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		fmt.Println(data)
		tmpl.Execute(w, data)

	} else {
		fmt.Println("I got to the else block")

		r.ParseMultipartForm(32 << 20)
		// r.ParseForm()
		file, handler, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)

		// r.ParseForm()

		// r.ParseMultipartForm(32 << 20)
		// fmt.Println("I am testing the file name")
		// fmt.Println(r.FormFile("image"))

		// var Buf bytes.Buffer
		// // in your case file would be fileupload
		// file, header, err := r.FormFile("image")
		// if err != nil {
		// 	panic(err)
		// }
		// defer file.Close()
		// file_name := strings.Split(header.Filename, ".")
		// fmt.Printf("File name %s\n", file_name[0])
		// Copy the file data to my buffer
		// io.Copy(&Buf, file)
		// do something with the contents...
		// I normally have a struct defined and unmarshal into a struct, but this will
		// work as an example
		// contents := Buf.String()
		// fmt.Println(contents)
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		// Buf.Reset()
		// do something else
		// etc write header

		// logic part of log in
		name := r.FormValue("firstname")
		address := r.FormValue("address")
		mobileno := r.FormValue("mobileno")
		religion := r.FormValue("religion")
		fmt.Println("Firstname:", name)
		fmt.Println("address:", address)
		fmt.Println("mobileno:", mobileno)
		fmt.Println("gender:", religion)

		student := &models.Student{}
		student.Name = name
		student.Address = address
		student.PhoneNo = mobileno
		student.Religion = religion

		fmt.Println(student.ID)
		student.Create()
		// fmt.Println("Student struct:", student)

		fmt.Println("Time to redirect")
		http.Redirect(w, r, "/students", http.StatusTemporaryRedirect)

	}
}

func AddSubjectHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "add-subject.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "add-teacher.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func AddClassHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: Definitely refactor this part of the code
	// TODO: Fix The duplicate submission issue
	if r.Method == "GET" {
		files := []string{
			filepath.Join(templatesDir, "add-class.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		tmpl.Execute(w, nil)

	} else {
		r.ParseForm()
		// logic part of log in
		name := r.FormValue("name")
		classCoordinator := r.FormValue("classCoordinator")

		fmt.Println("Firstname:", name)
		fmt.Println("address:", classCoordinator)

		class := &models.Class{}
		class.Name = name

		class.ClassCoordinator = classCoordinator
		class.NoOfStudent = 0

		class.Create()
		// fmt.Println("Student struct:", student)
		files := []string{
			filepath.Join(templatesDir, "add-class.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		tmpl.Execute(w, nil)

	}
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "add-student.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}
