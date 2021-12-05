package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/ironstone95/FlashQudoV2/authentication"
	"github.com/ironstone95/FlashQudoV2/database"
	"github.com/ironstone95/FlashQudoV2/handler"
)

func main() {
	l := log.New(os.Stdout, "BACKEND ", log.Default().Flags())
	l.SetFlags(log.LstdFlags | log.Lshortfile)
	db := database.ConnectDB(l, true, true, true, true)
	ctxApp := context.Background()
	app, err := firebase.NewApp(ctxApp, nil)
	if err != nil {
		log.Fatal(err)
	}
	ctxClient := context.Background()
	authClient, err := app.Auth(ctxClient)
	if err != nil {
		log.Fatal(err)
	}
	auth := authentication.NewAuthenticator(l, db, authClient, true)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "COMING SOON")
	})

	// GET
	gr := router.Methods(http.MethodGet).Subrouter()
	gr.Use(auth.AuthMW)
	gh := handler.NewGetHandler(l, db, true)

	// Authenticated Access
	gr.HandleFunc("/users/{username}", gh.GetUser)
	gr.HandleFunc("/groups/{groupID}", gh.GetGroup)

	// Authorized Access
	gr.Handle("/groups/{groupID}/users", auth.AuthGroupMemberMW(http.HandlerFunc(gh.GetGroupUsers)))
	gr.Handle("/groups/{groupID}/bundles", auth.AuthGroupMemberMW(http.HandlerFunc(gh.GetGroupBundles)))
	gr.Handle("/bundles/{bundleID}", auth.AuthBundleGroupMemberMW(http.HandlerFunc(gh.GetBundle)))
	gr.Handle("/bundles/{bundleID}/cards", auth.AuthBundleGroupMemberMW(http.HandlerFunc(gh.GetBundleCards)))
	gr.Handle("/users/{username}/groups", auth.AuthUser(http.HandlerFunc(gh.GetUserGroups)))

	// POST
	pr := router.Methods(http.MethodPost).Subrouter()
	pr.Use(auth.AuthMW)
	ph := handler.NewPostHandler(l, db, true)

	// Authenticated Access
	pr.HandleFunc("/groups", ph.InsertGroup)
	pr.HandleFunc("/users", ph.InsertUser)

	// Authorized Access
	pr.Handle("/groups/{groupID}/bundles", auth.AuthGroupAdminMW(http.HandlerFunc(ph.InsertBundle)))
	pr.Handle("/bundles/{bundleID}/cards", auth.AuthBundleGroupAdminMW(http.HandlerFunc(ph.InsertCard)))
	pr.Handle("/groups/{groupID}/users", auth.AuthGroupAdminMW(http.HandlerFunc(ph.InsertMember)))

	// Patch
	paR := router.Methods(http.MethodPatch).Subrouter()
	paR.Use(auth.AuthMW)
	paH := handler.NewPatchHandler(l, db, true)

	// Authenticated Access
	paR.HandleFunc("/users/{username}", paH.PatchUser)

	// Authorized Access
	paR.Handle("/bundles/{bundleID}", auth.AuthBundleGroupAdminMW(http.HandlerFunc(paH.PatchBundle)))
	paR.Handle("/cards/{cardID}", auth.AuthCardGroupAdminMW(http.HandlerFunc(paH.PatchCard)))
	paR.Handle("/groups/{groupID}", auth.AuthGroupAdminMW(http.HandlerFunc(paH.PatchGroup)))
	paR.Handle("/groups/{groupID}/users/{userID}", auth.AuthGroupAdminMW(http.HandlerFunc(paH.PatchMember)))

	// Delete
	dr := router.Methods(http.MethodDelete).Subrouter()
	dh := handler.NewDeleteHandler(l, db, true)

	// Authorized Access
	dr.Handle("/groups/{groupID}", auth.AuthGroupAdminMW(http.HandlerFunc(dh.DeleteGroup)))
	dr.Handle("/bundles/{bundleID}", auth.AuthBundleGroupAdminMW(http.HandlerFunc(dh.DeleteBundle)))
	dr.Handle("/bundles/{bundleID}/cards", auth.AuthBundleGroupAdminMW(http.HandlerFunc(dh.DeleteBundleCards)))
	dr.Handle("/cards/{cardID}", auth.AuthCardGroupAdminMW(http.HandlerFunc(dh.DeleteCard)))
	dr.Handle("/groups/{groupID}/users/{userID}", auth.AuthGroupAdminMW(http.HandlerFunc(dh.DeleteMember)))

	server := http.Server{
		Addr:     ":5000",
		ErrorLog: l,
		Handler:  router,
	}

	err = server.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}
