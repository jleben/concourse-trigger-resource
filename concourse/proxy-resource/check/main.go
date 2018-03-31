package main

import (
    "encoding/json"
    //"io"
    //"ioutil"
    "os"
    //"os/exec"
    "fmt"
    "errors"
    "strings"
    //"net/http"
    "github.com/jleben/trigger-resource/protocol"
    "github.com/nlopes/slack"
)

func main() {

    var request protocol.CheckRequest

    var err error

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("Parsing request.", err)
	}

	if len(request.Source.Token) == 0 {
        fatal1("Missing source field: token.")
    }

    if len(request.Source.Channel) == 0 {
        fatal1("Missing source field: channel.")
    }

    if len(request.Source.Command) == 0 {
        fatal1("Missing source field: command.")
    }


    slack_client := slack.New(request.Source.Token)

    channel_id := get_channel_id(request, slack_client)

    params := slack.NewHistoryParameters()

    if len(request.Version.Request) != 0 {
        params.Oldest = request.Version.Request
    }

    params.Inclusive = true
    params.Count = 20

    var history *slack.History
    history, err = slack_client.GetChannelHistory(channel_id, params)

    response := protocol.CheckResponse{}

    for i := len(history.Messages)-1; i >= 0; i-- {

        msg := history.Messages[i]

        text := msg.Msg.Text
        ts := msg.Msg.Timestamp
        fmt.Fprintf(os.Stderr, "Message %s: %s \n", ts, text)

        target_version := process_message(text, request)

        if len(target_version) == 0 { continue }

        var version protocol.Version
        version.Request = ts
        version.Target = target_version

        response = append(response, version)
    }

    json.NewEncoder(os.Stdout).Encode(&response)
}

type Channel struct {
    id string
    name string
}

type ChannelsMeta struct {
    next_cursor string
}

type Channels struct {
    ok bool
    channels []Channel
    meta ChannelsMeta
}

func get_channel_id(request protocol.CheckRequest, slack_client *slack.Client) string {

    channels, get_err := slack_client.GetChannels(false)
    if get_err != nil {
        fatal("getting channels", get_err)
    }

    fmt.Fprintf(os.Stderr, "Looking for channel '%s'\n", request.Source.Channel)

    var channel_id string
    for _, channel := range channels {
        if channel.Name == request.Source.Channel {
            fmt.Fprintf(os.Stderr, "Channel: name = '%s', id = '%s'\n", channel.Name, channel.ID)
            channel_id = channel.ID
        }
    }

    if len(channel_id) == 0 {
        fatal("finding channel", errors.New("Channel name not found.") )
    }

    fmt.Fprintf(os.Stderr, "Found channel ID '%s'\n", channel_id)

    return channel_id
}

func process_message(text string, request protocol.CheckRequest) map[string]string {

    version := make(map[string]string)

    prefix := "@" + request.Source.Command

    if !strings.HasPrefix(text, prefix) {
        fmt.Fprintf(os.Stderr, "Prefix '%s' not found.\n", prefix)
        return version
    }

    version_text := strings.Trim(text[len(prefix):], " ")

    version_parts := strings.Split(version_text, ":")

    if len(version_parts) != 2 {
        fmt.Fprintf(os.Stderr, "Invalid version format: '%s'.\n", version_text)
        return version
    }

    key := version_parts[0]
    value := version_parts[1]

    fmt.Fprintf(os.Stderr, "Parsed command for version: %s: %s\n", key, value)

    version[key] = value

    return version
}

/*
func get_channel_id(request protocol.CheckRequest) {

    var done = false
    var cursor string

    for !done {
        channels := get_channels(cursor)

        for _, channel := range channels.channels {
            fmt.Fprintf(os.Stderr, "Channel: %s %s\n", channel.id, channel.name)
        }

        cursor = channels.meta.next_cursor
        done = len(cursor) == 0
    }
}

func get_channels(cursor string) (Channels) {
    url = "https://slack.com/api/channels.list?" +
        "token=" + request.Source.Token +
        "&exclude_archived=true" +
        "&exclude_members=true"

    if len(cursor) > 0 {
        url += "&cursor=" + cursor
    }

    resp, get_err := http.Get(url)
    if get_err != nil { fatal("getting channels", get_err) }

    body, read_err := ioutil.ReadAll(resp.Body)
    if read_err != nil { fatal("getting channels - reading response body", read_err) }

    var channels Channels
    parse_err := json.Unmarshall(body, &channels)
    if parse_err != nil { fatal("getting channels - parsing response body", parse_err) }

    return channels
}

func get_history(request protocol.CheckRequest, channel_id) {
    url = "https://slack.com/api/channels.history?" +
        "token=" + request.Source.Token
    response, err := http.Get()
}
*/

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}

func fatal1(reason string) {
    fmt.Fprintf(os.Stderr, reason + "\n")
    os.Exit(1)
}
