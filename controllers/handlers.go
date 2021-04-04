package controllers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"sms-webapp/models"
	"sms-webapp/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const (
	templatesDir = "templates"
)

// TODO: Not necessarily a todo just the guide for the implementation
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
		MaxAge:   60 * 60, // TODO: Update the time here to something reasonable like a day
		HttpOnly: true,
	}

	gob.Register(models.User{})

	// tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func getUser(s *sessions.Session) models.User {
	val := s.Values["user"]
	var user = models.User{}
	user, ok := val.(models.User)
	if !ok {
		// return models.User{Authenticated: false}
		return models.User{}
	}
	return user
}

// func getUser2(s *sessions.Session) User {
// 	val := s.Values["user"]
// 	var user = User{}
// 	user, ok := val.(User)
// 	if !ok {
// 		return User{Authenticated: false}
// 	}
// 	return user
// }

func StudentLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/student-login.html"))
		tmpl.Execute(w, nil)
	} else {
		session, err := store.Get(r, "cookie-name")
		r.ParseMultipartForm(32 << 20)
		email := r.FormValue("surname")
		id := r.FormValue("password")
		intId, err := strconv.Atoi(id)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("name is and id is ", email, id)

		// TODO: Update the login here to return quickly for failed login
		if models.StudentLogin(email, uint(intId)) {
			currentStudent := models.GetSingleStudentById(uint(intId))
			fmt.Println("Login SuccessFul")
			session.Values["user"] = currentStudent.ID
			session.Values["authenticated"] = true
			session.Values["user"] = currentStudent.Name

			fmt.Println("The user has name of ", currentStudent.Name)

			err = session.Save(r, w)
			if err != nil {
				session.AddFlash("Error when trying to log in")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/student-dashboard", http.StatusFound)

		} else {
			fmt.Println("Login Failed")
			session.AddFlash("Login Failed")

		}

	}

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
			currentUser := models.GetUserByEmail(email)
			fmt.Println("Login SuccessFul")
			session.Values["user"] = email
			session.Values["authenticated"] = true
			session.Values["user"] = currentUser

			fmt.Println("The user has name of ", currentUser.Name)

			err = session.Save(r, w)
			if err != nil {
				session.AddFlash("Error when trying to log in")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/your-subject", http.StatusFound)

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

//AuthRequired is a middleware which will be used for each httpHandler to check if there is any active session
func AuthRequired(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler(w, r)
	}
}

//PrincipalRoleRequired is a middleware which will be used for each httpHandler to check if there is any active session
func PrincipalRoleRequired(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", 302)
			return
		}
		currentUser := getUser(session)

		if currentUser.Role != "principal" {
			http.Redirect(w, r, "/your-subject", 302)
			return
		}
		handler(w, r)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = models.User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)

}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// TODO: Uncomment the following below before deployment
	// session, _ := store.Get(r, "cookie-name")
	// if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	// http.Error(w, "Forbidden", http.StatusForbidden)
	// 	return
	// }

	// currentUser := getUser(session)

	// if currentUser.Role != "principal" {
	// 	// TODO: This should ideally move to the forbidden page
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	// http.Error(w, "Forbidden", http.StatusForbidden)
	// 	return
	// }

	type DashboardVariables struct {
		AmountOwed              float64
		NumberOfStudentInSchool int64
		AmountRealizedToday     float64
		AmountReaizedPastWeek   float64
		AmountRealizedThisTerm  float64
	}

	amountOWED := models.GetTotalAmountOwed()
	fmt.Println("Amount owed is ", amountOWED)

	noOfStudentsInSchool := models.GetNoOfStudentInSchool()
	fmt.Println("The number of student in school is ", noOfStudentsInSchool)

	amountRealizedInCurrentTerm := models.GetAmountPaidInCurrentTerm()

	dashboardPagevariable := DashboardVariables{
		AmountOwed:              amountOWED,
		NumberOfStudentInSchool: noOfStudentsInSchool,
		AmountRealizedThisTerm:  amountRealizedInCurrentTerm,
	}

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

	tmpl.Execute(w, &dashboardPagevariable)
}

func ViewAllTeacherHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	currentUser := getUser(session)

	// currentRole := currentUser.Role
	type TeacherPageVariable struct {
		CurrentTeacher models.User
		Teachers       []models.User
	}

	// FIXME: Currently some users are principal and teachers. So, returning all roles
	// teachers := models.GetAllUserByRole("teacher")
	teachers := models.GetAllUsers()

	fmt.Println("The current user is ", currentUser)
	// pVariables :=

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

func ViewAllFeesHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		data := models.GetAllStudents()
		files := []string{
			filepath.Join(templatesDir, "all-fees-details.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		tmpl.Execute(w, &data)
	}

}

func ViewYourSubjectHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	currentUser := getUser(session)

	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		type TeacherSubjectVariable struct {
			SubjectClass *models.SubjectClass
			Subject      *models.Subject
		}

		var SubjectsDetails []TeacherSubjectVariable
		type TeacherSubjectsPageVariable struct {
			Teacher             models.User
			SubjectClassDetails []TeacherSubjectVariable
			AllowGrading        bool
		}
		teacherSubjectClass := models.GetSubjectClassForTeacher(currentUser.ID)

		for _, singleSubjectClass := range teacherSubjectClass {

			currentSubject := models.GetSubjectById(singleSubjectClass.SubjectID)
			SubjectsDetails = append(SubjectsDetails, TeacherSubjectVariable{
				SubjectClass: singleSubjectClass,
				Subject:      currentSubject,
			})
		}
		finalPVariables := TeacherSubjectsPageVariable{
			Teacher:             currentUser,
			SubjectClassDetails: SubjectsDetails,
			AllowGrading:        true,
		}

		files := []string{
			filepath.Join(templatesDir, "all-teachers-subjects.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		// tmpl := template.Must(template.
		// 	ParseFiles(files...))

		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			panic(err.Error())
		}

		tmpl.Execute(w, &finalPVariables)
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

func add(x, y int) int {
	return x + y
}

func GradeStudentsHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: Fix the hack for loading static files for path variables
	session, _ := store.Get(r, "cookie-name")
	currentUser := getUser(session)

	requestParams := mux.Vars(r)
	id, err := strconv.Atoi(requestParams["id"])

	if err != nil {
		panic(err.Error())
	}

	currentSubjectClass := models.GetSubjectClassById(uint(id))

	if currentSubjectClass.UserID != currentUser.ID {
		// TODO: Think through how unauthorized users should be handled
		http.Redirect(w, r, "/login", http.StatusFound)
		// http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	studentsForSubject := models.GetStudentSubjectsClassBySubjectClassID(currentSubjectClass.ID)

	type StudentSubjectDetails struct {
		Student        *models.Student
		StudentSubject *models.StudentSubjectClass
	}

	var PageVariables []StudentSubjectDetails

	for _, singleStudentSubject := range studentsForSubject {
		currentStudent := models.GetSingleStudentById(singleStudentSubject.StudentID)
		PageVariables = append(PageVariables, StudentSubjectDetails{
			Student:        currentStudent,
			StudentSubject: singleStudentSubject,
		})

	}

	files := []string{
		filepath.Join(templatesDir, "result-page.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	// funcs := template.FuncMap{"add": add}

	tmpl := template.Must(template.
		ParseFiles(files...))

	// tmpl := template.Must(template.New("result-page.html").Funcs(funcs).
	// 	ParseFiles(files...))

	tmpl.Execute(w, PageVariables)
}

func ViewAllGradeHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "student-grade-old.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func ViewResultHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join(templatesDir, "result-page-old.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl := template.Must(template.
		ParseFiles(files...))

	tmpl.Execute(w, nil)
}

func UpdateGradeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The request method is ", r.Method)
	r.ParseForm()

	fieldName := r.FormValue("name")
	scoreValue := r.FormValue("value")
	studentSubjectClassID := r.FormValue("pk")
	// stringSlice := strings.Split(r.FormValue("pk"), "-")

	// studentID := stringSlice[1]
	// subjectClassID := stringSlice[0]

	fmt.Println(fieldName)
	fmt.Println(studentSubjectClassID)
	// fmt.Println(studentID)
	// fmt.Println(subjectClassID)
	fmt.Println(scoreValue)

	intStudentSubjectClassID, _ := strconv.ParseUint(studentSubjectClassID, 10, 32)
	intScoreValue, _ := strconv.ParseUint(scoreValue, 10, 32)

	studentSubjectClass := models.UpdateStudentscore(uint(intStudentSubjectClassID), fieldName, float32(intScoreValue))

	fmt.Println(studentSubjectClass)

	files := []string{
		filepath.Join(templatesDir, "result-page-old.html"),
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

// UpdateSubjectHandler - Update the optional subject for student
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

			fmt.Println("The current subject is ", currentSubject)
			fmt.Println("The single subject id is ", singleSubjectClass.SubjectID)

			// We are Only considering subjects that are not compulsory. Compulsory subjects are automatically added
			if currentSubject.IsCompulsory == "NO" {
				PageVariables = append(PageVariables, SubjectClassDetails{
					SubjectClass: singleSubjectClass,
					Subject:      currentSubject,
				})
			}

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

		fmt.Println("The subject Class Details is ", finalPVariables.SubjectClassDetails)
		// fmt.Println("The subject Class Details is ", finalPVariables.SubjectClassDetails[0].Subject.Name)

		// tmpl := template.Must(template.
		// 	ParseFiles(files...))

		// tmpl.Execute(w, nil)

		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			panic(err.Error())
		}

		tmpl.Execute(w, finalPVariables)
	} else {

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Got to the else part of assigning subjects")
		r.ParseMultipartForm(32 << 20)
		studentID := id

		class := r.Form["subjectClass"]
		fmt.Println("The class is", class)

		for _, singleClass := range class {
			intSingleClass, _ := strconv.ParseUint(singleClass, 10, 32)
			studentSubjectClass := &models.StudentSubjectClass{}

			studentSubjectClass.IsActive = true
			studentSubjectClass.StudentID = uint(studentID)
			studentSubjectClass.SubjectClassID = uint(intSingleClass)
			studentSubjectClass.Create()

			// intSingleClass, _ := strconv.ParseUint(singleClass, 10, 32)

			// respData := models.UpdateSubjectClassTeacher(uint(intSingleClass), uint(teacher))

			// fmt.Println(respData)
		}
		redirectURL := fmt.Sprintf("/student-profile/%d", studentID)
		http.Redirect(w, r, redirectURL, http.StatusFound)

	}

}

func PrintReceiptForStudnet(w http.ResponseWriter, r *http.Request) {

	type ReceiptStruct struct {
		Payment         *models.FeesPayment
		Student         *models.Student
		AmountInString  string
		CurrentDate     string
		BeingPaymentFor string
	}

	requestParams := mux.Vars(r)
	id, err := strconv.Atoi(requestParams["id"])

	if err != nil {
		panic(err.Error())
	}

	feePaid := models.GetFeePaidById(uint(id))

	amountInString := utils.Num2words(int(feePaid.Amount))
	fmt.Println(amountInString)

	student := models.GetSingleStudentById(feePaid.StudentID)

	beingPaymentString := ""
	if student.PresentTermBalance > 0 {
		// FIXME: This should ideally pick the name of the current term.
		// FIXME: This would be updated every term
		// FIXME: Also consider the string to return for outstanding debt

		beingPaymentString = "Second term part payment"
	} else {
		beingPaymentString = "Second term full payment"
	}

	currentTime := time.Now()
	receiptPageVariable := ReceiptStruct{
		Payment:         feePaid,
		Student:         student,
		AmountInString:  amountInString,
		CurrentDate:     currentTime.Format("02/01/2006"),
		BeingPaymentFor: beingPaymentString,
	}

	// fmt.Println("Current Time in String: ", currentTime.String())
	// fmt.Println("MM-DD-YYYY : ", currentTime.Format("01-02-2006"))
	// fmt.Println("DD|MM|YYYY : ", currentTime.Format("02|01|2006"))

	files := []string{
		filepath.Join(templatesDir, "student-receipt.html"),
		filepath.Join(templatesDir, "base.html"),
	}

	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		panic(err.Error())
	}

	tmpl.Execute(w, &receiptPageVariable)

}

func UpdateOutstandingDebt(w http.ResponseWriter, r *http.Request) {

	// GetAllSubjectsDetails
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}

		student := models.GetSingleStudentById(uint(id))

		files := []string{
			filepath.Join(templatesDir, "student-outstanding-debt.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		tmpl.Execute(w, student)
	} else {

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}

		student := models.GetSingleStudentById(uint(id))
		oustandingDebtStr := r.FormValue("amount")
		outstandingDebt, _ := strconv.ParseFloat(oustandingDebtStr, 64)
		student.OutstandingDebt = outstandingDebt

		models.GetDB().Save(&student)

		http.Redirect(w, r, "/fees", http.StatusFound)
	}

}
func GetRemarkFromScore(score float64) string {
	// TODO: Update this condition
	if score >= 70 {
		return "EXCELLENT"
	} else if score >= 60 && score < 70 {
		return "VERY GOOD"
	} else if score >= 50 && score < 60 {
		return "GOOD"
	} else if score >= 40 && score < 50 {
		return "PASS"
	} else if score >= 30 && score < 40 {
		return "POOR"
	} else {
		return "WEAK"
	}
}

func GetGradeFromScore(score float64) string {
	if score >= 70 {
		return "A"
	} else if score >= 60 && score < 70 {
		return "B"
	} else if score >= 50 && score < 60 {
		return "C"
	} else if score >= 45 && score < 50 {
		return "D"
	} else if score >= 40 && score < 45 {
		return "E"
	} else {
		return "F"
	}
}

func GetTeacherRemarkFromPercentage(percentage float64) string {
	if percentage >= 70 {
		return "A splendid result. Increase your academic tempo. The sky is the beginning."
	} else if percentage >= 60 && percentage < 70 {
		return "A good result. Do not relent in your efforts."
	} else if percentage >= 50 && percentage < 60 {
		return "An average result. Work hard and don't be left behind."
	} else if percentage >= 40 && percentage < 50 {
		return "Below average is not good for you. Work harder next term."
	} else {
		return "Poor result. Work hard or you will be left behind."
	}

}

func GetPrincipalRemarkFromPercentage(percentage float64) string {
	if percentage >= 70 {
		return "An excellent performance. keep it up."
	} else if percentage >= 60 && percentage < 70 {
		return "A good performance, however, there is still room for improvement next term."
	} else if percentage >= 50 && percentage < 60 {
		return "An average performance. Concentrate more on your weak subjects."
	} else if percentage >= 40 && percentage < 50 {
		return "A below average performance. Put in more effort in your academics."
	} else {
		return "A very poor performance. You need to put in extra effort next term."
	}

}

func DeleteSingleStudentSubjectHandler(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	subject_id, err := strconv.Atoi(requestParams["subject_id"])
	student_id, err := strconv.Atoi(requestParams["student_id"])
	if err != nil {
		panic(err.Error())
	}

	// TODO: What happens if the class is updated
	student := models.GetSingleStudentById(uint(student_id))
	subjectClass := models.GetSubjectClassBySubjectIdAndClass(uint(subject_id), student.ClassText)
	// studentSubjectClass := models.GetStudentSubjectClassById(uint(id))
	// studentID := studentSubjectClass.StudentID

	models.GetDB().Where("subject_class_id =? AND student_id=? ", subjectClass.ID, student.ID).Delete(&models.StudentSubjectClass{})
	// err := GetDB().Table("student_subject_classes").Where("id = ?", id).First(studentSubjectClass).Error
	// if err != nil || err == gorm.ErrRecordNotFound {
	// 	return nil
	// }
	redirectURL := fmt.Sprintf("/student-profile/%d", student_id)
	fmt.Println(redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)

}

func ViewSingleStudentResultHandler(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	id, err := strconv.Atoi(requestParams["id"])

	if err != nil {
		panic(err.Error())
	}

	type StudentResultPageData struct {
		StudentSubjectClass *models.StudentSubjectClass
		Subject             *models.Subject
	}

	var StudentSubjectScoreDetails []StudentResultPageData

	type StudentDetailsResultData struct {
		Student                     *models.Student
		SubjectScoreDetails         []StudentResultPageData
		TotalObtainable             int
		StudentPercentage           float64
		StudentTotalScore           float32
		TeacherRemarks              string
		PrincipalRemarks            string
		OverallRemarkFromPercentage string
		OverallGradeFromPercentage  string
	}

	GetSubjectListRow := models.GetStudentSubjectListRowByStudentID(uint(id))
	fmt.Println(GetSubjectListRow)

	studentSubjectClass := models.GetStudentSubjectsClassByStudentID(uint(id))
	fmt.Println(studentSubjectClass)
	// teachers := models.GetAllUserByRole("teacher")

	// FIXME: This was commented to update for second term
	// ************************************************************************************
	// TODO:  VERY IMPORTANT
	/*
			TODO: Ensure that the second term first term result  is gotten from the backup bcos
		// Initially models.UpdateStudentSubject(singleStudentSubjectDetail.ID, totalScore, grade, remark) was overwrting first term results during second term.
		// The backup is the source of truth for first term
		****************************************************
	*/

	for _, singleStudentSubjectDetail := range studentSubjectClass {
		totalScore := singleStudentSubjectDetail.FirstCA + singleStudentSubjectDetail.SecondCA + singleStudentSubjectDetail.FirstExam
		grade := GetGradeFromScore(float64(totalScore))
		remark := GetRemarkFromScore(float64(totalScore))
		updatedStudentSubjectClass := models.UpdateStudentSubject(singleStudentSubjectDetail.ID, totalScore, grade, remark)
		fmt.Println(updatedStudentSubjectClass)
	}

	// for _, singleStudentSubjectDetail := range studentSubjectClass {
	// 	totalScore := singleStudentSubjectDetail.SFirstCA + singleStudentSubjectDetail.SSecondCA + singleStudentSubjectDetail.SecondExam
	// 	grade := GetGradeFromScore(float64(totalScore))
	// 	remark := GetRemarkFromScore(float64(totalScore))
	// 	updatedStudentSubjectClass := models.UpdateStudentSubject(singleStudentSubjectDetail.ID, totalScore, grade, remark)
	// 	fmt.Println(updatedStudentSubjectClass)
	// }

	// TODO: So many expensive db calls here
	studentSubjectClass = models.GetStudentSubjectsClassByStudentID(uint(id))
	studentDetails := models.GetSingleStudentById(uint(id))

	println("I got the student details also ", studentDetails)

	var studentTotal float32

	for _, singleStudentSubjectDetail := range studentSubjectClass {
		studentTotal += singleStudentSubjectDetail.TotalFirst
		// studentTotal += singleStudentSubjectDetail.TotalSecond
		currentSubject := models.GetSubjectBySubjectClassId(singleStudentSubjectDetail.SubjectClassID)
		fmt.Println("The name of the subject is", currentSubject.Name)
		StudentSubjectScoreDetails = append(StudentSubjectScoreDetails, StudentResultPageData{
			StudentSubjectClass: singleStudentSubjectDetail,
			Subject:             currentSubject,
		})
	}

	numberOfSubjectOffered := len(studentSubjectClass)
	var studentPercentage float64
	studentPercentage = float64(studentTotal) / float64(len(studentSubjectClass))

	studentPercentage = math.Round(studentPercentage*10) / 10
	totalScoreObtainable := numberOfSubjectOffered * 100

	principalRemarks := GetPrincipalRemarkFromPercentage(studentPercentage)
	teacherRemarks := GetTeacherRemarkFromPercentage(studentPercentage)

	fmt.Println("The percentage is ", studentPercentage)
	fmt.Println("The total score obtainable is ", totalScoreObtainable)
	fmt.Println("The number of subject offered is ", numberOfSubjectOffered)

	pVariables := StudentDetailsResultData{
		Student:                     studentDetails,
		SubjectScoreDetails:         StudentSubjectScoreDetails,
		TotalObtainable:             totalScoreObtainable,
		StudentPercentage:           studentPercentage,
		StudentTotalScore:           studentTotal,
		PrincipalRemarks:            principalRemarks,
		TeacherRemarks:              teacherRemarks,
		OverallRemarkFromPercentage: GetRemarkFromScore(studentPercentage),
		OverallGradeFromPercentage:  GetGradeFromScore(studentPercentage),
	}

	// files := []string{
	// 	filepath.Join(templatesDir, "student-grade.html"),
	// 	filepath.Join(templatesDir, "base.html"),
	// }

	// tmpl, err := template.ParseFiles(files...)
	// if err != nil {
	// 	panic(err.Error())
	// }

	tmpl := template.Must(template.ParseFiles("templates/student-grade2.html"))
	// tmpl := template.Must(template.ParseFiles("templates/student-grade_uthman.html"))
	// tmpl := template.Must(template.ParseFiles("templates/student-grade_misturah.html"))
	// tmpl := template.Must(template.ParseFiles("templates/student-grade_ama.html"))
	// tmpl := template.Must(template.ParseFiles("templates/student-grade_bashirudeen.html"))

	// tmpl := template.Must(template.
	// 	ParseFiles(files...))

	tmpl.Execute(w, &pVariables)

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
			Subject []*models.Subject
			Payment []*models.FeesPayment
			// SubjectClass []*models.SubjectClass
		}

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		student := models.GetSingleStudentById(uint(id))

		fees_payment := models.GetPaymentsForStudent(uint(id))

		fmt.Println("The class is ", student.ClassText)

		classSubject := models.GetSubjectsClassForStudentByID(uint(id))

		StudentSubjects := make([]*models.Subject, 0)
		for _, singleClassSubject := range classSubject {
			singleSubject := models.GetSubjectBySubjectClassId(singleClassSubject.SubjectClassID)

			StudentSubjects = append(StudentSubjects, singleSubject)

		}

		fmt.Println(StudentSubjects)

		// var StudentSubjects []models.Subject

		// for _, singleClassSubject := range classSubject {
		// 	StudentSubjects = append(StudentSubjects, )
		// }

		// fmt.Println("The type of classsubject is ", reflect.TypeOf(classSubject))
		fmt.Println("The type of data is ", reflect.TypeOf(student))

		// fmt.Println(classSubject)

		pVariables := PageVariables{Student: student, Subject: StudentSubjects, Payment: fees_payment}
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

func TeacherSubjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		type TeacherSubjectVariable struct {
			SubjectClass *models.SubjectClass
			Subject      *models.Subject
		}

		var SubjectsDetails []TeacherSubjectVariable
		type TeacherSubjectsPageVariable struct {
			Teacher             models.User
			SubjectClassDetails []TeacherSubjectVariable
			AllowGrading        bool
		}
		teacherSubjectClass := models.GetSubjectClassForTeacher(uint(id))
		currentUser := *models.GetUserByID(uint(id))

		for _, singleSubjectClass := range teacherSubjectClass {

			currentSubject := models.GetSubjectById(singleSubjectClass.SubjectID)
			SubjectsDetails = append(SubjectsDetails, TeacherSubjectVariable{
				SubjectClass: singleSubjectClass,
				Subject:      currentSubject,
			})
		}
		finalPVariables := TeacherSubjectsPageVariable{
			Teacher:             currentUser,
			SubjectClassDetails: SubjectsDetails,
			AllowGrading:        false,
		}

		files := []string{
			filepath.Join(templatesDir, "all-teachers-subjects.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		// tmpl := template.Must(template.
		// 	ParseFiles(files...))

		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			panic(err.Error())
		}

		tmpl.Execute(w, &finalPVariables)
	}

}

func StudentPaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		student := models.GetSingleStudentById(uint(id))

		files := []string{
			filepath.Join(templatesDir, "student-payment.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		tmpl.Execute(w, student)

	} else {
		session, _ := store.Get(r, "cookie-name")
		// if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// 	http.Redirect(w, r, "/login", http.StatusFound)
		// 	// http.Error(w, "Forbidden", http.StatusForbidden)
		// 	return
		// }

		currentUser := getUser(session)

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		// UpdateStudentSubject
		student := models.GetSingleStudentById(uint(id))
		amountPaidStr := r.FormValue("amount")
		amountPaid, _ := strconv.ParseFloat(amountPaidStr, 64)
		outstandingDebt := student.OutstandingDebt
		IsFeeCompleted := false

		if amountPaid-outstandingDebt < 0 {
			student.OutstandingDebt = outstandingDebt - amountPaid
		} else if amountPaid-outstandingDebt > 0 {
			student.OutstandingDebt = 0
			currenTermPayment := student.PresentTermPayment + amountPaid - outstandingDebt
			student.PresentTermPayment = currenTermPayment
			student.PresentTermBalance = student.PresentTermFees - currenTermPayment
			// student.PresentTermBalance =
			if currenTermPayment >= student.PresentTermFees {
				IsFeeCompleted = true
				student.IsFeeCompleted = IsFeeCompleted
			}

		} else {
			// This definitely means payment made was just to clear the outstanding
			student.OutstandingDebt = 0
		}
		// studentModel := &models.Student{}
		// models.GetDB().Model(&studentModel).Updates(student)
		models.GetDB().Save(&student)

		feesPayment := &models.FeesPayment{}
		feesPayment.StudentID = student.ID
		feesPayment.Amount = amountPaid
		feesPayment.UserID = currentUser.ID
		feesPayment.IsComplete = IsFeeCompleted
		feesPayment.Create()
		http.Redirect(w, r, "/fees", http.StatusFound)
	}

}

func StudentDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		student := models.GetSingleStudentById(uint(id))

		files := []string{
			filepath.Join(templatesDir, "student-payment.html"),
			filepath.Join(templatesDir, "base.html"),
		}

		tmpl := template.Must(template.
			ParseFiles(files...))

		tmpl.Execute(w, student)

	} else {
		session, _ := store.Get(r, "cookie-name")
		// if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// 	http.Redirect(w, r, "/login", http.StatusFound)
		// 	// http.Error(w, "Forbidden", http.StatusForbidden)
		// 	return
		// }

		currentUser := getUser(session)

		requestParams := mux.Vars(r)
		id, err := strconv.Atoi(requestParams["id"])

		if err != nil {
			panic(err.Error())
		}
		// UpdateStudentSubject
		student := models.GetSingleStudentById(uint(id))
		amountPaidStr := r.FormValue("amount")
		amountPaid, _ := strconv.ParseFloat(amountPaidStr, 64)
		outstandingDebt := student.OutstandingDebt
		IsFeeCompleted := false

		if amountPaid-outstandingDebt < 0 {
			student.OutstandingDebt = outstandingDebt - amountPaid
		} else if amountPaid-outstandingDebt > 0 {
			student.OutstandingDebt = 0
			currenTermPayment := student.PresentTermPayment + amountPaid - outstandingDebt
			student.PresentTermPayment = currenTermPayment
			student.PresentTermBalance = student.PresentTermFees - currenTermPayment
			// student.PresentTermBalance =
			if currenTermPayment >= student.PresentTermFees {
				IsFeeCompleted = true
				student.IsFeeCompleted = IsFeeCompleted
			}

		} else {
			// This definitely means payment made was just to clear the outstanding
			student.OutstandingDebt = 0
		}
		// studentModel := &models.Student{}
		// models.GetDB().Model(&studentModel).Updates(student)
		models.GetDB().Save(&student)

		feesPayment := &models.FeesPayment{}
		feesPayment.StudentID = student.ID
		feesPayment.Amount = amountPaid
		feesPayment.UserID = currentUser.ID
		feesPayment.IsComplete = IsFeeCompleted
		feesPayment.Create()
		http.Redirect(w, r, "/fees", http.StatusFound)
	}

}
