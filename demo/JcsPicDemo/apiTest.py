import os
from JointCloudStorage import *

if __name__ == '__main__':
    ak = "b8dda4fa857544729aa87164ca7c5869"
    sk = "884abbf78a9c4ea7b1ff73e3977b1d95"
    endpoint_hohhot = "http://jsi-aliyun-hohhot.jointcloudstorage.cn/"
    auth = Auth(ak, sk)
    bucket = Bucket(auth, endpoint_hohhot)
    state = State(auth, endpoint_hohhot)

    info = state.get_server_info()
    print("服务器状态：", info)

    storage_state = state.get_storage_state()
    print("存储使用情况：", storage_state)

    storage_plan = state.get_storage_plan()
    print("存储方案：", storage_plan)

    files = bucket.get_object_list("/")  # 获取文件列表
    print(files)

    with open("./requirements.txt") as file:
        bucket.put_object("/python/r.txt", file)  # 上传文件
        file.seek(0, 0)
        task_id = bucket.put_object_async("/python/r2.txt", file)  # 异步上传
        print("上传 task id:", task_id)

    files = bucket.get_object_list("/python/")  # 获取以 /python/ 为前缀的文件对象
    print(files)

    if not os.path.exists("./tmp/"):
        os.makedirs("./tmp/")

    c = bucket.get_object("/python/r.txt")  # 下载文件
    with open("./tmp/t.txt", "wb+") as f:
        f.write(c)

    bucket.delete_object("/python/r.txt")  # 删除文件
    files = bucket.get_object_list("/python/")
    print(files)