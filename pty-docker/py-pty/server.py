import pty
import os
import select
import time

# 启动子进程
# pid = os.fork()

# if pid == 0:
#     # 在子进程中执行命令
#     pass


# else:
#     print('Child pid is %s' % pid)
#     os.waitpid(pid, 0)

with open('./log', 'w') as script, open('./stdin', 'r') as f:
    def master_read(fd):
        time.sleep(0.1)
        data = os.read(fd, 1)
        script.write(data.decode())
        script.flush()
        return data

    def stdin_read(fd):
        data = f.readline().encode()
        return data

    pty.spawn('sh', master_read, stdin_read)
