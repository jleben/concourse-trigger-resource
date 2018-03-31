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

    if len(request.Source.Channel) == 0 && len(request.Source.ChannelId) == 0 {
        fatal1("Missing source field: channel or channel_id.")
    }

    if len(request.Source.Command) == 0 {
        fatal1("Missing source field: command.")
    }


    slack_client := slack.New(request.Source.Token)

    if len(request.Source.ChannelId) == 0 {
        request.Source.ChannelId = get_channel_id(request, slack_client)
    }

    params := slack.NewHistoryParameters()

    if request_version, ok := request.Version["request"]; ok {
        params.Oldest = request_version
        fmt.Fprintf(os.Stderr, "Request version: %s\n", request_version)
    }

    params.Inclusive = true
    params.Count = 100

    var history *slack.History
    history, err = slack_client.GetChannelHistory(request.Source.ChannelId, params)
    if err != nil {
		fatal("getting messages.", err)
	}

    response := protocol.CheckResponse{}

    for i := len(history.Messages)-1; i >= 0; i-- {

        msg := history.Messages[i]

        version := process_message(&msg, request, slack_client)

        if version != nil {
            response = append(response, version)
        }
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

func process_message(message *slack.Message, request protocol.CheckRequest, slack_client *slack.Client) protocol.Version {

    is_reply := len(message.Msg.ThreadTimestamp) > 0 &&
        message.Msg.ThreadTimestamp != message.Msg.Timestamp

    if is_reply {
        fmt.Fprintf(os.Stderr, "Message %s is a reply. Skipping.\n", message.Msg.Timestamp)
        return nil
    }

    is_by_bot := message.Msg.SubType == "bot_message" || len(message.Msg.User) == 0
    if is_by_bot {
        fmt.Fprintf(os.Stderr, "Message %s is by a bot. Skipping.\n", message.Msg.Timestamp)
        return nil
    }

    text := message.Msg.Text
    ts := message.Msg.Timestamp
    fmt.Fprintf(os.Stderr, "Message %s: %s \n", ts, text)

    /*
    if message_has_reply(message) {
        fmt.Fprintf(os.Stderr, "Message already processed previously.\n", ts)
        return nil
    }
    */

    prefix := "@" + request.Source.Command

    if !strings.HasPrefix(text, prefix) {
        fmt.Fprintf(os.Stderr, "Prefix '%s' not found.\n", prefix)
        return nil
    }

    version_text := strings.Trim(text[len(prefix):], " ")

    version_parts := strings.Split(version_text, ":")

    if len(version_parts) != 2 {
        fmt.Fprintf(os.Stderr, "Invalid version format: '%s'.\n", version_text)
        return nil
    }

    version_key := version_parts[0]
    version_value := version_parts[1]

    fmt.Fprintf(os.Stderr, "Parsed command for version: %s: %s\n", version_key, version_value)

    version := protocol.Version{
        version_key: version_value,
        "request": ts,
    }

    reply(message, version_text, request, slack_client)

    return version
}

/*
func message_has_reply(message *slack.Message) bool {
    if message.Msg.ReplyCount == 0 {
        return false
    }

    for _, reply := range message.Msg.Replies {

    }
}
*/

func reply(message *slack.Message, target_version string, request protocol.CheckRequest, slack_client *slack.Client) {
    params := slack.NewPostMessageParameters()
    params.ThreadTimestamp = message.Msg.Timestamp

    text := fmt.Sprintf("@%s %s queued.", request.Source.Command, target_version)

    _, _, err := slack_client.PostMessage(request.Source.ChannelId, text, params)
    if err != nil {
        fatal("replying", err)
    }
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
    fmt.Fprintf(os.Stderr, "error " + doing + ": " + err.Error() + "\n")
	os.Exit(1)
}

func fatal1(reason string) {
    fmt.Fprintf(os.Stderr, reason + "\n")
    os.Exit(1)
}
