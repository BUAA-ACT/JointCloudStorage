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



send_index = 0

app_id = "24595163"
ak = "e3mXihhfyQPjY0gNoK97fj6v"
sk = "CczWdEiQs2gFd8w0cHn1EYeI1G4zCorL"
aip_client = AipImageProcess(app_id, ak, sk)


class Node(threading.Thread):
    def __init__(self, task_type, ak, sk, endpoint, interval, upstream_dict, output_dict, fallback_endpoint):
        threading.Thread.__init__(self)
        self.auth = Auth(ak, sk)
        self.bucket = Bucket(self.auth, endpoint)
        self.state = State(self.auth, endpoint)
        self.task_type = task_type
        self.interval = interval
        self.fail_times = 0
        self.upstream_dict = upstream_dict
        self.output_dict = output_dict
        self.fallback_endpoint = fallback_endpoint
        self.fallback_index = -1
        try:
            info = self.state.get_server_info()
            logger.warning(f"计算节点 {task_type} 初始化成功，存储接入点：{endpoint}, 节点信息：{info}")
        except Exception as e:
            logger.error(f"计算节点 {task_type} 初始化失败，存储接入点：{endpoint}, ")

    def run(self):
        logger.info("开始执行工作线程："+self.task_type)
        times = 0
        while True:
            #self.switch[self.task_type](self, self.bucket)
            self.work()
            times += 1
            # logger.info(f"{self.task_type} 成功执行 第 {times} 次")
            time.sleep(self.interval)

    def send(self):
        global send_index
        try:
            file_list = os.listdir("picture")
            with open("picture/" + file_list[send_index], 'rb') as f:
                dt = datetime.datetime.now()
                if f.name.endswith(".jpg"):
                    if self.bucket.put_object(self.output_dict + "pic" + dt.strftime("%Y-%m-%d-%H-%M-%S") + ".jpg", f):
                        logger.warning(f" STEP 1 : 文件上传成功 file upload OK!")
                    else:
                        print("upload Fail!")
            send_index = (send_index + 1) % len(file_list)
        except Exception as e:
            logger.warning("step1:" + str(e))
        pass

    def work(self):
        if self.task_type == "send":
            return self.send()
        try:
            file_list1 = self.bucket.get_object_list(self.upstream_dict)
            file_list2 = self.bucket.get_object_list(self.output_dict)
            set1 = set()
            set2 = set()
            if file_list1:
                set1 = set(self.getFileNames(file_list1))
            if file_list2:
                set2 = set(self.getFileNames(file_list2))
            difference = list(set1.difference(set2))
            if difference:
                logger.info(f"{self.task_type} 节点检测到 {len(difference)} 张待处理图片")
                for filename in difference:
                    c = self.bucket.get_object(self.upstream_dict + filename)
                    try:
                        if self.task_type == "colorize":
                            res = aip_client.colourize(c)
                        elif self.task_type == "lar_en":
                            res = aip_client.imageQualityEnhance(c)
                        elif self.task_type == "con_en":
                            res = aip_client.contrastEnhance(c)
                        else:
                            res = {"image": str(base64.b64encode(c), "utf-8")}
                    except Exception as e:
                        logger.error(f"图像处理时错误，Error: {e}")
                    self.sendBytes(res, self.output_dict + filename)
                    logger.info(f"{self.task_type} 节点处理 {filename} 结果上传成功")
        except Exception as e:
            logger.error(f"云际存储连接出错 Error:{e}, 错误次数 {self.fail_times}")
            self.fail_times += 1
            if self.fail_times % 3 == 0:
                self.fallback_index += 1
                self.bucket = Bucket(self.auth, self.fallback_endpoint[self.fallback_index])
                self.state = State(self.auth, self.fallback_endpoint[self.fallback_index])
                logger.error(f"切换到备份节点：{self.fallback_endpoint[self.fallback_index]}")



    def sendBytes(self, res, path):
        if "image" in res:
            img = base64.b64decode(res['image'].encode())
            # 写入step2
            # with open("tmp/tmp.jpg",'rwb') as tmp:
            #     tmp.write(img)
            #     dt=datetime.datetime.now()
            #     bucket.put_object(dict2+"pic"+dt.strftime("%Y-%m-%d-%H-%M-%S"),tmp)
            dt = datetime.datetime.now()
            self.bucket.put_object(path, img)
        else:
            logger.warning(res['error_code'] + ":" + res['error_msg'] + " path: " + path)
        pass

    def getFileNames(self, fileList):
        names = []
        for file in fileList:
            fullName = file['Filename']
            name = os.path.basename(fullName)
            names.append(name)
            pass
        return names
        pass

    def clear_all(self):
        bucket = self.bucket
        files = bucket.get_object_list("")
        for f in files:
            fileName = f['Filename']
            logger.info("删除文件：" + fileName)
            bucket.delete_object(f['Filename'])
