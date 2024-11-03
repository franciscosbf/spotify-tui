package api

import "fmt"

type ErrResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e ErrResponse) Error() string {
	return fmt.Sprintf("%d %s", e.Status, e.Message)
}

type UserProfile struct {
	Name      string `json:"display_name"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
}
