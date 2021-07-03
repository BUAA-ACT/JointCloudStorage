## JcsPan 通用 API 规范 v1.0

## Object 相关

### PutObject 上传文件

**Request Syntax**

```
PUT /Key+ HTTP/1.1
Host: 
AccessKey: ak
SecretKey: sk
Authorization: JSI lksdfj2u3q894rf:fa;kldsfj28394urkdlvksadfkj
Content-Length: ContentLength
Content-MD5: ContentMD5
Time: 

Body
```
**Reply Syntax**

```http request
HTTP/1.1 200
```

### DeleteObject 删除文件

```
DELETE /Key+ HTTP/1.1
Host: 
AccessKey: ak
SecretKey: sk
Authorization: 认证字符串
```

### ListObject 获取文件列表

```
GET /KeyPrefix?max-keys=MaxKeys&start-after=StartAfter HTTP/1.1
Host: 
AccessKey: ak
SecretKey: sk
Authorization: 认证字符串
```

### GetObject 获取文件

```
GET /Key HTTP/1.1 
Host: 
AccessKey: ak
SecretKey: sk
Authorization: 认证字符串
```

```http request
HTTP/1.1 200
Content-Length: ContentLength
Last-Modified: LastModified

Body
```