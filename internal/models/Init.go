package models

import (
	"github.com/jeevic/lego/components/mongo"
	"github.com/jeevic/lego/pkg/app"
)

////初始化mongo
func InitMongo() {
	cfg := app.App.GetConfiger()

	//判断是否有配置
	if !cfg.IsSet("mongo") {
		return
	}
	var instances map[string]interface{}
	var prefix string
	//判断是否多实例
	if cfg.IsSet("mongo.type") && app.IsMultiInstance(cfg.GetString("mongo.type")) {
		instances = cfg.GetStringMap("mongo.instance")
		prefix = "mongo.instance."
	}

	for instance := range instances {
		pre := prefix + instance + "."
		setting := mongo.Setting{}
		if cfg.IsSet(pre + "uri") {
			setting.Uri = cfg.GetString(pre + "uri")
		}
		if cfg.IsSet(pre + "hosts") {
			setting.Hosts = cfg.GetString(pre + "hosts")
		}
		if cfg.IsSet(pre + "replset") {
			setting.ReplSet = cfg.GetString(pre + "replset")
		}
		if cfg.IsSet(prefix + "username") {
			setting.Username = cfg.GetString(pre + "username")
		}
		if cfg.IsSet(prefix + "password") {
			setting.Password = cfg.GetString(pre + "password")
		}
		if cfg.IsSet(prefix + "max_pool_size") {
			setting.MaxPoolSize = cfg.GetUint64(pre + "max_pool_size")
		}
		if cfg.IsSet(prefix + "min_pool_size") {
			setting.MinPoolSize = cfg.GetUint64(pre + "min_pool_size")
		}
		if cfg.IsSet(prefix + "max_idle_time") {
			setting.MaxIdleTime = cfg.GetInt(prefix + "max_idle_time")
		}
		if cfg.IsSet(prefix + "read_preference") {
			setting.ReadPreference = cfg.GetString(pre + "read_preference")
		}

		err := mongo.Register(instance, setting)
		if err != nil {
			app.App.GetLogger().Fatalf("[init] mongo instance:%s  error:%s", instance, err.Error())
			continue
		}

		app.App.GetLogger().Infof("[init] mongo instance:%s set !", instance)
	}
	app.App.GetLogger().Info("[init] mongo component complete !")
}
