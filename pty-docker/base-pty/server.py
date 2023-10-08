import os
import select
import threading

from base_pty_pb2 import RunCmdResponse
from base_pty_pb2_grpc import BasePtyServicer as UnimplementedBasePtyServicer
from base_pty import BasePty


class BasePtyServicer(UnimplementedBasePtyServicer):
    def RunCmd(self, request_iterator, context):
        # 创建新Pty
        base_pty = BasePty()

        # 创建线程处理命令输入
        def send_cmds(request_iterator, pty_fd):
            for data in request_iterator:
                cmd = data.cmd + '\n'
                print('command: ' + cmd)
                pty_fd.write(cmd.encode('utf-8'))
                pty_fd.flush()

        threading.Thread(
            target=send_cmds, args=(request_iterator, base_pty.pty_fd)
        ).start()

        # 读取伪终端的输出并发送回客户端
        while True:
            r, _, _ = select.select([base_pty.pty_fd], [], [], 0.1)
            if base_pty.pty_fd in r:
                res = base_pty.pty_fd.read().decode('utf-8')
                if res:
                    print('result: ' + res)
                    yield RunCmdResponse(result=res)

            # 检查子进程是否退出
            if os.waitpid(base_pty.pid, os.WNOHANG) != (0, 0):
                break
