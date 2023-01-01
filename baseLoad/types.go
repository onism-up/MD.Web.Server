package baseLoad

type MDInterface struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	CreateTime    string `json:"create_time"`
	LastChangTime string `json:"lastChangTime"`
	IsPublic      bool   `json:"is_public"`
}

type MDBaseType = map[string]MDInterface

type UserInterface struct {
	Name  string `json:"name"`
	Pwd   string `json:"pwd"`
	Token string `json:"token"`
	Init  bool   `json:"init"`
}
