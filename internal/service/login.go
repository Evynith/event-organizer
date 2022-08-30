package service

import (
	auth "main/internal/repository/auth"
)

func LoginUser(email string, password string) bool {
	user, err := auth.SearchUser(email)

	if err != nil {
		return false
	}
	return user.Password == password
}

func SaveToken(username string, token string) bool {
	err := auth.PersistToken(username, token)
	if err != nil {
		return false
	}
	return true
}

func ExistsToken(token string) bool {
	return auth.ExistsToken(token)
}

func TypeUser(username string) string {
	u, err := auth.SearchUser(username)
	if err != nil {
		return ""
	}
	return u.Type
}
