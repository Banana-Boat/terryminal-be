import fcntl
import os
import pty


class BasePty:
    def __init__(self) -> None:
        master, slave = pty.openpty()
        self.pid = os.fork()

        if self.pid == 0:
            # 子进程
            os.close(master)
            os.dup2(slave, 0)
            os.dup2(slave, 1)
            os.dup2(slave, 2)

            os.chdir('/')
            os.execvp('/bin/bash', ['/bin/bash'])
        else:
            # 父进程
            os.close(slave)
            self.pty_fd = os.fdopen(
                master, 'rb+', buffering=0)

            # 设置非阻塞模式
            flags = fcntl.fcntl(self.pty_fd, fcntl.F_GETFL)
            fcntl.fcntl(self.pty_fd, fcntl.F_SETFL, flags | os.O_NONBLOCK)

            print('====base pty created successfully====')

    def __del__(self):
        self.pty_fd.close()
        print('====base pty closed successfully====')
