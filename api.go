package slack

const (
	apiBaseURL = "https://slack.com/api/"
)

type API struct {
	Token string
}

func NewAPI(token string) API {
	return API{token}
}
