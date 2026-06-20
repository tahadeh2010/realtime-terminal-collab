package domain

type Session struct {
	ID      string
	Clients []*Client
	Host    *Client
}
