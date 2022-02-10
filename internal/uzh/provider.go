package uzh

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lukasmoellerch/mensa-cli/internal/base"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

const (
	deDinnerSuffix = " - Abendessen"
	enDinnerSuffix = " - Dinner"

	deLunchSuffix = " - Mittag"
	enLunchSuffix = " - Lunch"
)

var _ base.Provider = (*Provider)(nil)

type Provider struct{}

func (p *Provider) Id() string {
	return "uzh"
}

func (p *Provider) Label() string {
	return "UZH"
}

func (p *Provider) FetchMenus(ctx context.Context, caanteens []base.CanteenRef, date string, daytime string, lang string) ([]base.CanteenMenu, error) {
	eg, ctx := errgroup.WithContext(ctx)
	menus := make([]base.CanteenMenu, len(caanteens))
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	inputYear, inputWeek := parsedDate.ISOWeek()
	nowYear, nowWeek := now.ISOWeek()
	if inputYear != nowYear || inputWeek != nowWeek {
		return nil, fmt.Errorf("cannot fetch menu for %s, only current week is supported", date)
	}

	weekday := parsedDate.Weekday()
	var strWeekday string
	if weekday == time.Monday {
		strWeekday = "montag"
	} else if weekday == time.Tuesday {
		strWeekday = "dienstag"
	} else if weekday == time.Wednesday {
		strWeekday = "mittwoch"
	} else if weekday == time.Thursday {
		strWeekday = "donnerstag"
	} else if weekday == time.Friday {
		strWeekday = "freitag"
	} else {
		return nil, fmt.Errorf("cannot fetch menu for %s, only monday to friday are supported", date)
	}

	for i, ref := range caanteens {
		i, ref := i, ref
		eg.Go(func() error {
			var slug string
			var ok bool
			if daytime == "lunch" || daytime == "dinner" {
				slug, ok = ref.Meta[daytime]
				if !ok {
					return nil
				}
			} else {
				return fmt.Errorf("invalid daytime: %s", daytime)
			}
			menu, err := fetchMenuUzh(ctx, ref.ID, slug, lang, strWeekday)
			if err != nil {
				return err
			}
			menus[i] = menu
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return menus, nil
}

func (p *Provider) FetchCanteens(ctx context.Context, lang string) ([]base.CanteenMetadata, error) {
	url := fmt.Sprintf("http://www.mensa.uzh.ch/%s/menueplaene.html", lang)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	list := doc.Find("ul[role=navigation] > li > ul > li > a")

	linkPrefix := fmt.Sprintf("/%s/menueplaene/", lang)
	linkSuffix := ".html"

	lunchMap := make(map[string]string)
	dinnerMap := make(map[string]string)
	list.EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, ok := s.Attr("href")
		if !ok {
			return false
		}
		name := strings.TrimSpace(s.Text())

		hasLunch := false
		hasDinner := false

		if lang == "de" {
			if strings.HasSuffix(name, deLunchSuffix) {
				hasLunch = true
				name = strings.TrimSuffix(name, deLunchSuffix)
			} else if strings.HasSuffix(name, deDinnerSuffix) {
				hasDinner = true
				name = strings.TrimSuffix(name, deDinnerSuffix)
			} else {
				hasLunch = true
			}
		} else {
			if strings.HasSuffix(name, enLunchSuffix) {
				hasLunch = true
				name = strings.TrimSuffix(name, enLunchSuffix)
			} else if strings.HasSuffix(name, enDinnerSuffix) {
				hasDinner = true
				name = strings.TrimSuffix(name, enDinnerSuffix)
			} else {
				hasLunch = true
			}
		}
		href = strings.TrimPrefix(href, linkPrefix)
		href = strings.TrimSuffix(href, linkSuffix)

		if hasLunch {
			lunchMap[name] = href
		} else if hasDinner {
			dinnerMap[name] = href
		}

		return true
	})

	canteenMap := make(map[string]base.CanteenMetadata)
	for name, href := range lunchMap {
		canteenMap[name] = base.CanteenMetadata{
			ID:    "L:" + href,
			Label: name,
			Meta: map[string]string{
				"lunch": href,
			},
		}
	}
	for name, href := range dinnerMap {
		existing, ok := canteenMap[name]
		if ok {
			existing.ID = existing.ID + ",D:" + href
			existing.Meta["dinner"] = href
			canteenMap[name] = existing
		} else {
			canteenMap[name] = base.CanteenMetadata{
				ID:    "D:" + href,
				Label: name,
				Meta: map[string]string{
					"dinner": href,
				},
			}
		}
	}

	result := make([]base.CanteenMetadata, 0, len(canteenMap))
	return result, nil
}

func fetchMenuUzh(ctx context.Context, id string, slug string, lang string, weekday string) (base.CanteenMenu, error) {
	url := fmt.Sprintf("http://www.mensa.uzh.ch/%s/menueplaene/%s/%s.html", lang, slug, weekday)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return base.CanteenMenu{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	divs := doc.Find(".newslist-description > div")
	if divs.Length() != 1 {
		panic("no div")
	}

	meals := []base.Meal{}

	numParagraphs := 0
	var label string
	var description []string
	var studentPrice string
	var staffPrice string
	var externPrice string

	var loopErr error
	divs.Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		if goquery.NodeName(s) == "h3" {
			// Reset state
			numParagraphs = 0
			if label != "" {
				parsedStudentPrice, err := base.ParsePrice(studentPrice)
				if err != nil {
					loopErr = err
					return false
				}
				parsedStaffPrice, err := base.ParsePrice(staffPrice)
				if err != nil {
					loopErr = err
					return false
				}
				parsedExternPrice, err := base.ParsePrice(externPrice)
				if err != nil {
					loopErr = err
					return false
				}
				meals = append(meals, base.Meal{
					Label:       label,
					Description: description,
					Prices: base.MealPrices{
						Student: parsedStudentPrice,
						Staff:   parsedStaffPrice,
						Extern:  parsedExternPrice,
					},
				})
				description = []string{}
				studentPrice = ""
				staffPrice = ""
				externPrice = ""
			}

			// Split by | - separating the name and the price
			parts := strings.Split(s.Text(), " | ")
			if len(parts) != 2 {
				panic("invalid h3 tag")
			}
			label = parts[0]
			price := parts[1]

			// Remove leading "CHF" from price
			price = strings.TrimPrefix(price, "CHF ")
			// Separate different price categories by splitting using " / "
			priceParts := strings.Split(price, " / ")

			for i, pricePart := range priceParts {
				pricePart = strings.TrimSpace(pricePart)
				if i == 0 {
					studentPrice = pricePart
				} else if i == 1 {
					staffPrice = pricePart
				} else if i == 2 {
					externPrice = pricePart
				}
			}
		}
		if goquery.NodeName(s) == "p" {
			lines := make([]string, 0)

			for n := s.Get(0).FirstChild; n != nil; n = n.NextSibling {
				if n.Type == html.TextNode {
					text := strings.TrimSpace(n.Data)
					if text == "" {
						continue
					}
					lines = append(lines, text)
				}
			}

			if numParagraphs == 0 {
				description = lines
			}

			numParagraphs++
		}
		return true
	})
	if loopErr != nil {
		return base.CanteenMenu{}, loopErr
	}

	if label != "" {
		parsedStudentPrice, err := base.ParsePrice(studentPrice)
		if err != nil {
			return base.CanteenMenu{}, err
		}
		parsedStaffPrice, err := base.ParsePrice(staffPrice)
		if err != nil {
			return base.CanteenMenu{}, err
		}
		parsedExternPrice, err := base.ParsePrice(externPrice)
		if err != nil {
			return base.CanteenMenu{}, err
		}
		meals = append(meals, base.Meal{
			Label:       label,
			Description: description,
			Prices: base.MealPrices{
				Student: parsedStudentPrice,
				Staff:   parsedStaffPrice,
				Extern:  parsedExternPrice,
			},
		})
	}

	return base.CanteenMenu{
		Canteen: id,
		Meals:   meals,
	}, nil
}
