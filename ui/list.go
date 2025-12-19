package ui

import (
	"fmt"

	"tv-shows-manager/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	show models.Show
}

func (i item) Title() string {
	return fmt.Sprintf("%s - S%s (%s)", i.show.Show, i.show.Season, i.show.YearWatched)
}

func (i item) Description() string {
	return fmt.Sprintf("Source: %s | Rating: %s | Kind: %s", i.show.Source, i.show.SeasonRating, i.show.Kind)
}

func (i item) FilterValue() string {
	return i.show.Show
}

func showsToItems(shows []models.Show) []list.Item {
	items := make([]list.Item, len(shows))
	for i, show := range shows {
		items[i] = item{show: show}
	}
	return items
}
