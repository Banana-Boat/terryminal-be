import { ServerDuplexStream } from "@grpc/grpc-js";
import { BasePty } from "./base-pty.js";
import {
  RunCmdRequest,
  RunCmdResponse,
  UnimplementedBasePtyService,
} from "./pb/base_pty.js";

export class BasePtyService extends UnimplementedBasePtyService {
  RunCmd(call: ServerDuplexStream<RunCmdRequest, RunCmdResponse>): void {
    let basePty = new BasePty(call);
    call.on("data", (chunk) => {
      const { cmd } = chunk;
      if (cmd) basePty.runCmd(cmd);
    });
  }
}
