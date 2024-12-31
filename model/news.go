package model

import (
	"fmt"
	"magmar/model/dao"
	"strings"
)

// News ...
type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Newses ...
type Newses []*News

// NewNewses ...
func NewNewses(articles dao.News) Newses {
	newses := make(Newses, len(articles.Articles))
	for i, article := range articles.Articles {
		newses[i] = &News{
			Title:       article.Title,
			Description: article.Description,
		}
	}
	return newses
}

// ToPromptData ...
func (ns Newses) ToPromptData() string {
	var result strings.Builder
	for i, news := range ns {
		comma := ""
		if i != len(ns)-1 {
			comma = ","
		}
		result.WriteString(fmt.Sprintf("{title: %s, description: %s}%s\n", news.Title, news.Description, comma))
	}
	return result.String()
}
