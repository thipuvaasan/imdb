package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/olivere/elastic.v5"

	"github.com/raazcrzy/imdb/models"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

/* addUserHandler handles the incoming requests to create a new user
The expected request body structure is:

type User struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	CreatedAt    int64  `json:"created_at"`
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

If role is specified as admin, the the request maker must be an existing admin/super admin.
*/
func addUserHandler(w http.ResponseWriter, r *http.Request) {
	var returnMsg map[string]interface{}
	var err error
	if r.Method != "POST" {
		returnMsg = map[string]interface{}{
			"message": "Invalid HTTP method, allowed POST",
			"status":  http.StatusBadRequest,
		}
		writeBack(w, returnMsg, err)
		return
	}
	email, reqCategory, ok, err := basicAuth(r)
	if err != nil {
		Log.Errorln(err)
		returnMsg = map[string]interface{}{
			"message": "Internal server error",
			"status":  http.StatusInternalServerError,
		}
		writeBack(w, returnMsg, err)
		return
	}
	if ok {
		if reqCategory != "users" {
			returnMsg := map[string]interface{}{
				"message": "Unauthorized",
				"status":  http.StatusUnauthorized,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	}
	d := json.NewDecoder(r.Body)
	body := models.User{}
	err = d.Decode(&body)
	if err != nil {
		Log.Errorln("decoding err: ", err)
		returnMsg = map[string]interface{}{
			"message": "Unable to decode request body",
			"status":  400,
		}
		writeBack(w, returnMsg, nil)
		return
	}

	// check if request make is admin before creating another admin user
	ok = (isAdmin(email) || isSuperAdmin(email) || body.Role != "admin")
	if body.Role == "" {
		body.Role = "user"
	}
	if ok {
		if body.Email == "" || body.UserName == "" || body.UserPassword == "" || body.Name == "" {
			returnMsg = map[string]interface{}{
				"message": "one or more fields missing in request body",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		if !emailRegexp.MatchString(body.Email) {
			returnMsg = map[string]interface{}{
				"message": "invalid email present in the request body",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		if len(body.UserName) > 32 || len(body.UserPassword) > 32 {
			returnMsg = map[string]interface{}{
				"message": "user_name and password has a max limit of 32 characters",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		if !(body.Role == "admin" || body.Role == "user") {
			returnMsg = map[string]interface{}{
				"message": "invalid role provided, valid roles: admin, user",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		body.CreatedAt = time.Now().Unix()
		returnMsg, err = createUser(body)
	} else {
		returnMsg = map[string]interface{}{
			"message": "Not Authorized",
			"status":  401,
		}
	}
	writeBack(w, returnMsg, err)
}

// removeUserHandler deletes an existing user from the database
func removeUserHandler(w http.ResponseWriter, r *http.Request) {
	var returnMsg map[string]interface{}
	var err error
	if r.Method != "DELETE" {
		returnMsg = map[string]interface{}{
			"message": "Invalid HTTP method, allowed DELETE",
			"status":  http.StatusBadRequest,
		}
		writeBack(w, returnMsg, err)
		return
	}
	email, reqCategory, ok, err := basicAuth(r)
	if err != nil {
		Log.Errorln(err)
		returnMsg = map[string]interface{}{
			"message": "Internal server error",
			"status":  http.StatusInternalServerError,
		}
		writeBack(w, returnMsg, err)
		return
	}
	if ok {
		if reqCategory != "users" {
			returnMsg := map[string]interface{}{
				"message": "Unauthorized",
				"status":  http.StatusUnauthorized,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	}
	userID, _, _ := r.BasicAuth()
	d := json.NewDecoder(r.Body)
	var body struct {
		Email string `json:"email"`
	}
	err = d.Decode(&body)
	if err != nil {
		Log.Errorln("decoding err: ", err)
		returnMsg = map[string]interface{}{
			"message": "Unable to decode request body",
			"status":  400,
		}
		writeBack(w, returnMsg, nil)
		return
	}
	if body.Email == "" {
		returnMsg = map[string]interface{}{
			"message": "one or more fields missing in request body",
			"status":  400,
		}
		writeBack(w, returnMsg, nil)
		return
	}
	if !emailRegexp.MatchString(body.Email) {
		returnMsg = map[string]interface{}{
			"message": "invalid email present in the request body",
			"status":  400,
		}
		writeBack(w, returnMsg, nil)
		return
	}
	ok = (isAdmin(email) || isSuperAdmin(email) || isAuthorizedUser(body.Email, userID))
	if ok {
		returnMsg, err = deleteUser(body.Email)
	} else {
		returnMsg = map[string]interface{}{
			"message": "Not Authorized",
			"status":  401,
		}
	}
	writeBack(w, returnMsg, err)
}

// addMovieHandler adds a new movie in the existing set of movies in elasticsearch
func addMovieHandler(w http.ResponseWriter, r *http.Request) {
	var returnMsg map[string]interface{}
	var err error
	if r.Method != "POST" {
		returnMsg = map[string]interface{}{
			"message": "Invalid HTTP method, allowed POST",
			"status":  http.StatusBadRequest,
		}
		writeBack(w, returnMsg, err)
		return
	}
	email, reqCategory, ok, err := basicAuth(r)
	if err != nil {
		Log.Errorln(err)
		returnMsg = map[string]interface{}{
			"message": "Internal server error",
			"status":  http.StatusInternalServerError,
		}
		writeBack(w, returnMsg, err)
		return
	}
	if ok {
		if reqCategory != "users" {
			returnMsg := map[string]interface{}{
				"message": "Unauthorized",
				"status":  http.StatusUnauthorized,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	}
	ok = (isAdmin(email) || isSuperAdmin(email))
	if ok {
		d := json.NewDecoder(r.Body)
		body := models.Movie{}
		err := d.Decode(&body)
		if err != nil {
			Log.Errorln("decoding err: ", err)
			returnMsg = map[string]interface{}{
				"message": "Unable to decode request body",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		if body.Name == "" || body.Director == "" || len(body.Genre) == 0 {
			returnMsg = map[string]interface{}{
				"message": "one or more fields missing in request body, required fields: name, director, genre",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		returnMsg, err = addMovie(body)
	} else {
		returnMsg = map[string]interface{}{
			"message": "Not Authorized",
			"status":  401,
		}
	}
	writeBack(w, returnMsg, err)
}

// removeMovieHandler deletes a movie from the existing set of movies
func removeMovieHandler(w http.ResponseWriter, r *http.Request) {
	var returnMsg map[string]interface{}
	var err error
	if r.Method != "DELETE" {
		returnMsg = map[string]interface{}{
			"message": "Invalid HTTP method, allowed DELETE",
			"status":  http.StatusBadRequest,
		}
		writeBack(w, returnMsg, err)
		return
	}
	email, reqCategory, ok, err := basicAuth(r)
	if err != nil {
		Log.Errorln(err)
		returnMsg = map[string]interface{}{
			"message": "Internal server error",
			"status":  http.StatusInternalServerError,
		}
		writeBack(w, returnMsg, err)
		return
	}
	if ok {
		if reqCategory != "users" {
			returnMsg := map[string]interface{}{
				"message": "Unauthorized",
				"status":  http.StatusUnauthorized,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	}
	ok = (isAdmin(email) || isSuperAdmin(email))
	if ok {
		movieID := r.URL.Query().Get("movie_id")
		if movieID == "" {
			returnMsg = map[string]interface{}{
				"message": "movie_id required as URL param",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		returnMsg, err = deleteMovie(movieID)
	} else {
		returnMsg = map[string]interface{}{
			"message": "Not Authorized",
			"status":  401,
		}
	}
	writeBack(w, returnMsg, err)
}

// updateMovieHandler updates a movie content
func updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var returnMsg map[string]interface{}
	var err error
	if r.Method != "PUT" {
		returnMsg = map[string]interface{}{
			"message": "Invalid HTTP method, allowed PUT",
			"status":  http.StatusBadRequest,
		}
		writeBack(w, returnMsg, err)
		return
	}
	email, reqCategory, ok, err := basicAuth(r)
	if err != nil {
		Log.Errorln(err)
		returnMsg = map[string]interface{}{
			"message": "Internal server error",
			"status":  http.StatusInternalServerError,
		}
		writeBack(w, returnMsg, err)
		return
	}
	if ok {
		if reqCategory != "users" {
			returnMsg := map[string]interface{}{
				"message": "Unauthorized",
				"status":  http.StatusUnauthorized,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	}
	ok = (isAdmin(email) || isSuperAdmin(email))
	if ok {
		d := json.NewDecoder(r.Body)
		body := models.Movie{}
		err := d.Decode(&body)
		if err != nil {
			Log.Errorln("decoding err: ", err)
			returnMsg = map[string]interface{}{
				"message": "Unable to decode request body",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		if body.Name == "" || body.Director == "" || len(body.Genre) == 0 || body.ID == "" {
			returnMsg = map[string]interface{}{
				"message": "one or more fields missing in request body, required fields: name, director, genre",
				"status":  400,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		returnMsg, err = editMovie(body)
	} else {
		returnMsg = map[string]interface{}{
			"message": "Not Authorized",
			"status":  401,
		}
	}
	writeBack(w, returnMsg, err)
}

// getMovieHandler fetches the list of movies matching the queries
func getMovieHandler(w http.ResponseWriter, r *http.Request) {
	var returnMsg map[string]interface{}
	var err error
	if r.Method != "GET" {
		returnMsg = map[string]interface{}{
			"message": "Invalid HTTP method, allowed GET",
			"status":  http.StatusBadRequest,
		}
		writeBack(w, returnMsg, err)
		return
	}
	user, reqCategory, ok, err := basicAuth(r)
	if err != nil {
		Log.Errorln(err)
		returnMsg = map[string]interface{}{
			"message": "Internal server error",
			"status":  http.StatusInternalServerError,
		}
		writeBack(w, returnMsg, err)
		return
	}
	Log.Infoln("user: ", user, reqCategory)
	if ok {
		if reqCategory != "users" {
			returnMsg := map[string]interface{}{
				"message": "Unauthorized",
				"status":  http.StatusUnauthorized,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	}
	searchQuery := elastic.NewBoolQuery()
	foundFilters := 0
	movieName := r.URL.Query().Get("name")
	if movieName != "" {
		searchQuery.Should(elastic.NewMatchQuery("name", fmt.Sprint(movieName)))
		foundFilters++
	}
	directorName := r.URL.Query().Get("director")
	if directorName != "" {
		searchQuery.Should(elastic.NewMatchQuery("director", fmt.Sprint(directorName)))
		foundFilters++
	}
	popularity := r.URL.Query().Get("99popularity")
	if popularity != "" {
		searchQuery.Should(elastic.NewMatchQuery("99popularity", popularity))
		foundFilters++
	}
	IMDBScore := r.URL.Query().Get("imdb_score")
	if IMDBScore != "" {
		searchQuery.Should(elastic.NewMatchQuery("imdb_score", IMDBScore))
		foundFilters++
	}
	genre := r.URL.Query().Get("genre")
	if genre != "" {
		searchQuery.Should(elastic.NewMatchPhraseQuery("genre", genre))
		foundFilters++
	}
	fromString := r.URL.Query().Get("from")
	var from, size int
	if fromString != "" {
		from, err = strconv.Atoi(fromString)
		if err != nil {
			Log.Errorln("Unable to parse from value: ", err)
			returnMsg := map[string]interface{}{
				"message": "from value must be an integer",
				"status":  http.StatusBadRequest,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		if from < 0 {
			Log.Errorln("Unable to parse size value: ", err)
			returnMsg := map[string]interface{}{
				"message": "from value must greater than -1",
				"status":  http.StatusBadRequest,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	} else {
		from = 0
	}
	sizeString := r.URL.Query().Get("size")
	if sizeString != "" {
		size, err = strconv.Atoi(sizeString)
		if err != nil {
			Log.Errorln("Unable to parse size value: ", err)
			returnMsg := map[string]interface{}{
				"message": "size value must be an integer",
				"status":  http.StatusBadRequest,
			}
			writeBack(w, returnMsg, nil)
			return
		}
		if size > 100 || size < 1 {
			Log.Errorln("Unable to parse size value: ", err)
			returnMsg := map[string]interface{}{
				"message": "size value must be greater than 0 and less than 100",
				"status":  http.StatusBadRequest,
			}
			writeBack(w, returnMsg, nil)
			return
		}
	} else {
		size = 20
	}

	src, err := searchQuery.Source()
	if err != nil {
		Log.Errorln(err)
	}
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		Log.Errorln(err)
	}
	Log.Infoln(string(data))
	searchQuery.MinimumNumberShouldMatch(foundFilters)
	returnMsg, err = listMovies(searchQuery, foundFilters, from, size)
	writeBack(w, returnMsg, err)
}

func writeBack(w http.ResponseWriter, returnMsg map[string]interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(500)
		e.Encode(err.Error())
	} else {
		statusCode, ok := returnMsg["status"]
		if ok {
			w.WriteHeader(statusCode.(int))
			delete(returnMsg, "status")
			e.Encode(returnMsg)
		} else {
			e.Encode(returnMsg)
		}
	}
}

// basicAuth fetches email, user category from the request
func basicAuth(r *http.Request) (string, string, bool, error) {
	_, _, ok := r.BasicAuth()
	if !ok {
		return "", "", false, nil
	}

	ctxEmail := r.Context().Value("email")
	if ctxEmail == nil {
		return "", "", false, fmt.Errorf("cannot fetch email from request context")
	}
	email, ok := ctxEmail.(string)
	if !ok {
		return "", "", false, fmt.Errorf("cannot cast context email %v to string", ctxEmail)
	}

	ctxCredentials := r.Context().Value("category")
	if ctxCredentials == nil {
		return "", "", false, fmt.Errorf("cannot fetch context credentials from request context")
	}
	credentials, ok := ctxCredentials.(string)
	if !ok {
		return "", "", false, fmt.Errorf("cannot cast context credentials %v to string", ctxCredentials)
	}

	return email, credentials, true, nil
}
