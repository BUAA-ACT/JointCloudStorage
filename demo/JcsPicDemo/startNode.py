from node import Node

class NodesRunner(object):

    def __init__(self):
        self.nodes = []

    def start_nodes(self, is_test=False):
        ak = "c4470e6c2b28433f88cff9642429684e"
        sk = "adc52bedd6824fad83b9e8235b867b37"
        endpoint_hohhot = "http://jsi-aliyun-hohhot.jointcloudstorage.cn/"
        endpoint_qingdao = "http://jsi-aliyun-qingdao.jointcloudstorage.cn/"
        endpoint_hangzhou = "http://jsi-aliyun-hangzhou.jointcloudstorage.cn/"
        endpoint_chengdu = "http://jsi-txyun-chengdu.jointcloudstorage.cn/"
        endpoint_guangzhou = "http://jsi-bdyun-guangzhou.jointcloudstorage.cn/"

        fallback_endpoint = [endpoint_chengdu, endpoint_guangzhou, endpoint_hohhot, endpoint_qingdao]

        dict1 = "/step1/"
        dict2 = "/step2/"
        dict3 = "/step3/"
        dict4 = "/step4/"

        # 初始化 呼和浩特 节点，用于文件上传
        upload_node = Node("send", ak, sk, endpoint_hohhot, 1.0, "local", dict1, fallback_endpoint)
        upload_node.clear_all()
        upload_node.start()

        if is_test:
            print("测试环境运行")

        # 初始化 青岛 节点，图片彩色化
        if is_test:
            colorize_node = Node("test", ak, sk, endpoint_qingdao, 1.0, dict1, dict2, fallback_endpoint)
        else:
            colorize_node = Node("colorize", ak, sk, endpoint_qingdao, 1.0, dict1, dict2, fallback_endpoint)
        colorize_node.start()

        # 初始化 杭州 节点，图像增强
        if is_test:
            contrast_node = Node("test", ak, sk, endpoint_qingdao, 1.0, dict2, dict3, fallback_endpoint)
        else:
            contrast_node = Node("lar_en", ak, sk, endpoint_hangzhou, 1.0, dict2, dict3, fallback_endpoint)
        contrast_node.start()
        # large_node = Node("lar_en", ak, sk, endpoint, 1.0)
        # large_node.start()
        self.nodes = [upload_node, colorize_node, contrast_node]


