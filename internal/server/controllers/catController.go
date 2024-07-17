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

type CatController struct {
	CatService services.CatService
	errorLog   *log.Logger
}

func NewCatController(catService services.CatService, errorLog *log.Logger) *CatController {
	return &CatController{
		CatService: catService,
		errorLog:   errorLog,
	}
}

func (c *CatController) HireCat(ctx *gin.Context) {
	var catInfo models.Cat

	if err := ctx.ShouldBindJSON(&catInfo); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.CatService.HireCat(catInfo)
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

type FireCatRequest struct {
	CatID uint `json:"cat_id" binding:"required,numeric,gt=0"`
}

func (c *CatController) FireCat(ctx *gin.Context) {
	var req FireCatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.CatService.FireCat(req.CatID)
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

type UpdateSalaryRequest struct {
	Salary float64 `json:"salary" binding:"required,numeric,gt=0"`
	CatId  uint    `json:"cat_id" binding:"required,numeric,gt=0"`
}

func (c *CatController) UpdateSalary(ctx *gin.Context) {
	var req UpdateSalaryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	err := c.CatService.UpdateSalary(req.CatId, req.Salary)
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

type ListCatsResponse struct {
	List []models.Cat `json:"list"`
}

func (c *CatController) ListCats(ctx *gin.Context) {
	var resp ListCatsResponse

	list, err := c.CatService.ListCats()

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

type GetCatRequest struct {
	CatID uint `json:"cat_id" binding:"required,numeric,gt=0"`
}

func (c *CatController) GetCat(ctx *gin.Context) {
	var req GetCatRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := "Couldn't bind request:" + err.Error()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errMsg,
		})
		c.errorLog.Println("Couldn't bind request")
		return
	}

	cat, err := c.CatService.GetCat(req.CatID)

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

	ctx.JSON(http.StatusOK, cat)
}
