package splat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// App is the container for Slack config and commands
type App struct {
	SigningSecret string
	commands      map[string]Command
	actions       []Action
}

// New creates a new SlackApp object (please don't kill me for using 'object')
func New(secret string) *App {
	return &App{secret, make(map[string]Command), make([]Action, 5)}
}

// RegisterCommand creates a new command to be executed when it is called from Slack
func (s *App) RegisterCommand(name string, handler func(*Payload) *Response) {
	s.commands[name] = Command{name, handler}
}

// RegisterAction creates actions on the Slack app
func (s *App) RegisterAction(callbackID, endpoint string, handler func(*ActionPayload)) error {
	if len(s.actions) == 5 {
		return errors.New("cannot add another action (exceeded limit of 5)")
	}

	s.actions[len(s.actions)-1] = Action{callbackID, endpoint, handler}
	return nil
}

// Open starts HTTP server listening on addr at endpoint
func (s *App) Open(addr string, endpoint string) error {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		payload, err := s.fromRequest(r)
		if err != nil {
			log.Println(err)
			return
		}

		response := new(Response)

		for k, v := range s.commands {
			if k == payload.Command {
				response = v.Execute(payload)
				break
			}
		}
		if response != nil {
			res, err := json.Marshal(&response)
			if err != nil {
				log.Println("Splat:", err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	for _, a := range s.actions {
		http.HandleFunc(endpoint+"/"+a.Endpoint, func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			log.Println(string(body))

			err = s.checkSignature(r, body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			payload := new(ActionPayload)

			err = json.Unmarshal(body, &payload)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

			a.Execute(payload)
		})
	}

	return http.ListenAndServe(addr, nil)
}

// Parses the body data into a program readable format
func (s *App) fromRequest(r *http.Request) (p *Payload, err error) {

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

	p = new(Payload)
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
func (s *App) checkSignature(r *http.Request, body []byte) (err error) {

	timestamp := r.Header.Get("X-Slack-Request-Timestamp")
	// TODO: Check timestamp age for possible replay attack

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
