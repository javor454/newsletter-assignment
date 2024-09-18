package response

type HealthStatus struct {
	Status     string      `json:"status"`
	Indicators []Indicator `json:"indicators"`
}

type Indicator struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
