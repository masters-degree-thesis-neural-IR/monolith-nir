package repositories

type IndexMemoryRepository interface {
	FindByTerm(term string) ([]string, error)
	Update(term string, documents []string) error
	Save(term string, documents []string) error
}
