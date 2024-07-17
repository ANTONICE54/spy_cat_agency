package services

import (
	"net/http"
	appErrors "spy_cat_agency/internal/appErorrs"
	"spy_cat_agency/internal/models"
)

type IMissionDao interface {
	AddMission(mission models.Mission) error
	Assign(missionId, catId uint) error
	GetMissionByID(id uint) (*models.Mission, error)
	GetMissionByCatID(catID uint) (*models.Mission, error)
	DeleteMission(id uint) error
	ListMissions() ([]models.Mission, error)
	UpdateMission(id uint, completed bool) error
	GetTarget(id uint) (*models.Target, error)
	DeleteTarget(id uint) error
	AddTarget(missionId uint, target models.Target) error
	CompleteTarget(id uint) error
	UpdateTargetNotes(id uint, notes string) error
}

type MissionService struct {
	MissionDao IMissionDao
}

func NewMissionService(missionDao IMissionDao) *MissionService {
	return &MissionService{
		MissionDao: missionDao,
	}
}

func (s *MissionService) AddMission(mission models.Mission) error {

	if len(mission.TargetList) > 3 || len(mission.TargetList) < 1 {
		return appErrors.NewHttpError("Target limit exceeded", http.StatusBadRequest, map[string]interface{}{"error": "mission can only have from 1 to 3 targets!"})

	}
	err := s.MissionDao.AddMission(mission)
	return err
}

func (s *MissionService) Assign(missionId, catId uint) error {
	mission, err := s.MissionDao.GetMissionByID(missionId)
	if err != nil && mission == nil {
		return appErrors.NewHttpError("There is no mission with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no mission with such id"})
	} else if err != nil {
		return err
	}

	catMission, err := s.MissionDao.GetMissionByCatID(catId)
	if err != nil {

		return err
	}

	if catMission.ID != 0 {
		return appErrors.NewHttpError("This cat has already been assigned a mission", http.StatusBadRequest, map[string]interface{}{"error": "this cat has already been assigned a mission"})
	}

	if mission.CatId != nil {
		return appErrors.NewHttpError("Already assigned", http.StatusBadRequest, map[string]interface{}{"error": "this mission is already assigned to a cat"})
	}

	err = s.MissionDao.Assign(missionId, catId)

	return err
}

func (s *MissionService) GetMission(id uint) (*models.Mission, error) {
	mission, err := s.MissionDao.GetMissionByID(id)
	return mission, err
}

func (s *MissionService) DeleteMission(id uint) error {

	mission, err := s.MissionDao.GetMissionByID(id)
	if err != nil && mission == nil {
		return appErrors.NewHttpError("There is no mission with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no mission with such id"})
	} else if err != nil {
		return err
	}

	if mission.CatId != nil {
		return appErrors.NewHttpError("Assigned mission cannot be deleted", http.StatusBadRequest, map[string]interface{}{"error": "assigned mission cannot be deleted"})
	}

	err = s.MissionDao.DeleteMission(id)

	return err
}

func (s *MissionService) ListMissions() ([]models.Mission, error) {
	list, err := s.MissionDao.ListMissions()
	return list, err
}

func (s *MissionService) UpdateMission(id uint, completed bool) error {
	mission, err := s.MissionDao.GetMissionByID(id)

	if err != nil && mission == nil {
		return appErrors.NewHttpError("There is no mission with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no mission with such id"})
	} else if err != nil {
		return err
	}

	if mission.IsCompleted {
		return appErrors.NewHttpError("Completed mission cannot be updated", http.StatusInternalServerError, map[string]interface{}{"error": "completed mission cannot be updated"})
	}
	err = s.MissionDao.UpdateMission(id, completed)
	return err
}

func (s *MissionService) GetTarget(id uint) (*models.Target, error) {
	target, err := s.MissionDao.GetTarget(id)

	return target, err
}

func (s *MissionService) DeleteTarget(id uint) error {
	target, err := s.MissionDao.GetTarget(id)
	if err != nil && target == nil {
		return appErrors.NewHttpError("There is no target with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no target with such id"})
	} else if err != nil {
		return err
	}
	if target.IsCompleted {
		return appErrors.NewHttpError("Completed target cannot be deleted", http.StatusInternalServerError, map[string]interface{}{"error": "completed target cannot be deleted"})
	}
	mission, err := s.MissionDao.GetMissionByID(target.MissionID)

	if err != nil {
		return err
	}

	if mission.IsCompleted {
		return appErrors.NewHttpError("Targets cannot be deleted from completed missions", http.StatusInternalServerError, map[string]interface{}{"error": "targets cannot be deleted from completed missions"})
	}

	if len(mission.TargetList) == 1 {
		return appErrors.NewHttpError("Target limit exceeded", http.StatusBadRequest, map[string]interface{}{"error": "you cannot delete the last target of the mission!"})

	}

	err = s.MissionDao.DeleteTarget(id)

	return err
}
func (s *MissionService) AddTarget(missionId uint, target models.Target) error {
	mission, err := s.MissionDao.GetMissionByID(missionId)

	if err != nil && mission == nil {
		return appErrors.NewHttpError("There is no mission with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no mission with such id"})
	} else if err != nil {
		return err
	}

	if mission.IsCompleted {
		return appErrors.NewHttpError("Completed mission cannot be updated with new targets", http.StatusInternalServerError, map[string]interface{}{"error": "completed mission cannot be updated with new targets"})
	}
	if len(mission.TargetList) == 3 {
		return appErrors.NewHttpError("Target limit exceeded", http.StatusBadRequest, map[string]interface{}{"error": "mission can only have from 1 to 3 targets!"})
	}

	err = s.MissionDao.AddTarget(missionId, target)

	return err
}

func (s *MissionService) CompleteTarget(id uint) error {
	var allTargetsCompleted = true

	target, err := s.MissionDao.GetTarget(id)
	if err != nil && target == nil {
		return appErrors.NewHttpError("There is no target with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no target with such id"})
	} else if err != nil {
		return err
	}

	if target.IsCompleted {
		return appErrors.NewHttpError("Completed target cannot be updated", http.StatusInternalServerError, map[string]interface{}{"error": "Completed target cannot be updated"})
	}

	mission, err := s.MissionDao.GetMissionByID(target.MissionID)

	if err != nil && mission == nil {
		return appErrors.NewHttpError("There is no mission with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no mission with such id"})
	} else if err != nil {
		return appErrors.NewHttpError("Internal server error", http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	if mission.IsCompleted {
		return appErrors.NewHttpError("Target cannot be updated", http.StatusInternalServerError, map[string]interface{}{"error": "Target of completed mission cannot be updated"})
	}
	err = s.MissionDao.CompleteTarget(id)
	if err != nil {
		return err
	}

	for _, v := range mission.TargetList {
		if v.ID == id {
			continue
		}
		if !v.IsCompleted {
			allTargetsCompleted = false
			break
		}
	}

	if allTargetsCompleted {
		err = s.MissionDao.UpdateMission(mission.ID, true)
		return err
	}

	return nil
}

func (s *MissionService) UpdateTargetNotes(id uint, notes string) error {
	target, err := s.MissionDao.GetTarget(id)
	if err != nil && target == nil {
		return appErrors.NewHttpError("There is no target with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no target with such id"})
	} else if err != nil {
		return err
	}

	if target.IsCompleted {
		return appErrors.NewHttpError("Completed target cannot be updated", http.StatusInternalServerError, map[string]interface{}{"error": "Completed target cannot be updated"})
	}

	mission, err := s.MissionDao.GetMissionByID(target.MissionID)

	if err != nil && mission == nil {
		return appErrors.NewHttpError("There is no mission with such id", http.StatusBadRequest, map[string]interface{}{"error": "there is no mission with such id"})
	} else if err != nil {
		return appErrors.NewHttpError("Internal server error", http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	if mission.IsCompleted {
		return appErrors.NewHttpError("Target cannot be updated", http.StatusInternalServerError, map[string]interface{}{"error": "Target of completed mission cannot be updated"})
	}
	err = s.MissionDao.UpdateTargetNotes(id, notes)

	return err

}
