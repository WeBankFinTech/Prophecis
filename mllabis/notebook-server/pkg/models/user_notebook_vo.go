package models

type UserNotebookVO struct {
	Role          string   `json:"role,omitempty"`
	NamespaceList []string `json:"namespaceList,omitempty"`
}

type PageListVO struct {
	List  interface{} `json:"list"`
	Pages int         `json:"pages"`
	Total int         `json:"total"`
}
