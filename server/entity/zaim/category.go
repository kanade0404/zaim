package zaim

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Mode     string `json:"mode"`
	Sort     int    `json:"sort"`
	Active   int    `json:"active"`
	Modified string `json:"modified"`
}
