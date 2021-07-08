## JcsPan 通用 API 规范 v1.0

云际网盘定义了一套 REST 风格的对象操作 API，通过 `Key` 来声明待操作的存储对象，使用鉴权秘钥对请求进行签名，生成鉴权 Token，来认证用户身份。

云际网盘内各节点采用对等方案设计，任意节点均提供了注册、存储方案构建、鉴权秘钥生成接口。

## Object 相关

### PutObject 上传文件

**请求语法**

```http
PUT /Key+ HTTP/1.1
Host: 云际网盘节点 Host
Authorization: 云际网盘鉴权 Token
Content-Length: 内容长度
Content-MD5: 内容 MD5 校验码
Time: 鉴权 Token 生成时间

Body 上传文件内容
```
**响应语法**

```http request
HTTP/1.1 200
```

### DeleteObject 删除文件

**请求语法**

```http
DELETE /Key+ HTTP/1.1
Host: 云际网盘节点 Host
Authorization: 云际网盘鉴权 Token
Time: 鉴权 Token 生成时间
```

**响应语法**

```http request
HTTP/1.1 200
```

### ListObject 获取文件列表

**请求语法**

```http
GET /KeyPrefix?max-keys=MaxKeys&start-after=StartAfter HTTP/1.1
Host: 云际网盘节点 Host
Authorization: 云际网盘鉴权 Token
Time: 鉴权 Token 生成时间
```

**响应语法**

```http request
HTTP/1.1 200

[
{
	"Filename": "文件 key",
	"Owner": "文件所有者",
	"Size": 文件大小,
	"LastModified": "2021-07-04T03:07:14.044Z", // 最后修改时间
	"SyncStatus": "Done", // 同步状态
},
{
 ...
}
]
```

### GetObject 获取文件

**请求语法**

```http
GET /Key HTTP/1.1 
Host: 云际网盘节点 Host
Authorization: 云际网盘鉴权 Token
Time: 鉴权 Token 生成时间
```

**响应语法**

```http
HTTP/1.1 200
Content-Length: 正文内容长度
Last-Modified: 最后修改时间

Body
```

