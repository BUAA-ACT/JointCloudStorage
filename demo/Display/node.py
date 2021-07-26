import argparse
import functools
import os
import datetime

from JointCloudStorage import *

from pic_haddle import AipClient
from utils import add_arguments
from aip import AipImageProcess

parser = argparse.ArgumentParser(description=__doc__)
add_arg = functools.partial(add_arguments, argparser=parser)
# yapf: disable
add_arg('type', str, "send", "Node type:[send,colorize,con_en,lar_en]")
add_arg('ak', str, None, "access key")
add_arg('sk', str, None, "secret key")
add_arg('endpoint', str, None, "service address,ip:port")

dict1 = "/step1/"
dict2 = "/step2/"
dict3 = "/step3/"
dict4 = "/step4/"

send_index = 0

appid = "24595163"
ak = "e3mXihhfyQPjY0gNoK97fj6v"
sk = "CczWdEiQs2gFd8w0cHn1EYeI1G4zCorL"
aip_client = AipImageProcess(appid, ak, sk)


def send(bucket):
    buck = bucket
    try:
        fileList = os.listdir("picture")
        with open("picture/" + fileList[send_index], 'rb') as f:
            dt = datetime.datetime.now()
            if buck.put_object(dict1 + "pic" + dt.strftime("%Y-%m-%d-%H-%M-%S") + ".jpg", f):
                print("upload OK!")
            else:
                print("upload Fail!")
        sendindex = (send_index + 1) / len(fileList)
    except Exception as e:
        logger.warning("step1:", e)
    pass


def colorize(bucket):
    buck = bucket
    try:
        fileList1 = bucket.get_object_list(dict1)
        fileList2 = bucket.get_object_list(dict2)
        set1=()
        set2=()
        if fileList1:
            set1 = set(getFileNames(fileList1))
        if fileList2:
            set2 = set(getFileNames(fileList2))
        difference = list(set1.difference(set2))
        if difference:
            filename = difference[0]
            c = buck.get_object(filename)
            res = aip_client.colourize(c)
            # 发送文件
            rel_filename=filename.split('/')[-1]
            sendBytes(res, dict2 + rel_filename)
    except Exception as e:
        logger.warning("step2:", e)
    pass


def constrastenhance(bucket):
    buck = bucket
    try:
        fileList1 = bucket.get_object_list(dict2)
        fileList2 = bucket.get_object_list(dict3)
        set1=()
        set2=()

        if fileList1:
            set1 = set(getFileNames(fileList1))
        if fileList2:
            set2 = set(getFileNames(fileList2))

        difference = list(set1.difference(set2))
        if difference:
            filename = difference[0]
            c = buck.get_object(filename)
            res = aip_client.contrastEnhance(c)
            # 发送文件
            rel_filename=filename.split('/')[-1]
            sendBytes(res, dict3 + rel_filename)
    except Exception as e:
        logger.warning("step3:", e)
    pass


def largeenhance(bucket):
    buck = bucket
    try:
        fileList1 = bucket.get_object_list(dict3)
        fileList2 = bucket.get_object_list(dict4)
        set1=()
        set2=()

        if fileList1:
            set1 = set(getFileNames(fileList1))
        if fileList2:
            set2 = set(getFileNames(fileList2))

        difference = list(set1.difference(set2))
        if difference:
            filename = difference[0]
            c = buck.get_object(filename)
            res = aip_client.contrastEnhance(c)
            # 发送文件
            rel_filename=filename.split('/')[-1]
            sendBytes(res, dict4 + rel_filename)
    except Exception as e:
        logger.warning("step4:", e)
    pass


def sendBytes(res, path):
    if "image" in res:
        img = base64.b64decode(res['image'].encode())
        # 写入step2
        # with open("tmp/tmp.jpg",'rwb') as tmp:
        #     tmp.write(img)
        #     dt=datetime.datetime.now()
        #     bucket.put_object(dict2+"pic"+dt.strftime("%Y-%m-%d-%H-%M-%S"),tmp)
        dt = datetime.datetime.now()
        bucket.put_object(path, img)
    else:
        logger.warning(res['error_code'] + ":" + res['error_msg'])
    pass

def getFileNames(fileList):
    names=[]
    for file in fileList:
        names.append(file['Filename'])
        pass
    return names
    pass


switch = {
    "send": send,
    "colorize": colorize,
    "con_en": constrastenhance,
    "lar_en": largeenhance,
}

if __name__ == "__main__":
    args = parser.parse_args()
    print(args)
    auth = Auth(args.ak, args.sk)
    bucket = Bucket(auth, args.endpoint)
    while True:
        switch[args.type](bucket)
        print("succeed!")
        time.sleep(5.0)
