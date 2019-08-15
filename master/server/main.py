import threading
import requests
import json
import time
import os


class server:
    def __init__(self, p):
        self.port = p
        ts = threading.Thread(
            target=os.system,
            args=("python server\\server.py " + str(self.port),)
        )
        ts.start()
        time.sleep(1)
        if not self.selfTest():
            self.closeServer()
            os._exit(-1)
        tu = threading.Thread(
            target=os.system,
            args=("python server\\frontend\\main.py " + str(self.port),))
        tu.start()
        while tu.isAlive() and ts.isAlive():
            pass
        self.closeServer()
        os._exit(0)


    def selfTest(self):
        try:
            res = requests.get(url="http://localhost:" + str(self.port) + "/selfTest")
            resjson = json.loads(res.text)
            if "stat" in resjson and resjson["stat"] == "ok":
                print("Server Online")
                return True
            else:
                print("Err: Server Launch fail")
                return False
        except:
            print("Err: Server Launch fail")
            return False

    def closeServer(self):
        try:
            requests.get(url="http://localhost:" + str(self.port) + "/end")
        except:
            print("Server Closed")
