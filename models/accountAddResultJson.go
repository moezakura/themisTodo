package models

type AccountAddResultJson struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}
