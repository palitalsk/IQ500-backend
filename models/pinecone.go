package models

// request structure for Pinecone upsert
type PineconeUpsertRequest struct {
	Vectors   []PineconeVector `json:"vectors"`
	Namespace string           `json:"namespace,omitempty"`
}

// a single vector in Pinecone
type PineconeVector struct {
	ID       string                 `json:"id"`
	Values   []float64              `json:"values"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// request structure for Pinecone query
type PineconeQueryRequest struct {
	Vector          []float64 `json:"vector"`
	TopK            int       `json:"topK"`
	IncludeMetadata bool      `json:"includeMetadata"`
	Namespace       string    `json:"namespace,omitempty"`
}

// response structure from Pinecone query
type PineconeQueryResponse struct {
	Matches []PineconeMatch `json:"matches"`
}

// a single match from Pinecone query
type PineconeMatch struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}
