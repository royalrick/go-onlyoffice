package models

type Config struct {
	Type         string       `json:"type"`
	Document     Document     `json:"document"`
	DocumentType string       `json:"documentType"`
	EditorConfig EditorConfig `json:"editorConfig"`
	Token        string       `json:"token,omitempty"`
}

type Document struct {
	FileType      string        `json:"fileType"`
	Key           string        `json:"key,omitempty"`
	Title         string        `json:"title"`
	Url           string        `json:"url"`
	Info          MetaInfo      `json:"info"`
	Version       string        `json:"version,omitempty"`
	Permissions   Permissions   `json:"permissions,omitempty"`
	ReferenceData ReferenceData `json:"referenceData"`
}

type MetaInfo struct {
	Author   string `json:"owner"`
	Created  string `json:"uploaded"`
	Favorite any    `json:"favorite,omitempty"`
}

type Permissions struct {
	Chat           bool           `json:"chat"`
	Comment        bool           `json:"comment,omitempty"`
	Download       bool           `json:"download"`
	Edit           bool           `json:"edit"`
	FillForms      bool           `json:"fillForms,omitempty"`
	Print          bool           `json:"print,omitempty"`
	Review         bool           `json:"review,omitempty"`
	Protect        bool           `json:"protect,omitempty"`
	RewiewGroups   []string       `json:"reviewGroups,omitempty"`
	UserInfoGroups []string       `json:"userInfoGroups,omitempty"`
	CommentGroups  map[string]any `json:"commentGroups,omitempty"`
}

type ReferenceData struct {
	FileKey string `json:"fileKey"`
	Link    any    `json:"link,omitempty"`
}

type EditorConfig struct {
	User          UserInfo      `json:"user"`
	CallbackUrl   string        `json:"callbackUrl"`
	Customization Customization `json:"customization,omitempty"`
	Lang          string        `json:"lang,omitempty"`
	Mode          string        `json:"mode,omitempty"`
	Templates     []Template    `json:"templates,omitempty"`
}

type UserInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image,omitempty"`
}

type Goback struct {
	RequestClose bool `json:"requestClose"`
}

type Customization struct {
	About      bool           `json:"about"`
	Comments   bool           `json:"comments,omitempty"`
	Feedback   bool           `json:"feedback"`
	Forcesave  bool           `json:"forcesave,omitempty"`
	SubmitForm bool           `json:"submitForm,omitempty"`
	Goback     Goback         `json:"goback,omitempty"`
	Close      map[string]any `json:"close,omitempty"`
}

type Template struct {
	Image string `json:"image,omitempty"`
	Title string `json:"title,omitempty"`
	Url   string `json:"url,omitempty"`
}

type EditorParams struct {
	Filename    string
	Mode        string
	Type        string
	Language    string
	UserId      string
	UserName    string
	UserEmail   string
	CallbackUrl string
	CanEdit     bool
	CanDownload bool
	ReadOnly    bool
}
