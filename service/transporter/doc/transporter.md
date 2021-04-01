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

## 测试环境

- 在开发过程中可以使用测试环境来对 transporter 提供的接口进行测试
- 在进行测试前，请确保请求发起方位于 ACT 网络内，通过 ping 192.168.105.2 来验证请求发起方与 transporter 的连通性
- 测试环境地址

  1. http://transporter.act.sumblog.cn/
  2. http://192.168.105.2:9648/
- 测试接口 1 仅为接口 2 的反向代理，两者并无区别，若无法使用，请联系张俊华
- 测试时，cloudID 请填写 `aliyun-beijing` 

## 历史版本

- [v0.3.24-transporter-alpha](http://gitlab.act.buaa.edu.cn/jointcloudstorage/jcspan/-/releases/v0.3.24-transporter-alpha)

  首个 transporter alpha 版本

  - 支持 Uplaod、Download、Sync、Delete 任务
  - 支持用户端通过表单上传文件
  - 生成有超时时间的临时下载链接
  - 所有用户文件共用临时测试 bucket
  - 采用 jwt 生成签名 token 认证用户身份
  - 从配置文件读取配置
  - 使用 MongoDB 进行持久化存储 (MongoDB 服务器硬编码，暂不能修改）
  - 受限的云存储支持，不能容忍云服务失效

## 接口文档

### 1. Upload

- 用户文件上传接口，用于通过 HTML 表单上传的方式将文件上传到用户网盘
- Upload 操作通过验证 post 表单中的 token 来判断操作的合法性

请求语法：

```
POST /upload/Path/to/MyFilename.jpg HTTP/1.1
Host: transporter.host
Content-Length：ContentLength
Content-Type: multipart/form-data; boundary=9431149156168
--9431149156168
Content-Disposition: form-data; name="file"; filename="MyFilename.jpg"
Content-Type: image/jpeg
file_content
--9431149156168
Content-Disposition: form-data; name="token"

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXRoIjoiIiwiVWlkIjoidGVzdGVyIiwiVGlkIjoiNjA1YjI4YWRkOTNmZGM5ZjBlOWY4YTBlIiwiZXhwIjoxNjE2NjczMzI1LCJpc3MiOiJ0cmFuc3BvcnRlciJ9.MsaNNkZhKTV9e4Ek5DxuhVky8qtAqywWSdDYeH4527s
--9431149156168
```

参数说明:

> Upload 的消息实体通过多重表单格式（multipart/form-data）编码，在 Upload 操作中，参数通过消息体中的表单域传递
>
> 文件的上传路径直接通过 POST 的请求路径指定，如希望上传到用户网盘的 /Path/to/MyFilename.jpg 路径，则需要将 POST 的请求路径设置为：`/upload/Path/to/MyFilename.jpg`

| 名称  | 类型   | 是否必选 | 描述                                                         |
| ----- | ------ | -------- | ------------------------------------------------------------ |
| file  | 字符串 | 必选     | 文件或文本内容。浏览器会自动根据文件类型来设置Content-Type，一次只能上传一个文件 |
| token | 字符串 | 必选     | 文件上传  token， 该 token 通过创建上传task获得              |

返回示例：

```
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 0  
Server: JcsPanTransporter
```

### 2. GetFileTmpUrl

- 用户获取文件临时下载链接接口
- **该接口准备弃用，下载文件前端请与 httpserver 请求，由 httpserver 通过创建task的方式获取下载链接**

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

- 任务创建接口，供其他后端服务调用 **前端不应调用此接口**

请求语法：

```
POST /task HTTP/1.1
Host: transporter.host
Content-Length：ContentLength

{
  "TaskType": "Upload",
   "Uid": "tester",
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

参数说明：

> ~~CreateTask 的消息实体通过多重表单格式（multipart/form-data）编码，在 Upload 操作中，参数通过消息体中的表单域传递~~
>
> ~~创建的任务类型直接通过 POST 的请求路径指定，如创建类型为 `taskType` 的任务，则请求路径应该为 `/task/taskType`~~
>
> **所有参数直接通过 Json 结构体传递**

- Task Json 格式

  ```json
  {
    "TaskType": "Task 类型",
     "Uid": "用户 id",
     "Sid": "用户鉴权 sid",
     "SourcePath":"任务作用路径",
     "DestinationPath":"任务目的路径",
     "SourceStoragePlan":{},
     "DestinationStoragePlan":{}
  }
  ```

  | 名称                   | 值                                           | 是否必选 | 描述                                                         |
  | ---------------------- | -------------------------------------------- | -------- | ------------------------------------------------------------ |
  | TaskType               | "Upload" \| "Download" \| "Sync" \| "Delete" | 必选     | 任务类型，当前实现的任务类型有 Upload、Download、Sync、Delete |
  | SourcePath             | string                                       | 可选     | 任务的作用路径，创建的任务对 SourcePath 路径下（前缀）的文件进行处理 |
  | DestinationPath        | string                                       | 可选     | 任务的目的路径，任务完成后，生成的目标文件在网盘中的存储位置 |
  | Uid                    | string                                       | 必选     | 用户 user id，用户指定特定的用户                             |
  | SourceStoragePlan      | storagePlanStruct                            | 可选     | 任务作用路径的存储方案                                       |
  | DestinationStoragePlan | storagePlanStruct                            | 可选     | 任务路径的存储方案                                           |

- storagePlanStruct：

  ```json
  {
      "StorageMode": "EC",
      "Clouds": [],
      "N": 3,
      "K": 1
  }
  ```

  | 名称        | 值                | 是否必选 | 描述                                      |
  | ----------- | ----------------- | -------- | ----------------------------------------- |
  | StorageMode | "Replica" \| "EC" | 必选     | 存储方案类型  多副本：Replica、纠删吗：EC |
  | Clouds      | [] cloudStruct    | 必选     | 存储方案云服务提供商                      |
  | N           | int               | 可选     | 文件数据块分块数量                        |
  | K           | int               | 可选     | 文件校验块数量                            |

- cloudStruct：

  ```json
  {
  	"ID": "txyun-chongqing"	
  }
  ```

  | 名称 | 值     | 是否必选 | 描述          |
  | ---- | ------ | -------- | ------------- |
  | ID   | string | 必选     | 存储提供商 ID |



**请求示例：**

- 请求创建上传文件 Task (多副本模式）：

  ```json
  {
    "TaskType": "Upload",
     "Uid": "12",
     "Sid": "tttteeeesssstttt",
     "DestinationPath":"/path/to/upload/",
     "DestinationStoragePlan":{
        "StorageMode": "Replica",
        "Clouds": [
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-shanghai"
           }
        ]
     }
  }
  ```

- 请求创建上传文件 Task (纠删码模式）：

  ```json
  {
    "TaskType": "Upload",
     "Uid": "12",
     "DestinationPath":"/path/to/upload/",
     "DestinationStoragePlan":{
        "StorageMode": "EC",
        "Clouds": [
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-shanghai"
           },
           {
              "ID": "txyun-beijing"
           }
        ],
        "N": 2,
        "K": 1
     }
  }
  ```

- 请求创建下载文件 Task (纠删码模式）：

  ```json
  {
    "TaskType": "Download",
     "Uid": "tester",
     "SourcePath":"/path/to/jcspantest.txt",
     "SourceStoragePlan":{
        "StorageMode": "EC",
        "Clouds": [
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-shanghai"
           },
           {
              "ID": "txyun-beijing"
           }
        ],
        "N": 2,
        "K": 1
     }
  }
  ```

  请注意，纠删码模式下，云服务提供商数量必须等于 N+K，且按照云服务商在 Clouds 中出现的顺序，依次存储数据分块以及校验分块

- 创建同步任务 （从纠删码同步到多副本）：

  ```json
  {
    "TaskType": "Sync",
     "Uid": "tester",
     "DestinationPath":"tmp/test/sync/",
     "SourcePath": "tmp/test/sync/",
     "SourceStoragePlan":{
        "StorageMode": "EC",
        "Clouds": [
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           }
        ],
        "N": 2,
        "K": 1
     },
     "DestinationStoragePlan":{
        "StorageMode": "Replica",
        "Clouds": [
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           }
        ]
     }
  }
  ```

- 创建删除任务（删除整个目录）：

  ```json
  {
    "TaskType": "Delete",
     "Uid": "tester",
     "SourcePath": "tmp/test/sync/",
     "SourceStoragePlan":{
        "StorageMode": "Replica",
        "Clouds": [
           {
              "ID": "aliyun-beijing"
           },
           {
              "ID": "aliyun-beijing"
           }
        ]
     }
  }
  ```

**返回值示例：**

- 异步任务：如Download task (EC 模式)、Sync 

  ```
  {
  	"Code": 200,
  	"Msg": "Generate download task OK",
  	"Data": {
  		"Type": "tid",
  		"Result": "6061c72fc85f9f22ebac91e7"
  	}
  }
  ```

  唯一标识的 task id，用于跟踪任务进度

- 同步任务：

  如 Download （Replica 模式）

  ```
  {
  	"Code": 200,
	"Msg": "Generate download url OK",
  	"Data": {
		"Type": "url",
  		"Result": "http://jcspan-aliyun-bj-test.oss-cn-beijing.aliyuncs.com/tester/path/to/jcspantest.txt?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=LTAI4G3PCfrg7aXQ6EvuDo25%2F20210329%2Fus-east-1%2Fs3%2Faws4_request\u0026X-Amz-Date=20210329T122646Z\u0026X-Amz-Expires=1800\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=a27c353b0f0e0bcfd1b30435bd39e8110ba731a636667c0ef544df761670515b"
  	}
  }
  ```
  
  如 Upload 生成token：
  
  ```
  {
  	"Code": 200,
  	"Msg": "Create upload task OK",
  	"Data": {
  		"Type": "token",
  		"Result": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXRoIjoiIiwiVWlkIjoidGVzdGVyIiwiVGlkIjoiNjA2MWM2ODg1YTQ3YTlkYzNlNWMyMmJkIiwiZXhwIjoxNjE3MTA2OTUzLCJpc3MiOiJ0cmFuc3BvcnRlciJ9.ntEnf79hlKA7TF0pDvIm-kBZo8gPsoZ-eZio0aGsnCA"
  	}
  }
  ```
  

### 4. DownloadCacheFile

- 下载重建完成的本地文件

**请求语法：**

- GET 模式

```
GET /cache_file?token=<token> HTTP/1.1
Host: transporter.host
```

参数说明:

> 使用 token 同时进行用户身份和下载文件的指定，token 采用 JWT 方式进行编码，由 httpserver 从 file 数据库中读取

| 名称  | 类型   | 是否必选 | 描述                                          |
| ----- | ------ | -------- | --------------------------------------------- |
| token | 字符串 | 必选     | 文件下载  token， 该 token 由 httpserver 获取 |



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

## Docker 配置
账号：robot$pub+gitlab-ci
令牌：tVGegelWtspeQFGZYPivzIWP8w2vy8fj

```shell
docker build -t transporter:v1-dev .
```

在项目中标记镜像：
```shell
docker tag SOURCE_IMAGE[:TAG] harbor.sumblog.cn/pub/REPOSITORY[:TAG]
```
推送镜像到harbor：
```shell
docker push harbor.sumblog.cn/pub/REPOSITORY[:TAG]
```

运行 docker transporter
```shell
docker volume create transporter-vol
docker run -d -p 8083:8083 \
  --name transporter-dev-v1 \
  --mount source=transporter-vol,target=/data/transporter \
  --mount type=bind,source=/home/jcspan/transporter/transporter_config.json.docker,readonly \
  transporter[:TAG]
```
