package response

type HealthStatus struct {
	Status     string      `json:"status" example:"healthy"`
	Indicators []Indicator `json:"indicators"`
}

type Indicator struct {
	Name   string `json:"name" example:"postgres"`
	Status string `json:"status" example:"healthy"`
}
