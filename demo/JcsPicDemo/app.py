from flask import Flask
from node import Node
import logging
from startNode import NodesRunner
import multiprocessing

logging.basicConfig(level=logging.DEBUG)

app = Flask(__name__)


@app.route('/')
def hello_world():
    return 'Hello World!'


@app.route("/start")
def start_run():
    runner = NodesRunner()
    background_process = multiprocessing.Process(name="nodes", target=runner.start_nodes, args=(False, ))
    background_process.daemon = True
    background_process.start()
    return hello_world()


if __name__ == '__main__':
    pass