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


### 鉴权Authorization生成方法

#### 用户加工request：(accessKey, secretKey, http.request) -> (http.request, error)

1. 用户传入的request，需要包含url，method，对post还需要放好body，secretKey与accessKey为32字节，为十六进制数字

2. 将目前时间取出，记作time，request的url和method取出，将其拼凑成为stringToSign，格式如下，其中url不包含域名和端口

   ```http
   Method:GET\r\n
   URL:/photos/puppy.jpg\r\n
   Time:1982738932754\r\n
   ```

3. 将stringToSign和secretKey进行sha_3加密，其中sha_3中先放入secretKey，并在吸收态进行轮转算法，然后再放入stringToSign，继续吸收并轮转，最终输出48字节的digest

4. 将digest转换为base64(url)格式，记作signature，为64字节

5. 生成auth字段，格式为`JCS(accessKey:signature)`

6. 在request中放入两个header，格式如下：

   |      KEY      | VALUE |
   | :-----------: | :---: |
   |     Time      | time  |
   | Authorization | auth  |



#### 请求验证:

1. 得到的请求格式如下：

   ```http
   GET /photos/puppy.jpg HTTP/1.1
   Time: 1982738932754
   Authorization: JCS(accessKey:signature)
   ```

2. 将Authorization中的accessKey和signature解析出来，其中accessKey为32字节uuid格式，signature为64字节base64(url)格式

3. 从数据库中根据accessKey取出对应的secretKey，并将signature用base64(url)解码为originSign

4. 将method，url，time取出，按照用户发出的方法拼凑为stringToSign，同样的sha_3方法加密输出摘要，记作sign

5. 将originSign与sign相比，若一致则通过，否则丢弃



#### 示例请求：

|            AccessKey             |            secretKey             |
| :------------------------------: | :------------------------------: |
| 41b320e3e75b4adbbeeda73540133757 | f841583fae584a4faf39ad769687f78c |

```http
PUT /jsiTest.txt HTTP/1.1
Time: 1625746465
Authorization: JCS(41b320e3e75b4adbbeeda73540133757:5plNpuX5c-Xb-N5DAoJ6-irPi6xKPDQHxI4wTnhWPrVLvWXyNk9NvH8ngRe5IvuG)
```

