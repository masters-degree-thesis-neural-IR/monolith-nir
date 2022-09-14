package score

import (
	"monolith-nir/service/application/domain"
)

func BM25(query []string, document *domain.NormalizedDocument, idf map[string]float64, corpusSize int, b float64, k1 float64) float64 {

	var score = 0.0
	docLength := float64(document.Length)
	frequencies := document.Tf
	avgDocLen := docLength / float64(corpusSize)

	for _, term := range query {

		if frequencies[term] == 0 {
			continue
		}

		freq := float64(frequencies[term])
		numerator := idf[term] * freq * (k1 + 1)
		denominator := freq + k1*(1-b+b*docLength/avgDocLen)
		score += numerator / denominator
	}

	return score
}

func BM25plus(query []string, document *domain.NormalizedDocument, idf map[string]float64, corpusSize int, b float64, k1 float64) float64 {

	var score = 0.0
	docLength := float64(document.Length)
	frequencies := document.Tf
	avgDocLen := docLength / float64(corpusSize)
	delta := 1.0

	for _, term := range query {

		if frequencies[term] == 0 {
			continue
		}
		//1f / (k1 * ((1 - b) + b * LENGTH_TABLE[i] / avgdl));
		freq := float64(frequencies[term])
		numerator := idf[term] * (delta + (freq * (k1 + 1)))
		denominator := k1 * ((1 - b + b*docLength/avgDocLen) + freq)
		score += numerator / denominator
	}

	return score
}

func BM25L(query []string, document *domain.NormalizedDocument, idf map[string]float64, corpusSize int, b float64, k1 float64) float64 {

	var score = 0.0
	docLength := float64(document.Length)
	frequencies := document.Tf
	avgDocLen := docLength / float64(corpusSize)
	delta := 1.0

	for _, term := range query {

		if frequencies[term] == 0 {
			continue
		}

		//ctd = q_freq / (1 - self.b + self.b * doc_len / self.avgdl)
		//score += (self.idf.get(q) or 0) * (self.k1 + 1) * (ctd + self.delta) / \
		//(self.k1 + ctd + self.delta)

		freq := float64(frequencies[term])
		ctd := freq / (1 - b + b*docLength/avgDocLen)
		numerator := idf[term] * (k1 + 1) * (ctd + delta)
		denominator := k1 + ctd + delta

		//numerator := idf[term] * freq * (k1 + 1)
		//denominator := freq + k1*(1-b+b*docLength/avgDocLen)
		score += numerator / denominator
	}

	return score
}
