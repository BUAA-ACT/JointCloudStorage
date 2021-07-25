from JointCloudStorage import *
import logging

access_key = "47bed2954a3647fdbc7a3364778c388f"
secret_key = "2d61f1ec5d834af6ad7150d756fa33e5"
auth = Auth(access_key, secret_key)
bucket = Bucket(auth, "http://192.168.105.13:8085")
state = State(auth, "http://192.168.105.13:8085")


def test_put_object():
    with open("../requirements.txt") as file:
        bucket.put_object("/python/r.txt", file)
    logging.info("put object success")


def test_get_object_list():
    objects = bucket.get_object_list("/python/")
    logging.info(objects)
    for obj in objects:
        logging.info(obj)


def test_get_object():
    test_put_object()
    r = bucket.get_object("/python/r.txt")
    logging.info(r)


def test_put_object_async():
    with open("../requirements.txt") as file:
        bucket.put_object_async("/python/r.txt", file)
    logging.info("put object success")


def test_state():
    r = state.get_storage_state()
    logging.info(r)
    r = state.get_storage_plan()
    logging.info(r)
    r = state.get_server_info()
    logging.info(r)

