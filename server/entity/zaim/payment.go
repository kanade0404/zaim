package zaim

type PaymentParameter struct {
	CategoryID    string `json:"category_id"`
	GenreID       string `json:"genre_id"`
	Amount        string `json:"amount"`
	Date          string `json:"date"`
	FromAccountID string `json:"from_account_id"`
	Name          string `json:"name"`
	Place         string `json:"place"`
	Comment       string `json:"comment"`
}

type Money struct {
	ID       int     `json:"id"`
	PlaceUID *string `json:"place_uid"`
	Modified string  `json:"modified"`
}
type User struct {
	InputCount   int    `json:"input_count"`
	RepeatCount  int    `json:"repeat_count"`
	DayCount     int    `json:"day_count"`
	DataModified string `json:"data_modified"`
}
type Place struct {
	ID                int    `json:"id"`
	UserID            int    `json:"user_id"`
	GenreID           int    `json:"genre_id"`
	AccountID         int    `json:"account_id"`
	TransferAccountID int    `json:"transfer_account_id"`
	Mode              string `json:"mode"`
	PlaceUID          string `json:"place_uid"`
	Service           string `json:"service"`
	Name              string `json:"name"`
	OriginalName      string `json:"original_name"`
	Tel               string `json:"tel"`
	Count             int    `json:"count"`
	PlacePatternID    int    `json:"place_pattern_id"`
	CalcFlag          int    `json:"calc_flag"`
	EditFlag          int    `json:"edit_flag"`
	Active            int    `json:"active"`
	Modified          string `json:"modified"`
	Created           string `json:"created"`
}
type PaymentResponse struct {
	Stamps    *string  `json:"stamps"`
	Banners   []string `json:"banners"`
	Money     Money    `json:"money"`
	User      User     `json:"user"`
	Place     Place    `json:"place"`
	Requested int      `json:"requested"`
}
type Order string

const OrderDate Order = "date"

type Mode string

const Payment Mode = "payment"

type ListPaymentParameter struct {
	CategoryID *int    `json:"category_id"`
	GenreID    *int    `json:"genre_id"`
	Mode       Mode    `json:"mode"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	Order      *Order  `json:"order"`
	Page       *int    `json:"page"`
	Limit      *int    `json:"limit"`
	GroupBy    *string `json:"group_by"`
}

type ListPaymentResponse struct {
	Money []struct {
		ID            int    `json:"id"`
		Mode          Mode   `json:"mode"`
		UserID        int    `json:"user_id"`
		Date          string `json:"date"`
		CategoryID    int    `json:"category_id"`
		GenreID       int    `json:"genre_id"`
		ToAccountID   int    `json:"to_account_id"`
		FromAccountID int    `json:"from_account_id"`
		Amount        int    `json:"amount"`
		Comment       string `json:"comment"`
		Active        int    `json:"active"`
		Name          string `json:"name"`
		ReciptID      int    `json:"recipt_id"`
		Place         string `json:"place"`
		Created       string `json:"created"`
		CurrencyCode  string `json:"currency_code"`
	} `json:"money"`
}
