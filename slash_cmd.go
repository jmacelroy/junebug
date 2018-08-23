package junebug

import (
	"fmt"
	"net/http"

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

	// room := jitsi.RandomName()
	// teamName := r.PostFormValue("team_domain")
	// if teamName == "" {
	// 	http.Error(w, "expected team domain in request form", http.StatusBadRequest)
	// 	return
	// }
	// meetingURL := fmt.Sprintf(
	// 	"https://meet.jit.si/%s/%s",
	// 	strings.ToLower(teamName),
	// 	room,
	// )
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, prompt)
}
