from aip import AipImageProcess

import threading


class NodeClient(threading.Thread):
    def __init__(self, bucket, up_dict, out_dict, filename, state, file_lock, num_lock, state_lock):
        self.bucket = bucket
        self.up_dict = up_dict
        self.out_dict = out_dict
        self.filename = filename
        self.state = state
        self.file_lock = file_lock
        self.num_lock = num_lock
        self.state_lock = state_lock

    def run(self):
        c = self.bucket.get_object(self.dict + self.filename)
        file_lock.
        self.node_state.file_processing = filename
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
        self.sendBytes(res, self.output_dict + filename)
        logger.info(f"{self.task_type} 节点处理 {filename} 结果上传成功")
        self.node_state.finish_num += 1
        self.node_state.state = "OK"
