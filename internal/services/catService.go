package services

import (
	"net/http"
	appErrors "spy_cat_agency/internal/appErorrs"
	"spy_cat_agency/internal/models"

	"github.com/lib/pq"
)

type ICatDao interface {
	Add(cat models.Cat) error
	Delete(id uint) error
	Update(id uint, salary float64) error
	List() ([]models.Cat, error)
	Get(id uint) (*models.Cat, error)
}

type CatService struct {
	CatDao ICatDao
}

func NewCatService(catDao ICatDao) *CatService {
	return &CatService{
		CatDao: catDao,
	}
}

func (s *CatService) HireCat(cat models.Cat) error {
	err := s.CatDao.Add(cat)

	return err
}

func (s *CatService) FireCat(id uint) error {

	cat, _ := s.CatDao.Get(id)

	if cat == nil {
		return appErrors.NewHttpError("There is no cat with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no cat with such id"})
	}

	err := s.CatDao.Delete(id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":

				return appErrors.NewHttpError("You cannot fire cat, while it is on mission", http.StatusBadRequest, map[string]interface{}{"error": "you cannot fire cat, while it is on mission"})
			}
		}

		return appErrors.ErrDatabase
	}

	return nil

}

func (s *CatService) UpdateSalary(id uint, salary float64) error {
	cat, _ := s.CatDao.Get(id)

	if cat == nil {
		return appErrors.NewHttpError("There is no cat with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no cat with such id"})
	}

	err := s.CatDao.Update(id, salary)

	return err

}

func (s *CatService) ListCats() ([]models.Cat, error) {
	list, err := s.CatDao.List()

	return list, err
}
func (s *CatService) GetCat(id uint) (*models.Cat, error) {
	cat, err := s.CatDao.Get(id)

	if cat == nil {
		return nil, appErrors.NewHttpError("There is no cat with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no cat with such id"})
	}

	return cat, err
}
