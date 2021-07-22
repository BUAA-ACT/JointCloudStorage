import requests


class Bucket(object):
    def __init__(self, auth, endpoint: str):
        self.auth = auth
        self.endpoint = endpoint.rstrip("/")

    def _gen_request_url(self, category, key: str):
        url = self.endpoint + "/" + category + "/"
        if len(key) > 0:
            url += key.rstrip("/")
        return url

    def put_object(self, key, value):
        pass

    def get_object(self, key):
        pass

    def delete_object(self, key):
        pass

    def get_object_list(self, prefix):
        req = requests.Request('GET', url=self._gen_request_url()
        pass
