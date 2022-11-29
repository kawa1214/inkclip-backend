package api

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	db "github.com/bookmark-manager/bookmark-manager/db/sqlc"
	"github.com/bookmark-manager/bookmark-manager/token"
	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createWebRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// omitempty 空の場合はレスポンスに含めない
type webResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func newWebResponse(web db.Web) webResponse {
	return webResponse{
		ID:           web.ID,
		UserID:       web.UserID,
		URL:          web.Url,
		Title:        web.Title,
		ThumbnailURL: web.ThumbnailUrl,
		CreatedAt:    web.CreatedAt,
	}
}

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
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
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
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

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

	ctx.JSON(http.StatusOK, resWebs)
}

type deleteWebRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

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
