import base64
from aip import AipImageProcess


class AipClient(object):
    """
    AipClient for handle
    """

    def __init__(self, appId, ak, sk):
        """
        :param appId: string App_ID
        :param ak: string Access_Key
        :param sk: string Secret_Key
        """
        self.appId = appId
        self.accessKey = ak
        self.secretKey = sk
        self.client = AipImageProcess(self.appId, self.accessKey, self.secretKey)
        pass

    def newClient(self, appId, ak, sk):
        """
        update the client field of the object
        :param appId: string App_ID
        :param ak: string Access_Key
        :param sk: string Secret_Key
        :return: nil
        """
        self.appId = appId
        self.accessKey = ak
        self.secretKey = sk
        self.client = AipImageProcess(self.appId, self.accessKey, self.secretKey)
        pass

    def _base64encode_(self, imgPath):
        """
        :param imgPath:string 图像路径
        :return: string 图像的base64编码
        """
        try:
            with open(imgPath,'rb') as f:
                return f.read()
        except:
            print("can't open and encode image.")
            exit(1)
        pass

    def _base64decode_(self, b64_code, newPath):
        """
        decode the base64 code into image and write into a jpg
        :param img: base64 code of the image
        :param newPath: path of new image
        :return: succeed or not
        """
        with open(newPath, 'wb') as f:
            img = base64.b64decode(b64_code.encode())
            f.write(img)

        pass

    def enlargeAndQualityEnhance(self, path, newpath):
        """
        enlarge the picture and enhance it's quality
        :param path: the path of picture
        :return: success or not
        """
        img = self._base64encode_(path)
        if img == None:
            return False
        else:
            res = self.client.imageQualityEnhance(img)
            if "image" in res:
                return self._base64decode_(res['image'],newpath)
            else:
                print(res['error_code'],":",res['error_msg'])
                return False
            pass
        pass

    def contrastEnhance(self, path, newpath):
        """
        enlarge contrast of the picture
        :param path: the path of picture
        :return: success or not
        """
        img = self._base64encode_(path)
        if img == None:
            return False
        else:
            res = self.client.imageQualityEnhance(img)
            if "image" in res:
                return self._base64decode_(res['image'], newpath)
            else:
                print(res['error_code'], ":", res['error_msg'])
                return False
            pass
        pass

    def definitionEnhance(self, path, newpath):
        """
        enlarge definition of the picture
        :param path: the path of picture
        :return: success or not
        """
        img = self._base64encode_(path)
        if img == None:
            return False
        else:
            res = self.client.imageDefinitionEnhance(img)
            if "image" in res:
                return self._base64decode_(res['image'], newpath)
            else:
                print(res['error_code'], ":", res['error_msg'])
                return False
            pass
        pass

    def colourize(self,path,newpath):
        """
        colorrize an image
        :param path: the path of picture
        :param newpath : the path of new picture
        :return: success or not
        """
        img = self._base64encode_(path)
        if img == None:
            return False
        else:
            res = self.client.colourize(img)
            if "image" in res:
                return self._base64decode_(res['image'], newpath)
            else:
                print(res['error_code'], ":", res['error_msg'])
                return False
            pass
        pass


if __name__ == "__main__":
    appid="24595163"
    ak="e3mXihhfyQPjY0gNoK97fj6v"
    sk="CczWdEiQs2gFd8w0cHn1EYeI1G4zCorL"
    aipclient=AipClient(appid,ak,sk)
    aipclient.colourize("picture/1.jpg","col1.jpg")
    #aipclient.definitionEnhance("test.png","de_test.png")
    #aipclient.enlargeAndQualityEnhance("test.png","enl_test.png")
