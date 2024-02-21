package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/snykk/go-rest-boilerplate/internal/constants"
	"github.com/snykk/go-rest-boilerplate/internal/datasources/model"
	"github.com/snykk/go-rest-boilerplate/pkg/logger"
)

type api struct {
	router *gin.RouterGroup
	db     *sqlx.DB
}

type Response struct {
	Err    error `json:"err"`
	Status bool  `json:"status"`
	Data   any   `json:"data"`
}

func Api(router *gin.RouterGroup, db *sqlx.DB) *api {
	return &api{router: router, db: db}
}

func (api *api) Routes() {
	api.router.GET("/:uuid", api.GetInvite)
	api.router.POST("/:uuid", api.UpdateInvite)
}

func (api *api) GetInvite(ctx *gin.Context) {
	newUuid := ctx.Param("newUuid")

	invite, err := api.GetInviteByUUID(newUuid)
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status": false,
			"err":    err,
			"data":   invite,
		})
		return
	}

	ctx.JSON(http.StatusOK, invite)
}

func (api *api) UpdateInvite(ctx *gin.Context) {

	var requestBody model.Request

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Err: err, Status: false})
		return
	}

	if err := api.SetRsvpForInvite(requestBody.RSVP, requestBody.InviteId); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Err: err, Status: false})
		return
	}

	for _, attendee := range requestBody.Data {
		if err := api.SaveAttendee(attendee); err != nil {
			ctx.JSON(http.StatusBadRequest, Response{Err: err, Status: false})
			return
		}
	}

	ctx.JSON(http.StatusOK, Response{Status: true})
}

func (api *api) GetInviteByUUID(uuid string) (*model.Invite, error) {

	var invite model.Invite
	if err := api.db.Get(&invite, `SELECT * FROM invite WHERE id = $1`, uuid); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", invite), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return nil, err
	}

	logger.Info(fmt.Sprintf("%#v \n", invite), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})

	// if we get this data for the first time, we set the first opened at, otherwise the last opened at
	if invite.FirstOpenedAt == nil {
		if err := api.SetFirstOpenedAtForInvite(invite.Id); err != nil {
			logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
			return nil, err
		}
	} else {
		if err := api.SetLastOpenedAtForInvite(invite.Id); err != nil {
			logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
			return nil, err
		}
	}

	return &invite, nil
}

func (api *api) SetRsvpForInvite(rsvp bool, inviteId string) error {
	query := `
	UPDATE invite SET rsvp = $1, updated_at = $2 WHERE invite_id = $3`
	if _, err := api.db.Exec(query, rsvp, time.Now(), inviteId); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}

func (api *api) SetFirstOpenedAtForInvite(inviteId string) error {
	query := `
	UPDATE invite SET first_opened_at = $1 WHERE invite_id = $2`
	if _, err := api.db.Exec(query, time.Now(), inviteId); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}

func (api *api) SetLastOpenedAtForInvite(inviteId string) error {
	query := `
	UPDATE invite SET last_opened_at = $1 WHERE invite_id = $2`
	if _, err := api.db.Exec(query, time.Now(), inviteId); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}

func (api *api) SaveAttendee(attendee model.Attendee) error {
	query := `
	INSERT INTO attendee (
		id,
		invite_id,
		name,
		email,
		phone,
		age,
		created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (invite_id) DO UPDATE SET
		name = $3,
	    email = $4,
		phone = $5,
		age = $6,
	    updated_at = $7
	`
	if _, err := api.db.Exec(
		query,
		uuid.New(),
		attendee.InviteId,
		attendee.Name,
		attendee.Email,
		attendee.Phone,
		attendee.Age,
		time.Now(),
	); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}
