package protocol

type Source struct {
    Channel string `json:"channel"`
    Token string `json:"token"`
    Command string `json:"command"`
    Target interface{} `json:"target"`
}

type Version struct {
    Request string `json:"request"`
    Target interface{} `json:"target"`
}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type TargetInRequest struct {
    Source  interface{}  `json:"source"`
	Version interface{} `json:"version"`
}

type TargetInResponse struct {
    Version interface{} `json:"version"`
    Metadata interface{} `json:"metadata"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckResponse []Version
