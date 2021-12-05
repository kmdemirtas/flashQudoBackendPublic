package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/ironstone95/FlashQudoV2/database"
)

type GetHandler struct {
	l        *log.Logger
	db       *database.Database
	debugLog *log.Logger
}

func NewGetHandler(l *log.Logger, db *database.Database, fullLog bool) *GetHandler {
	gh := new(GetHandler)
	gh.l = l
	gh.db = db
	if fullLog {
		gh.debugLog = log.New(os.Stdout, "[GetHandler] ", 0)
	}
	return gh
}

func (gh *GetHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	username, err := getParam("username", r)
	if err != nil {
		gh.log("GetUser", err.Error())
		SendError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO dont use map
	user, err := gh.db.GetUser(map[string]interface{}{"username": username})
	if err != nil {
		gh.log("GetUser", err.Error())
		SendError(rw, "cannot find", http.StatusNotFound)
		return
	}
	enc := json.NewEncoder(rw)
	if err := enc.Encode(user); err != nil {
		gh.log("GetUser", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	gh.log("GetUser", "SUCCESS")
}

func (gh *GetHandler) GetGroup(rw http.ResponseWriter, r *http.Request) {
	groupID, err := getParam("groupID", r)
	if err != nil {
		gh.log("GetGroup", err.Error())
		SendError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO dont use map
	group, err := gh.db.GetGroup(map[string]interface{}{"id": groupID})
	if err != nil {
		gh.log("GetGroup", err.Error())
		SendError(rw, "cannot find", http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(group); err != nil {
		gh.log("GetGroup", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	gh.log("GetGroup", "SUCCESS")
}

func (gh *GetHandler) GetUserGroups(rw http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	p := newPaging(r)
	groups, err := gh.db.GetUserGroups(userID, p.page, p.limit)
	if err != nil {
		gh.log("GetUserGroups dbGet", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(groups); err != nil {
		gh.log("GetUserGroups encode", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
}

func (gh *GetHandler) GetGroupUsers(rw http.ResponseWriter, r *http.Request) {
	groupID, err := getParam("groupID", r)
	if err != nil {
		gh.log("GetGroupUsers", err.Error())
		SendError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	p := newPaging(r)
	users, err := gh.db.GetGroupMembers(groupID, p.page, p.limit)
	if err != nil {
		gh.log("GetGroupUsers", err.Error())
		SendError(rw, "cannot find", http.StatusNotFound)
		return
	}
	enc := json.NewEncoder(rw)
	if err := enc.Encode(&users); err != nil {
		gh.log("GetGroupUsers", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	gh.log("GetGroupUsers", "SUCCESS")
}

func (gh *GetHandler) GetGroupBundles(rw http.ResponseWriter, r *http.Request) {
	groupID, err := getParam("groupID", r)
	if err != nil {
		gh.log("GetGroupBundles", err.Error())
		SendError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	p := newPaging(r)
	bundles, err := gh.db.GetGroupBundles(groupID, p.page, p.limit)
	if err != nil {
		gh.log("GetGroupBundles", err.Error())
		SendError(rw, "cannot find", http.StatusNotFound)
		return
	}
	enc := json.NewEncoder(rw)
	if err := enc.Encode(bundles); err != nil {
		gh.log("GetGroupBundles", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	gh.log("GetGroupBundles", "SUCCESS")
}

// TODO ONLY FOR MEMBERS!! TEST
func (gh *GetHandler) GetBundle(rw http.ResponseWriter, r *http.Request) {
	bundleID, err := getParam("bundleID", r)
	if err != nil {
		gh.log("GetBundle", err.Error())
		SendError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	bundle, err := gh.db.GetBundle(bundleID)
	if err != nil {
		gh.log("GetBundle", err.Error())
		SendError(rw, "cannot find", http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(bundle); err != nil {
		gh.log("GetBundle", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	gh.log("GetBundle", "SUCCESS")
}

// TODO ONLY FOR MEMBERS!! TEST
func (gh *GetHandler) GetBundleCards(rw http.ResponseWriter, r *http.Request) {

	bundleID, err := getParam("bundleID", r)
	if err != nil {
		gh.log("GetBundleCards", err.Error())
		SendError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	p := newPaging(r)
	cards, err := gh.db.GetBundleCards(bundleID, p.page, p.limit)
	if err != nil {
		gh.log("GetBundleCards", err.Error())
		SendError(rw, "cannot find", http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(&cards); err != nil {
		gh.log("GetBundleCards", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	gh.log("GetBundleCards", "SUCCESS")
}

func (gh *GetHandler) log(prefix, msg string) {
	if gh.debugLog != nil {
		gh.debugLog.Printf("[%s] %s\n", prefix, msg)
	}
}
