package nlp

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"math"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/score"
	"monolith-nir/service/application/stopwords"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

//func NotContains(document domain.NormalizedDocument, documents []domain.NormalizedDocument) bool {
//
//	for _, doc := range documents {
//		if doc.Id == document.Id {
//			return false
//		}
//	}
//
//	return true
//}

func NotContains(documentID string, documents []string) bool {

	for _, id := range documents {
		if documentID == id {
			return false
		}
	}

	return true
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func CleanSpecialCharacters(word string) string {

	reg, err := regexp.Compile("[\\p{P}\\p{S}]+")

	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(word, "")
}

func Tokenizer(document string, normalize bool) []string {

	fields := strings.Fields(document)

	if normalize {

		var localSlice = make([]string, 0)
		for _, token := range fields {
			var tempToken = strings.ToLower(RemoveAccents(CleanSpecialCharacters(token)))
			tempToken = strings.TrimSpace(tempToken)
			if len(tempToken) > 2 {
				localSlice = append(localSlice, tempToken)
			}
		}

		return localSlice
	}

	return fields

}

func StopWordLang(lang string) (map[string]bool, error) {

	if lang == "en" {
		return stopwords.English, nil
	}

	if lang == "pt" {
		return stopwords.Portuguese, nil
	}

	return nil, exception.ThrowValidationError("Not found language from stop word")
}

func RemoveStopWords(tokens []string, lang string) ([]string, error) {

	stopWordLang, err := StopWordLang(lang)

	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return make([]string, 0), nil
	}

	var localSlice = make([]string, 0)

	for _, token := range tokens {
		if !stopWordLang[token] {
			localSlice = append(localSlice, token)
		}
	}

	return localSlice, nil

}

func TermFrequency(tokens []string) map[string]int {

	localMap := make(map[string]int)

	for _, token := range tokens {

		if localMap[token] == 0 {
			localMap[token] = 1
		} else {
			localMap[token] = localMap[token] + 1
		}
	}

	return localMap

}

func CalcIdf(df map[string]int, corpusSize int) map[string]float64 {

	idf := make(map[string]float64)

	for term, frequency := range df {
		//idf[term] = math.log(1 + (corpus_size - freq + 0.5) / (freq + 0.5))
		//freq := float64(frequency) + 0.5
		//corpusSize := float64(corpusSize)
		//idf[term] = math.Log(1 + (corpusSize-freq)/freq)

		//Novo modelo para teste
		//idf = math.log((self.corpus_size + 1) / freq)
		freq := float64(frequency)
		corpusSize := float64(corpusSize + 1)
		idf[term] = math.Log(corpusSize / freq)
	}

	return idf

}

func ScoreBM25(query []string, invertedIndex *domain.InvertedIndex) []domain.QueryResult {

	queryResults := make([]domain.QueryResult, invertedIndex.CorpusSize)

	var i = 0
	for _, doc := range invertedIndex.NormalizedDocumentFound {

		score := score.BM25(query, &doc, invertedIndex.Idf, invertedIndex.CorpusSize, 0.75, 1.5)

		queryResults[i] = domain.QueryResult{
			Similarity:         score,
			NormalizedDocument: doc,
		}
		i++

	}

	return queryResults
}

func SortDesc(results []domain.QueryResult, top int) []domain.QueryResult {

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	var maxScore = make([]domain.QueryResult, 0)

	count := 0
	for _, document := range results {
		if count == top {
			break
		}
		maxScore = append(maxScore, document)
		count++
	}

	return maxScore
}
