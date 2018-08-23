package splat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// App is the container for Slack config and commands
type App struct {
	SigningSecret string
	commands      map[string]Command
}

// New creates a new SlackApp object (please don't kill me for using 'object')
func New(secret string) *App {
	return &App{secret, make(map[string]Command)}
}

// RegisterCommand creates a new command to be executed when it is called from Slack
func (s *App) RegisterCommand(name string, handler func(*SlashRequest)) {
	s.commands["/"+name] = Command{"/" + name, handler}
}

// Open starts HTTP server listening on addr at endpoint
func (s *App) Open(addr string, endpoint string) error {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		payload, err := s.fromRequest(r)
		if err != nil {
			log.Println(err)
			return
		}

		for k, v := range s.commands {
			if k == payload.Command {
				go v.Execute(payload)
				break
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(addr, nil)
}

// Parses the body data into a program readable format
func (s *App) fromRequest(r *http.Request) (p *SlashRequest, err error) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	err = s.checkSignature(r, body)
	if err != nil {
		return nil, err
	}

	// I can't find a good way to do this. Please start using JSON, Slack
	kv := make(map[string]string)
	split := strings.Split(string(body), "&")
	for _, v := range split {
		t := strings.Split(v, "=")
		t[0], _ = url.QueryUnescape(t[0])
		t[1], _ = url.QueryUnescape(t[1])
		kv[t[0]] = t[1]
	}

	// Map all body parameters to a SlashRequest struct
	p = new(SlashRequest)
	if val, ok := kv["token"]; ok {
		p.Token = val
	}
	if val, ok := kv["team_id"]; ok {
		p.TeamID = val
	}
	if val, ok := kv["team_domain"]; ok {
		p.TeamDomain = val
	}
	if val, ok := kv["enterprise_id"]; ok {
		p.EnterpriseID = val
	}
	if val, ok := kv["enterprise_name"]; ok {
		p.EnterpriseName = val
	}
	if val, ok := kv["channel_id"]; ok {
		p.ChannelID = val
	}
	if val, ok := kv["channel_name"]; ok {
		p.ChannelName = val
	}
	if val, ok := kv["user_id"]; ok {
		p.UserID = val
	}
	if val, ok := kv["user_name"]; ok {
		p.UserName = val
	}
	if val, ok := kv["command"]; ok {
		p.Command = val
	}
	if val, ok := kv["text"]; ok {
		p.Text = val
	}
	if val, ok := kv["response_url"]; ok {
		p.ResponseURL = val
	}
	if val, ok := kv["trigger_id"]; ok {
		p.TriggerID = val
	}

	return p, nil
}

// Validates the request signature from Slack
func (s *App) checkSignature(r *http.Request, body []byte) error {

	timestamp := r.Header.Get("X-Slack-Request-Timestamp")
	time := time.Now().Unix()
	stamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return err
	}

	// Check if request is over 5 minutes old and reject if it is. Prevents replay attack
	if time-stamp > 300 {
		return errors.New("invalid timestamp in request header")
	}

	base := "v0:" + timestamp + ":" + string(body)
	sign := r.Header.Get("X-Slack-Signature")
	key := []byte(s.SigningSecret)

	h := hmac.New(sha256.New, key)
	h.Write([]byte(base))
	result := "v0=" + hex.EncodeToString(h.Sum(nil))

	if sign == result {
		return nil
	}

	return errors.New("invalid request signature")
}
