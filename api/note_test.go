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
	n := 5
	webs := make([]db.Web, n)
	webIds := make([]uuid.UUID, n)
	bodyWebIds := make([]string, n)
	for i := 0; i < n; i++ {
		webs[i] = randomWeb(t, user.ID)
		webIds[i] = webs[i].ID
		bodyWebIds[i] = webs[i].ID.String()
	}
	result := db.TxCreateNoteResult{
		Note: note,
		Webs: webs,
	}

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
				"web_ids": bodyWebIds,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TxCreateNoteParams{
					CreateNoteParams: db.CreateNoteParams{
						UserID:  user.ID,
						Title:   note.Title,
						Content: note.Content,
					},
					WebIds: webIds,
				}
				store.EXPECT().
					TxCreateNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNote(t, recorder.Body, note, webs)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"title":   note.Title,
				"content": note.Content,
				"web_ids": bodyWebIds,
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
				"web_ids": bodyWebIds,
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
				"web_ids": bodyWebIds,
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
				"web_ids": bodyWebIds,
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
				"web_ids": bodyWebIds,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TxCreateNoteParams{
					CreateNoteParams: db.CreateNoteParams{
						UserID:  user.ID,
						Title:   note.Title,
						Content: note.Content,
					},
					WebIds: webIds,
				}
				store.EXPECT().
					TxCreateNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TxCreateNoteResult{}, sql.ErrConnDone)
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
			name:   "ErrGetNoteDB",
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
					Return([]db.Web{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "ErrListWebsDB",
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

func TestListNote(t *testing.T) {
	user, _ := randomUser(t)

	noteN := 5
	notes := make([]db.Note, noteN)
	for i := 0; i < noteN; i++ {
		note := randomNote(t, user.ID)
		notes[i] = note
	}
	webN := 5
	webs := make([]db.Web, noteN*webN)
	webRows := make([]db.ListWebByNoteIdsRow, noteN*webN)
	for ni, note := range notes {
		for i := 0; i < webN; i++ {
			web := randomWeb(t, user.ID)
			webs[ni*noteN+i] = web
			webRows[ni*noteN+i] = db.ListWebByNoteIdsRow{
				ID:           web.ID,
				UserID:       web.UserID,
				Url:          web.Url,
				Title:        web.Title,
				ThumbnailUrl: web.ThumbnailUrl,
				NoteID:       note.ID,
				CreatedAt:    web.CreatedAt,
			}

		}
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: noteN,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListNotesByUserIdParams{
					UserID: user.ID,
					Limit:  int32(noteN),
					Offset: 0,
				}
				store.EXPECT().ListNotesByUserId(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(notes, nil)
				noteIDs := util.Select(notes, func(note db.Note) uuid.UUID {
					return note.ID
				})
				store.EXPECT().
					ListWebByNoteIds(gomock.Any(), gomock.InAnyOrder(noteIDs)).
					Times(1).
					Return(webRows, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchNotes(t, recorder.Body, notes, webs)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				pageID:   1,
				pageSize: noteN,
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
			name: "InvalidPageId",
			query: Query{
				pageID:   0,
				pageSize: noteN,
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
			name: "InvalidMaxPageSize",
			query: Query{
				pageID:   1,
				pageSize: 11,
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
			name: "InvalidMinPageSize",
			query: Query{
				pageID:   1,
				pageSize: 4,
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
			name: "ListNote DBError",
			query: Query{
				pageID:   1,
				pageSize: noteN,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListNotesByUserIdParams{
					UserID: user.ID,
					Limit:  int32(noteN),
					Offset: 0,
				}
				store.EXPECT().ListNotesByUserId(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Note{}, sql.ErrConnDone)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "ListWeb DBError",
			query: Query{
				pageID:   1,
				pageSize: noteN,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListNotesByUserIdParams{
					UserID: user.ID,
					Limit:  int32(noteN),
					Offset: 0,
				}
				store.EXPECT().ListNotesByUserId(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(notes, nil)
				noteIDs := util.Select(notes, func(note db.Note) uuid.UUID {
					return note.ID
				})
				store.EXPECT().
					ListWebByNoteIds(gomock.Any(), gomock.InAnyOrder(noteIDs)).
					Times(1).
					Return([]db.ListWebByNoteIdsRow{}, sql.ErrConnDone)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := "/notes"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchNotes(t *testing.T, body *bytes.Buffer, notes []db.Note, webs []db.Web) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var res []noteResponse
	err = json.Unmarshal(data, &res)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	for _, noteRes := range res {
		require.NotEmpty(t, noteRes)
		require.NotEmpty(t, noteRes.ID)
		require.NotEmpty(t, noteRes.UserID)
		require.NotEmpty(t, noteRes.Title)
		require.NotEmpty(t, noteRes.Content)
		require.Equal(t, len(webs)/len(notes), len(noteRes.Webs))
	}
}

func TestDeleteNoteAPI(t *testing.T) {
	user, _ := randomUser(t)
	note := randomNote(t, user.ID)
	n := 5
	webs := make([]db.Web, n)
	webIds := make([]uuid.UUID, n)
	bodyWebIds := make([]string, n)
	for i := 0; i < n; i++ {
		webs[i] = randomWeb(t, user.ID)
		webIds[i] = webs[i].ID
		bodyWebIds[i] = webs[i].ID.String()
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
					TxDeleteNote(gomock.Any(), gomock.Eq(db.TxDeleteNoteParams{
						NoteID: note.ID,
					})).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				require.Equal(t, data, []byte("{}"))
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
			name:   "ErrGetQuery",
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
			name:   "ErrDeleteQuery",
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
					TxDeleteNote(gomock.Any(), gomock.Eq(db.TxDeleteNoteParams{
						NoteID: note.ID,
					})).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "RequestFromUnauthorizedUser",
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

			url := fmt.Sprintf("/notes/%v", tc.noteID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
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
	require.NotEmpty(t, res)
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
