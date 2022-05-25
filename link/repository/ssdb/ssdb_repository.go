package ssdb

import (
	"github.com/LittleLe6owski/LinkShortener/domain"
	"github.com/seefan/gossdb/v2/pool"
)

type LinkSSDBRepo struct {
	repo   domain.LinkSSDBRepo
	client *pool.Client
}

func NewLinkRepository(client *pool.Client) domain.LinkSSDBRepo {
	return &LinkSSDBRepo{client: client}
}

func (l LinkSSDBRepo) GetLink(shortLink, mapName string) (domain.Link, error) {
	link, err := l.client.HGet(mapName, shortLink)
	if err != nil {
		return domain.Link{}, err
	}
	return domain.Link{FullVersion: link.String(), ShortVersion: shortLink, ExpiringDate: mapName}, nil
}

func (l LinkSSDBRepo) DelLink(shortLink, mapName string) error {
	return l.client.HDel(mapName, shortLink)
}

func (l LinkSSDBRepo) SaveLink(link domain.Link) error {
	return l.client.HSet(link.ExpiringDate, link.ShortVersion, link.FullVersion)
}

func (l LinkSSDBRepo) GetMaps(nameExpiringToken ...string) ([]string, error) {
	if len(nameExpiringToken) != 0 {
		return l.client.HList("", nameExpiringToken[0], -1)
	}
	return l.client.HList("", "", -1)
}

func (l LinkSSDBRepo) GetKeys(mName string) ([]string, error) {
	return l.client.HKeys(mName, "", "", -1)
}
