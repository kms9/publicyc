# 日志库

## 使用方式

// 单例
log := New("info", "developemnt")

// 每次请求实例化一个
l := log.With(&logger.BaseContentInfo{
    TraceID: "39b99823-9f9b-4afc-a9a2-2ae33d90a690",
    UID:     "5eb5200ec23e490001cf405a",
})

l.Warn("这是测试", "warn")
l.Warnf("这是测试 %s %d", "warnf", 1)
l.Info("这是测试 info")
l.Infof("这是测试 %s", "infof")
l.Error("这是测试", "error")
l.Errorf("这是测试 %s %d", "errorf", 1)
l.Debug("这是测试 debug")
l.Debugf("这是测试 %s", "debugf")
l.Notice(&hook.DingMsg{
    Msgtype: "text",
    Text:    &hook.DingMsgText{Content: "test-3214321432-342143214321-43214231"},
}) // 不配置不会真实通知
l.Panic("这是测试 panic")

// GinLogMiddle gin 日志中间件
func GinLogMiddle(appName string, goEnv string) gin.HandlerFunc

// GetBaseByContext 通过gin context 获取基础内容
func GetBaseByContext(c *gin.Context) *logger.BaseContentInfo
