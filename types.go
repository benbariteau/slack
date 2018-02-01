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

var _ Response = &BasicResponse{}

func (r *BasicResponse) Ok() bool {
	return r.OK
}

func (r *BasicResponse) Error() error {
	if r.OK {
		return nil
	}
	return SlackError{r.ErrorString}
}

type User struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Deleted           bool        `json"deleted"`
	Color             string      `json:"color"`
	Profile           UserProfile `json:"profile"`
	IsAdmin           bool        `json:"is_admin"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	Has2FA            bool        `json:"has_2fa"`
	HasFiles          bool        `json:"has_files"`
}

type UserProfile struct {
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	RealName    string `json:"real_name"`
	Email       string `json:"email"`
	Skype       string `json:"skype"`
	Phone       string `json:"phone"`
	Image24     string `json:"image_24"`
	Image32     string `json:"image_32"`
	Image48     string `json:"image_48"`
	Image72     string `json:"image_72"`
	Image192    string `json:"image_192"`
}

type Channel struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	IsChannel          string                 `json:"is_channel"`
	IsGroup            string                 `json:"is_group"`
	Created            int                    `json:"created"`
	Creator            string                 `json:"creator"`
	IsArchived         bool                   `json:"is_archived"`
	IsGeneral          bool                   `json:"is_general"`
	Members            []string               `json:"members"`
	Topic              Topic                  `json:"topic"`
	Purpose            Topic                  `json:"purpose"`
	IsMember           bool                   `json:"is_member"`
	LastRead           string                 `json:"last_read"`
	UnreadCount        int                    `json:"unread_count"`
	UnreadCountDisplay int                    `json:"unread_count_display"`
	Latest             map[string]interface{} `json:"latest"`
}

type Topic struct {
	Value   string `json:"value"`
	Creator string `json:"creator"`
	LastSet int    `json:"last_set"`
}
