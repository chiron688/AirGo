package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ppoonk/AirGo/api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/service"
	"github.com/ppoonk/AirGo/utils/response"
	"time"
)

// 获取用户动态路由
func GetMenuList(ctx *gin.Context) {
	uIdInt, _ := api.GetUserIDFromGinContext(ctx)
	//查cache
	userRouteList, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_USER_MENU_LIST_BY_ID, uIdInt))
	if ok {
		response.OK("GetMenuList success", userRouteList, ctx)
		return
	}
	//查询uId对应的角色
	roleIds, err := service.RoleSvc.FindRoleIdsByuId(uIdInt)
	if err != nil {
		global.Logrus.Error(err.Error())
		response.Fail("GetMenuList error:"+err.Error(), nil, ctx)
		return
	}
	menuIds, err := service.MenuSvc.GetMenuIdsByRoleIds(roleIds)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenuIdsByRoleIds error:"+err.Error(), nil, ctx)
		return
	}
	menuSlice, err := service.MenuSvc.GetMenusByMenuIds(menuIds)
	if err != nil {
		global.Logrus.Error(err)
		response.Fail("GetMenusByMenuIds error:"+err.Error(), nil, ctx)
		return
	}
	route := service.MenuSvc.GetMenus(menuSlice)
	global.LocalCache.Set(fmt.Sprintf("%s%d",
		constant.CACHE_USER_MENU_LIST_BY_ID, uIdInt),
		*route,
		constant.CAHCE_USER_MENUS_TIMEOUT*time.Minute)
	response.OK("GetMenuList success", route, ctx)
}
