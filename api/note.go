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

type listNoteRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listNote(ctx *gin.Context) {
	var req listNoteRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListNotesByUserIdParams{
		UserID: authPayload.UserID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	notes, err := server.store.ListNotesByUserId(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	noteIDs := make([]uuid.UUID, len(notes))
	for i := range notes {
		noteIDs[i] = notes[i].ID
	}

	webRows, err := server.store.ListWebByNoteIds(ctx, noteIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := make([]noteResponse, len(notes))
	for i, note := range notes {
		var websFiltterByNote []db.Web
		for _, row := range webRows {
			if row.NoteID == note.ID {
				websFiltterByNote = append(websFiltterByNote, db.Web{
					ID:           row.ID,
					UserID:       row.UserID,
					Url:          row.Url,
					Title:        row.Title,
					ThumbnailUrl: row.ThumbnailUrl,
					CreatedAt:    row.CreatedAt,
				})
			}
		}
		res[i] = newNoteResponse(note, websFiltterByNote)
	}

	ctx.JSON(http.StatusOK, res)
}

type deleteNoteRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (server *Server) deleteNote(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req deleteNoteRequest
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
	if note.UserID != authPayload.UserID {
		err := errors.New("web doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if err := server.store.TxDeleteNote(ctx, db.TxDeleteNoteParams{
		NoteID: id,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type putNoteRequest struct {
	// ID      string   `uri:"id" binding:"required,uuid"`
	Title   string   `form:"title" binding:"required"`
	Content string   `form:"content" binding:"required"`
	WebIDs  []string `json:"web_ids" binding:"min=1,max=5,dive,uuid"`
}

func (server *Server) putNote(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	var req putNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	note, err := server.store.GetNote(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if note.UserID != authPayload.UserID {
		err := errors.New("note doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	webIDs := make([]uuid.UUID, len(req.WebIDs))
	for i, id := range req.WebIDs {
		webIDs[i], _ = uuid.Parse(id)
	}
	updateNoteArg := db.TxUpdateNoteParams{
		UpdateNoteParams: db.UpdateNoteParams{
			ID:      id,
			Title:   req.Title,
			Content: req.Content,
		},
		WebIds: webIDs,
	}
	result, err := server.store.TxUpdateNote(ctx, updateNoteArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newNoteResponse(result.Note, result.Webs))
}
