import { ServerDuplexStream } from "@grpc/grpc-js";
import pty from "node-pty";
import { RunCmdRequest, RunCmdResponse } from "./pb/base_pty.js";
import { sleep } from "./util.js";

export class BasePty {
  ptyProcess: pty.IPty;

  constructor(call: ServerDuplexStream<RunCmdRequest, RunCmdResponse>) {
    this.ptyProcess = pty.spawn("bash", [], {
      cols: 300, // 影响传输次数，待查明！！
      rows: 30,
      cwd: "/",
    });

    this.ptyProcess.onData((data) => {
      console.log(`result: ${data}`);
      call.write(new RunCmdResponse({ result: data }));
    });

    sleep(500);
  }

  close() {
    this.ptyProcess.kill();
  }

  runCmd(cmd: string) {
    // 判断命令，注意多命令连接符，待实现...

    this.ptyProcess.write(cmd + "\r");
  }
}
