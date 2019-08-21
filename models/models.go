package models

type Version struct {
	ID string `json:"id"`
}

type MetadataEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Source struct {
	Host     string   `json:"host"`
	Password *string  `json:"password"`
	Key      string   `json:"key"`
	Fields   []string `json:"fields"`
}

type GetParams struct {
}

type PutParams struct {
}

type CheckRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
}

type CheckResponse []Version

type GetRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type GetResponse struct {
	Version  Version         `json:"version"`
	Metadata []MetadataEntry `json:"metadata"`
}

type PutRequest struct {
}

type PutResponse struct {
}
