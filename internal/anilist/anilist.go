package anilist

import (
	"time"
)

type MediaSeason string

const (
	MediaSeasonWinter MediaSeason = "WINTER"
	MediaSeasonSpring MediaSeason = "SPRING"
	MediaSeasonSummer MediaSeason = "SUMMER"
	MediaSeasonFall   MediaSeason = "FALL"
)

type MediaSort string

const (
	MediaSortTrendingDesc   MediaSort = "TRENDING_DESC"
	MediaSortPopularityDesc MediaSort = "POPULARITY_DESC"
	MediaSortScoreDesc      MediaSort = "SCORE_DESC"
)

type ListMedia struct {
	Id    int
	Score int
}

type List struct {
	UserName       string
	Name           string
	IsCustom       bool
	MediaIds       []int
	ScoreByMediaId map[int]int
}

func (l List) GetId() string {
	return l.UserName + ":" + l.Name
}

type getUserAnimeListQuery struct {
	MediaListCollection struct {
		Lists []struct {
			Name         string
			IsCustomList bool
			Entries      []struct {
				Score     int `graphql:"score(format: POINT_100)"`
				MediaList struct {
					Media struct {
						Id int
					}
				} `graphql:"... on MediaList"`
			}
		}
	} `graphql:"MediaListCollection(userName: $userName, type: ANIME)"`
}

func FetchUserList(userName, name string) (*List, error) {
	return nil, nil
}

const searchAnimeListMaxPage = 4
const searchAnimeListPerPage = 50
const searchAnimeListQuery = `query (
  $page: Int!
  $season: MediaSeason
  $seasonYear: Int
  $sort: [MediaSort]
) {
  Page(page: $page, perPage: 50) {
		media(type: ANIME, season: $season, seasonYear: $seasonYear, sort: $sort) {
      id
    }
  }
}`

type SearchAnimeListData struct {
	Page struct {
		Media []struct {
			Id int `json:"id"`
		} `json:"media"`
	} `json:"Page"`
}

func getSeason(month time.Month) MediaSeason {
	switch month {
	case time.January, time.February, time.March:
		return MediaSeasonWinter
	case time.April, time.May, time.June:
		return MediaSeasonSpring
	case time.July, time.August, time.September:
		return MediaSeasonSummer
	case time.October, time.November, time.December:
		return MediaSeasonFall
	}
	panic("unreachable")
}

type searchListMeta struct {
	getInput func(page int) map[string]any
	name     string
}

var searchListQueryInputByName = map[string]searchListMeta{
	"trending": {
		name: "Trending",
		getInput: func(page int) map[string]any {
			return map[string]any{
				"page": page,
				"sort": []MediaSort{MediaSortTrendingDesc, MediaSortPopularityDesc},
			}
		},
	},
	"this-season": {
		name: "Popular This Season",
		getInput: func(page int) map[string]any {
			t := time.Now()
			return map[string]any{
				"page":       page,
				"season":     getSeason(t.Month()),
				"seasonYear": t.Year(),
				"sort":       []MediaSort{MediaSortPopularityDesc, MediaSortScoreDesc},
			}
		},
	},
	"next-season": {
		name: "Upcoming Next Season",
		getInput: func(page int) map[string]any {
			t := time.Now().AddDate(0, 3, 0)
			return map[string]any{
				"page":       page,
				"season":     getSeason(t.Month()),
				"seasonYear": t.Year(),
				"sort":       []MediaSort{MediaSortPopularityDesc, MediaSortScoreDesc},
			}
		},
	},
	"popular": {
		name: "All Time Popular",
		getInput: func(page int) map[string]any {
			return map[string]any{
				"page": page,
				"sort": []MediaSort{MediaSortPopularityDesc},
			}
		},
	},
	"top-100": {
		name: "Top 100",
		getInput: func(page int) map[string]any {
			return map[string]any{
				"page": page,
				"sort": []MediaSort{MediaSortScoreDesc},
			}
		},
	},
}

func IsValidSearchList(name string) bool {
	_, ok := searchListQueryInputByName[name]
	return ok
}

func FetchSearchList(name string) (*List, error) {
	return nil, nil
}

type fetchMediasQuery struct {
	Page struct {
		Media []struct {
			Id     int
			IdMal  int
			Format string
			Title  struct {
				English string
				Romaji  string
			}
			Description string
			BannerImage string
			CoverImage  struct {
				ExtraLarge string
			}
			Duration  int
			IsAdult   bool
			Genres    []string
			StartDate struct {
				Year int
			}
		} `graphql:"media(id_in: $ids)"`
	} `graphql:"Page(page: $page, perPage: 50)"`
}

type Media struct {
	Id          int
	IdMal       int
	Format      MediaFormat
	Title       string
	Description string
	BannerImage string
	CoverImage  string
	Duration    int
	IsAdult     bool
	Genres      []string
	StartYear   int
}

func FetchMedias(mediaIds []int) ([]Media, error) {
	if len(mediaIds) == 0 {
		return nil, nil
	}

	medias := []Media{}
	return medias, nil
}

type MediaFormat string

const (
	MediaFormatTV      MediaFormat = "TV"
	MediaFormatTVShort MediaFormat = "TV_SHORT"
	MediaFormatMovie   MediaFormat = "MOVIE"
	MediaFormatSpecial MediaFormat = "SPECIAL"
	MediaFormatOVA     MediaFormat = "OVA"
	MediaFormatONA     MediaFormat = "ONA"
	MediaFormatMusic   MediaFormat = "MUSIC"
	MediaFormatManga   MediaFormat = "MANGA"
	MediaFormatNovel   MediaFormat = "NOVEL"
	MediaFormatOneShot MediaFormat = "ONE_SHOT"
)

func (mf MediaFormat) ToSimple() string {
	switch mf {
	case MediaFormatTV, MediaFormatTVShort, MediaFormatOVA, MediaFormatONA:
		return "series"
	case MediaFormatMovie, MediaFormatSpecial, MediaFormatOneShot:
		return "movie"
	case MediaFormatMusic, MediaFormatManga, MediaFormatNovel:
		return ""
	default:
		return ""
	}
}

type fetchAnimeMediaFormatInfoQuery struct {
	Page struct {
		Media []struct {
			Id     int
			IdMal  int
			Format string
		} `graphql:"media(type: ANIME, id_in: $ids)"`
	} `graphql:"Page(page: $page, perPage: 50)"`
}

type MediaFormatInfo struct {
	Id     int
	IdMal  int
	Format MediaFormat
}

func FetchAnimeMediaFormatInfo(mediaIds []int) ([]MediaFormatInfo, error) {
	if len(mediaIds) == 0 {
		return nil, nil
	}

	infos := []MediaFormatInfo{}
	return infos, nil
}
