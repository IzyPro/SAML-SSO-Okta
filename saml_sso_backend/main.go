package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/sso"
)

const PORT = ":3001"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/authenticate", authenticate)
	fmt.Println("Server is listening on port", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	PROJECT_ID := os.Getenv("STYTCH_PROJECT_ID")
	SECRET_KEY := os.Getenv("STYTCH_SECRET_KEY")

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
	http.Redirect(w, r, "http://localhost:3000/success", http.StatusSeeOther)
	fmt.Println(resp.StatusCode)
}

func httpError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
