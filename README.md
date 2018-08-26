# Splat - a small Slack app library
Currently Splat just uses the slash commands API for Slack. If you are for some reason considering using this library, don't. This library exists as a result of a 'reinvent the wheel' situation of me wanting to practice API implementation.

### Usage
This example registers a command `/test` (the '/' is omitted in the register call) and listens for requests on port `:3000` at the `/slackcommands` endpoint
```go
app := splat.New("YourSigningSecret")

app.RegisterCommand("test", func(r *splat.SlashRequest) {
    r.Write(&splat.Response{Text: "Hello, World!"})
})

log.Fatal(app.Open(":3000", "/slackcommands"))
```