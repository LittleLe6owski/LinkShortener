package usecase

import (
	"github.com/LittleLe6owski/LinkShortener/domain"
	"github.com/robfig/cron/v3"
	"github.com/speps/go-hashids/v2"
	"log"
	"sync"
	"time"
)

type LinkUseCase struct {
	pgRepo   domain.LinkPgRepo
	ssdbRepo domain.LinkSSDBRepo
}

func NewLinkUseCase(pgRepo domain.LinkPgRepo, ssdbRepo domain.LinkSSDBRepo) domain.LinkUseCase {
	return &LinkUseCase{
		pgRepo:   pgRepo,
		ssdbRepo: ssdbRepo,
	}
}

func (l LinkUseCase) GenShortLink(link string) (code string, err error) {
	code, err = l.genCode(domain.Link{
		FullVersion:  link,
		ExpiringDate: time.Now().Format("25.05"),
	})
	if err != nil {
		return
	}
	err = l.ssdbRepo.SaveLink(domain.Link{
		ExpiringDate: time.Now().AddDate(0, 0, 6).Format("25.05"),
		FullVersion:  link,
		ShortVersion: code,
	})
	if err != nil {
		log.Println("failed to save code and origin url")
		return
	}
	return
}

func (l LinkUseCase) GetOriginalLink(shortLink string) (string, error) {
	var link = domain.Link{}
	maps, err := l.ssdbRepo.GetMaps()
	if err != nil {
		return "", err
	}
	for i := range maps {
		link, err = l.ssdbRepo.GetLink(shortLink, maps[i])
		if err != nil {
			log.Printf("link getting error, map = %s", maps[i])
			continue
		}
		if link.FullVersion != "" {
			return link.FullVersion, nil
		}
	}
	return link.FullVersion, nil
}

func (l LinkUseCase) InitRestoreLinks() {
	c := cron.New()
	_, err := c.AddFunc("@daily", l.restoreLinks)
	if err != nil {
		log.Fatalf("restorer did not work out because %s", err)
	}
	c.Start()
}

func (l LinkUseCase) restoreLinks() {
	var wg sync.WaitGroup
	jobs := make(chan domain.Link)

	wg.Add(6)
	for i := 0; i < 6; i++ {
		go l.restoreWorker(&wg, jobs)
	}
	maps, err := l.ssdbRepo.GetMaps(time.Now().Format("25.05"))
	if err != nil {
		log.Printf("error getting map with keys because %s", err)
		return
	}
	for _, mapName := range maps {
		keys, err := l.ssdbRepo.GetKeys(mapName)
		if err != nil {
			log.Printf("error getting keys because %s", err)
			return
		}
		for _, key := range keys {
			jobs <- domain.Link{
				ExpiringDate: mapName,
				ShortVersion: key,
			}
		}
	}
}

func (l LinkUseCase) restoreWorker(wg *sync.WaitGroup, links <-chan domain.Link) {
	for key := range links {
		link, err := l.ssdbRepo.GetLink(key.ShortVersion, key.ExpiringDate)
		if err != nil {
			log.Printf("failed to get full url because %s", err)
		}
		err = l.ssdbRepo.DelLink(key.ShortVersion, key.ExpiringDate)
		if err != nil {
			log.Printf("failed to delete full url because %s", err)
		}
		err = l.pgRepo.SaveLink(link)
		if err != nil {
			log.Printf("failed to save link because %s", err)
		}
		wg.Done()
	}
}

func (l LinkUseCase) genCode(link domain.Link) (code string, err error) {
	hd := hashids.NewData()
	hd.Salt = link.FullVersion
	hd.MinLength = 8

	hashId, err := hashids.NewWithData(hd)
	if err != nil {
		return
	}
	return hashId.Encode(link.ConvertDateToListNumber())
}
