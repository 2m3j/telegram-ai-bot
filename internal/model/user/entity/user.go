package entity

import (
	"time"
)

type User struct {
	ID           uint64
	AIPlatform   string
	AIModel      string
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func NewUser(
	id uint64,
	aIPlatform string,
	aIModel string,
	username string,
	firstName string,
	lastName string,
	languageCode string,
) *User {
	now := time.Now()
	return &User{
		ID:           id,
		AIPlatform:   aIPlatform,
		AIModel:      aIModel,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		LanguageCode: languageCode,
		UpdatedAt:    now,
		CreatedAt:    now,
	}
}

func (u *User) IsEqualUserInfo(username string, firstName string, lastName string, languageCode string) bool {
	return u.Username == username && u.FirstName == firstName && u.LastName == lastName && u.LanguageCode == languageCode
}

func (u *User) ChangeAISettings(aiPlatform string, aiModel string) {
	u.AIPlatform = aiPlatform
	u.AIModel = aiModel
	u.UpdatedAt = time.Now()
}

func (u *User) ChangeUserInfo(username string, firstName string, lastName string, languageCode string) {
	u.Username = username
	u.FirstName = firstName
	u.LastName = lastName
	u.LanguageCode = languageCode
	u.UpdatedAt = time.Now()
}
