package server

import (
	"context"
	. "github.com/cyberax/go-dd-service-base/visibility"
	"github.com/jmoiron/sqlx"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"xyzzy/gen/db"
	. "xyzzy/gen/xyzzy"
)

type XyzzyApi struct {
	Q        *db.Queries
	Database *sqlx.DB
}

var _ XyzzyData = &XyzzyApi{}

func (t *XyzzyApi) Ping(ctx context.Context, _ *PingRequest) (*PingOk, error) {
	_, err := t.Q.Ping(ctx)
	if err != nil {
		CL(ctx).Warn("healthcheck failed", zap.Error(err))
		return nil, twirp.InternalError("unhealthy")
	}

	return &PingOk{}, nil
}

func (t *XyzzyApi) translateErr(err error) error {
	if err.Error() == "sql: no rows in result set" {
		return twirp.NotFoundError("element not found")
	}
	return err
}
