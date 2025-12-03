package dto

type ResponseLov struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type ResponseFile struct {
	Metadata map[string]interface{} `json:"metadata"`
	FilePath []string               `json:"file_paths"`
}
