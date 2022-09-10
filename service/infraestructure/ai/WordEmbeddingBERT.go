package ai

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"monolith-nir/service/application/ports"
	"net/http"
)

type WordEmbeddingBERT struct {
}

func NewWordEmbeddingBERT() ports.WordEmbedding {
	return WordEmbeddingBERT{}
}

type Response struct {
	StatusCode int       `json:"statusCode"`
	Embedding  []float64 `json:"embedding"`
}

func (w WordEmbeddingBERT) Generate(sentence string) []float64 {

	postBody, _ := json.Marshal(map[string]string{
		"sentence": sentence,
	})

	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://127.0.0.1:5000/sentence-embedding", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	var response Response
	json.Unmarshal([]byte(sb), &response)

	return response.Embedding
}
