package authentication

import (
	"context"
	"net/http"

	"github.com/ironstone95/FlashQudoV2/handler"

	"github.com/gorilla/mux"
)

// AuthMW authenticates users with their tokens, uses database and firebase
func (a *Authenticator) AuthMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		if err := a.AuthToken(r.Header.Get("X-Auth-Token")); err != nil {
			a.log("AuthMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		} else {
			a.log("AuthMW", "SUCCESS")
			next.ServeHTTP(rw, r)
		}
	})
}

// AuthGroupMemberMW authorizes the user if he/she is member of the group
func (a *Authenticator) AuthGroupMemberMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := a.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
		if err != nil {
			a.log("AuthGroupMemberMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		groupID := mux.Vars(r)["groupID"]
		if len(groupID) == 0 {
			a.log("AuthGroupMemberMW", "groupID not in URL")
			handler.SendError(rw, "bad request", http.StatusBadRequest)
			return
		}
		isMember, err := a.db.IsGroupMember(groupID, userID)
		if err != nil {
			a.log("AuthGroupMemberMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		if !isMember {
			a.log("AuthGroupMemberMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		a.log("AuthGroupMemberMW", "SUCCESS")
		next.ServeHTTP(rw, r)
	})
}

// AuthBundleGroupMemberMW authorizes the user if he/she member of the group
func (a *Authenticator) AuthBundleGroupMemberMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := a.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
		if err != nil {
			a.log("AuthBundleGroupMemberMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		bundleID := mux.Vars(r)["bundleID"]
		canSeeBundle, err := a.db.CanSeeBundle(bundleID, userID)
		if err != nil {
			a.log("AuthBundleGroupMemberMW", err.Error())
			handler.SendError(rw, "bad request", http.StatusBadRequest)
			return
		} else if !canSeeBundle {
			a.log("AuthBundleGroupMemberMW", "user not authorized for this action")
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		a.log("AuthBundleGroupMemberMW", "SUCCESS")
		next.ServeHTTP(rw, r)
	})
}

// AuthGroupAdminMW authorizes the user if he/she the admin of the group param group
func (a *Authenticator) AuthGroupAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := a.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
		if err != nil {
			a.log("AuthGroupAdminMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		groupID := mux.Vars(r)["groupID"]
		if isAdmin, err := a.db.IsAdmin(groupID, userID); err != nil {
			a.log("AuthGroupAdminMW", err.Error())
			handler.SendError(rw, "bad request", http.StatusBadRequest)
			return
		} else if !isAdmin {
			a.log("AuthGroupAdminMW", "FAIL")
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		a.log("AuthGroupAdminMW", "SUCCESS")
		next.ServeHTTP(rw, r)
	})
}

// AuthBundleGroupAdminMW authorizes the user if he/she the admin of the group of the param bundle
func (a *Authenticator) AuthBundleGroupAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := a.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
		if err != nil {
			a.log("AuthBundleGroupAdminMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		bundleID := mux.Vars(r)["bundleID"]
		if canEditBundle, err := a.db.CanEditBundle(bundleID, userID); err != nil {
			a.log("AuthBundleGroupAdminMW", err.Error())
			handler.SendError(rw, "bad request", http.StatusBadRequest)
			return
		} else if !canEditBundle {
			a.log("AuthBundleGroupAdminMW", "FAIL")
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}
		a.log("AuthBundleGroupAdminMW", "SUCCESS")
		next.ServeHTTP(rw, r)
	})
}

// AuthCardGroupAdminMW authorizes the user if he/she is a member of the group of the param card
func (a *Authenticator) AuthCardGroupAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := a.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
		if err != nil {
			a.log("AuthCardGroupAdminMW", err.Error())
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}

		cardID := mux.Vars(r)["cardID"]
		if canEditCard, err := a.db.CanEditCard(cardID, userID); err != nil {
			a.log("AuthCardGroupAdminMW", err.Error())
			handler.SendError(rw, "bad request", http.StatusBadRequest)
			return
		} else if !canEditCard {
			a.log("AuthCardGroupAdminMW", "FAIL")
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}

		a.log("AuthCardGroupAdminMW", "SUCCESS")
		next.ServeHTTP(rw, r)
	})
}

// AuthUser authrozies if the user id of param username and header token belong each other.
func (a *Authenticator) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]
		userID, err := a.db.GetIDFromUsername(username)
		if err != nil {
			a.log("AuthUser getUserID", err.Error())
			handler.SendError(rw, "bad request", http.StatusBadRequest)
			return
		}

		tokenID, err := a.db.GetIDFromToken(r.Header.Get("X-Auth-Token"))
		if err != nil || userID != tokenID {
			if err != nil {
				a.log("AuthUser getUserTokenID", err.Error())
			} else {
				a.log("AuthUser getUserTokenID", "userID, tokenID not same")
			}
			handler.SendError(rw, "forbidden", http.StatusForbidden)
			return
		}

		ctxRaw := context.Background()
		ctx := context.WithValue(ctxRaw, "id", userID)
		r2 := r.Clone(ctx)
		next.ServeHTTP(rw, r2)
	})
}
