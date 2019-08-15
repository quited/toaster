import os
import subprocess
import sys

from flask import Flask
import flask_restful as restful
from flask_apscheduler import APScheduler
import datetime

workloads = []


def scanWorkLoads():
    for root, dirs, files in os.walk(r"..\workloads"):
        for d in dirs:
            lp = os.path.join(os.path.join(root, d), "launcher.exe")
            if os.path.exists(lp):
                sub = subprocess.Popen(lp + " -s", stdout=subprocess.PIPE)
                sub.wait()
                stat = sub.stdout.read().decode("utf-8")
                sub = subprocess.Popen(lp + " -n", stdout=subprocess.PIPE)
                sub.wait()
                name = sub.stdout.read().decode("utf-8").split(":")
                if (name[0] == "toast_workload" and len(name) == 2) and \
                        (stat == "running" or stat == "stop"):
                    if not (name[1], lp) in workloads:
                        workloads.append((name[1], lp))
    print(workloads)


class selfTest(restful.Resource):
    def get(self):
        scheduler.add_job(
            func=scanWorkLoads, trigger='date', id="0",
            next_run_time=datetime.datetime.now() + datetime.timedelta(seconds=1)
        )
        return {"stat": "ok"}


class end(restful.Resource):
    def get(self):
        print("Server End Called")
        os._exit(0)


if __name__ == "__main__":
    app = Flask(__name__)
    api = restful.Api(app)
    api.add_resource(selfTest, "/selfTest")
    api.add_resource(end, "/end")
    scheduler = APScheduler()
    scheduler.init_app(app)
    scheduler.start()
    print(len(sys.argv))
    if len(sys.argv) == 2:
        app.run(
            port=int(sys.argv[1])
        )
    else:
        app.run()
