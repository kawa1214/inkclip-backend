package api

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/inkclip/backend/db/sqlc"
	"github.com/inkclip/backend/token"
	"github.com/lib/pq"
)

type createWebRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// omitempty 空の場合はレスポンスに含めない
type webResponse struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	UserID       uuid.UUID `json:"user_id" binding:"required"`
	URL          string    `json:"url" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty" binding:"required"`
	HTML         string    `json:"html" binding:"required"`
	CreatedAt    time.Time `json:"created_at" binding:"required"`
}

func newWebResponse(web db.Web) webResponse {
	return webResponse{
		ID:           web.ID,
		UserID:       web.UserID,
		URL:          web.Url,
		Title:        web.Title,
		ThumbnailURL: web.ThumbnailUrl,
		HTML:         web.Html,
		CreatedAt:    web.CreatedAt,
	}
}

// @Param request body api.createWebRequest true "query params"
// @Success 200 {object} api.webResponse
// @Router /webs [post]
// @Tags web
// @Security AccessToken
func (server *Server) createWeb(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req createWebRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res, err := http.Get(req.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Fatal("Error closing resp:", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	og := opengraph.NewOpenGraph()
	err = og.ProcessHTML(strings.NewReader(string(body)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var thumbnailURL string
	if len(og.Images) != 0 {
		thumbnailURL = og.Images[0].URL
	}

	arg := db.CreateWebParams{
		UserID:       authPayload.UserID,
		Url:          req.URL,
		Title:        og.Title,
		ThumbnailUrl: thumbnailURL,
		Html:         string(body),
	}

	if arg.Title == "" {
		arg.Title = req.URL
	}

	web, err := server.store.CreateWeb(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newWebResponse(web))
}

type getWebRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// @Param id path string true "Web ID"
// @Success 200 {object} api.webResponse
// @Router /webs/{id} [get]
// @Tags web
// @Security AccessToken
func (server *Server) getWeb(ctx *gin.Context) {
	var req getWebRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, _ := uuid.Parse(req.ID)

	web, err := server.store.GetWeb(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if web.UserID != authPayload.UserID {
		err := errors.New("web doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newWebResponse(web))
}

type listWebRequest struct {
	PageID   int32 `json:"page_id" form:"page_id" binding:"required,min=1"`
	PageSize int32 `json:"page_size" form:"page_size" binding:"required,min=5,max=10"`
}

type listWebResponse struct {
	Webs []webResponse `json:"webs"`
}

// @Param request query api.listWebRequest true "query params"
// @Success 200 {object} api.listWebResponse
// @Router /webs [get]
// @Tags web
// @Security AccessToken
func (server *Server) listWeb(ctx *gin.Context) {
	var req listWebRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListWebsByUserIdParams{
		UserID: authPayload.UserID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	webs, err := server.store.ListWebsByUserId(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resWebs := []webResponse{}
	for _, web := range webs {
		resWebs = append(resWebs, newWebResponse(web))
	}

	res := listWebResponse{
		Webs: resWebs,
	}

	ctx.JSON(http.StatusOK, res)
}

type deleteWebRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// @Param id path string true "Web ID"
// @Success 200 {} {}
// @Router /webs/{id} [delete]
// @Tags web
// @Security AccessToken
func (server *Server) deleteWeb(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req deleteWebRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, _ := uuid.Parse(req.ID)

	web, err := server.store.GetWeb(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if web.UserID != authPayload.UserID {
		err := errors.New("web doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := server.store.DeleteWeb(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
