package controllers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"sms-webapp/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const (
	templatesDir = "templates"
)

// store will hold all session data
var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(models.User{})

	// tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
	} else {
		session, err := store.Get(r, "cookie-name")
		r.ParseMultipartForm(32 << 20)
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Println("Email is and password is ", email, password)
		// TODO: Update the login here to return quickly for failed login
		if models.Login(email, password) {
			fmt.Println("Login SuccessFul")
			session.Values["user"] = email
			session.Values["authenticated"] = true

			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)

		} else {
			fmt.Println("Login Failed")
			session.AddFlash("Login Failed")

		}

	}

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
	} else {
		r.ParseMultipartForm(32 << 20)
		name := r.FormValue("fullName")
		phoneNo := r.FormValue("phoneNo")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user := &models.User{}
		user.Name = name
		user.PhoneNo = phoneNo
		user.Email = email
		user.Password = password
		user.Role = "teacher"

		TeacherID := user.Create()

		fmt.Println("The teacher Id is ", TeacherID)

		http.Redirect(w, r, "/teachers", http.StatusTemporaryRedirect)

	}

}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// TODO: Uncomment the following below before deployment
	// session, _ := store.Get(r, "cookie-name")
	// if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
	// 	http.Error(w, "Forbidden", http.StatusForbidden)
	// 	return
	// }

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
	teachers := models.GetAllUserByRole("teacher")

	files := []string{
		filepath.Join(templatesDir, "all-teachers.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		panic(err.Error())
	}

	// tmpl := template.Must(template.
	// 	ParseFiles(files...))

	tmpl.Execute(w, &teachers)

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

		classIdMap := make(map[string]uint)
		classIdMap["jss1"] = 1
		classIdMap["jss2"] = 2
		classIdMap["jss3"] = 3
		classIdMap["sss1"] = 4
		classIdMap["sss2"] = 5
		classIdMap["sss3"] = 6

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
		student.ClassID = classIdMap[class]
		student.ClassText = class
		// student.Class = class
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
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		data := models.GetAllSubjects()
		fmt.Print(data[0].Name)

		files := []string{
			filepath.Join(templatesDir, "all-subjects.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		// tmpl := template.Must(template.
		// 	ParseFiles(files...))

		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			panic(err.Error())
		}

		tmpl.Execute(w, &data)
	} else {
		fmt.Println("Got to the else part of viewing subjects")
		r.ParseMultipartForm(32 << 20)
		name := r.FormValue("Subject")
		// class := r.FormValue("class")
		class := r.Form["class"]
		isCompulsory := r.FormValue("isCompulsory")
		fmt.Println("The name is ", name)
		fmt.Println("The class is ", class)
		fmt.Println("The isCompulsory is ", isCompulsory)

		// TODO: Include the image for the subject also
		subject := &models.Subject{}
		subject.Name = name
		subject.IsCompulsory = isCompulsory
		subjectID := subject.Create()

		for _, singleClass := range class {
			subjectClass := &models.SubjectClass{}
			subjectClass.Class = singleClass
			subjectClass.SubjectID = subjectID
			subjectClass.IsCompulsory = false
			subjectClass.Create()
		}

		http.Redirect(w, r, "/subjects", http.StatusFound)

	}
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

func ViewAllGradeHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "student-grade.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

// AssignSubjectHandler - Assign subject to teacher handler
func AssignSubjectHandler(w http.ResponseWriter, r *http.Request) {
	// GetAllSubjectsDetails
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {

		type AssignSubjectPageVariable struct {
			Teacher     *models.User
			AllSubjects []*models.Subject
		}

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		currentTeacher := models.GetUserByID(uint(id))

		// data := models.GetAllSubjects()
		data := models.GetAllSubjectsDetails()

		fmt.Println(data)

		files := []string{
			filepath.Join(templatesDir, "assign-subjects.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		// tmpl := template.Must(template.
		// 	ParseFiles(files...))

		// tmpl.Execute(w, nil)

		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			panic(err.Error())
		}
		pVariables := AssignSubjectPageVariable{Teacher: currentTeacher, AllSubjects: data}

		tmpl.Execute(w, &pVariables)
	} else {
		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Got to the else part of assigning subjects")
		r.ParseMultipartForm(32 << 20)
		// teacher := r.FormValue("")
		teacher := id

		fmt.Println("The teacher id is ", teacher)

		class := r.Form["subjectClass"]
		fmt.Println("The class is", class)

		for _, singleClass := range class {

			intSingleClass, _ := strconv.ParseUint(singleClass, 10, 32)

			respData := models.UpdateSubjectClassTeacher(uint(intSingleClass), uint(teacher))
			fmt.Println(respData)
		}

		http.Redirect(w, r, "/teachers", http.StatusTemporaryRedirect)
	}

}

// UpdateSubjectHandler - Update the optional subject for the current user
func UpdateSubjectHandler(w http.ResponseWriter, r *http.Request) {
	// GetAllSubjectsDetails
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {

		type SubjectClassDetails struct {
			SubjectClass *models.SubjectClass
			Subject      *models.Subject
		}

		var PageVariables []SubjectClassDetails
		type StudentSubjectClassPageVariable struct {
			Student             *models.Student
			SubjectClassDetails []SubjectClassDetails
		}

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		student := models.GetSingleStudentById(uint(id))

		// data := models.GetAllSubjects()
		// data := models.GetSubjectsDetailsForClass(student.ClassText)
		studentClass := student.ClassText
		data := models.GetSubjectsForClass(studentClass)

		for _, singleSubjectClass := range data {

			// TODO: Check if the subject is compulsory and drop from the list of struct with the corresponding subjectclass
			fmt.Println(singleSubjectClass)

			currentSubject := models.GetSubjectById(singleSubjectClass.SubjectID)

			PageVariables = append(PageVariables, SubjectClassDetails{
				SubjectClass: singleSubjectClass,
				Subject:      currentSubject,
			})

		}

		finalPVariables := StudentSubjectClassPageVariable{
			Student:             student,
			SubjectClassDetails: PageVariables,
		}

		fmt.Println(data)

		files := []string{
			filepath.Join(templatesDir, "edit-student-subjects.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		// tmpl := template.Must(template.
		// 	ParseFiles(files...))

		// tmpl.Execute(w, nil)

		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			panic(err.Error())
		}

		tmpl.Execute(w, finalPVariables)
	} else {
		fmt.Println("Got to the else part of assigning subjects")
		r.ParseMultipartForm(32 << 20)
		// teacher := r.FormValue("")
		teacher := 1

		class := r.Form["subjectClass"]
		fmt.Println("The class is", class)

		for _, singleClass := range class {

			intSingleClass, _ := strconv.ParseUint(singleClass, 10, 32)

			respData := models.UpdateSubjectClassTeacher(uint(intSingleClass), uint(teacher))
			fmt.Println(respData)
		}

	}

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

		type PageVariables struct {
			Student *models.Student
			// SubjectClass []*models.SubjectClass
		}
		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		student := models.GetSingleStudentById(uint(id))

		fmt.Println("The class is ", student.ClassText)

		// classSubject := models.GetSubjectsForClass(student.ClassText)

		// fmt.Println("The type of classsubject is ", reflect.TypeOf(classSubject))
		fmt.Println("The type of data is ", reflect.TypeOf(student))

		// fmt.Println(classSubject)

		pVariables := PageVariables{Student: student}
		// pVariables := PageVariables{Student: student, SubjectClass: classSubject}

		fmt.Println("The page variables are ", pVariables)

		files := []string{
			filepath.Join(templatesDir, "student-profile.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		// fmt.Println(data)
		tmpl.Execute(w, pVariables)

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
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		files := []string{
			filepath.Join(templatesDir, "add-subject.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		tmpl.Execute(w, nil)
	}

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
