package api

type Params struct {
	Args interface{}
	Request *Request
	User *User
}

type Request struct {
	IPAdress string
	Token string
}