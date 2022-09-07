package dto

type Result struct {
	Total        int           `json:"total"`
	Duration     string        `json:"duration"`
	QueryResults []QueryResult `json:"queryResults"`
}
