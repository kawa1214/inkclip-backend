package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inkclip/backend/config"
	db "github.com/inkclip/backend/db/sqlc"
	"github.com/inkclip/backend/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := config.Config{
		TokenSecretKey:      util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
