package ogin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kms9/publicyc/pkg/conf"
)

// YAPIMock yapi mock数据
func (s *Server) YAPIMock() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := fmt.Sprintf("https://xxx.xxxx.tv/mock/%s%s", conf.Detail().GetString("yc.server.http.yapiProject"), c.Request.URL.Path)
		req, _ := http.NewRequest(c.Request.Method, url, nil)
		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var data interface{}
		_ = json.Unmarshal(body, &data)
		c.JSON(200, data)
	}
}
