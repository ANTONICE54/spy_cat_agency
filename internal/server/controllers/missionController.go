package controllers

import (
	"errors"
	"log"
	"net/http"
	appErrors "spy_cat_agency/internal/appErorrs"
	"spy_cat_agency/internal/models"
	"spy_cat_agency/internal/services"

	"github.com/gin-gonic/gin"
)

type MissionController struct {
	MissionService services.MissionService
	errorLog       *log.Logger
}

func NewMissionController(missionService services.MissionService, errorLog *log.Logger) *MissionController {
	return &MissionController{
		MissionService: missionService,
		errorLog:       errorLog,
	}
}

func (c *MissionController) AddMission(ctx *gin.Context) {
	var missionInfo models.Mission

	if err := ctx.ShouldBindJSON(&missionInfo); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.AddMission(missionInfo)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)

}

type AssignRequest struct {
	MissionID uint `json:"mission_id" binding:"required,numeric,gt=0"`
	CatID     uint `json:"cat_id" binding:"required,numeric,gt=0"`
}

func (c *MissionController) Assign(ctx *gin.Context) {
	var req AssignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.Assign(req.MissionID, req.CatID)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)

}

type GetMissionRequest struct {
	MissionID uint `json:"mission_id" binding:"required,numeric,gt=0"`
}

func (c *MissionController) GetMission(ctx *gin.Context) {
	var req GetMissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	mission, err := c.MissionService.GetMission(req.MissionID)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.JSON(http.StatusOK, mission)
}

type DeleteMissionRequest struct {
	MissionID uint `json:"mission_id" binding:"required,numeric,gt=0"`
}

func (c *MissionController) DeleteMission(ctx *gin.Context) {
	var req DeleteMissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.DeleteMission(req.MissionID)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)
}

type ListMissionsResponse struct {
	List []models.Mission `json:"list"`
}

func (c *MissionController) ListMissions(ctx *gin.Context) {
	var resp ListMissionsResponse

	list, err := c.MissionService.ListMissions()

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	resp.List = list
	ctx.JSON(http.StatusOK, list)
}

type UpdateMissionRequest struct {
	MissionID   uint `json:"mission_id" binding:"required,numeric,gt=0"`
	IsCompleted bool `json:"is_completed" binding:"required"`
}

func (c *MissionController) UpdateMission(ctx *gin.Context) {
	var req UpdateMissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.UpdateMission(req.MissionID, req.IsCompleted)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)
}

type GetTargetRequest struct {
	TargetID uint `json:"target_id" binding:"required,numeric,gt=0"`
}

func (c *MissionController) GetTarget(ctx *gin.Context) {
	var req GetTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	target, err := c.MissionService.GetTarget(req.TargetID)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)

			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.JSON(http.StatusOK, target)
}

type DeleteTargetRequest struct {
	TargetID uint `json:"target_id" binding:"required,numeric,gt=0"`
}

func (c *MissionController) DeleteTarget(ctx *gin.Context) {
	var req DeleteTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.DeleteTarget(req.TargetID)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)
}

type AddTargetRequest struct {
	MissionID uint          `json:"mission_id" binding:"required,numeric,gt=0"`
	TargetObj models.Target `json:"target" binding:"required"`
}

func (c *MissionController) AddTarget(ctx *gin.Context) {
	var req AddTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.AddTarget(req.MissionID, req.TargetObj)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)
}

type CompleteTargetRequest struct {
	TargetID uint `json:"target_id" binding:"required,numeric,gt=0"`
}

func (c *MissionController) CompleteTarget(ctx *gin.Context) {
	var req CompleteTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.CompleteTarget(req.TargetID)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)
}

type UpdateTargetNotesRequest struct {
	TargetID uint   `json:"target_id" binding:"required,numeric,gt=0"`
	Notes    string `json:"notes" binding:"required"`
}

func (c *MissionController) UpdateTargetNotes(ctx *gin.Context) {
	var req UpdateTargetNotesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.MissionService.UpdateTargetNotes(req.TargetID, req.Notes)

	if err != nil {
		var httpErr *appErrors.HttpError
		if errors.As(err, &httpErr) {
			ctx.JSON(httpErr.StatusCode, httpErr.JSONResponse)
			c.errorLog.Println(httpErr.Message)
			return
		}
		ctx.JSON(appErrors.ErrInternalServer.StatusCode, appErrors.ErrInternalServer.JSONResponse)
		c.errorLog.Println(appErrors.ErrInternalServer.Message)
		return
	}

	ctx.Status(http.StatusOK)
}
