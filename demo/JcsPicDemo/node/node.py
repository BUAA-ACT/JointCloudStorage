import argparse
import functools
import os
import datetime
import threading
import base64
import time

import JointCloudStorage

from .pic_haddle import AipClient
from .utils import add_arguments
from aip import AipImageProcess

from JointCloudStorage import Auth, Bucket, State
from logzero import logger

parser = argparse.ArgumentParser(description=__doc__)
add_arg = functools.partial(add_arguments, argparser=parser)
# yapf: disable
add_arg('type', str, "send", "Node type:[send,colorize,con_en,lar_en]")
add_arg('ak', str, None, "access key")
add_arg('sk', str, None, "secret key")
add_arg('endpoint', str, None, "service address,ip:port")

dict1 = "/step1/"
dict2 = "/step2/"
dict3 = "/step3/"
dict4 = "/step4/"

output_dicts = [dict1, dict2, dict3, dict4]

send_index = 0

app_id = "24595163"
ak = "e3mXihhfyQPjY0gNoK97fj6v"
sk = "CczWdEiQs2gFd8w0cHn1EYeI1G4zCorL"
aip_client = AipImageProcess(app_id, ak, sk)


def send(bucket):
    buck = bucket
    global send_index
    try:
        file_list = os.listdir("picture")
        with open("picture/" + file_list[send_index], 'rb') as f:
            dt = datetime.datetime.now()
            if buck.put_object(dict1 + "pic" + dt.strftime("%Y-%m-%d-%H-%M-%S") + ".jpg", f):
                logger.warning(f" STEP 1 : 文件上传成功 file upload OK!")
            else:
                print("upload Fail!")
        send_index = (send_index + 1) % len(file_list)
    except Exception as e:
        logger.warning("step1:" + str(e))
    pass


def colorize(bucket):
    buck = bucket
    try:
        fileList1 = bucket.get_object_list(dict1)
        fileList2 = bucket.get_object_list(dict2)
        set1=set()
        set2=set()
        if fileList1:
            set1 = set(getFileNames(fileList1))
        if fileList2:
            set2 = set(getFileNames(fileList2))
        difference = list(set1.difference(set2))
        if difference:
            logger.info(f"colorize 节点检测到 {len(difference)} 张待处理图片")
            for filename in difference:
                c = buck.get_object(dict1 + filename)
                res = aip_client.colourize(c)
                # 发送文件
                sendBytes(res, dict2 + filename, bucket)
                logger.warning(f" STEP 2 : 彩色化图片 处理成功")
            logger.info(f"colorize 节点完成处理")
    except Exception as e:
        logger.warning("step2:" + str(e))
    pass


def contrast_enhance(bucket):
    buck = bucket
    try:
        fileList1 = bucket.get_object_list(dict2)
        fileList2 = bucket.get_object_list(dict3)
        set1=set()
        set2=set()

        if fileList1:
            set1 = set(getFileNames(fileList1))
        if fileList2:
            set2 = set(getFileNames(fileList2))

        difference = list(set1.difference(set2))
        if difference:
            logger.info(f"contrast 节点检测到 {len(difference)} 张待处理图片")
            for filename in difference:
                c = buck.get_object(dict2 + filename)
                res = aip_client.contrastEnhance(c)
                # 发送文件
                sendBytes(res, dict3 + filename, bucket)
                logger.warning(f" STEP 3 : 图像增强 处理成功")
            logger.info(f"contrast 节点完成处理")
    except Exception as e:
        logger.warning("step3:" + str(e))
    pass


def large_enhance(bucket):
    buck = bucket
    try:
        fileList1 = bucket.get_object_list(dict3)
        fileList2 = bucket.get_object_list(dict4)
        set1=set()
        set2=set()

        if fileList1:
            set1 = set(getFileNames(fileList1))
        if fileList2:
            set2 = set(getFileNames(fileList2))

        difference = list(set1.difference(set2))
        if difference:
            logger.info(f"large 节点检测到 {len(difference)} 张待处理图片")
            for filename in difference:
                c = buck.get_object(dict3 + filename)
                res = aip_client.imageQualityEnhance(c)
                # 发送文件
                sendBytes(res, dict4 + filename, bucket)
    except Exception as e:
        logger.warning("step4:", e)
    pass


def sendBytes(res, path, bucket):
    if "image" in res:
        img = base64.b64decode(res['image'].encode())
        # 写入step2
        # with open("tmp/tmp.jpg",'rwb') as tmp:
        #     tmp.write(img)
        #     dt=datetime.datetime.now()
        #     bucket.put_object(dict2+"pic"+dt.strftime("%Y-%m-%d-%H-%M-%S"),tmp)
        dt = datetime.datetime.now()
        bucket.put_object(path, img)
    else:
        logger.warning(res['error_code'] + ":" + res['error_msg'])
    pass

def getFileNames(fileList):
    names=[]
    for file in fileList:
        fullName = file['Filename']
        name = os.path.basename(fullName)
        names.append(name)
        pass
    return names
    pass


def clear_all(node):
    bucket = node.bucket
    for d in output_dicts:
        files = bucket.get_object_list(d)
        for f in files:
            fileName = f['Filename']
            logger.info("删除文件：" + fileName)
            bucket.delete_object(f['Filename'])



switch = {
    "send": send,
    "colorize": colorize,
    "con_en": contrast_enhance,
    "lar_en": large_enhance,
}


class Node(threading.Thread):
    def __init__(self, task_type, ak, sk, endpoint, interval):
        threading.Thread.__init__(self)
        self.auth = Auth(ak, sk)
        self.bucket = Bucket(self.auth, endpoint)
        self.state = State(self.auth, endpoint)
        self.task_type = task_type
        self.interval = interval
        info = self.state.get_server_info()
        logger.warning(f"计算节点 {task_type} 初始化成功，存储接入点：{endpoint}, 节点信息：{info}")



    def run(self):
        logger.info("开始执行工作线程："+self.task_type)
        times = 0
        while True:
            switch[self.task_type](self.bucket)
            times += 1
            # logger.info(f"{self.task_type} 成功执行 第 {times} 次")
            time.sleep(self.interval)

