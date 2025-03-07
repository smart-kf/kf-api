package index

type AiMessage struct {
	Id      int64  `json:"id"`
	CardId  string `json:"cardId"`
	Keyword string `json:"keyword"`
	Message string `json:"message"`
}
