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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/inkclip/backend/db/mock"
	db "github.com/inkclip/backend/db/sqlc"
	"github.com/inkclip/backend/token"
	"github.com/inkclip/backend/util"
	"github.com/jarcoal/httpmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateWebAPI(t *testing.T) {
	user, _ := randomUser(t)
	web := randomWeb(t, user.ID)

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
				"url": web.Url,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				httpmock.RegisterResponder("GET", web.Url,
					httpmock.NewStringResponder(
						http.StatusOK,
						web.Html,
					),
				)

				arg := db.CreateWebParams{
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Title,
					ThumbnailUrl: web.ThumbnailUrl,
					Html:         web.Html,
				}
				store.EXPECT().
					CreateWeb(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(web, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWeb(t, recorder.Body, web)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"url": web.Url,
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
			name: "ErrDB",
			body: gin.H{
				"url": web.Url,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				httpmock.RegisterResponder("GET", web.Url,
					httpmock.NewStringResponder(
						http.StatusOK,
						web.Html,
					),
				)

				arg := db.CreateWebParams{
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Title,
					ThumbnailUrl: web.ThumbnailUrl,
					Html:         web.Html,
				}
				store.EXPECT().
					CreateWeb(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Web{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "CannotGetOGP",
			body: gin.H{
				"url": web.Url,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				httpmock.RegisterResponder("GET", web.Url,
					httpmock.NewStringResponder(
						http.StatusOK,
						"",
					),
				)

				arg := db.CreateWebParams{
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Url,
					ThumbnailUrl: "",
				}
				expectWeb := db.Web{
					ID:           web.ID,
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Url,
					ThumbnailUrl: "",
				}
				store.EXPECT().
					CreateWeb(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(expectWeb, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				matchWeb := db.Web{
					ID:           web.ID,
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Url,
					ThumbnailUrl: "",
				}
				requireBodyMatchWeb(t, recorder.Body, matchWeb)
			},
		},
		{
			name: "InvalidURL",
			body: gin.H{
				"url": "invalid url",
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
			name: "URLRequestFailed",
			body: gin.H{
				"url": web.Url,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateURLAndUserID",
			body: gin.H{
				"url": web.Url,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				httpmock.RegisterResponder("GET", web.Url,
					httpmock.NewStringResponder(
						http.StatusOK,
						web.Html,
					),
				)

				arg := db.CreateWebParams{
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Title,
					ThumbnailUrl: web.ThumbnailUrl,
					Html:         web.Html,
				}
				store.EXPECT().
					CreateWeb(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Web{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			url := "/webs"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetWebAPI(t *testing.T) {
	user, _ := randomUser(t)
	web := randomWeb(t, user.ID)

	testCases := []struct {
		name          string
		webID         string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(web, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWeb(t, recorder.Body, web)
			},
		},
		{
			name:  "Unauthorized",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "InvalidID",
			webID: "invalid",
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
			name:  "ErrDB",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(db.Web{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				user2, _ := randomUser(t)
				user2Web := randomWeb(t, user2.ID)
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(user2Web, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "RequestFromUnauthorizedUser",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(db.Web{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprintf("/webs/%v", tc.webID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListWebAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	webs := make([]db.Web, n)
	// rows := make([]db.ListWebByNoteIdsRow, n)
	for i := 0; i < n; i++ {
		webs[i] = randomWeb(t, user.ID)
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
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListWebsByUserIdParams{
					UserID: user.ID,
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListWebsByUserId(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(webs, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchWebs(t, recorder.Body, webs)
			},
		},
		// {
		// 	name: "InvalidPageID",
		// 	query: Query{
		// 		pageID:   0,
		// 		pageSize: n,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidMaxPageSize",
		// 	query: Query{
		// 		pageID:   1,
		// 		pageSize: 11,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidMinPageSize",
		// 	query: Query{
		// 		pageID:   1,
		// 		pageSize: 4,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "Unauthorized",
		// 	query: Query{
		// 		pageID:   1,
		// 		pageSize: n,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "DBError",
		// 	query: Query{
		// 		pageID:   1,
		// 		pageSize: n,
		// 	},
		// 	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
		// 		addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.ListWebsByUserIdParams{
		// 			UserID: user.ID,
		// 			Limit:  int32(n),
		// 			Offset: 0,
		// 		}
		// 		store.EXPECT().
		// 			ListWebsByUserId(gomock.Any(), gomock.Eq(arg)).
		// 			Times(1).
		// 			Return([]db.Web{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
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

			url := "/webs"
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

func TestDeleteWebAPI(t *testing.T) {
	user, _ := randomUser(t)
	web := randomWeb(t, user.ID)

	testCases := []struct {
		name          string
		webID         string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(web, nil)

				store.EXPECT().
					DeleteWeb(gomock.Any(), gomock.Eq(web.ID)).
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
			name:  "Unauthorized",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:  "InvalidID",
			webID: "invalid",
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
			name:  "NotFound",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(db.Web{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "ErrGetQuery",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(db.Web{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "ErrDeleteQuery",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(web, nil)

				store.EXPECT().
					DeleteWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "RequestFromUnauthorizedUser",
			webID: web.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				user2, _ := randomUser(t)
				user2Web := randomWeb(t, user2.ID)
				store.EXPECT().
					GetWeb(gomock.Any(), gomock.Eq(web.ID)).
					Times(1).
					Return(user2Web, nil)
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

			url := fmt.Sprintf("/webs/%v", tc.webID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomWeb(t *testing.T, userID uuid.UUID) db.Web {
	id, err := uuid.NewRandom()
	require.NoError(t, err)
	title := util.RandomName()
	thumbnailURL := util.RandomThumbnailURL()
	return db.Web{
		ID:           id,
		UserID:       userID,
		Url:          util.RandomURL(),
		Title:        title,
		ThumbnailUrl: thumbnailURL,
		Html: fmt.Sprintf(`<html>
		<head>
		<meta property="og:title" content="%s" />
		<meta property="og:image" content="%s" />
		</head>
		<body></body>
		</html>`,
			title,
			thumbnailURL,
		),
	}
}

func requireBodyMatchWeb(t *testing.T, body *bytes.Buffer, web db.Web) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotWeb db.Web
	err = json.Unmarshal(data, &gotWeb)
	require.NoError(t, err)
	require.Equal(t, web.ID, gotWeb.ID)
	require.Equal(t, web.UserID, gotWeb.UserID)
	require.Equal(t, web.Url, gotWeb.Url)
	require.Equal(t, web.Title, gotWeb.Title)
	require.Equal(t, web.ThumbnailUrl, gotWeb.ThumbnailUrl)
}

func requireBodyMatchWebs(t *testing.T, body *bytes.Buffer, webs []db.Web) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got listWebResponse
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	require.Equal(t, len(got.Webs), len(webs))
}
