
class JcsError(Exception):
    def __init__(self, status, headers, body, details):
        #: HTTP 状态码
        self.status = status

        #: HTTP响应体（部分）
        self.body = body

        #: 详细错误信息，是一个string
        self.details = details

    def __str__(self):
        error = {'status': self.status,
                 'details': self.details}
        return str(error)


def make_exception(resp):
    status = resp.status_code
    headers = resp.headers
    body = resp.text
    return JcsError(status, headers, body, body)
