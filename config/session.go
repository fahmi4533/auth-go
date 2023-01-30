package config

import "github.com/gorilla/sessions"

const SEESSION_ID = "go_auth_sess"

var Store = sessions.NewCookieStore([]byte("adkaksdwenfwkefowk"))
