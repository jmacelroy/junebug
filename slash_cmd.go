package junebug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	jitsi "github.com/jitsi/jitsi-slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/rs/zerolog/hlog"
)

// JitsiTokenGenerator provides an interface for creating Jitsi Meet
// video conference authenticated access via JWT.
type JitsiTokenGenerator interface {
	CreateJWT(tenantID, tenantName, roomClaim, userID, userName, avatarURL string) (string, error)
}

// Slash handles /junebug Slack commands.
func Slash(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		hlog.FromRequest(r).Error().
			Err(err).
			Msg("unable to parse form data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	callbackID := uuid.New().String()
	fmt.Fprintf(w, prompt, callbackID, callbackID, callbackID)
}

// Meeting is the configuration for a meeting.
type Meeting struct {
	ConferenceType string
	TrackingType   string
}

// InteractionStateStore stores state for meeting setup interactions.
type InteractionStateStore struct {
	state map[string]*Meeting
	mux   sync.Mutex
}

// NewInteractionStateStore provides a clean instance of the
// state store.
func NewInteractionStateStore() *InteractionStateStore {
	return &InteractionStateStore{state: make(map[string]*Meeting)}
}

func (x *InteractionStateStore) startMeetingMsg(callbackID string) string {
	var (
		meetingURL = fmt.Sprintf(
			"https://meet.jit.si/atlassian/%s",
			jitsi.RandomName(),
		)
		meeting *Meeting
		ok      bool
		msg     string
	)

	if meeting, ok = x.state[callbackID]; !ok {
		// user hit start w/o making a selection so default to a Jitsi meeting ;-)
		return fmt.Sprintf("Meeting starting in %s\n", meetingURL)
	}

	if meeting.ConferenceType == "jitsi" {
		msg = fmt.Sprintf("Meeting starting in %s\n", meetingURL)
	}
	if meeting.TrackingType == "confluence" {
		msg = msg + fmt.Sprintf("Take notes with %s", "https://www.google.com")
	}
	fmt.Printf("\n%s\n", msg)
	return msg
}

// SlashInteraction is
func (x *InteractionStateStore) SlashInteraction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		hlog.FromRequest(r).Error().
			Err(err).
			Msg("unable to parse form data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(r.PostFormValue("payload"))

	var payload slackevents.MessageAction
	err = json.NewDecoder(strings.NewReader(r.PostFormValue("payload"))).Decode(&payload)
	if err != nil {
		hlog.FromRequest(r).Error().
			Err(err).
			Msg("parsing interaction")
	}

	if len(payload.Actions) != 1 {
		panic("too many actions")
	}

	x.mux.Lock()
	switch payload.Actions[0].Name {
	case "video":
		if meeting, ok := x.state[payload.CallbackId]; ok {
			meeting.ConferenceType = payload.Actions[0].Value
		} else {
			x.state[payload.CallbackId] = &Meeting{
				ConferenceType: payload.Actions[0].Value,
			}
		}
		w.WriteHeader(http.StatusOK)
		x.mux.Unlock()
		return
	case "tracking":
		if meeting, ok := x.state[payload.CallbackId]; ok {
			meeting.TrackingType = payload.Actions[0].Value
		} else {
			x.state[payload.CallbackId] = &Meeting{
				TrackingType: payload.Actions[0].Value,
			}
		}
		w.WriteHeader(http.StatusOK)
		x.mux.Unlock()
		return
	case "start":
		fmt.Print("\n start \n")
		msg := x.startMeetingMsg(payload.CallbackId)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, msg)
		x.mux.Unlock()
		return
	default:
		hlog.FromRequest(r).Error().
			Msg("unexpected action name")
		w.WriteHeader(http.StatusInternalServerError)
		x.mux.Unlock()
		return
	}
}
