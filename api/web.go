package api

import (
	"fmt"
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

func newWeb(web db.Web) webResponse {
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
	payloadValue, exists := ctx.Get(authorizationPayloadKey)
	payload, isPyalod := payloadValue.(*token.Payload)

	if !exists || !isPyalod {
		err := fmt.Errorf("authorization payload not found")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

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
		UserID:       payload.UserID,
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

	ctx.JSON(http.StatusOK, newWeb(web))
}
