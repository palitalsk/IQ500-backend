package models

type GenerateRequest struct {
	Context string `json:"context"`
	Query   string `json:"query"`
	Gender  string `json:"gender,omitempty"`
}

type GenerateResponse struct {
	Answer string `json:"answer"`
}

type ChatRequest struct {
	Query     string `json:"query"`
	Namespace string `json:"namespace"`
	TopK      int    `json:"topK"`
	Gender    string `json:"gender,omitempty"`
}

type ChatResponse struct {
	Results string `json:"results"`
}

type UploadResponse struct {
	Message   string `json:"message"`
	Chunks    int    `json:"chunks"`
	Namespace string `json:"namespace"`
}
