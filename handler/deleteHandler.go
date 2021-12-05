package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ironstone95/FlashQudoV2/database"
)

type DeleteHandler struct {
	l        *log.Logger
	db       *database.Database
	debugLog *log.Logger
}

func NewDeleteHandler(l *log.Logger, db *database.Database, fullLog bool) *DeleteHandler {
	dh := new(DeleteHandler)
	dh.l = l
	dh.db = db
	if fullLog {
		dh.debugLog = log.New(os.Stdout, "[GetHandler] ", 0)
	}
	return dh
}

func (dh *DeleteHandler) DeleteGroup(rw http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["groupID"]
	if err := dh.db.DeleteGroup(groupID); err != nil {
		dh.log("DeleteGroup delete", err.Error())
		if errors.Is(err, database.ErrMembersExists) {
			dh.log("DeleteGroup database delete", err.Error())
			SendError(rw, "cannot delete group with members", http.StatusBadRequest)

		} else {
			SendError(rw, "bad request", http.StatusBadRequest)
		}
		return
	}

	_, err := fmt.Fprint(rw, "group deleted")
	if err != nil {
		dh.log("DeleteGroup response", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
}

func (dh *DeleteHandler) DeleteBundle(rw http.ResponseWriter, r *http.Request) {
	bundleID := mux.Vars(r)["bundleID"]
	if err := dh.db.DeleteBundle(bundleID); err != nil {
		dh.log("DeleteBundle delete", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	_, err := fmt.Fprint(rw, "bundle delete")
	if err != nil {
		dh.log("DeleteBundle response", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
}

func (dh *DeleteHandler) DeleteBundleCards(rw http.ResponseWriter, r *http.Request) {
	bundleID := mux.Vars(r)["bundleID"]
	if err := dh.db.DeleteBundleCards(bundleID); err != nil {
		dh.log("DeleteBundleCards delete", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	_, err := fmt.Fprint(rw, "bundle cards delete")
	if err != nil {
		dh.log("DeleteBundleCards response", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
}

func (dh *DeleteHandler) DeleteCard(rw http.ResponseWriter, r *http.Request) {
	cardID := mux.Vars(r)["cardID"]
	if err := dh.db.DeleteCard(cardID); err != nil {
		dh.log("DeleteCard delete", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	_, err := fmt.Fprint(rw, "card delete")
	if err != nil {
		dh.log("DeleteCard response", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
}

func (dh *DeleteHandler) DeleteMember(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if isAdmin, err := dh.db.IsAdmin(vars["groupID"], vars["userID"]); err != nil {
		dh.log("DeleteMember admin check", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	} else if isAdmin {
		requesterID, _ := dh.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
		if requesterID != vars["userID"] {
			dh.log("DeleteMember adminCheck", "admin is not deleted member")
			SendError(rw, "an admin cannot change other admins status", http.StatusForbidden)
			return
		}
	}

	if err := dh.db.DeleteMember(vars["groupID"], vars["userID"]); err != nil {
		dh.log("DeleteMember dbDelete", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	_, err := fmt.Fprint(rw, "Member deleted")
	if err != nil {
		dh.log("Delete Member Response", err.Error())
	}
}

func (dh *DeleteHandler) log(prefix, msg string) {
	if dh.debugLog != nil {
		dh.debugLog.Printf("[%s] %s\n", prefix, msg)
	}
}
