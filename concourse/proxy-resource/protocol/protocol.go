package protocol

import (
    "strings"
)

type Source struct {
    Channel string `json:"channel"`
    ChannelId string `json:"channel_id"`
    BotId string `json:"bot_id"`
    Token string `json:"token"`
    Context string `json:"context"`
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

func ParseSlackRequest(text string, source *Source) *SlackRequest {

    text = strings.Trim(text, " ")

    bot_mention := "<@" + source.BotId + ">"
    if !strings.HasPrefix(text, bot_mention) { return nil }

    words := strings.Split(text, " ")
    if len(words) < 3 { return nil }

    if words[1] != source.Context { return nil }

    request := new(SlackRequest)
    request.Context = source.Context
    request.Version = Version{}

    versions := words[2:]

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
    text := request.Context
    for key, value := range request.Version {
        text += " " + key + ":" + value
    }
    return text
}
