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
from node_client import NodeClient

parser = argparse.ArgumentParser(description=__doc__)
add_arg = functools.partial(add_arguments, argparser=parser)
# yapf: disable
add_arg('type', str, "send", "Node type:[send,colorize,con_en,lar_en]")
add_arg('ak', str, None, "access key")
add_arg('sk', str, None, "secret key")
add_arg('endpoint', str, None, "service address,ip:port")

send_index = 0


class NodeState(object):
    def __init__(self, finish_num, fail_num, endpoint_name, endpoint_address, state):
        self.finish_num = finish_num
        self.fail_num = fail_num
        self.endpoint_name = endpoint_name
        self.endpoint_address = endpoint_address
        self.file_processing = ""


class Node(threading.Thread):
    def __init__(self, task_type, ak, sk, endpoint, interval, upstream_dict, output_dict, fallback_endpoint,
                 endpoint_name_dict):
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
        self.endpoint_name_dict = endpoint_name_dict
        self.node_state = NodeState(0, 0, endpoint_name_dict[endpoint], endpoint, "init")

        try:
            info = self.state.get_server_info()
            logger.warning(f"计算节点 {task_type} 初始化成功，存储接入点：{endpoint}, 节点信息：{info}")
        except Exception as e:
            logger.error(f"计算节点 {task_type} 初始化失败，存储接入点：{endpoint}, ")

    def get_state(self) -> NodeState:
        return self.node_state

    def run(self):
        logger.info("开始执行工作线程：" + self.task_type)
        times = 0
        while True:
            self.work()
            times += 1
            time.sleep(self.interval)

    def send(self):
        global send_index
        try:
            file_list = os.listdir("picture")
            with open("picture/" + file_list[send_index], 'rb') as f:
                dt = datetime.datetime.now()
                if f.name.endswith(".jpg"):
                    self.node_state.file_processing = f.name
                    if self.bucket.put_object(self.output_dict + "pic" + dt.strftime("%Y-%m-%d-%H-%M-%S") + ".jpg", f):
                        logger.warning(f" STEP 1 : 文件上传成功 file upload OK!")
                    else:
                        print("upload Fail!")
                        self.node_state.fail_num += 1
            send_index = (send_index + 1) % len(file_list)
            self.node_state.finish_num += 1
        except Exception as e:
            logger.error(f"云际存储连接出错 Error:{e}, 错误次数 {self.fail_times}")
            self.fail_times += 1
            if self.fail_times % 3 == 0:
                self.fallback_index += 1
                self.node_state.fail_num += 1
                self.state = "ERROR"
                self.bucket = Bucket(self.auth, self.fallback_endpoint[self.fallback_index])
                self.state = State(self.auth, self.fallback_endpoint[self.fallback_index])
                self.node_state.endpoint_address = self.fallback_endpoint[self.fallback_index]
                self.node_state.endpoint_name = self.endpoint_name_dict[self.node_state.endpoint_address]
                logger.error(f"切换到备份节点：{self.fallback_endpoint[self.fallback_index]}")

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
                    file_lock = threading.Lock()
                    num_lock = threading.Lock()
                    state_lock = threading.Lock()
                    client = NodeClient(bucket=self.bucket, task_type=self.task_type, up_dict=self.upstream_dict,
                                        out_dict=self.output_dict, file_name=filename, state=self.state,
                                        file_lock=file_lock, num_lock=num_lock, state_lock=state_lock)
                    client.run()
                    client.join()
        except Exception as e:
            logger.error(f"云际存储连接出错 Error:{e}, 错误次数 {self.fail_times}")
            self.fail_times += 1
            if self.fail_times % 3 == 0:
                self.fallback_index += 1
                self.node_state.fail_num += 1
                self.state = "ERROR"
                self.bucket = Bucket(self.auth, self.fallback_endpoint[self.fallback_index])
                self.state = State(self.auth, self.fallback_endpoint[self.fallback_index])
                logger.error(f"切换到备份节点：{self.fallback_endpoint[self.fallback_index]}")

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
