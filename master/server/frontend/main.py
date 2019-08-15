import wx
import wx.lib.scrolledpanel as scrolled
from wx.html2 import WebView
import requests
import sys
import json
import time
import threading
import os

port = 5000
if len(sys.argv) == 2:
    port = int(sys.argv[1])


class winFrame(wx.Frame):
    def __init__(self, title):
        wx.Frame.__init__(self, None, -1, title, style=wx.DEFAULT_FRAME_STYLE ^ wx.RESIZE_BORDER)

        mainBoxSizer = wx.BoxSizer(wx.HORIZONTAL)
        scrollFlexSizer = wx.FlexGridSizer(100, 1, 10, 0)
        scrollPanel = scrolled.ScrolledPanel(self, -1)
        infoFlexSizer = wx.FlexGridSizer(2, 1, 0, 0)
        infoPanel = scrolled.ScrolledPanel(self, -1)

        clickBox1 = wx.CheckBox(scrollPanel, label="Pornhub")
        clickBox2 = wx.CheckBox(scrollPanel, label="QQ")
        clickBox3 = wx.CheckBox(scrollPanel, label="WeChat")
        clickBox4 = wx.CheckBox(scrollPanel, label="NetEase")

        scrollFlexSizer.AddMany([
            (clickBox1, 0, wx.SHAPED),
            (clickBox2, 0, wx.SHAPED),
            (clickBox3, 0, wx.SHAPED),
            (clickBox4, 0, wx.SHAPED),
        ])
        scrollFlexSizer.AddGrowableCol(0)

        scrollPanel.SetMaxSize((200, 500))
        scrollPanel.SetMinSize((200, 500))
        scrollPanel.SetSizerAndFit(scrollFlexSizer)
        scrollPanel.SetAutoLayout(1)
        scrollPanel.SetupScrolling()

        # textBox1 = wx.TextCtrl(infoPanel, -1, size=(600, 400), style=wx.TE_READONLY | wx.TE_MULTILINE)
        self.textBox2 = wx.TextCtrl(infoPanel, -1, size=(600, 100),
                                    style=wx.TE_MULTILINE | wx.TE_RICH2 | wx.TE_READONLY)
        wv = WebView.New(infoPanel, size=(600, 400))
        wv.LoadURL("https://www.baidu.com/")

        # textBox1.write("hello\nworld")
        # textBox1.Clear()
        self.textBox2.Bind(wx.EVT_TEXT_ENTER, self.onTextEnter)

        infoFlexSizer.AddMany([
            # (textBox1, 0, wx.SHAPED),
            (wv, 0, wx.SHAPED),
            (self.textBox2, 0, wx.SHAPED),
        ])

        infoPanel.SetSizerAndFit(infoFlexSizer)
        infoPanel.SetAutoLayout(1)
        infoPanel.SetupScrolling()

        mainBoxSizer.Add(scrollPanel, 1, wx.ALIGN_LEFT, wx.ALL)
        mainBoxSizer.Add(infoPanel, 4, wx.ALIGN_RIGHT, wx.ALL)

        self.SetSizerAndFit(mainBoxSizer)
        self.SetSize((700, 500))

        self.Bind(wx.EVT_CHECKBOX, self.onChecked)

    def onChecked(self, e):
        cb = e.GetEventObject()
        tc = self.textBox2
        tc.SetDefaultStyle(wx.TextAttr(wx.RED))
        tc.AppendText(u"【系统】")
        tc.SetDefaultStyle(wx.TextAttr(wx.BLUE))
        tc.AppendText(cb.GetLabel() + " is clicked to " + str(cb.GetValue()) + "\n")
        tc.SetDefaultStyle(wx.TextAttr(wx.BLACK))

    def onTextEnter(self, e):
        tc = e.GetEventObject()
        tc.SetDefaultStyle(wx.TextAttr(wx.RED))
        tc.AppendText(u"【系统】")
        tc.SetDefaultStyle(wx.TextAttr(wx.BLUE))
        tc.AppendText(u"无交互式命令行\n")
        tc.SetDefaultStyle(wx.TextAttr(wx.BLACK))


def closeServer():
    try:
        requests.get(url="http://localhost:" + str(port) + "/end")
    except:
        print("Server Closed")


def selfTest():
    try:
        res = requests.get(url="http://localhost:" + str(port) + "/selfTest")
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


def heartBeatCheck():
    while 1:
        try:
            if not selfTest():
                os._exit(0)
        except:
            os._exit(0)
    time.sleep(2)


if __name__ == "__main__":
    app = wx.App()
    winFrame("frame Test").Show()
    tHbc = threading.Thread(target=heartBeatCheck)
    tHbc.start()
    app.MainLoop()
    closeServer()
