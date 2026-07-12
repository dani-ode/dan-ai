package schema

type Chunk struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Keywords []string `json:"keywords"`
}

type KnowledgeBuilderResponse struct {
	Chunks []Chunk `json:"chunks"`
}
