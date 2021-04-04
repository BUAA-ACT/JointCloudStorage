# 部署&运维文档

## transporter 
### transporter 构建方法

1. 进入 transporter 目录

   ```
   cd service/transporter
   ```

2. 拉取 go mod 依赖

   ```
   go mod tidy
   ```

3. 使用 make 进行构建

   ```
   make build
   ```

构建生成的可执行文件为：`/service/transporter/build/bin/transporter`

### transporter 执行

1. 在命令行中启动 transporter：

   ```
   ./transporter -c path/to/transporter_config.json
   ```

   启动时需要从硬盘中读取 json 格式的配置文件，使用 -c 参数指定，若不指定，则默认从执行目录下读取 transporter_config.json

   在仓库中给出了示例配置文件 `transporter_config.json.sample`

   运行后，程序会在前台驻留，输出日志

2. 访问 ip:port 查看 transporter 运行状态，正常运行时应该输出 transporter 的版本号：

    ![image-20210404111237161](http://media.sumblog.cn/uPic/2021-04-04image-20210404111237161mp9sy0.png)

