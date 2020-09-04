
开发日志

- [x] 命令行程序
- [x] Postgres数据库对接
- [x] 配置数据读取和入库
- [x] 爬虫引擎程序
- [x] 缓存程序
- [x] 分布式调度程序
- [x] 日志处理
- [ ] 丰富内置插件：加密，
- [ ] 结果推送webhook
- [ ] 消息通知：钉钉，微信，QQ，企业微信，邮箱，Webhook等
- [ ] Block检测机制
- [ ] 统计功能
- [ ] 仪表盘
- [x] 前端UI
- [x] 前端功能细化完善
- [x] 使用文档编写
- [x] 插件编写指南
- [x] 函数插件开发
- [x] 配置调试功能
- [x] ajax配置支持
- [ ] 链接，结果去重
- [x] 结束状态判定
- [x] 任务状态控制
- [x] 配置可打包上传下载
- [x] 插件表移除slot
- [x] project节点亲和标签
- [x] redis订阅事件
- [x] env配置
- [x] 配置校验
- [x] 并发控制
- [x] 处理中的queue状态中断处理
- [x] 删除前置判断
- [x] docker集成
- [x] 配置travisCI
- [x] 任务完成清理缓存
- [ ] stage增加remark
- [x] field html类型
- [x] worker精简，去掉redis和postgres
- [x] worker支持travisCI识别和定时关闭
- [x] managerUrl的schema处理
- [ ] t_queue和t_result的分库分表
- [x] worker自动将自己的id加入label
- [ ] 模板引擎渲染结果
- [ ] 资源下载
- [ ] http代理
- [ ] 根据sql展示task进度图select a.stage_name, count(*) from t_queue a where a.task_id = 31 and status=0 group by a.stage_name
- [x] 任务删除功能，同时清理数据
- [ ] postgres压力太大，想办法降低
- [x] 压缩导出文件
- [ ] 用户管理
- [x] 日志乱码
- [x] node_affinity改为list类型
- [x] 无法完成任务：task already shutdown   无法完成任务：task not exists
- [ ] restapi打印错误日志
- [x] 支持xpath
- [ ] 任务调度节奏性卡顿问题













