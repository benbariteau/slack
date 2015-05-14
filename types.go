package slack

type SlackError struct {
	ErrorString string
}

func (err SlackError) Error() string {
	return err.ErrorString
}

type Response interface {
	Ok() bool
	Error() error
}

type BasicResponse struct {
	OK          bool   `json:"ok"`
	ErrorString string `json:"error"`
}

var _ Response = BasicResponse{}

func (r BasicResponse) Ok() bool {
	return r.OK
}

func (r BasicResponse) Error() error {
	return SlackError{r.ErrorString}
}
