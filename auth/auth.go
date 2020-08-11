package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nashl/online-store-server/graph/model"
	"log"
	"net/http"
	"time"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}
var writerCtxKey = &contextKey{"cookieWriter"}
type contextKey struct {
	name string
}

// HTTP is the struct used to inject the response writer and request http structs.
type HTTP struct {
	W *http.ResponseWriter
	R *http.Request
}


// A stand-in for our database backed user object
type User struct {
	UserID    string
	FullName  string
	Email     string
	//Name string
	//IsAdmin bool
}

type authResponseWriter struct {
	http.ResponseWriter
	userIDToResolver string
	userIDFromCookie string
}

func (w *authResponseWriter) Write(b []byte) (int, error) {
	if w.userIDToResolver != w.userIDFromCookie {
		fmt.Print("\n", "DIfferente", "\n")
		http.SetCookie(w, &http.Cookie{
			Name:     "gonlineToke",
			Value:    w.userIDToResolver,
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
			SameSite: http.SameSiteLaxMode,
		})
	}
	return w.ResponseWriter.Write(b)
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			httpContext := HTTP{
				W: &w,
				R: r,
			}

			ctx := context.WithValue(r.Context(),writerCtxKey, httpContext)
			r = r.WithContext(ctx)

			fmt.Print("\n r.Headers: ", r.Header, "\n")
			c, err := r.Cookie("onlineStore")


			// Allow unauthenticated users in
			if err != nil || c == nil {
				fmt.Print("\n", " $$$ unauthenticated user (NO COOKIE) $$$ ", "\n")
				next.ServeHTTP(w, r)
				return
			}
			//
			userId, err := validateAndGetUserID(c)
			if err != nil {
				fmt.Print("\n INVALID COOKIE!! \n")
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}
			//
			// get the user from the database
			user := getUserByID(db, userId)
			//
			fmt.Print("\nUser got from database", user, "\n")
			//
			// put it in context
			ctx = context.WithValue(ctx, userCtxKey, user)

			r = r.WithContext(ctx)
			fmt.Print("\n", "BEFORe SErVEHTtP", "\n")
			next.ServeHTTP(w, r)

		})
	}
}

func getUserByID(db *sql.DB, userId string) *model.User {
	var u model.User
	err := db.QueryRow("SELECT userId, email, fullName FROM `users` WHERE userId = ? ", userId).Scan(&u.UserID, &u.Email, &u.FullName)

	if err != nil {
		log.Fatal(err)
	}
	return &u
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *model.User {
	//fmt.Print("**ctx.value.User", ctx.Value(userCtxKey).(*model.User), "\n")
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	fmt.Print("USER FROM CONteXT: ", raw, "\n")
	return raw
}
// ForContext finds the HTTP Struct from the context. REQUIRES Middleware to have run.
func SaveAuthCookie(ctx context.Context) {
	fmt.Print("\nInsideSaveAUthCookie--0..\n")
	fmt.Print("ctxValue =>", ctx.Value(writerCtxKey), "\n")
	raw := ctx.Value(writerCtxKey).(HTTP)
	fmt.Print("--1--\n")
	authCookie := &http.Cookie{
		Name: "onlineStore",
		Value: "SoyUnTokenDifficultDeAdivinar",
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	fmt.Print("--1--\n")
	http.SetCookie(*(raw.W), authCookie)
	fmt.Print("--1--\n")
}

func ReadAuthCookie(ctx context.Context) string{
	fmt.Print("\nInsideReadAuthCookie--0..\n")
	raw := ctx.Value(writerCtxKey).(HTTP)
	fmt.Print("--R1--\n")
	c, err := raw.R.Cookie("onlineStore")
	fmt.Print("--R2--\n")
	if err != nil {
		log.Fatal("ReadAuthCookie ERROR", err)
	}
	fmt.Print("--R3--C:", c , "\n")
	//http.SetCookie(*(raw.W), authCookie)
	return c.Value
}

func RemoveAuthCookie(ctx context.Context) {
	fmt.Print("\nInsideRemoveAuthCookie--0..\n")
	expire := time.Now().Add(-7 * 24 * time.Hour)
	raw := ctx.Value(writerCtxKey).(HTTP)
	authCookie := &http.Cookie{
		Name: "onlineStore",
		Value: "",
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires: expire,
	}
	http.SetCookie(*(raw.W), authCookie)
}


func validateAndGetUserID(c *http.Cookie) (string, error){
	fmt.Print("\n COokie:", c, "\n")
	return "1", nil
}
