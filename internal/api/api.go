package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/peti562/wedding/internal/constants"
	"github.com/peti562/wedding/internal/datasources/model"
	"github.com/peti562/wedding/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
	inviteId := ctx.Param("uuid")

	invite, err := api.GetInviteByUUID(inviteId)
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status": false,
			"err":    err,
			"data":   invite,
		})
		return
	}

	attendees, err := api.GetAttendeesForInvite(inviteId)
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status": false,
			"err":    err,
		})
		return
	}
	logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
	logger.Info(fmt.Sprintf("%#v \n", attendees), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})

	for _, att := range attendees {
		if invite.Attendees == nil {
			invite.Attendees = map[string]*model.Attendee{
				att.Id: att,
			}
		}

		invite.Attendees[att.Id] = att
	}

	ctx.JSON(http.StatusOK, invite)
}

func (api *api) UpdateInvite(ctx *gin.Context) {

	var invite model.Invite

	if err := ctx.ShouldBindJSON(&invite); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Err: err, Status: false})
		return
	}

	if err := api.SetRsvpForInvite(invite.RSVP, invite.Id); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{Err: err, Status: false})
		return
	}
	if invite.RSVP {
		for _, attendee := range invite.Attendees {
			if err := api.SaveAttendee(attendee, invite.Id); err != nil {
				ctx.JSON(http.StatusBadRequest, Response{Err: err, Status: false})
				return
			}
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

func (api *api) GetAttendeesForInvite(inviteId string) ([]*model.Attendee, error) {
	q := `SELECT id, invite_id, name, email, phone, age, is_child FROM attendee WHERE invite_id = $1`

	var attendees []*model.Attendee

	rows, err := api.db.Query(q, inviteId)
	if err != nil {
		logger.Info(fmt.Sprintf("%#v \n", attendees), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var att model.Attendee
		if err := rows.Scan(
			&att.Id,
			&att.InviteId,
			&att.Name,
			&att.Email,
			&att.Phone,
			&att.Age,
			&att.IsChild,
		); err != nil {
			logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})

			return nil, err
		}

		attendees = append(attendees, &att)
	}

	return attendees, rows.Err()
}

func (api *api) SetRsvpForInvite(rsvp bool, inviteId string) error {
	query := `UPDATE invite SET rsvp = $1, updated_at = $2 WHERE id = $3`
	if _, err := api.db.Exec(query, rsvp, time.Now(), inviteId); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}

func (api *api) SetFirstOpenedAtForInvite(inviteId string) error {
	query := `UPDATE invite SET first_opened_at = $1 WHERE id = $2`
	if _, err := api.db.Exec(query, time.Now(), inviteId); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}

func (api *api) SetLastOpenedAtForInvite(inviteId string) error {
	query := `UPDATE invite SET last_opened_at = $1 WHERE id = $2`
	if _, err := api.db.Exec(query, time.Now(), inviteId); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}

func (api *api) SaveAttendee(attendee *model.Attendee, inviteId string) error {
	query := `
	INSERT INTO attendee (
		invite_id,
		name,
		email,
		phone,
		age,
		is_child,
		active
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO UPDATE SET
		name = $2,
	    email = $3,
		phone = $4,
		age = $5,
	    is_child = $6,
	    active = $7,
	    updated_at = $8 
	`
	if _, err := api.db.Exec(
		query,
		inviteId,
		attendee.Name,
		attendee.Email,
		attendee.Phone,
		attendee.Age,
		attendee.IsChild,
		attendee.Active,
		time.Now(),
	); err != nil && err != sql.ErrNoRows {
		logger.Info(fmt.Sprintf("%#v \n", err), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryDatabase})
		return err
	}
	return nil
}
