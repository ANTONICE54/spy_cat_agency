package database

import (
	"database/sql"
	appErrors "spy_cat_agency/internal/appErorrs"
	"spy_cat_agency/internal/models"
)

type CatRepository struct {
	*sql.DB
}

func NewCatRepository(db *sql.DB) *CatRepository {
	return &CatRepository{
		db,
	}
}

func (db *CatRepository) Add(cat models.Cat) error {
	query := "INSERT INTO cats(name, years_of_experience, breed, salary ) VALUES ($1, $2, $3, $4);"

	_, err := db.Exec(query, cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary)

	if err != nil {
		return appErrors.ErrDatabase
	}

	return nil
}

func (db *CatRepository) Delete(id uint) error {
	query := "DELETE FROM cats WHERE id = $1;"

	_, err := db.Exec(query, id)

	return err
}

func (db *CatRepository) Update(id uint, salary float64) error {
	query := "UPDATE cats SET salary = $1 WHERE id = $2;"

	_, err := db.Exec(query, salary, id)

	if err != nil {
		return appErrors.ErrDatabase
	}

	return nil
}

func (db *CatRepository) List() ([]models.Cat, error) {
	list := []models.Cat{}
	query := "SELECT * FROM cats;"

	rows, err := db.Query(query)

	if err != nil {
		return nil, appErrors.ErrDatabase
	}

	defer rows.Close()

	for rows.Next() {
		var cat models.Cat
		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.YearsOfExperience,
			&cat.Breed,
			&cat.Salary,
			&cat.CreatedAt,
		); err != nil {
			return nil, appErrors.ErrDatabase
		}
		list = append(list, cat)
	}

	if err := rows.Close(); err != nil {
		return nil, appErrors.ErrDatabase
	}

	if err := rows.Err(); err != nil {
		return nil, appErrors.ErrDatabase
	}

	return list, nil
}

func (db *CatRepository) Get(id uint) (*models.Cat, error) {
	var res models.Cat
	query := "SELECT * FROM cats WHERE id = $1;"

	row := db.QueryRow(query, id)

	err := row.Scan(
		&res.ID,
		&res.Name,
		&res.YearsOfExperience,
		&res.Breed,
		&res.Salary,
		&res.CreatedAt,
	)

	if err != nil {
		return nil, appErrors.ErrDatabase
	}

	return &res, nil

}
