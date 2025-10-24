package app

type AppService struct{}

func NewService() *AppService {
	return &AppService{}
}

func (s *AppService) Hello() string {
	return "Hello, Wold!"
}
