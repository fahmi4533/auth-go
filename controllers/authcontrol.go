package controllers

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/makcik45/jwt-go/config"
	"github.com/makcik45/jwt-go/entities"
	"github.com/makcik45/jwt-go/libraries"
	"github.com/makcik45/jwt-go/models"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	UserName string
	Password string
}

var usermodel = models.NuserModel()
var validation = libraries.NewValid()

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SEESSION_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	} else {
		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {

			data := map[string]interface{}{
				"nama_lengkap": session.Values["nama_lengkap"],
			}

			temp, err := template.ParseFiles("views/index.html")
			if err != nil {
				panic(err)
			}
			temp.Execute(w, data)

		}

	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/login.html")

		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		userInput := &UserInput{
			UserName: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		}

		vError := validation.Struct(userInput)

		if vError != nil {
			data := map[string]interface{}{
				"validation": vError,
			}
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else {
			var user entities.User
			usermodel.Where(&user, "username", userInput.UserName)

			var message error
			if user.Username == "" {
				// tidak ditemukan di database
				message = errors.New("Username atau Password salah!")
			} else {
				// pengecekan passwprd
				errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
				if errPass != nil {
					message = errors.New("Username atau Password salah woy!")

				}
			}
			if message != nil {

				data := map[string]interface{}{
					"error": message,
				}

				temp, _ := template.ParseFiles("views/login.html")
				temp.Execute(w, data)
			} else {
				//set session
				session, _ := config.Store.Get(r, config.SEESSION_ID)

				session.Values["loggedIn"] = true
				session.Values["email"] = user.Email
				session.Values["username"] = user.Username
				session.Values["nama_lengkap"] = user.NamaLengkap

				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SEESSION_ID)

	// delete session

	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		temp, _ := template.ParseFiles("views/register.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		//  melakukan proses registrasi

		// mengambil inputan form
		r.ParseForm()

		user := entities.User{
			NamaLengkap: r.Form.Get("nama_lengkap"),
			Email:       r.Form.Get("email"),
			Username:    r.Form.Get("username"),
			Password:    r.Form.Get("password"),
			Cpassword:   r.Form.Get("cpassword"),
		}
		vError := validation.Struct(user)

		if vError != nil {
			data := map[string]interface{}{
				"validation": vError,
				"user":       user,
			}
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)
		} else {
			// ketika sudah validasi kita memasukan ke database dengan haspassword
			hashPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			user.Password = string(hashPass)

			// kita lakukan insert ke database
			usermodel.Create(user)

			data := map[string]interface{}{
				"pesan": "Registrasi Berhasil!!",
			}
			temp, _ := template.ParseFiles("views/register.html")
			temp.Execute(w, data)

		}

	}
}
