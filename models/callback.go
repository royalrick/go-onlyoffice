package models

type Callback struct {
	Actions     []Action `json:"actions"`
	ChangesUrl  string   `json:"changesurl"`
	History     History  `json:"history"`
	Key         string   `json:"key"`
	Status      int      `json:"status"`
	Users       []string `json:"users"`
	Url         string   `json:"url"`
	FileId      string   `json:"-"`
	Token       string   `json:"token,omitempty"`
	Filename    string   `json:"filename,omitempty"`
	UserAddress string   `json:"userAddress,omitempty"`
}

type Action struct {
	Type   int    `json:"type"`
	UserID string `json:"userid"`
}

type History struct {
	Changes       []Change `json:"changes,omitempty"`
	ServerVersion string   `json:"serverVersion,omitempty"`
	Created       string   `json:"created,omitempty"`
	Key           string   `json:"key,omitempty"`
	User          *User    `json:"user,omitempty"`
	Version       int      `json:"version,omitempty"`
}

type Change struct {
	Created string `json:"created"`
	User    User   `json:"user"`
}

type User struct {
	Id                string         `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	Group             string         `json:"group,omitempty"`
	ReviewGroups      []string       `json:"reviewGroups,omitempty"`
	CommentGroups     map[string]any `json:"commentGroups,omitempty"`
	UserInfoGroups    []string       `json:"userInfoGroups,omitempty"`
	Favorite          int            `json:"favorite,omitempty"`
	DeniedPermissions []string       `json:"deniedPermissions,omitempty"`
	Description       []string       `json:"description,omitempty"`
	Templates         bool           `json:"templates,omitempty"`
	Avatar            bool           `json:"avatar,omitempty"`
}
