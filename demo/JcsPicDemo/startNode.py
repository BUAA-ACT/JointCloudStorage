from node import Node, clear_all

if __name__ == '__main__':
    ak = "c4470e6c2b28433f88cff9642429684e"
    sk = "adc52bedd6824fad83b9e8235b867b37"
    endpoint_hohhot = "http://jsi-aliyun-hohhot.jointcloudstorage.cn/"
    endpoint_qingdao = "http://jsi-aliyun-qingdao.jointcloudstorage.cn/"
    endpoint_hangzhou = "http://jsi-aliyun-hangzhou.jointcloudstorage.cn/"

    # 初始化 呼和浩特 节点，用于文件上传
    upload_node = Node("send", ak, sk, endpoint_hohhot, 3.0)
    clear_all(upload_node)
    upload_node.start()

    # 初始化 青岛 节点，图片彩色化
    colorize_node = Node("colorize", ak, sk, endpoint_qingdao, 1.0)
    colorize_node.start()

    # 初始化 杭州 节点，图像增强
    contrast_node = Node("lar_en", ak, sk, endpoint_hangzhou, 1.0)
    contrast_node.start()
    # large_node = Node("lar_en", ak, sk, endpoint, 1.0)
    # large_node.start()

    upload_node.join()
    colorize_node.join()
    contrast_node.join()
