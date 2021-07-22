class UserInfo(object):
    accessKey = ""
    accessSecret = ""

    def __init__(self, ak, sk):
        self.accessKey = ak
        self.accessSecret = sk


class AuthBase(object):
    """
    提供认证基础服务
    """

    def __init__(self, user_info):
        self.user_info = user_info
        pass

    def _sign_url(self, method, url) -> str:
        return ""

    def _sign_request(self, req):
        pass


class Auth(AuthBase):
    """
    jcs 认证类
    """

    def __init__(self, ak, sk):
        user_info = UserInfo(ak, sk)
        super().__init__(user_info)

    def _sign_url(self, method, url) -> str:
        """
        为指定的 method 和 url 生成签名字符串 access_token
        :param method: 如 “PUT”， “GET”
        :param url: 如 /path/to/file?param=xxx
        :return: 完整签名 token
        """
        # todo 待编写
        return ""

    def _sign_request(self, req):
        """
        为传入的 request 对象，设置签名好的请求头
        :param req: 待签名 request 对象
        :return: 签名好的 request
        """
        # todo 调用 _sign_url 方式获取签名 token，并设置到请求头内
        pass
