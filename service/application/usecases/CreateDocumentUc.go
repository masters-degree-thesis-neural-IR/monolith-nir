package usecases

type DocumentUc interface {
	Create(id string, title string, body string) error
}
