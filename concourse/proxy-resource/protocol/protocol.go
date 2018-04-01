package protocol

import (
    "strings"
)

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
	Params interface{} `json:"params"`
}

type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type TargetInRequest struct {
    Source  interface{}  `json:"source"`
	Version interface{} `json:"version"`
	Params interface{} `json:"params"`
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

type SlackRequest struct {
    Context string
    Version Version
}

func ParseSlackRequest(text string, context string) *SlackRequest {
    if !strings.HasPrefix(text, "@" + context) { return nil }
    words := strings.Split(text, " ")

    if len(words) < 2 { return nil }

    request := new(SlackRequest)
    request.Context = context
    request.Version = Version{}

    versions := words[1:]

    for _, version := range versions {
        parts := strings.Split(version, ":")
        if len(parts) != 2 { return nil }
        key := parts[0]
        value := parts[1]
        request.Version[key] = value
    }

    return request
}

func (request *SlackRequest) String() string {
    text := "@" + request.Context
    for key, value := range request.Version {
        text += " " + key + ":" + value
    }
    return text
}
