import hashlib
import base64
import time

HttpHeaderKeyAuthorization = "Authorization"
HttpHeaderKeyTime = "Time"
HttpHeaderMethod = "Method"
HttpHeaderURL = "URL"


def url_suffix(url):
    return "/" + list(reversed(url.split("//")))[0].split('/', 1)[1].rstrip('/')


def sha3_encode(secret_key, sign):
    h = hashlib.new('sha3_384')
    h.update(bytes(secret_key, encoding='utf-8'))
    h.update(bytes(sign, encoding='utf-8'))
    return base64.urlsafe_b64encode(h.digest())


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

    def _sign_url(self, method, url, unix_time_now) -> str:
        return ""

    def _sign_request(self, req):
        pass


class Auth(AuthBase):
    """
    jcs 认证类
    """
    keys = [HttpHeaderMethod, HttpHeaderURL, HttpHeaderKeyTime]

    def __init__(self, ak, sk):
        user_info = UserInfo(ak, sk)
        super().__init__(user_info)

    def _sign_url(self, method, url, unix_time_now) -> str:
        """
        为指定的 method 和 url 生成签名字符串 access_token
        :param method: 如 “PUT”， “GET”
        :param url: 如 /path/to/file?param=xxx
        :param unix_time_now: 如 1626938559 -> unix timestamp
        :return: 完整签名 token
        """

        value_dict = {
            HttpHeaderMethod: method,
            HttpHeaderURL: url,
            HttpHeaderKeyTime: unix_time_now,
        }
        # generate the origin sign
        sign = ""
        for key in self.keys:
            sign += key + ":" + str(value_dict[key]) + "\r\n"
        # encode with sha-3-384
        encode_sign = sha3_encode(self.user_info.accessSecret, sign)
        authorization = "JCS(" + self.user_info.accessKey + ":" + str(encode_sign, encoding="utf-8") + ")"
        return authorization

    def _sign_request(self, req):
        """
        为传入的 request 对象，设置签名好的请求头
        :param req: 待签名 request 对象
        :return: 签名好的 request
        """
        # get time
        time_now = time.time()

        # use _sign_url method to get authorization token
        authorization = self._sign_url(method=req.method, url=url_suffix(req.url), unix_time_now=int(time_now))
        # set authorization token in header
        req.headers[HttpHeaderKeyAuthorization] = authorization
        # set time in header
        req.headers[HttpHeaderKeyTime] = str(int(time_now))
