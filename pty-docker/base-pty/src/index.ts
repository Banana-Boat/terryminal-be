import { BasePtyService } from "./server.js";
import grpc from "@grpc/grpc-js";
import { UnimplementedBasePtyService } from "./pb/base_pty.js";

const HOST = "0.0.0.0";
const PORT = "8081";

(() => {
  const basePtyService = new BasePtyService();
  const server = new grpc.Server();

  server.addService(UnimplementedBasePtyService.definition, basePtyService);

  server.bindAsync(
    `${HOST}:${PORT}`,
    grpc.ServerCredentials.createInsecure(),
    (err, port) => {
      if (err) {
        console.log(err);
        return;
      }
      server.start();
      console.log(`Server running on port: ${port}`);
    }
  );
})();
