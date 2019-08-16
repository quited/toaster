from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse
import json
import random
import string
import os
import threading
import requests
import time
import wx
from wx.html2 import WebView
import subprocess
import ctypes
import sys
import logging as log

RESTfulHost = ('localhost', 8888)
RESTfulKeyPairs = {}
loopBackRESTfulKey = ''.join(random.sample(string.ascii_letters + string.digits, 8))
RESTfulKeyPairs[loopBackRESTfulKey] = "toast_master:Sv"

workloads = []
workloadPages = {}
workloadKeys = {}


def is_admin():
    try:
        return ctypes.windll.shell32.IsUserAnAdmin()
    except Exception as e:
        return False


def wlExec(workload, cmd):
    pwd = os.getcwd()
    path = os.path.join(os.getcwd(), workload[1][:-12]).replace('\\', '/')
    os.chdir(path)
    sub = subprocess.Popen(
        "launcher.exe " + cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    sub.wait()
    os.chdir(pwd)
    return sub


def scanWorkLoads():
    for root, dirs, files in os.walk(r"workloads"):
        for d in dirs:
            lp = os.path.join(os.path.join(root, d), "launcher.exe")
            if os.path.exists(lp):
                name = wlExec(("", lp), "-n").stdout.read().decode("utf-8") \
                    .replace("\r", "").replace("\n", "").split(":")
                if name[0] == "toast_workload" and len(name) == 2:
                    if not (name[1], lp) in workloads:
                        workloads.append((name[1], lp))


def resetWorkLoads():
    for wl in workloads:
        wlExec(wl, "-e")
        if wlExec(wl, "-s").stdout.read().decode("utf-8")\
                .replace("\r", "").replace("\n", "") != "stopped":
            log.error(u"停止工作负载 %s 失败" % wl[0])


def initWorkLoads():
    wl_cfg_dict = {}
    try:
        with open("workloads.json", "r") as wl_cfg:
            wl_cfg_dict = json.load(wl_cfg)
        for wl in workloads:
            if wl[0] in wl_cfg_dict:
                if wl_cfg_dict[wl[0]] == "new":
                    wlExec(wl, "-i")
                    log.info(u"初始化工作负载 %s" % wl[0])
                elif wl_cfg_dict[wl[0]] == "inited":
                    log.info(u"已初始化工作负载 %s" % wl[0])
            else:
                wl_cfg_dict[wl[0]] = "new"
                wlExec(wl, "-i")
                log.info(u"初始化工作负载 %s" % wl[0])
                wl_cfg_dict[wl[0]] = "inited"
        with open("workloads.json", "w") as wl_cfg:
            json.dump(wl_cfg_dict, wl_cfg)
    except Exception as e:
        log.error(u"初始化工作负载失败： %s" % e)


class RESTfulHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        details = urlparse(self.path)
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        pathItems = details.path.split('/')
        if pathItems[1] in RESTfulKeyPairs:
            workloadName = RESTfulKeyPairs[pathItems[1]]
            if pathItems[2] == "selfT":
                if workloadName == "toast_master:Sv":
                    data = {'stat': 'ok'}
                    self.wfile.write(json.dumps(data).encode())
            if pathItems[2] == "shutdown":
                if workloadName == "toast_master:Sv":
                    resetWorkLoads()
                    os._exit(0)


def serverTest():
    try:
        res = requests.get(url="http://%s:%s/%s/selfT" % (RESTfulHost + (loopBackRESTfulKey,)))
        resJson = json.loads(res.text)
        if "stat" in resJson and resJson["stat"] == "ok":
            return True
        else:
            return False
    except Exception as e:
        log.error(u"服务器状态异常： %s" % e)
        return False


class wlPage(wx.Panel):
    def __init__(self, parent, link):
        wx.Panel.__init__(self, parent)
        self.wv = WebView.New(self, size=(600, 400))
        self.wv.LoadURL(link)
        self.stat = wx.TextCtrl(self, style=wx.TE_READONLY | wx.TE_RICH2)
        page_box_sizer = wx.BoxSizer(wx.VERTICAL)
        page_box_sizer.Add(self.wv, 1, wx.ALIGN_CENTER, wx.ALL)
        page_box_sizer.Add(self.stat, 0, wx.ALIGN_LEFT | wx.EXPAND, wx.ALL)
        self.SetSizerAndFit(page_box_sizer)

    def showStat(self, pub, msg):
        self.stat.Clear()
        self.stat.SetDefaultStyle(wx.TextAttr(wx.RED))
        self.stat.AppendText(u"【" + pub + u"】")
        self.stat.SetDefaultStyle(wx.TextAttr(wx.BLUE))
        self.stat.AppendText(msg)


class wlFrame(wx.Frame):
    def __init__(self):
        wx.Frame.__init__(
            self, None, -1, title="Toaster",
            style=wx.DEFAULT_FRAME_STYLE ^ wx.RESIZE_BORDER
        )
        nb = wx.Notebook(self)
        for wl in workloads:
            wlKey = loopBackRESTfulKey
            while wlKey == loopBackRESTfulKey:
                wlKey = ''.join(random.sample(string.ascii_letters + string.digits, 8))
            wlEp = "http://%s:%s/%s/" % (RESTfulHost + (wlKey,))
            sub = wlExec(wl, "-b " + wlEp)
            page = wlPage(nb, sub.stdout.read().decode("utf-8"))
            if wlExec(wl, "-s").stdout.read().decode("utf-8")\
                    .replace("\r", "").replace("\n", "") == "running":
                nb.AddPage(page, wl[0])
                page.showStat(u"系统", u"工作负载：" + wl[0])
                workloadPages[wl[0]] = page
            else:
                log.error(u"启动负载 %s 失败" % wl[0])
        self.SetSize(620, 500)


def uiLauncher():
    time.sleep(1)
    if not serverTest():
        resetWorkLoads()
        os._exit(0)
    log.info(u"本地服务器已上线")
    app = wx.App()
    wlFrame().Show()
    app.MainLoop()
    resetWorkLoads()
    os._exit(0)


if __name__ == '__main__':
    log.basicConfig(
        level=log.DEBUG,
        format="%(asctime)s - %(name)s- %(levelname)s - %(message)s",
        datefmt="%a, %d %b %Y %H:%M:%S",
        handlers={log.FileHandler(filename="toaster.log", mode='a', encoding="utf-8")}
    )
    log.critical("----------------------------")
    whnd = ctypes.windll.kernel32.GetConsoleWindow()
    if whnd != 0:
        log.info(u"隐藏命令行窗口")
        ctypes.windll.user32.ShowWindow(whnd, 0)
    ctypes.windll.kernel32.CloseHandle(whnd)
    if not is_admin():
        log.info(u"需要申请管理员权限")
        ctypes.windll.shell32.ShellExecuteW(None, "runas", sys.executable, __file__, None, 1)
    if is_admin():
        log.info(u"已获取管理员权限")
        RESTfulServer = HTTPServer(RESTfulHost, RESTfulHandler)
        scanWorkLoads()
        resetWorkLoads()
        initWorkLoads()
        UIt = threading.Thread(target=uiLauncher)
        UIt.start()
        print("RESTful: %s:%s" % RESTfulHost)
        print("Master Key: %s" % loopBackRESTfulKey)
        RESTfulServer.serve_forever()
        resetWorkLoads()
    log.info(u"获取管理员权限失败")
