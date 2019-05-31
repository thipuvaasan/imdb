package models

type User struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	CreatedAt    int64  `json:"created_at"`
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

type Movie struct {
	ID         string   `json:"movie_id,omitempty"`
	Name       string   `json:"name"`
	Popularity float32  `json:"99popularity"`
	Director   string   `json:"director"`
	Genre      []string `json:"genre"`
	IMDBScore  float32  `json:"imdb_score"`
}
