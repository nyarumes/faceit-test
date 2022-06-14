package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// HttpRouter declares http routes
func HttpRouter(handler UniqueHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		var (
			offset   = r.URL.Query().Get("offset")
			count    = r.URL.Query().Get("count")
			nickname = r.URL.Query().Get("nickname")
		)

		if len(offset) == 0 {
			offset = "0"
		}

		if len(count) == 0 {
			count = "10"
		}

		offsetParse, err := strconv.Atoi(offset)
		if err != nil {
			log.Error().Err(err).Msg("can't parse offset value from request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		countParse, err := strconv.Atoi(count)
		if err != nil {
			log.Error().Err(err).Msg("can't parse count value from request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		users, err := handler.FindAllWithPaginationAndFilter(offsetParse, countParse, map[string]string{
			FilterUserNickname: nickname,
		})
		if err != nil {
			log.Error().Err(err).Msg("can't find users")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(users)
		if err != nil {
			log.Error().Err(err).Msg("can't marshal users response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}).Methods("GET")

	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		var request Request

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Error().Err(err).Msg("can't parse request payload")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.Add(request.FirstName, request.LastName, request.Nickname, request.Password, request.Email, request.Country)
		if err != nil {
			log.Error().Err(err).Msg("can't add a new user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	r.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		var (
			request Request
			vars    = mux.Vars(r)
		)

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Error().Err(err).Msg("can't parse request payload")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(vars["id"])
		if err != nil {
			log.Error().Err(err).Msg("can't parse id as uuid")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.Update(id, request.FirstName, request.LastName, request.Nickname, request.Country)
		if err != nil {
			log.Error().Err(err).Msg("can't update existing user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("PUT")

	r.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		var vars = mux.Vars(r)

		id, err := uuid.Parse(vars["id"])
		if err != nil {
			log.Error().Err(err).Msg("can't parse id as uuid")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.Remove(id)
		if err != nil {
			log.Error().Err(err).Msg("can't update existing user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	log.Debug().Msg("User http server is starting")

	return r
}
