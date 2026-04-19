package model

type valid_status int

const (
	Aberta valid_status = iota
	Fechada
)

func (s valid_status) String() string {
	return [...]string{"Aberta", "Fechada"}[s]
}

type Invoice struct {
	ID             int    `json:"id"`
	CURRENT_STATUS string `json:"current_status"`
}
