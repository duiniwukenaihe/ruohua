package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chenhqchn/ruohua/server/utils/config"
	"github.com/pkg/errors"
)

var Enforcer *casbin.Enforcer

func InitEnforcer() error {
	adapter, err := gormadapter.NewAdapterByDB(config.DG())
	if err!= nil || adapter == nil{
		return errors.Wrap(err, "Failed to initing adapter DB")
	}

	Enforcer, err = casbin.NewEnforcer("config/casbin_model.conf", adapter)
	if err != nil {
		return errors.Wrap(err, "failed to creating enforcer")
	}
	return nil
}

