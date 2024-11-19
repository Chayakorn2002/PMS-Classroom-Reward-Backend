package classroom

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/constants"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

func GetGoogleClassroomClient(ctx context.Context) (*classroom.Service, error) {
	// Setup Google OAuth client
	googleCredentialFile := constants.GoogleCredentialFilePath
	data, err := os.ReadFile(googleCredentialFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read google credential file: %v", err)
	}

	googleCredentialConfig, err := google.ConfigFromJSON(
		data,
		classroom.ClassroomCoursesReadonlyScope,
		classroom.ClassroomRostersReadonlyScope,
		classroom.ClassroomProfileEmailsScope,
		classroom.ClassroomStudentSubmissionsStudentsReadonlyScope,
		classroom.ClassroomStudentSubmissionsMeReadonlyScope,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create google config from json: %v", err)
	}
	googleCredentialConfig.RedirectURL = "http://localhost:8080/oauth2callback"

	client, err := getClient(ctx, googleCredentialConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get google client: %v", err)
	}

	srv, err := classroom.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create classroom service: %v", err)
	}

	slog.InfoContext(ctx, "🔌 Connected to Google Classroom API")
	return srv, nil
}

// Retrieves an authenticated HTTP client using a refresh token
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	tokenFile := constants.GoogleOAuthTokenFilePath

	// Check if the token already exists
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		// If the token doesn't exist, initiate the OAuth flow
		token, err = getTokenFromWeb(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("failed to get token from web: %v", err)
		}

		err = saveToken(ctx, tokenFile, token)
		if err != nil {
			return nil, fmt.Errorf("failed to save token: %v", err)
		}
	}

	// Create a TokenSource using the refresh token
	tokenSource := config.TokenSource(context.Background(), token)
	return oauth2.NewClient(context.Background(), tokenSource), nil
}

// Request a token from the web (one-time setup)
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	slog.InfoContext(ctx, fmt.Sprintf("Go to the following link in your browser and enter the authorization code:\n%v\n", authURL))

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return token, nil
}

// Retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// Saves a token to a file path
func saveToken(ctx context.Context, file string, token *oauth2.Token) error {
	slog.InfoContext(ctx, fmt.Sprintf("Saving credential file to: %s\n", file))
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}
