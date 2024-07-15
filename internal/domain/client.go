package domain

type Client struct {
	ClientName   string `form:"client-name" binding:"required"`
	ExchangeName string `form:"exchange" binding:"required"`
	Label        string `form:"label" binding:"required"`
	Pair         string `form:"pair" binding:"required"`
}
