package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/golang/glog"
)

type PageParse struct {
	dataCh chan *PageData
}

func NewPageParse() *PageParse {
	parse := new(PageParse)
	parse.dataCh = make(chan *PageData)
	return parse
}

func (p *PageParse) Run() {
	select {
	case data := <-p.dataCh:
		p.parsePageData(data)
	}
}

func (p *PageParse) SendPageData(data *PageData) {
	p.dataCh <- data
}

func (p *PageParse) parsePageData(data *PageData) {
	filmNameReg := regexp.MustCompile(`<img width="100" alt="(?s:(.*?))"`)
	filmNames := filmNameReg.FindAllStringSubmatch(data.data, -1)
	glog.Infoln(filmNames)

	filmScoreReg := regexp.MustCompile(`<span class="rating_num" property="v:average">(.*)</span>`)
	filmScores := filmScoreReg.FindAllStringSubmatch(data.data, -1)
	glog.Infoln(filmScores)

	filmScoreNumReg := regexp.MustCompile(`<span>(.*)人评价</span>`)
	filmScoreNum := filmScoreNumReg.FindAllStringSubmatch(data.data, -1)
	glog.Infoln(filmScoreNum)

	filmCommentReg := regexp.MustCompile(`<span class="inq">(.*)</span>`)
	filmComments := filmCommentReg.FindAllStringSubmatch(data.data, -1)
	glog.Infoln(filmComments)

	filmContentReg := regexp.MustCompile(`(?s)<p class="">(.*?)</p>`)
	filmContents := filmContentReg.FindAllStringSubmatch(data.data, -1)
	glog.Infoln(filmContents)

	films := make([]FilmData, len(filmNames))
	for i := 0; i < len(films); i++ {
		if len(filmComments) <= i {
			filmComments = append(filmComments, []string{"", ""})
		}
		films[i] = FilmData{
			Name:    filmNames[i][1],
			Score:   filmScores[i][1],
			People:  filmScoreNum[i][1],
			Comment: filmComments[i][1],
			Content: filmContents[i][1],
		}
	}

	err := p.save2File(films, data.index)
	if err != nil {
		glog.Errorf("%d page data save to File error %e", data.index, err)
	}
}

func (p *PageParse) save2File(films []FilmData, page uint) error {
	f, err := os.Create(fmt.Sprintf("result %d.json", page))
	if err != nil {
		return err
	}

	defer f.Close()

	bytes, err := json.MarshalIndent(films, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(bytes))
	if err != nil {
		return err
	}

	return nil
}
