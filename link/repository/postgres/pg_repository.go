package postgres

import (
	"github.com/LittleLe6owski/LinkShortener/domain"
	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
)

const linkTable = "link"

type LinkPGRepo struct {
	repo domain.LinkPgRepo
	sess *dbr.Session
}

func NewLinkRepository(sess *dbr.Session) domain.LinkPgRepo {
	return &LinkPGRepo{
		sess: sess,
	}
}

func (l LinkPGRepo) GetLink(shortLink string) (link domain.Link, err error) {
	_, err = l.sess.Select("full_version").
		From(linkTable).
		Where("shortLink = ?", shortLink).
		Load(link)
	if err != nil {
		return
	}
	return
}

func (l LinkPGRepo) SaveLink(link domain.Link) (err error) {
	_, err = l.sess.InsertInto(linkTable).
		Columns("short_version", "full_version").
		Record(&link).Exec()
	return
}
