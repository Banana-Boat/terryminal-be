import pty from "node-pty";
import { sleep } from "./util.js";

export class BasePty {
  ptyProcess: pty.IPty;

  constructor() {
    this.ptyProcess = pty.spawn("bash", [], {
      name: "xterm-color",
      cols: 80,
      rows: 30,
      cwd: "/",
    });

    this.ptyProcess.onData((data) => {
      process.stdout.write("terryminal: " + data);
    });

    sleep(1000);
  }

  close() {
    this.ptyProcess.kill();
  }

  runCmd(cmd: string) {
    // 判断命令，注意多命令连接符，待实现...

    this.ptyProcess.write(cmd + "\r");
  }
}
