package entities

type Context struct {
	IP  		string 	`json:"ip" valid:"required"`
	UserAgent 	string 	`json:"user_agent" valid:"required"`
}
