package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"monolith-nir/service/application/domain"
	"monolith-nir/service/application/exception"
	"monolith-nir/service/application/service"
	"monolith-nir/service/infraestructure/controller"
	"monolith-nir/service/infraestructure/dto"
	"monolith-nir/service/infraestructure/memory"
	"monolith-nir/service/infraestructure/mysq"
	"monolith-nir/service/infraestructure/sns"
	"net/http"
	"time"
)

func errorHandler(err error, c *gin.Context) {

	switch err.(type) {
	case *exception.ValidationError:
		err, _ := err.(*exception.ValidationError)
		c.JSON(err.StatusCode, gin.H{"error": err.Message})

	case *exception.UnexpectedError:
		err, _ := err.(*exception.UnexpectedError)
		c.JSON(err.StatusCode, gin.H{"error": err.Message})

	default:
		c.JSON(500, gin.H{"error": "Internal error"})
	}

}

func main() {

	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/nir-database")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Success!")

	documentEvent := sns.NewDocumentEvent(nil, "")
	memoryRepository := memory.NewMemoryIndexRepository()

	documentRepository := mysq.NewDocumentRepository(db)
	documentMetricsRepository := mysq.NewDocumentMetricsRepository(db)
	//indexRepository := mysq.NewIndexRepository(db)
	indexService := service.NewIndexService(documentMetricsRepository, memoryRepository)
	documentService := service.NewDocumentService(documentEvent, documentRepository)
	searchService := service.NewSearch(documentMetricsRepository, memoryRepository, documentRepository)
	controller := controller.NewController(documentService, indexService, searchService)

	r := gin.New()
	r.POST("/nir", func(c *gin.Context) {

		var document dto.Document
		if err := c.ShouldBindJSON(&document); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := controller.CreateDocument(document)

		if err != nil {
			errorHandler(err, c)
		}

		c.JSON(http.StatusCreated, gin.H{"success": 201})

	})

	r.GET("/nir", func(c *gin.Context) {

		start := time.Now()
		paramPairs := c.Request.URL.Query()
		results, err := controller.Search.SearchDocument(paramPairs.Get("query"))
		duration := time.Since(start)

		if err != nil {
			errorHandler(err, c)
		}

		body, err := makeBody(results, duration)
		c.JSON(http.StatusOK, body)
	})
	r.Run()
}

func makeBody(results []domain.QueryResult, duration time.Duration) (dto.Result, error) {

	total := len(results)

	rst := dto.Result{
		Total:        total,
		Duration:     duration.String(),
		QueryResults: make([]dto.QueryResult, total),
	}

	for i, result := range results {
		rst.QueryResults[i] = dto.QueryResult{
			Similarity: result.Similarity,
			Document: dto.Document{
				Id: result.NormalizedDocument.Id,
				//Title: result.Document.Title,
				//Body:  result.Document.Body,
			},
		}
	}

	//body, err := json.Marshal(rst)
	return rst, nil

}
