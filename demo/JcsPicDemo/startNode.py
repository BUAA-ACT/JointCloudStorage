import threading

from node import Node


class NodesRunner(threading.Thread):

    def __init__(self, is_test=False):
        threading.Thread.__init__(self)
        self.nodes = []
        self.nodes_num = 0
        self.state = "ready"
        self.is_test = is_test

    def run(self):
        self.start_nodes(self.is_test)

    def start_nodes(self, is_test=False):
        self.state = "running"
        ak = "b5d2504ba5804a5690c2a57b4660b4ea"
        sk = "9d83480b5d07431e97f49a8fce67ce23"
        endpoint_hohhot = "http://jsi-aliyun-hohhot.jointcloudstorage.cn/"
        endpoint_qingdao = "http://jsi-aliyun-qingdao.jointcloudstorage.cn/"
        endpoint_hangzhou = "http://jsi-aliyun-hangzhou.jointcloudstorage.cn/"
        endpoint_chengdu = "http://jsi-txyun-chengdu.jointcloudstorage.cn/"
        endpoint_guangzhou = "http://jsi-bdyun-guangzhou.jointcloudstorage.cn/"

        fallback_endpoint = [endpoint_chengdu, endpoint_guangzhou, endpoint_hohhot, endpoint_qingdao]
        endpoint_name_dict = {
            endpoint_hohhot: "呼和浩特",
            endpoint_qingdao: "青岛",
            endpoint_hangzhou: "杭州",
            endpoint_chengdu: "成都",
            endpoint_guangzhou: "广州",
        }

        dict1 = "/step1/"
        dict2 = "/step2/"
        dict3 = "/step3/"
        dict4 = "/step4/"

        # 初始化 呼和浩特 节点，用于文件上传
        upload_node = Node("send", ak, sk, endpoint_hohhot, 1.0, "local", dict1, fallback_endpoint, endpoint_name_dict)

        if is_test:
            print("测试环境运行")

        # 初始化 青岛 节点，图片彩色化
        if is_test:
            colorize_node = Node("test", ak, sk, endpoint_qingdao, 0.2, dict1, dict2, fallback_endpoint,
                                 endpoint_name_dict)
        else:
            colorize_node = Node("colorize", ak, sk, endpoint_qingdao, 0.2, dict1, dict2, fallback_endpoint,
                                 endpoint_name_dict)

        # 初始化 杭州 节点，图像增强
        if is_test:
            contrast_node = Node("test", ak, sk, endpoint_hangzhou, 0.2, dict2, dict3, fallback_endpoint,
                                 endpoint_name_dict)
        else:
            contrast_node = Node("lar_en", ak, sk, endpoint_hangzhou, 0.2, dict2, dict3, fallback_endpoint, endpoint_name_dict)
        # large_node = Node("lar_en", ak, sk, endpoint, 1.0)
        # large_node.start()
        self.nodes = [upload_node, colorize_node, contrast_node]
        self.nodes_num = 3

        upload_node.clear_all()
        upload_node.start()
        colorize_node.start()
        contrast_node.start()
