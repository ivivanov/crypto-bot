package response

type Error struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
	Status string `json:"status"`
}
