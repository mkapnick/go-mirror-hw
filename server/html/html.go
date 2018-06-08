package html

import (
	"fmt"
	"net/http"
)

const (
	landingHtml = `
	<h2>Log in or Sign up</h2>
	<br />
	<h4>Login</h4>
	<form action="/auth/login" method="POST">
		<input type="text" name="user" placeholder="username">
		<input type="password" name="pass" placeholder="password">
		<input type="submit">
	</form>
	<br />
	<br />
	<h4>Sign up</h4>
	<form action="/auth/create" method="POST">
		<input type="text" name="user" placeholder="username">
		<input type="password" name="pass" placeholder="password">
		<input type="submit">
	</form>
	`
)

// serve the form
func LandingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, landingHtml)
}
