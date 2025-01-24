package configuration

type OsekPaHesh struct {
	Osek struct {
		Name    string `json:"name"`
		Tz      string `json:"tz"`
		Address string `json:"address"`
		Email   string `json:"email"`
		Sign    string `json:"sign"`
		Account []struct {
			Id       int    `json:"id"`
			Currency string `json:"currency"`
			Number   string `json:"number"`
		} `json:"account"`
		Service []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"service"`
	} `json:"osek"`
	Client []struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
		Tz   string `json:"tz,omitempty"`
	} `json:"client"`
	Transaction []struct {
		Receipt uint64  `json:"receipt"`
		Date    string  `json:"date"`
		Service int     `json:"service"`
		Client  int     `json:"client"`
		Amount  int     `json:"amount"`
		Account int     `json:"account"`
		Rate    float64 `json:"rate,omitempty"`
		Total   float64 `json:"total"`
	} `json:"transaction"`
}
