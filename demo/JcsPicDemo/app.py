import jsonpickle
from flask import Flask, jsonify, Response, send_file, render_template
from node import Node
import logging
from startNode import NodesRunner
import multiprocessing
from node import NodeState


logging.basicConfig(level=logging.DEBUG)
RUNNER = NodesRunner(is_test=True)

app = Flask(__name__)


class Info(object):
    def __init__(self, runner_state: str, node_states: list[NodeState]):
        self.runner_state = runner_state
        self.node_states = node_states

    def to_json(self):
        return jsonpickle.encode(self)

    def to_json_response(self):
        return Response(self.to_json(), mimetype='application/json')


@app.route('/')
def hello_world():
    return render_template("index.html")


@app.route("/start")
def start_run():
    global RUNNER
    if RUNNER.state == "running":
        return hello_world()
    # background_process = multiprocessing.Process(name="nodes", target=RUNNER.start_nodes, args=(False,))
    # background_process.daemon = True
    # background_process.start()
    RUNNER.start()
    return hello_world()


@app.route("/info")
def get_info():
    global RUNNER
    if RUNNER.state != "running":
        return Info("流水线尚未启动", []).to_json_response()

    info = Info(RUNNER.state, [])
    for node in RUNNER.nodes:
        info.node_states.append(node.get_state())

    return info.to_json_response()

@app.route("/charts")
def charts():
    return app.send_static_file("charts.html")


if __name__ == '__main__':
    pass
