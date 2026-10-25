package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(app.cfg.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	loadAndSaveSession := func(fn func(http.ResponseWriter, *http.Request)) http.Handler {
		return app.sessionManager.LoadAndSave(http.HandlerFunc(fn))
	}

	mux.Handle("GET /{$}", loadAndSaveSession(app.home))
	mux.Handle("GET /snippet/view/{id}", loadAndSaveSession(app.snippetView))
	mux.Handle("GET /snippet/create", loadAndSaveSession(app.snippetCreate))
	mux.Handle("POST /snippet/create", loadAndSaveSession(app.snippetCreatePost))

	mux.Handle("GET /user/signup", loadAndSaveSession(app.userSignup))
	mux.Handle("POST /user/signup", loadAndSaveSession(app.userSignupPost))
	mux.Handle("GET /user/login", loadAndSaveSession(app.userLogin))
	mux.Handle("POST /user/login", loadAndSaveSession(app.userLoginPost))
	mux.Handle("POST /user/logout", loadAndSaveSession(app.userLogoutPost))

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
