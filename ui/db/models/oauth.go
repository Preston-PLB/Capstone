package models

import (
	"time"
)

type OauthCredential struct {
	AccessToken  string    `bson:"access_token,omitempty" json:"access_token,omitempty"`
	ExpiresIn    int       `bson:"expires_in,omitempty" json:"expires_in,omitempty"`
	ExpiresAt    time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	TokenType    string    `bson:"token_type,omitempty" json:"token_type,omitempty"`
	Scope        string    `bson:"scope,omitempty" json:"scope,omitempty"`
	RefreshToken string    `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
}

type OauthRefreshBody struct {
	ClientId     string `json:"cleint_id"`
	ClientSecret string `json:"cleint_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}


