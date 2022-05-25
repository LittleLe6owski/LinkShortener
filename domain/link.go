package domain

import (
	"strconv"
	"strings"
)

type Link struct {
	ExpiringDate string `json:"creation_date"`
	FullVersion  string `json:"full_version"`
	ShortVersion string `json:"short_version"`
}

func (l *Link) ConvertDateToListNumber() []int {
	date := make([]int, 2)
	for i, v := range strings.Split(l.ExpiringDate, ".") {
		number, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		date[i] = number
	}
	return date
}

type LinkUseCase interface {
	GenShortLink(string) (string, error)
	GetOriginalLink(string) (string, error)
	InitRestoreLinks()
}

type LinkPgRepo interface {
	GetLink(string) (Link, error)
	SaveLink(Link) error
}

type LinkSSDBRepo interface {
	GetLink(string, string) (Link, error)
	DelLink(string, string) error
	GetMaps(...string) ([]string, error)
	GetKeys(mName string) ([]string, error)
	SaveLink(Link) error
}
