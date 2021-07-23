import requests
import logging
from . import exceptions

logger = logging.getLogger(__name__)


class _BASE(object):
    def __init__(self, auth, endpoint: str):
        self.auth = auth
        self.endpoint = endpoint.rstrip("/")
        self.session = requests.session()

    def _gen_request_url(self, category, key: str):
        url = self.endpoint + "/" + category + "/"
        if len(key) > 0:
            url += key.strip("/")
        return url

    def _do(self, method, category, key, **kwargs):
        req = requests.Request(method, url=self._gen_request_url(category, key), **kwargs)
        pre_request = self.session.prepare_request(req)
        self.auth._sign_request(pre_request)
        resp = self.session.send(pre_request)
        if resp.status_code // 100 != 2:
            e = exceptions.make_exception(resp)
            logger.info("Exception: {0}".format(e))
            raise e
        return resp


class Bucket(_BASE):
    """

    """
    CATEGORY_OBJECT = "object"
    CATEGORY_STATE = "state"

    def __init__(self, auth, endpoint: str):
        super().__init__(auth, endpoint)

    def put_object(self, key, data):
        self._do('PUT', Bucket.CATEGORY_OBJECT, key=key, data=data)
        return True

    def get_object(self, key):
        result = self._do("GET", Bucket.CATEGORY_OBJECT, key=key)
        return result.content

    def delete_object(self, key):
        self._do("DELETE", Bucket.CATEGORY_OBJECT, key=key)
        return True

    def get_object_list(self, prefix):
        params = {
            'keyPrefix': prefix
        }
        req = self._do('GET', Bucket.CATEGORY_OBJECT, key="", params=params)
        result = req.json()
        return result

