package repository

import (
	"MusicLibrary_Test/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type SongRepository struct {
	DB *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) *SongRepository {
	return &SongRepository{
		DB: db,
	}
}

func (r *SongRepository) GetSongs(ctx context.Context, filters map[string]string, page int, limit int) ([]models.Song, error) {
	var songs []models.Song
	offset := (page - 1) * limit
	query := `SELECT * FROM songs WHERE 1=1`
	var args []interface{}
	argID := 1

	for key, value := range filters {
		query += `AND` + key + ` ILIKE '%' || $` + fmt.Sprint(argID) + ` || '%'`
		args = append(args, value)
		argID++
	}

	query += ` LIMIT $` + fmt.Sprint(argID) + ` OFFSET $` + fmt.Sprint(argID+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.Text, &song.DateRelease, &song.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (r *SongRepository) GetSongText(ctx context.Context, id int, page int, limit int) (string, error) {
	var text string
	err := r.DB.QueryRow(ctx, `SELECT text FROM songs WHERE id=$1`, id).Scan(&text)
	if err != nil {
		return "", err
	}
	verses := splitTextIntoVerses(text)
	start := (page - 1) * limit
	end := start + limit
	if start > len(verses) {
		return "", err
	}
	if end > len(verses) {
		end = len(verses)
	}
	return joinVerses(verses[start:end]), nil
}

func (r *SongRepository) AddSong(ctx context.Context, song models.Song) error {
	_, err := r.DB.Exec(ctx, `INSERT INTO songs (group_name, song, text, release_date, link) VALUES ($1, $2, $3, $4, $5)`,
		song.GroupName, song.SongName, song.Text, song.DateRelease, song.Link)
	return err
}

func (r *SongRepository) UpdateSong(ctx context.Context, id int, song models.Song) error {
	_, err := r.DB.Exec(ctx, `UPDATE songs SET group_name=$1, song=$2, text=$3, release_date=$4, link=$5 WHERE id=$6`,
		song.GroupName, song.SongName, song.Text, song.DateRelease, song.Link, song.ID)
	return err
}

func (r *SongRepository) DeleteSong(ctx context.Context, id int) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM songs WHERE id=$1`, id)
	return err
}

func splitTextIntoVerses(text string) []string {
	return strings.Split(text, "\n\n")
}

func joinVerses(verses []string) string {
	return strings.Join(verses, "\n\n")
}
