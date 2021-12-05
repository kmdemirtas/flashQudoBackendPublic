package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ironstone95/FlashQudoV2/database"
	"github.com/ironstone95/FlashQudoV2/handler/request"
)

type PostHandler struct {
	l        *log.Logger
	db       *database.Database
	debugLog *log.Logger
}

func NewPostHandler(l *log.Logger, db *database.Database, fullLog bool) *PostHandler {
	ph := new(PostHandler)
	ph.l = l
	ph.db = db
	if fullLog {
		ph.debugLog = log.New(os.Stdout, "[PostHandler] ", 0)
	}
	return ph
}

func (ph *PostHandler) InsertBundle(rw http.ResponseWriter, r *http.Request) {
	bpr := request.BundlePostRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&bpr); err != nil {
		ph.log("InsertBundle", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	b, err := bpr.CreateBundle()
	if err != nil {
		ph.log("InsertBundle", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	b.GroupID = mux.Vars(r)["groupID"]

	dbBundle, err := ph.db.InsertBundle(b)
	if err != nil {
		ph.log("InsertBundle", err.Error())
		SendError(rw, "insertion failed", http.StatusBadRequest)
		return
	}
	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbBundle); err != nil {
		ph.log("InsertBundle", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	ph.log("InsertBundle", "SUCCESS")
}

func (ph *PostHandler) InsertCard(rw http.ResponseWriter, r *http.Request) {
	cpr := request.CardPostRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&cpr); err != nil {
		ph.log("InsertCard", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	c, err := cpr.CreateCard()
	if err != nil {
		ph.log("InsertCard", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	c.BundleID = mux.Vars(r)["bundleID"]

	dbCard, err := ph.db.InsertCard(c)
	if err != nil {
		ph.log("InsertCard", err.Error())
		SendError(rw, "insertion failed", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbCard); err != nil {
		ph.log("InsertCard", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	ph.log("InsertCard", "SUCCESS")
}

func (ph *PostHandler) InsertGroup(rw http.ResponseWriter, r *http.Request) {
	gpr := request.GroupPostRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&gpr); err != nil {
		ph.log("InsertGroup", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	g, err := gpr.CreateGroup()
	if err != nil {
		ph.log("InsertGroup", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	userID, err := ph.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
	if err != nil {
		ph.log("InsertGroup", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	dbGroup, err := ph.db.InsertGroup(g, userID)
	if err != nil {
		ph.log("InsertGroup", err.Error())
		SendError(rw, "insertion failed", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbGroup); err != nil {
		ph.log("InsertGroup", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	ph.log("InsertGroup", "SUCCESS")
}

func (ph *PostHandler) InsertMember(rw http.ResponseWriter, r *http.Request) {
	mpr := request.MemberPostRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&mpr); err != nil {
		ph.log("InsertMember decode", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	userID, err := ph.db.GetIDFromUsername(*mpr.Username)
	if err != nil {
		ph.log("InsertMember getIDFromUsername", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}
	m, err := mpr.CreateMember(userID)
	if err != nil {
		ph.log("InsertMember mpr.CreateMember", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	dbMember, err := ph.db.InsertMember(m)
	if err != nil {
		ph.log("InsertMember DB insertion", err.Error())
		SendError(rw, "insertion failed", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbMember); err != nil {
		ph.log("InsertMember encoding", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	ph.log("InsertMember", "SUCCESS")
}

func (ph *PostHandler) InsertUser(rw http.ResponseWriter, r *http.Request) {
	urp := request.UserPostRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&urp); err != nil {
		ph.log("InsertUser", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	u, err := urp.CreateUser()
	if err != nil {
		ph.log("InsertUser", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	userID, err := ph.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
	if err != nil {
		SendError(rw, "forbidden", http.StatusForbidden)
		return
	}

	u.ID = userID
	dbUser, err := ph.db.InsertUser(u)

	if err != nil {
		ph.log("InsertUser", err.Error())
		SendError(rw, "insertion failed", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbUser); err != nil {
		ph.log("InsertUser", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	ph.log("InsertUser", "SUCCESS")
}

func (ph *PostHandler) log(prefix, msg string) {
	if ph.debugLog != nil {
		ph.debugLog.Printf("[%s] %s\n", prefix, msg)
	}
}
