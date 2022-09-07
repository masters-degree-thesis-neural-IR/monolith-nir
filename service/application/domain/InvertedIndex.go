package domain

type InvertedIndex struct {
	CorpusSize int
	Df         map[string]int
	Idf        map[string]float64
	//Terms                   map[string][]string
	NormalizedDocumentFound map[string]NormalizedDocument
}
