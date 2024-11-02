package api

type ErrResponse struct {
	Message string
	Code    int
}

func (e ErrResponse) Error() string {
	return e.Message
}

type UserProfile struct {
	Name      string `json:"display_name"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
}
