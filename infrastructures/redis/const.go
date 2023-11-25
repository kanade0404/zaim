package redis

type RequestSecret struct {
	Secret string `json:"secret"`
	User   string `json:"user"`
}

type OauthToken struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}
