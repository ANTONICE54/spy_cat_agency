package database

import (
	"database/sql"
	appErrors "spy_cat_agency/internal/appErorrs"
	"spy_cat_agency/internal/models"
)

type MissionRepository struct {
	*sql.DB
}

func NewMissionRepository(db *sql.DB) *MissionRepository {
	return &MissionRepository{
		db,
	}
}

func (db *MissionRepository) AddMission(mission models.Mission) error {

	var res models.Mission

	res.TargetList = make([]models.Target, 0)
	query := "INSERT INTO missions (name) VALUES($1) RETURNING id, name, cat_id, is_completed, created_at;"
	row := db.QueryRow(query, mission.Name)

	err := row.Scan(
		&res.ID,
		&res.Name,
		&res.CatId,
		&res.IsCompleted,
		&res.CreatedAt,
	)
	if err != nil {

		return appErrors.ErrDatabase
	}

	if mission.CatId != nil {
		err = db.Assign(res.ID, *mission.CatId)
		if err != nil {
			return appErrors.ErrDatabase
		}
	}

	for _, v := range mission.TargetList {
		query := "INSERT INTO targets (name, country, notes, mission_id) VALUES ($1, $2, $3, $4);"
		_, err = db.Exec(query, v.Name, v.Country, v.Notes, res.ID)
		if err != nil {
			return appErrors.ErrDatabase
		}
	}

	return nil
}

func (db *MissionRepository) Assign(missionId, catId uint) error {
	query := "UPDATE missions SET cat_id = $1 WHERE id = $2 ;"
	_, err := db.Exec(query, catId, missionId)

	if err != nil {

		return appErrors.ErrDatabase
	}

	return nil
}

func (db *MissionRepository) GetMissionByID(id uint) (*models.Mission, error) {
	var mission models.Mission
	mission.TargetList = make([]models.Target, 0)
	query := "SELECT * FROM missions WHERE id = $1;"
	row := db.QueryRow(query, id)

	err := row.Scan(
		&mission.ID,
		&mission.Name,
		&mission.CatId,
		&mission.IsCompleted,
		&mission.CreatedAt,
	)
	if err != nil {

		return nil, appErrors.ErrDatabase
	}

	query = "SELECT * FROM targets WHERE mission_id = $1;"
	rows, err := db.Query(query, mission.ID)

	if err != nil {
		return nil, appErrors.ErrDatabase
	}
	defer rows.Close()
	for rows.Next() {
		var target models.Target
		if err := rows.Scan(
			&target.ID,
			&target.MissionID,
			&target.Name,
			&target.Country,
			&target.Notes,
			&target.IsCompleted,
			&target.CreatedAt,
		); err != nil {
			return nil, appErrors.ErrDatabase
		}
		mission.TargetList = append(mission.TargetList, target)
	}

	return &mission, nil
}

func (db *MissionRepository) GetMissionByCatID(catID uint) (*models.Mission, error) {
	var mission models.Mission
	mission.TargetList = make([]models.Target, 0)
	query := "SELECT * FROM missions WHERE cat_id = $1;"
	row := db.QueryRow(query, catID)

	err := row.Scan(
		&mission.ID,
		&mission.Name,
		&mission.CatId,
		&mission.IsCompleted,
		&mission.CreatedAt,
	)
	if err != nil && err != sql.ErrNoRows {

		return nil, appErrors.ErrDatabase
	}

	query = "SELECT * FROM targets WHERE mission_id = $1;"
	rows, err := db.Query(query, mission.ID)

	if err != nil {

		return nil, appErrors.ErrDatabase
	}
	defer rows.Close()
	for rows.Next() {
		var target models.Target
		if err := rows.Scan(
			&target.ID,
			&target.MissionID,
			&target.Name,
			&target.Country,
			&target.Notes,
			&target.IsCompleted,
			&target.CreatedAt,
		); err != nil {

			return nil, appErrors.ErrDatabase
		}
		mission.TargetList = append(mission.TargetList, target)
	}

	return &mission, nil
}

func (db *MissionRepository) DeleteMission(id uint) error {
	query := "DELETE FROM missions WHERE id = $1;"
	_, err := db.Exec(query, id)
	if err != nil {
		return appErrors.ErrDatabase
	}

	query = "DELETE FROM targets WHERE mission_id = $1;"
	_, err = db.Exec(query, id)

	if err != nil {
		return appErrors.ErrDatabase
	}

	return nil
}

func (db *MissionRepository) ListMissions() ([]models.Mission, error) {
	res := make([]models.Mission, 0)
	query := "SELECT * FROM missions;"

	rows, err := db.Query(query)

	if err != nil {

		return nil, appErrors.ErrDatabase
	}
	defer rows.Close()

	for rows.Next() {
		var mission models.Mission
		mission.TargetList = make([]models.Target, 0)
		if err := rows.Scan(
			&mission.ID,
			&mission.Name,
			&mission.CatId,
			&mission.IsCompleted,
			&mission.CreatedAt,
		); err != nil {

			return nil, appErrors.ErrDatabase
		}
		res = append(res, mission)
	}

	query = "SELECT * FROM targets WHERE mission_id = $1"
	for i, mission := range res {
		rows, err := db.Query(query, mission.ID)
		if err != nil {

			return nil, appErrors.ErrDatabase
		}
		for rows.Next() {
			var target models.Target
			if err := rows.Scan(
				&target.ID,
				&target.MissionID,
				&target.Name,
				&target.Country,
				&target.Notes,
				&target.IsCompleted,
				&target.CreatedAt,
			); err != nil {

				return nil, appErrors.ErrDatabase
			}

			res[i].TargetList = append(res[i].TargetList, target)
			mission.TargetList = append(mission.TargetList, target)
		}
		rows.Close()
	}

	return res, nil
}

func (db *MissionRepository) UpdateMission(id uint, completed bool) error {
	query := "UPDATE missions SET is_completed = $1 WHERE id = $2;"
	_, err := db.Exec(query, completed, id)

	if err != nil {
		return appErrors.ErrDatabase
	}

	return nil
}

func (db *MissionRepository) GetTarget(id uint) (*models.Target, error) {
	var target models.Target
	query := "SELECT * FROM targets WHERE id = $1"
	row := db.QueryRow(query, id)
	err := row.Scan(
		&target.ID,
		&target.MissionID,
		&target.Name,
		&target.Country,
		&target.Notes,
		&target.IsCompleted,
		&target.CreatedAt,
	)

	if err != nil {

		return nil, appErrors.ErrDatabase
	}

	return &target, nil
}

func (db *MissionRepository) DeleteTarget(id uint) error {
	query := "DELETE FROM targets WHERE id = $1;"
	_, err := db.Exec(query, id)
	return err
}

func (db *MissionRepository) AddTarget(missionId uint, target models.Target) error {
	query := "INSERT INTO targets (name, country, notes, mission_id) VALUES ($1, $2, $3, $4);"
	_, err := db.Exec(query, target.Name, target.Country, target.Notes, missionId)
	if err != nil {
		return appErrors.ErrDatabase
	}
	return nil
}

func (db *MissionRepository) CompleteTarget(id uint) error {
	query := "UPDATE targets SET is_completed = TRUE WHERE id = $1;"
	_, err := db.Exec(query, id)
	if err != nil {
		return appErrors.ErrDatabase
	}
	return nil
}

func (db *MissionRepository) UpdateTargetNotes(id uint, notes string) error {
	query := "UPDATE targets SET notes = $1 WHERE id = $2;"
	_, err := db.Exec(query, notes, id)
	if err != nil {
		return appErrors.ErrDatabase
	}
	return nil
}
