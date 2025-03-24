package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// User represents the authenticated user
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// UserInfoResponse represents the OAuth2 user info response
type UserInfoResponse struct {
	ID            string `json:"id,omitempty"`
	Sub           string `json:"sub,omitempty"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture,omitempty"`
}

// AuthMiddleware handles OAuth2 authentication
type AuthMiddleware struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	JWTSecret    []byte
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(clientID, clientSecret, redirectURL, authURL, tokenURL, userInfoURL, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		UserInfoURL:  userInfoURL,
		JWTSecret:    []byte(jwtSecret),
	}
}

// GetLoginURL generates the OAuth2 provider login URL
func (m *AuthMiddleware) GetLoginURL(state string) string {
	return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=email%%20profile&state=%s",
		m.AuthURL, m.ClientID, m.RedirectURL, state)
}

// ExchangeCodeForToken exchanges the authorization code for an access token
func (m *AuthMiddleware) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	data := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, m.ClientID, m.ClientSecret, m.RedirectURL)

	req, err := http.NewRequest("POST", m.TokenURL, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange code: %s - %s", resp.Status, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetUserInfo retrieves user information using the access token
func (m *AuthMiddleware) GetUserInfo(accessToken string) (*UserInfoResponse, error) {
	req, err := http.NewRequest("GET", m.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s - %s", resp.Status, string(body))
	}

	var userInfo UserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// GenerateJWT generates a JWT token for the authenticated user
func (m *AuthMiddleware) GenerateJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(m.JWTSecret)
}

// ValidateJWT validates the JWT token
func (m *AuthMiddleware) ValidateJWT(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &User{
			ID:    claims["id"].(string),
			Email: claims["email"].(string),
			Name:  claims["name"].(string),
		}
		return user, nil
	}

	return nil, errors.New("invalid token")
}

// Authenticate verifies the user's token
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Check that it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		token := parts[1]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token cannot be empty"})
			return
		}

		// Validate the JWT token
		user, err := m.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set the user in the context
		c.Set("user", user)
		c.Next()
	}
}

// HandleCallback processes the OAuth callback
func (m *AuthMiddleware) HandleCallback(c *gin.Context) {
	// Get the authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is required"})
		return
	}

	// Exchange the code for a token
	tokenResp, err := m.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to exchange code: %v", err)})
		return
	}

	// Get user information
	userInfo, err := m.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get user info: %v", err)})
		return
	}

	// Create user object
	user := &User{
		ID:    userInfo.ID,
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}

	// Generate JWT
	jwtToken, err := m.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate JWT: %v", err)})
		return
	}

	// Return the JWT token
	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"user":  user,
	})
}

// InitiateLogin initiates the OAuth login flow
func (m *AuthMiddleware) InitiateLogin(c *gin.Context) {
	// Generate a random state
	state := fmt.Sprintf("%d", time.Now().UnixNano())
	
	// Store the state in the session
	// In a real application, you would use a session manager
	
	// Redirect to the OAuth provider's login page
	loginURL := m.GetLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, loginURL)
}