package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/bookmark-manager/bookmark-manager/db/sqlc"
	"github.com/bookmark-manager/bookmark-manager/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createNoteRequest struct {
	Title   string   `json:"title" binding:"required,min=1,max=100"`
	Content string   `json:"content" binding:"required,max=10000"`
	WebIDs  []string `json:"web_ids" binding:"min=1,max=5,dive,uuid"`
}

type noteResponse struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"user_id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	Webs      []webResponse `json:"webs"`
}

func newNoteResponse(note db.Note, webs []db.Web) noteResponse {
	webResponses := make([]webResponse, len(webs))
	for i := range webs {
		webResponses[i] = newWebResponse(webs[i])
	}
	return noteResponse{
		ID:        note.ID,
		UserID:    note.UserID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		Webs:      webResponses,
	}
}

func (server *Server) createNote(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req createNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	webIds := make([]uuid.UUID, len(req.WebIDs))
	for i, id := range req.WebIDs {
		webIds[i], _ = uuid.Parse(id)
	}

	arg := db.TxCreateNoteParams{
		CreateNoteParams: db.CreateNoteParams{
			UserID:  authPayload.UserID,
			Title:   req.Title,
			Content: req.Content,
		},
		WebIds: webIds,
	}

	txNote, err := server.store.TxCreateNote(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newNoteResponse(txNote.Note, txNote.Webs))
}

type getNoteRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getNote(ctx *gin.Context) {
	var req getNoteRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, _ := uuid.Parse(req.ID)
	note, err := server.store.GetNote(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if note.UserID != authPayload.UserID {
		err := errors.New("note doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	webs, err := server.store.ListWebByNoteId(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newNoteResponse(note, webs))
}
