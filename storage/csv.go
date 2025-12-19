package storage

import (
	"encoding/csv"
	"os"

	"tv-shows-manager/models"
)

type CSVStorage struct {
	filename string
}

func NewCSVStorage(filename string) *CSVStorage {
	return &CSVStorage{filename: filename}
}

func (s *CSVStorage) Load() ([]models.Show, error) {
	file, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []models.Show{}, nil
	}

	shows := make([]models.Show, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		if len(records[i]) < 7 {
			continue
		}
		shows = append(shows, models.Show{
			Show:         records[i][0],
			Season:       records[i][1],
			YearWatched:  records[i][2],
			Source:       records[i][3],
			TmdbID:       records[i][4],
			Kind:         records[i][5],
			SeasonRating: records[i][6],
		})
	}

	return shows, nil
}

func (s *CSVStorage) Save(shows []models.Show) error {
	file, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"show", "season", "yearWatched", "source", "tmdb_id", "kind", "season_rating"}); err != nil {
		return err
	}

	// Write records
	for _, show := range shows {
		record := []string{
			show.Show,
			show.Season,
			show.YearWatched,
			show.Source,
			show.TmdbID,
			show.Kind,
			show.SeasonRating,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}