package api

import (
	"bytes"
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
		checkResponse func(recorder *httptest.ResponseRecorder)
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
						fmt.Sprintf(`<html>
						<head>
						<meta property="og:title" content="%s" />
						<meta property="og:image" content="%s" />
						</head>
						<body></body>
						</html>`,
							web.Title,
							web.ThumbnailUrl,
						),
					),
				)

				arg := db.CreateWebParams{
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Title,
					ThumbnailUrl: web.ThumbnailUrl,
				}
				store.EXPECT().
					CreateWeb(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(web, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPostWeb(t, recorder.Body, web)
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				matchWeb := db.Web{
					ID:           web.ID,
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Url,
					ThumbnailUrl: "",
				}
				requireBodyMatchPostWeb(t, recorder.Body, matchWeb)
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "URLStatusBadRequest",
			body: gin.H{
				"url": web.Url,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				httpmock.RegisterResponder("GET", web.Url,
					httpmock.NewStringResponder(
						http.StatusBadRequest,
						"",
					),
				)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
						fmt.Sprintf(`<html>
						<head>
						<meta property="og:title" content="%s" />
						<meta property="og:image" content="%s" />
						</head>
						<body></body>
						</html>`,
							web.Title,
							web.ThumbnailUrl,
						),
					),
				)

				arg := db.CreateWebParams{
					UserID:       web.UserID,
					Url:          web.Url,
					Title:        web.Title,
					ThumbnailUrl: web.ThumbnailUrl,
				}
				store.EXPECT().
					CreateWeb(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Web{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			tc.checkResponse(recorder)
		})
	}

}

func randomWeb(t *testing.T, userID uuid.UUID) db.Web {
	id, err := uuid.NewRandom()
	require.NoError(t, err)

	ThumbnailURL := util.RandomThumbnailUrl()
	return db.Web{
		ID:           id,
		UserID:       userID,
		Url:          util.RandomUrl(),
		Title:        util.RandomName(),
		ThumbnailUrl: ThumbnailURL,
	}
}

func requireBodyMatchPostWeb(t *testing.T, body *bytes.Buffer, web db.Web) {
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
