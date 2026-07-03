package app

type OrderingService struct {
	repo OrderingRepository
}

func NewOrderingService(repo OrderingRepository) *OrderingService {
	return &OrderingService{repo: repo}
}
