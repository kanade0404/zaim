package zaim

type Account struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Modified  string `json:"modified"`
	Sort      int    `json:"sort"`
	Active    int    `json:"active"`
	WebsiteID int    `json:"website_id"`
}
