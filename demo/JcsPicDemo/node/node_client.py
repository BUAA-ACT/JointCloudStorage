import base64


import threading
from logzero import logger
from aip import AipImageProcess

app_id = "24595163"
ak = "e3mXihhfyQPjY0gNoK97fj6v"
sk = "CczWdEiQs2gFd8w0cHn1EYeI1G4zCorL"
aip_client = AipImageProcess(app_id, ak, sk)

class NodeClient(threading.Thread):
    def __init__(self, bucket, task_type, up_dict, out_dict, file_name, state, file_lock, num_lock, state_lock):
        threading.Thread.__init__(self)
        self.bucket = bucket
        self.task_type = task_type
        self.up_dict = up_dict
        self.out_dict = out_dict
        self.file_name = file_name
        self.state = state
        self.file_lock = file_lock
        self.num_lock = num_lock
        self.state_lock = state_lock

    def run(self):
        c = self.bucket.get_object(self.up_dict + self.file_name)
        self.file_lock.acquire()
        self.state.file_processing = self.file_name
        self.file_lock.release()
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
            res = {"image": str(base64.b64encode(c), "utf-8")}
            logger.error(f"图像处理时错误，Error: {e}")
        self.sendBytes(res, self.out_dict + self.file_name)
        logger.info(f"{self.task_type} 节点处理 {self.file_name} 结果上传成功")

        self.num_lock.acquire()
        self.state.finish_num += 1
        self.num_lock.release()

        self.state_lock.acquire()
        self.state.state = "OK"
        self.state_lock.release()

    def sendBytes(self, res, path):
        if "image" in res:
            img = base64.b64decode(res['image'].encode())
            self.bucket.put_object(path, img)
        else:
            logger.warning(res['error_code'] + ":" + res['error_msg'] + " path: " + path)
        pass
