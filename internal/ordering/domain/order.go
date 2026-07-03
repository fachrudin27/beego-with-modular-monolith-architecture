package domain

import "time"

type Order struct {
	Id        string    `json:"id"`
	UserId    int64     `json:"user_id"`
	OrderName string    `json:"order_name"`
	CreatedAt time.Time `json:"created_at"`
}
