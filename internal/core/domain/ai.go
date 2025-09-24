package domain

type AIPredictionResult struct {
	IsSlip    bool                   `json:"is_slip" bson:"is_slip"`
	OCRResult map[string]interface{} `json:"ocr_result" bson:"ocr_result"`
}
