package services

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
	supa "github.com/nedpals/supabase-go"

	bcrypt "golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Uid           string
	Username      string
	Password      string
	EmailAddress  string
	PhoneNumber   string
	Created_at    string
	LastLoginTime string
}

type LoginResult struct {
	Uid          string `json:"uid"`
	Username     string `json:"username"`
	EmailAddress string `json:"emailAddress"`
	LastLoggedIn string `json:"lastLoggedIn"`
}

type LoginState struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    LoginResult `json:"data"`
}

func supabase() *supa.Client {

	supabaseUrl := "https://fmipfsmmrgqidckaxbqk.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImZtaXBmc21tcmdxaWRja2F4YnFrIiwicm9sZSI6ImFub24iLCJpYXQiOjE2NzI0NTc5NjAsImV4cCI6MTk4ODAzMzk2MH0.b7saxxttkv7cRxDVqRIJ-tasRauZkNOQkGBdxX5seb0"

	supabase := supa.CreateClient(supabaseUrl, supabaseKey)

	return supabase

}

func RegisterUser(name string, username string, emailAddress string, password string, phoneNumber string) {

	id, idErr := gonanoid.New()

	hashedBytes, passErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if passErr != nil {
		panic(passErr)
	}

	if idErr != nil {
		panic(idErr)
	}

	payload := map[string]string{
		"uid":          id,
		"username":     username,
		"emailAddress": emailAddress,
		"password":     string(hashedBytes),
		"phoneNumber":  phoneNumber,
	}

	var results map[string]interface{}

	dbErr := supabase().DB.From("UserAuth").Insert(payload).Execute(&results)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
	}

	fmt.Println(results)
}

func TryLogin(email string, password string) LoginState {

	var usersCandidate []LoginResponse

	searchParam := [1]string{email}

	dbErr := supabase().DB.From("UserAuth").Select("*").In("emailAddress", searchParam[:]).Execute(&usersCandidate)

	if dbErr != nil {
		panic(dbErr)
	}

	if len(usersCandidate) == 0 {
		return LoginState{
			Error:   true,
			Message: "The email address or password you entered is invalid!",
		}
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(usersCandidate[0].Password), []byte(password))

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword {
		return LoginState{
			Error:   true,
			Message: "The email address or password you entered is invalid!",
		}
	} else if passwordErr != nil {
		panic(passwordErr)
	}

	return LoginState{
		Error:   false,
		Message: "Login successful!",
		Data: LoginResult{
			Uid:          usersCandidate[0].Uid,
			Username:     usersCandidate[0].Username,
			EmailAddress: usersCandidate[0].EmailAddress,
			LastLoggedIn: usersCandidate[0].LastLoginTime,
		},
	}

}
