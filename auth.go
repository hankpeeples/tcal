package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/gookit/slog"
	"golang.org/x/oauth2"
)

// GetClient returns a new client
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile("token.json")
	if err != nil {
		slog.Info("Getting new token.")
		getTokenFromWeb(config)
	} else {
		slog.Info("Using stored token.")
	}

	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) {
	// Start web server and listen to callback URL
	server := &http.Server{Addr: config.RedirectURL}

	// Add handler to grab auth code and close server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get auth code
		authCode := r.URL.Query().Get("code")
		if authCode == "" {
			slog.Error("Unable to get auth code from server...")
			io.WriteString(w, "Error: could not find auth code in this request...\n")

			cleanup(server)
			return
		}

		tok, err := config.Exchange(context.TODO(), authCode)
		if err != nil {
			slog.Fatalf("Unable to authorize token: %v", err)
		}

		saveToken("token.json", tok)

		// return an indication of success to the caller
		io.WriteString(w, `
		<html>
			<body>
				<h1>Login successful!</h1>
				<h2>You can now close this window.</h2>
			</body>
		</html>`)

		slog.Infof("Successfully authenticated with google calendar API.")

		// close the HTTP server
		cleanup(server)
	})

	// parse the redirect URL for the port number
	u, err := url.Parse(config.RedirectURL)
	if err != nil {
		slog.Fatalf("Bad redirect URL: %s\n", err)
	}

	// set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	l, err := net.Listen("tcp", port)
	if err != nil {
		slog.Fatalf("Can't listen to port %s: %s\n", port, err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Press [ enter ] to open the authorization link in your browser...")
	fmt.Scanln()
	if err := openInBrowser(authURL); err != nil {
		log.Fatalf("Unable to open the authorization link: %v", err)
	}

	server.Serve(l)
}

func cleanup(server *http.Server) {
	slog.Infof("Stopping and cleaning auth server...")
	go server.Close()
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	slog.Infof("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0600)
	if err != nil {
		slog.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func openInBrowser(url string) error {
	msg := "Opening browser for "
	var err error
	switch runtime.GOOS {
	case "linux":
		slog.Info("%s linux", msg)
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		slog.Info("%s windows", msg)
		err = exec.Command("rundll32",
			"url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		slog.Info("%s macOS", msg)
		err = exec.Command("open", url).Start()
	default:
		slog.Warnf("Sorry, I can't open your browser at this time. Please use the following link: %s\n", url)
	}
	if err != nil {
		return err
	}

	return nil
}
