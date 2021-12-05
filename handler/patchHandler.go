package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ironstone95/FlashQudoV2/database"
	"github.com/ironstone95/FlashQudoV2/handler/request"
	"log"
	"net/http"
	"os"
)

type PatchHandler struct {
	l        *log.Logger
	db       *database.Database
	debugLog *log.Logger
}

func NewPatchHandler(l *log.Logger, db *database.Database, fullLog bool) *PatchHandler {
	ph := new(PatchHandler)
	ph.l = l
	ph.db = db
	if fullLog {
		ph.debugLog = log.New(os.Stdout, "[GetHandler] ", 0)
	}
	return ph
}

func (ph *PatchHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	upr := request.UserPatchRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&upr); err != nil {
		ph.log("PatchUser decode", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}
	pv, err := upr.GetPatchValues()
	if err != nil {
		ph.log("PatchUser getPatchValues", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	userID, err := ph.db.GetIDFromUsername(mux.Vars(r)["username"])
	if err != nil {
		ph.log("PatchUser getIDFromUsername", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}
	tokenUserID, err := ph.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
	if err != nil || userID != tokenUserID {
		if err != nil {
			ph.log("PatchUser readToken", err.Error())
		} else {
			ph.log("PatchUser readToken", "token belongs to another user")
		}
		SendError(rw, "forbidden", http.StatusForbidden)
		return
	}

	dbUser, err := ph.db.UpdateUser(userID, pv)
	if err != nil {
		ph.log("PatchUser updateDB", err.Error())
		SendError(rw, "update error", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbUser); err != nil {
		ph.log("PatchUser encode", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

	ph.log("PatchUser", "SUCCESS")
}

func (ph *PatchHandler) PatchGroup(rw http.ResponseWriter, r *http.Request) {
	gpr := request.GroupPatchRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&gpr); err != nil {
		ph.log("PatchGroup decode", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	pv, err := gpr.GetPatchValues()
	if err != nil {
		ph.log("PatchGroup getPatchValues", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	dbGroup, err := ph.db.UpdateGroup(mux.Vars(r)["groupID"], pv)
	if err != nil {
		ph.log("PatchGroup updateDB", err.Error())
		SendError(rw, "update error", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbGroup); err != nil {
		ph.log("PatchGroup encode", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
	ph.log("PatchGroup", "SUCCESS")
}

func (ph *PatchHandler) PatchBundle(rw http.ResponseWriter, r *http.Request) {
	bpr := request.BundlePatchRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&bpr); err != nil {
		ph.log("PatchBundle decode", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	pv, err := bpr.GetPatchValues()
	if err != nil {
		ph.log("PatchBundle getPatchValues", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}
	dbBundle, err := ph.db.UpdateBundle(mux.Vars(r)["bundleID"], pv)
	if err != nil {
		ph.log("PatchBundle updateDB", err.Error())
		SendError(rw, "update error", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbBundle); err != nil {
		ph.log("PatchBundle encode", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
	ph.log("PatchBundle", "SUCCESS")
}

func (ph *PatchHandler) PatchCard(rw http.ResponseWriter, r *http.Request) {
	cpr := request.CardPatchRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&cpr); err != nil {
		ph.log("PatchCard decode", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	pv, err := cpr.GetPatchValues()
	if err != nil {
		ph.log("PatchCard getPatchValues", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	dbCard, err := ph.db.UpdateCard(mux.Vars(r)["cardID"], pv)
	if err != nil {
		ph.log("PatchCard updateDB", err.Error())
		SendError(rw, "update error", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbCard); err != nil {
		ph.log("PatchCard encode", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}

}

func (ph *PatchHandler) PatchMember(rw http.ResponseWriter, r *http.Request) {
	mpr := request.MemberPatchRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&mpr); err != nil {
		ph.log("PatchMember decode", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	}

	pv, err := mpr.GetPatchValues()
	if err != nil {
		ph.log("PatchMember getPatchValues", err.Error())
		SendError(rw, "missing field", http.StatusBadRequest)
		return
	}

	requesterID, err := ph.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
	if err != nil {
		ph.log("PatchMember Auth", err.Error())
		SendError(rw, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	if isAdmin, err := ph.db.IsAdmin(vars["groupID"], vars["userID"]); err != nil {
		ph.log("PatchMember isAdmin", err.Error())
		SendError(rw, "bad request", http.StatusBadRequest)
		return
	} else if isAdmin && vars["userID"] != requesterID {
		ph.log("PatchMember requesterCheck", "Requester is not same with admin")
		SendError(rw, "an admin cannot change other admins status", http.StatusForbidden)
		return
	}

	dbMember, err := ph.db.UpdateMember(vars["groupID"], vars["userID"], pv)
	if err != nil {
		ph.log("PatchMember updateDB", err.Error())
		SendError(rw, "update error", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dbMember); err != nil {
		ph.log("PatchMember encode", err.Error())
		SendError(rw, "server error", http.StatusInternalServerError)
		return
	}
}

func (ph *PatchHandler) log(prefix, msg string) {
	if ph.debugLog != nil {
		ph.debugLog.Printf("[%s] %s\n", prefix, msg)
	}
}
