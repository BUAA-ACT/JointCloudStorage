from flask import Flask
from node import Node
import logging
logging.basicConfig(level=logging.DEBUG)

app = Flask(__name__)


@app.route('/')
def hello_world():
    return 'Hello World!'


if __name__ == '__main__':
    ak = "bd05544ea474472c9deeffe1f917f239"
    sk = "b4d5a58c7a454c4f852667bafd6a4eac"
    endpoint = "http://192.168.105.13:8085"
    upload_node = Node("send", ak, sk, endpoint)
    upload_node.start()
    upload_node.join()
