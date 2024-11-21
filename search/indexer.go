package search

import "github.com/xelathan/mini_search_engine/db"

// in memory representation of the index. inverted index
type Index map[string][]string

func (idx Index) Add(docs []db.CrawledUrl) {
	for _, doc := range docs {
		for _, token := range analyze(doc.Url + " " + doc.PageTitle + " " + doc.PageDescription + " " + doc.Heading) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				continue
			}
			idx[token] = append(ids, doc.ID)
		}
	}
}
