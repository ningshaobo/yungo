package models

import (
	dmModels "com.my/dmSvrWeb/models"
	myutils "com.my/utils"
	"encoding/json"
	beegologs "github.com/astaxie/beego/logs"
	"net/http"
)

type UsersRespBody struct {
	Username string
}

/**
* @param  反向代理 业务逻辑
 */
func UserLoginModify(resp *http.Response) error {
	beegologs.Debug("proxy 后 modify 操作")

	if resp.StatusCode < 400 {
		buf, err := ProxyDrainBody(resp)
		if err != nil {
			return err
		}
		var user dmModels.User
		json.Unmarshal(buf.Bytes(), &user)
		beegologs.Debug("user %v", user)

		cluArray := GetAllClusters()
		var cluList []*dmModels.Cluster
		for _, clu := range cluArray {
			var cc dmModels.Cluster
			myutils.CommconvParamString(&clu, &cc)
			//			cc.Id = clu.Id
			//			cc.Name = clu.Name
			cluList = append(cluList, &cc)
		}
		for _, ttu := range user.TenantTypeUsers {
			ttu.Tenant.Clusters = cluList
		}

		newBody, err := json.Marshal(&user)
		if err != nil {
			return err
		}
		// 重新定义 response
		buf.Reset()
		buf.Write(newBody)

		// 重新赋值 response body
		ProxySetBody(resp, buf)
	}

	return nil
}
