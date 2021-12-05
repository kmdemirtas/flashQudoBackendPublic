package authentication

import (
	"context"
	"fmt"
	"log"
	"os"

	"firebase.google.com/go/auth"
	"github.com/ironstone95/FlashQudoV2/database"
	"github.com/ironstone95/FlashQudoV2/model"
)

type Authenticator struct {
	db         *database.Database
	l          *log.Logger
	authClient *auth.Client
	debugLog   *log.Logger
}

// NewAuthenticator creates a new Authenticator instance
func NewAuthenticator(l *log.Logger, db *database.Database, authClient *auth.Client, fullLog bool) *Authenticator {
	a := new(Authenticator)
	a.l = l
	a.db = db
	a.authClient = authClient
	if fullLog {
		a.debugLog = log.New(os.Stdout, "[Authenticator] ", 0)
	}
	return a
}

// AuthToken checks if given token is saved in the database, if not checks in Firebase
func (a *Authenticator) AuthToken(token string) error {
	log.Println() // to separate each request
	if len(token) == 0 {
		a.log("AuthToken", "Empty Token")
		return fmt.Errorf("token cannot be empty")
	}
	err := a.db.AuthToken(token)
	if err == nil {
		a.log("AuthToken", "SUCCESS")
		return nil
	}
	ctx := context.Background()
	authToken, err := a.authClient.VerifyIDToken(ctx, token)
	if err != nil {
		a.log("AuthToken", err.Error())
		return err
	}
	t := model.Token{Token: token, UserID: authToken.UID, Expires: authToken.Expires}
	_, err = a.db.InsertToken(t)
	if err != nil {
		a.log("AuthToken", err.Error())
		return err
	}
	a.log("AuthToken", "SUCCESS")
	return nil
}

// IsGroupMember checks if the passed token is the member of the group with id groupID
func (a *Authenticator) IsGroupMember(groupID, token string) bool {
	userID, err := a.db.GetIDFromToken(token)
	if err != nil {
		a.log("IsGroupMember", err.Error())
		return false
	}
	isMember, err := a.db.IsGroupMember(groupID, userID)
	if err != nil {
		a.log("IsGroupMember", err.Error())
		return false
	}
	if isMember {
		a.log("IsGroupMember", "SUCCESS")
	} else {
		a.log("IsGroupMember", "FAIL")
	}
	return isMember
}

// IsTheUser checks given token belongs to user with id userID
func (a *Authenticator) IsTheUser(userID, token string) bool {
	id, err := a.db.GetIDFromToken(token)
	if err != nil {
		a.log("IsTheUser", err.Error())
		return false
	}
	res := id == userID
	if res {
		a.log("IsTheUser", "SUCCESS")
	} else {
		a.log("IsTheUser", "FAIL")
	}
	return id == userID
}

// CanSeeBundle checks owner of the given token can see the bundle with id bundleID
func (a *Authenticator) CanSeeBundle(bundleID, token string) bool {
	userID, err := a.db.GetIDFromToken(token)
	if err != nil {
		a.log("CanSeeBundle", err.Error())
		return false
	}
	canSee, err := a.db.CanSeeBundle(bundleID, userID)
	if err != nil {
		a.log("CanSeeBundle", err.Error())
		return false
	}
	if canSee {
		a.log("CanSeeBundle", "SUCCESS")
	} else {
		a.log("CanSeeBundle", "FAIL")
	}
	return canSee
}

// TODO change it
// Default log implementation
func (a *Authenticator) log(prefix, msg string) {
	if a.debugLog != nil {
		a.debugLog.Printf("[%s] %s\n", prefix, msg)
	}
}
