import ctypes
from server.main import server

if __name__ == "__main__":
    whnd = ctypes.windll.kernel32.GetConsoleWindow()
    if whnd != 0:    ctypes.windll.user32.ShowWindow(whnd, 0)
    ctypes.windll.kernel32.CloseHandle(whnd)
    Sv = server(6000)
