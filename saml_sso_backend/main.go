package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso"
)

const PROJECT_ID = "YOUR_STYTCH_PROJECT_ID"
const SECRET_KEY = "YOUR_STYTCH_SECRET_KEY"
const PORT = ":3001"

func main() {
	http.HandleFunc("/authenticate", authenticate)
	fmt.Println("Server is listening on port", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		httpError(w, "Failed to authenticate user. Token is null")
		return
	}
	client, err := b2bstytchapi.NewClient(
		PROJECT_ID,
		SECRET_KEY,
	)
	if err != nil {
		httpError(w, fmt.Sprintf("Error instantiating API client %s\n", err))
		return
	}
	params := &sso.AuthenticateParams{
		SSOToken: token,
	}

	resp, err := client.SSO.Authenticate(context.Background(), params)
	if err != nil {
		httpError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("<h3>SAML SSO Login Process Complete!</h3><br/><p>You  have been successfully Logged in to the %s organization. Your role is %s and your session expires at %s.</p><br/><br/>Session Token: %s<br/>Session JWT: %s", resp.Organization.OrganizationName, resp.MemberSession.Roles[0], resp.MemberSession.ExpiresAt, resp.SessionToken, resp.SessionJWT)
	fmt.Fprint(w, message)
}

func httpError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
