## jcs-sdk Joint Cloud Storage python sdk

云际网盘 python sdk

Joint Cloud Storage 是一个分布式云际存储系统。多个云服务访问方式统一，实现互操作性。
在此基础上，通过统一定义的接口互相协作，形成云际存储系统。

### 使用方法：

```python
import os
from JointCloudStorage import *

if __name__ == '__main__':
    access_key = "47bed2954a3647fdbc7a3364778c388f"
    secret_key = "2d61f1ec5d834af6ad7150d756fa33e5"
    auth = Auth(access_key,secret_key)
    bucket = Bucket(auth, "http://127.0.0.1:8085") # 指定 endpoint 地址
    files = bucket.get_object_list("/") # 获取文件列表
    print(files)

    with open("../requirements.txt") as file:
        bucket.put_object("/python/r.txt", file) # 上传文件

    files = bucket.get_object_list("/python/") # 获取以 /python/ 为前缀的文件对象
    print(files)

    if not os.path.exists("./tmp/"):
        os.makedirs("./tmp/")

    c = bucket.get_object("/python/r.txt") # 下载文件
    with open("./tmp/t.txt", "wb+") as f:
        f.write(c)

    bucket.delete_object("/python/r.txt") # 删除文件
    files = bucket.get_object_list("/python/")
    print(files)
```