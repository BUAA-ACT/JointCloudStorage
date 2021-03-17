# Transporter 设计文档

## 功能描述

1. 用户上传文件暂存
   将用户待上传文件暂存到本地存储，并上传到用户对应的存储后端
   
2. 用户文件下载链接生成
   响应用户的文件下载请求，生成临时下载链接。对于普通文件，返回一个 s3 的临时 Presigned  GET Url，对于纠删码文件，返回一个由 Transporter 负责响应的 Url，GET 该 Url 返回已经拼装好的整个文件流

3. 存储迁移

   transporter 定时从数据库中获取任务，完成待迁移存储后端的数据迁移

## 系统架构

![2021-03-08未命名文件3pSlJy](http://media.sumblog.cn/uPic/2021-03-08未命名文件3pSlJy.png)

- `Router`: 对外提供 http API，响应调用请求
- `TaskDatabase`: 任务库，存储新建任务，跟踪任务状态
- `StorageDatabase`：存储库，将 JcsPan 对象路径，映射为可操作的 `storage client`
- `storage client`: 提供不同云存储客户端的统一抽象
- `minio client`: 使用 minio 库对 Amazon S3 类型的对象存储进行操作



## 接口文档

```json
{
  "TaskType": "Upload",
   "Uid": "12",
   "Sid": "sdjfsdjfsadjkf21",
   "DestinationPath":"/path/to/upload/",
   "StoragePlan":{
      "StorageMode": "Replica",
      "Clouds": [
         {
            "ID": "txyun-chongqing"
         },
         {
            "ID": "aliyun-beijing"
         }
      ]
   }
}
```

### 1. Upload

- 用户文件上传接口，用于通过 HTML 表单上传的方式将文件上传到用户网盘

- Upload 操作通过验证请求头中 cookies session id 来判断操作的合法性

请求语法：

```
POST /upload/Path/to/MyFilename.jpg HTTP/1.1
Host: transporter.host
Content-Length：ContentLength
Content-Type: multipart/form-data; boundary=9431149156168
Cookie: sid=usersid
--9431149156168
Content-Disposition: form-data; name="file"; filename="MyFilename.jpg"
Content-Type: image/jpeg
file_content
--9431149156168
```

参数说明:

> Upload 的消息实体通过多重表单格式（multipart/form-data）编码，在 Upload 操作中，参数通过消息体中的表单域传递
>
> 文件的上传路径直接通过 POST 的请求路径指定，如希望上传到用户网盘的 /Path/to/MyFilename.jpg 路径，则需要将 POST 的请求路径设置为：`/upload/Path/to/MyFilename.jpg`

| 名称 | 类型   | 是否必选 | 描述                                                         |
| ---- | ------ | -------- | ------------------------------------------------------------ |
| file | 字符串 | 必选     | 文件或文本内容。浏览器会自动根据文件类型来设置Content-Type，一次只能上传一个文件 |

返回示例：

```
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 0  
Server: JcsPanTransporter
x-transporter-task-id: 5C2CBC8655718B5911EF4535
```

响应头：

| 名称                  | 类型   | 示例值                   | 描述                                              |
| --------------------- | ------ | ------------------------ | ------------------------------------------------- |
| x-transporter-task-id | 字符串 | 5C2CBC8655718B5911EF4535 | 上传的文件同步到云存储的任务 id，用于跟踪任务进度 |

### 2. GetFileTmpUrl

- 用户获取文件临时下载链接接口

请求语法：

```
GET /jcspan/Path/to/MyFilename.jpg HTTP/1.1
Host: transporter.host
Cookie: sid=usersid
```

参数说明：

> 请求获取的文件直接通过 GET 的请求路径指定

返回示例：

```
HTTP/1.1 200 OK
Server: JcsPanTransporter
Conternt-Type: text/plain; charset=utf-8

https://jcspan-aliyun-bj-test.oss-cn-beijing.aliyuncs.com/Path/to/MyFilename.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=LTAI4G3PCfrg7aXQ6EvuDo25%2F20210305%2Foss-cn-beijing%2Fs3%2Faws4_request&X-Amz-Date=20210305T022014Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&response-content-disposition=attachment%3B%20filename%3D%22your-filename.txt%22&X-Amz-Signature=c705a726732e9dd650986cdd56e78226996d43df3872ee225e9424b2d0fbdebb
```

### 3. CreateTask

- 任务创建接口，供其他后端服务调用

请求语法：

```
POST /task/taskType HTTP/1.1
Host: transporter.host
Content-Length：ContentLength
Content-Type: multipart/form-data; boundary=9431149156168
--9431149156168
Content-Disposition: form-data; name="srcpath"

path/to/
--9431149156168
Content-Disposition: form-data; name="dstpath"

dst/to/
--9431149156168
Content-Disposition: form-data; name="sid"

user-session-id
--9431149156168--
```

参数说明：

> CreateTask 的消息实体通过多重表单格式（multipart/form-data）编码，在 Upload 操作中，参数通过消息体中的表单域传递
>
> 创建的任务流行径直接通过 POST 的请求路径指定，如创建类型为 `taskType` 的任务，则请求路径应该为 `/task/taskType`

| 名称     | 类型   | 是否必选 | 描述                                                         |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| srcpath  | string | 必选     | 任务的作用路径，创建的任务对 srcpath 路径下的文件进行处理    |
| dstpath  | string | 可选     | 任务的目的路径，任务完成后，生成的目标文件在网盘中的存储位置 |
| sid      | string | 必选     | 用户 session id，用于用户身份认证                            |
| taskType | string | 必选     | 该参数通过 POST 请求的路径传入，目前已经实现的 taskType 有： `simplesync` |

返回示例：

```
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 0  
Server: JcsPanTransporter
x-transporter-task-id: 5C2CBC8655718B5911EF4535
```

响应头：

| 名称                  | 类型   | 示例值                   | 描述                            |
| --------------------- | ------ | ------------------------ | ------------------------------- |
| x-transporter-task-id | 字符串 | 5C2CBC8655718B5911EF4535 | 生成的Task id，用于跟踪任务进度 |



## 数据表设计

### 1. Task

| 字段名          | 描述                     | 类型   | 示例                     |
| --------------- | ------------------------ | ------ | ------------------------ |
| tid             | 任务 id                  | int    | 385                      |
| taskType        | 任务类型                 | enum   | USER_UPLOAD_SIMPLE       |
| state           | 任务状态                 | enum   | WAITING                  |
| startTime       | 任务开始时间             | time   | 2021-03-06 18:00:01      |
| sid             | session id               | string | JIOWEJ238HFFQ89          |
| sourcePath      | 任务操作对象路径         | string | /path/to/file.txt        |
| destinationPath | 任务生成对象路径（可选） | string | /jcspan/path/to/file.txt |