import os

from JointCloudStorage import *
import requests

if __name__ == '__main__':
    access_key = "47bed2954a3647fdbc7a3364778c388f"
    secret_key = "2d61f1ec5d834af6ad7150d756fa33e5"
    auth = Auth(access_key,secret_key)
    bucket = Bucket(auth, "http://127.0.0.1:8085")
    files = bucket.get_object_list("/jsi")
    print(files)

    with open("../requirements.txt") as file:
        bucket.put_object("/python/r.txt", file)

    files = bucket.get_object_list("/python/")
    print(files)

    if not os.path.exists("./tmp/"):
        os.makedirs("./tmp/")

    c = bucket.get_object("/python/r.txt")
    with open("./tmp/t.txt", "wb+") as f:
        f.write(c)

    bucket.delete_object("/python/r.txt")
    files = bucket.get_object_list("/python/")
    print(files)



