package auth

import (
	"crypto/rsa"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	privKeyPath = "/keys/app.rsa"
	pubKeyPath  = "/keys/app.rsa.pub"
)

const (
	tokenName = "AccessToken"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// poor man's db for storing user credentials
var db = make(map[string]string)

// poor man's db for keeping track of clients
var dbClients = make(map[string]UserInfo)

// stored in jwt token
type UserInfo struct {
	UserId   string
	Username string
	// every user has their own channelId to allow direct messages
	ChannelId string
	jwt.StandardClaims
}

// session user is used throughout the app to keep track of the session
// user
var SessionUser UserInfo

// read the key files before starting http handlers
func init() {
	pwd, _ := os.Getwd()
	signBytes, err := ioutil.ReadFile(pwd + privKeyPath)
	if err != nil {
		fmt.Printf("%s", err.Error())
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		fmt.Printf("%s", err.Error())
		log.Fatal(err)
	}

	verifyBytes, err := ioutil.ReadFile(pwd + pubKeyPath)
	if err != nil {
		fmt.Printf("%s", err.Error())
		log.Fatal(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		fmt.Printf("%s", err.Error())
		log.Fatal(err)
	}
}

func LoginFilter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(tokenName)
		println(tokenCookie)

		// no cookie found, redirect to auth flow
		if err == http.ErrNoCookie {
			fmt.Println("No cookie found")
			http.Redirect(w, r, "/", 301)
			return
		}

		// some other error - yikes
		if err != nil {
			fmt.Println("Error occurred with cookie")
			http.Redirect(w, r, "/", 301)
			return
		}

		userInfo := UserInfo{}

		// validate token
		token, verr := jwt.ParseWithClaims(tokenCookie.Value, &userInfo, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

		if !token.Valid {
			fmt.Println("Error occurred with cookie")
			http.Redirect(w, r, "/", 301)
			return
		}

		if verr != nil {
			valerr := err.(*jwt.ValidationError)
			switch valerr.Errors {
			case jwt.ValidationErrorExpired:
				fmt.Println("Token expired")
				http.Redirect(w, r, "/", 301)
				return

			default:
				fmt.Println("Error parsing token")
				http.Redirect(w, r, "/", 301)
				return
			}
		}

		fmt.Println("userInfo.Username", userInfo.Username)
		fmt.Println("userInfo.UserId", userInfo.UserId)
		fmt.Println("userInfo.ChannelId", userInfo.ChannelId)

		// keep track of the session user
		SessionUser = userInfo

		// add the session user to the client pool
		dbClients[userInfo.UserId] = userInfo
		next(w, r)
	}
}

// read the form values and create jwt
func AuthLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("user")
	pass := r.FormValue("pass")

	if len(username) == 0 || len(pass) == 0 {
		fmt.Println("Username and password are required fields")
		http.Redirect(w, r, "/", 301)
	}

	if len(db[username]) == 0 {
		fmt.Println("user does not exist in db")
		http.Redirect(w, r, "/", 301)
	}

	if db[username] != pass {
		fmt.Println("invalid credentials")
		http.Redirect(w, r, "/", 301)
	}

	log.Printf("authenticate: username[%s] pass[%s]\n", username, pass)

	// create a signer for rsa 256
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), &UserInfo{
		// userId
		uuid.Must(uuid.NewV4()).String(),
		// username
		username,
		// channelId
		uuid.Must(uuid.NewV4()).String(),
		// expiration
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Printf("Token Signing error: %v\n", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:       tokenName,
		Value:      tokenString,
		Path:       "/",
		RawExpires: "0",
	})

	// valid login, redirect to `/chat`
	http.Redirect(w, r, "/chat", 301)
}

func AuthCreate(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("user")
	pass := r.FormValue("pass")

	if len(username) == 0 || len(pass) == 0 {
		fmt.Println("Username and password are required fields")
		http.Redirect(w, r, "/", 301)
	}

	// if the username exists, validate their login
	if len(db[username]) != 0 {
		AuthLogin(w, r)
		return
	}

	// set the user in the db
	db[username] = pass

	// create a jwt token for the user
	AuthLogin(w, r)
	return
}
