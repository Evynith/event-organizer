package service

import (
	auth "main/internal/repository/auth"
)

/*
Corrobora que el nombre de usuario pertenezca a un usuario guardado y la contraseña para dicho
usuario sea igual a la recibida
*/
func LoginUser(username string, password string) bool {
	user, err := auth.SearchUser(username)

	if err != nil {
		return false
	}
	return user.Password == password
}

/*
Persiste el token recibido en la base de datos asociandolo al usuario recibido
*/
func SaveToken(username string, token string) bool {
	err := auth.PersistToken(username, token)
	if err != nil {
		return false
	}
	return true
}

func ExistsToken(token string, id string) bool {
	return auth.ExistsToken(token, id)
}

func DataUser(username string) (string, string) {
	u, err := auth.SearchUser(username)
	if err != nil {
		return "", ""
	}
	return u.Type, u.ID.Hex()
}
