package main

import (
	ctx "context"
	"encoding/json"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/raazcrzy/imdb/models"
	"github.com/raazcrzy/imdb/utils"
)

// addMovie function adds a new movie to the elasticsearch index
func addMovie(movie models.Movie) (map[string]interface{}, error) {
	_, err := utils.Elasticconn.Index().Index(utils.MovieIndex).Type("imdb").BodyJson(movie).Do(ctx.Background())
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  400,
		}, nil
	}
	return map[string]interface{}{
		"message": "movie added successfully",
		"status":  201,
	}, nil
}

// deleteMovie function deletes a movie from the elasticsearch index
func deleteMovie(movieID string) (map[string]interface{}, error) {
	_, err := utils.Elasticconn.Delete().Index(utils.MovieIndex).Type("imdb").Id(movieID).Do(ctx.Background())
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  400,
		}, nil
	}
	return map[string]interface{}{
		"message": "movie deleted successfully",
		"status":  200,
	}, nil
}

// editMovie function edits an existing movie
func editMovie(movie models.Movie) (map[string]interface{}, error) {
	_, err := utils.Elasticconn.Update().Index(utils.MovieIndex).Type("imdb").Id(movie.ID).Doc(movie).Do(ctx.Background())
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  400,
		}, nil
	}
	return map[string]interface{}{
		"message": "movie updated successfully",
		"status":  200,
	}, nil
}

// listMovies function queries the elasticsearch with appropriate query and fetches the list of movies matching the query
func listMovies(query elastic.Query, filters int, from, size int) (map[string]interface{}, error) {
	if filters == 0 {
		query = elastic.NewMatchAllQuery()
	}
	src, err := query.Source()
	if err != nil {
		Log.Errorln(err)
	}
	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		Log.Errorln(err)
	}
	Log.Infoln("here:", string(data))
	movies := []models.Movie{}
	response, err := utils.Elasticconn.Search().Index(utils.MovieIndex).Type("imdb").Query(query).From(from).Size(size).Do(ctx.Background())
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  400,
		}, nil
	}
	Log.Infof("%#v", response, response.Hits.TotalHits)
	for i := range response.Hits.Hits {
		movie := models.Movie{}
		Log.Infoln(string(*response.Hits.Hits[i].Source))
		json.Unmarshal(*response.Hits.Hits[i].Source, &movie)
		movie.ID = response.Hits.Hits[i].Id
		movies = append(movies, movie)

	}
	return map[string]interface{}{
		"message": "request successful",
		"movies":  movies,
		"status":  200,
	}, nil
}
