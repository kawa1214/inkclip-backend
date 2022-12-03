package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/bookmark-manager/bookmark-manager/db/mock"
	db "github.com/bookmark-manager/bookmark-manager/db/sqlc"
	"github.com/bookmark-manager/bookmark-manager/token"
	"github.com/bookmark-manager/bookmark-manager/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestCeateNote(t *testing.T) {
	user, _ := randomUser(t)
	note := randomNote(t, user.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":   note.Title,
				"content": note.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNoteParams{
					UserID:  user.ID,
					Title:   note.Title,
					Content: note.Content,
				}
				store.EXPECT().
					CreateNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(note, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNote(t, recorder.Body, note, []db.Web{})
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"title":   note.Title,
				"content": note.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidTitleMax",
			body: gin.H{
				"title":   util.RandomString(101),
				"content": note.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTitleZero",
			body: gin.H{
				"title":   "",
				"content": note.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidContentMax",
			body: gin.H{
				"title":   note.Title,
				"content": util.RandomString(10001),
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DBErr",
			body: gin.H{
				"title":   note.Title,
				"content": note.Content,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNoteParams{
					UserID:  user.ID,
					Title:   note.Title,
					Content: note.Content,
				}
				store.EXPECT().
					CreateNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Note{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/notes"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetNote(t *testing.T) {
	user, _ := randomUser(t)
	note := randomNote(t, user.ID)
	n := 5
	webs := make([]db.Web, n)
	for i := 0; i < n; i++ {
		webs[i] = randomWeb(t, user.ID)
	}

	testCases := []struct {
		name          string
		noteID        string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			noteID: note.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetNote(gomock.Any(), gomock.Eq(note.ID)).
					Times(1).
					Return(note, nil)
				store.EXPECT().
					ListWebByNoteId(gomock.Any(), gomock.Eq(note.ID)).
					Times(1).
					Return(webs, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNote(t, recorder.Body, note, webs)
			},
		},
		{
			name:   "Unauthorized",
			noteID: note.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			noteID: "invalid",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "ErrDB",
			noteID: note.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetNote(gomock.Any(), gomock.Eq(note.ID)).
					Times(1).
					Return(db.Note{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			noteID: note.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetNote(gomock.Any(), gomock.Eq(note.ID)).
					Times(1).
					Return(db.Note{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "RequestFromAnotherUser",
			noteID: note.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				user2, _ := randomUser(t)
				user2Note := randomNote(t, user2.ID)
				store.EXPECT().
					GetNote(gomock.Any(), gomock.Eq(note.ID)).
					Times(1).
					Return(user2Note, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/notes/%s", tc.noteID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchNote(t *testing.T, body *bytes.Buffer, note db.Note, webs []db.Web) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var res noteResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, note.ID, res.ID)
	require.Equal(t, note.UserID, res.UserID)
	require.Equal(t, note.Title, res.Title)
	require.Equal(t, note.Content, res.Content)

	require.Equal(t, len(webs), len(res.Webs))
	for _, web := range res.Webs {
		require.NotEmpty(t, web)
		require.Equal(t, res.UserID, web.UserID)
	}
}

func randomNote(t *testing.T, userID uuid.UUID) db.Note {
	id, err := uuid.NewRandom()
	require.NoError(t, err)

	return db.Note{
		ID:      id,
		UserID:  userID,
		Title:   util.RandomString(6),
		Content: util.RandomString(100),
	}
}

// func randomNoteWeb(t *testing.T, noteID uuid.UUID, webID uuid.UUID) db.NoteWeb {
// 	return db.NoteWeb{
// 		NoteID: noteID,
// 		WebID:  webID,
// 	}
// }
