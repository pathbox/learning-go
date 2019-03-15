package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var cambridgeURL = "https://dictionary.cambridge.org/dictionary/english-chinese-simplified/"
var cambridgeEnURL = "https://dictionary.cambridge.org/dictionary/english/"

type cambridge struct{}

func (e cambridge) audio(word string, us bool) (mp3, ipa, def string, err error) {
	u := cambridgeURL + word
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("failed to get audio from cambridge: %v\n", err)
		return mp3, ipa, def, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response from cambridge: %v\n", err)
		return mp3, ipa, def, err
	}

	if us {
		val, ok := doc.Find(".us span.audio_play_button").Attr("data-src-mp3")
		if !ok {
			fmt.Printf("failed to get audio for %s\n", word)
			return mp3, ipa, def, errors.New("not found")
		}

		mp3 = "https://dictionary.cambridge.org" + val

		ipa = doc.Find(".us .pron span.ipa").First().Text()
	} else { //uk
		val, ok := doc.Find(".uk span.audio_play_button").Attr("data-src-mp3")
		if !ok {
			//fmt.Println("failed to get audio")
			return e.englishAudio(word, us)
		}

		mp3 = "https://dictionary.cambridge.org" + val

		ipa = doc.Find(".uk .pron span.ipa").First().Text()
	}

	doc.Find("span.def-body > span.trans").Each(func(_ int, s *goquery.Selection) {
		def = def + strings.Trim(s.Text(), " ") + "\n"
	})

	ipa = "[" + ipa + "]"
	return mp3, ipa, def, nil
}

func (e cambridge) englishAudio(word string, us bool) (mp3, ipa, def string, err error) {
	u := cambridgeEnURL + word
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("failed to get audio from cambridge: %v\n", err)
		return mp3, ipa, def, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response from cambridge: %v\n", err)
		return mp3, ipa, def, err
	}

	if us {
		val, ok := doc.Find(".us span.audio_play_button").Attr("data-src-mp3")
		if !ok {
			fmt.Println("failed to get audio")
			return mp3, ipa, def, errors.New("not found")
		}

		mp3 = "https://dictionary.cambridge.org" + val

		ipa = doc.Find(".us span.ipa").Text()
	} else { //uk
		val, ok := doc.Find(".uk span.audio_play_button").Attr("data-src-mp3")
		if !ok {
			fmt.Println("failed to get audio")
			return mp3, ipa, def, errors.New("not found")
		}

		mp3 = "https://dictionary.cambridge.org" + val

		ipa = doc.Find(".uk span.ipa").Text()
	}

	def = doc.Find("p.def-head b.def").Text()
	ipa = "[" + ipa + "]"

	return mp3, ipa, def, nil
}
