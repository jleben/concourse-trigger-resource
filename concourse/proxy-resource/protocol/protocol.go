package protocol

type Source struct {
    Channel string `json:"channel"`
    ChannelId string `json:"channel_id"`
    Token string `json:"token"`
    Command string `json:"command"`
    Target interface{} `json:"target"`
}

type Version map[string]string

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
    Metadata Metadata `json:"metadata"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckResponse []Version
