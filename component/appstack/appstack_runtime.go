// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package appstack

import (
	"context"
	"fmt"
	"os"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rulex/glogger"
	"github.com/hootrhino/rulex/typex"
)

var __DefaultAppStackRuntime *AppStackRuntime

func InitAppStack(re typex.RuleX) *AppStackRuntime {
	__DefaultAppStackRuntime = &AppStackRuntime{
		RuleEngine:   re,
		Applications: make(map[string]*typex.Application),
	}
	return __DefaultAppStackRuntime
}
func AppRuntime() *AppStackRuntime {
	return __DefaultAppStackRuntime
}

/*
*
* 加载本地文件到lua虚拟机, 但是并不执行
*
 */
func LoadApp(app *typex.Application) error {
	bytes, err := os.ReadFile(app.Filepath)
	if err != nil {
		return err
	}
	// 重新读
	app.VM().DoString(string(bytes))
	// 检查函数入口
	AppMainVM := app.VM().GetGlobal("Main")
	if AppMainVM == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if AppMainVM.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	// 抽取main
	fMain := *AppMainVM.(*lua.LFunction)
	app.SetMainFunc(&fMain)
	// 加载库
	LoadAppLib(app, __DefaultAppStackRuntime.RuleEngine)
	// 加载到内存里
	__DefaultAppStackRuntime.Applications[app.UUID] = app
	return nil
}

/*
* 此时才是真正的启动入口:
* 启动 function Main(args) --do-some-thing-- return 0 end
*
 */
func StartApp(uuid string) error {
	app, ok := __DefaultAppStackRuntime.Applications[uuid]
	if !ok {
		return fmt.Errorf("app not exists:%s", uuid)
	}
	if app.AppState == 1 {
		return fmt.Errorf("app not already started:%s", uuid)
	}
	// args := lua.LBool(false) // Main的参数，未来准备扩展
	ctx, cancel := context.WithCancel(typex.GCTX)
	app.SetCnC(ctx, cancel)
	go func(app *typex.Application) {
		defer func() {
			glogger.GLogger.Debug("App exit:", app.UUID)
			if err := recover(); err != nil {
				glogger.GLogger.Error("App recover:", err)
			}
		}()
		glogger.GLogger.Debugf("Ready to run app:%s", app.UUID)
		app.AppState = 1
		err := app.VM().CallByParam(lua.P{
			Fn:      app.GetMainFunc(),
			NRet:    1,
			Protect: true, // If ``Protect`` is false,
			// GopherLua will panic instead of returning an ``error`` value.
			Handler: &lua.LFunction{
				GFunction: func(*lua.LState) int {
					return 1
				},
			},
		}, lua.LBool(false))
		if err != nil {
			glogger.GLogger.Error("normal app.VM().CallByParam error:", err)
		}
		app.AppState = 0
	}(app)
	glogger.GLogger.Info("App started:", app.UUID)
	return nil
}

/*
*
* 从内存里面删除APP
*
 */
func RemoveApp(uuid string) error {
	if app, ok := __DefaultAppStackRuntime.Applications[uuid]; ok {
		app.Remove()
		delete(__DefaultAppStackRuntime.Applications, uuid)
	}
	glogger.GLogger.Info("App removed:", uuid)
	return nil
}

/*
*
* 停止应用并不删除应用, 将其进程结束，状态置0
*
 */
func StopApp(uuid string) error {
	if app, ok := __DefaultAppStackRuntime.Applications[uuid]; ok {
		app.Remove()
		delete(__DefaultAppStackRuntime.Applications, uuid)
	}
	glogger.GLogger.Info("App removed:", uuid)
	return nil
}

/*
*
* 更新应用信息
*
 */
func UpdateApp(app typex.Application) error {
	if oldApp, ok := __DefaultAppStackRuntime.Applications[app.UUID]; ok {
		oldApp.Name = app.Name
		oldApp.Version = app.Version
		glogger.GLogger.Info("App updated:", app.UUID)
		return nil
	}
	return fmt.Errorf("update failed, app not exists:%s", app.UUID)

}
func GetApp(uuid string) *typex.Application {
	if app, ok := __DefaultAppStackRuntime.Applications[uuid]; ok {
		return app
	}
	return nil
}

/*
*
* 获取列表
*
 */
func AllApp() []*typex.Application {
	return ListApp()
}
func ListApp() []*typex.Application {
	apps := []*typex.Application{}
	for _, v := range __DefaultAppStackRuntime.Applications {
		apps = append(apps, v)
	}
	return apps
}

func Stop() {
	for _, app := range __DefaultAppStackRuntime.Applications {
		glogger.GLogger.Info("Stop App:", app.UUID)
		app.Stop()
		glogger.GLogger.Info("Stop App:", app.UUID, " Successfully")
	}
	glogger.GLogger.Info("Appstack stopped")

}
func GetRuleX() typex.RuleX {
	return __DefaultAppStackRuntime.RuleEngine
}