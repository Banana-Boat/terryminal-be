import { BasePtyService } from "./server.js";
import grpc from "@grpc/grpc-js";
import { UnimplementedBasePtyService } from "./pb/base_pty.js";
import dotenv from "dotenv";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
dotenv.config({
  path: path.join(__dirname, "..", `.env`),
});

(() => {
  const { SERVICE_HOST, SERVICE_PORT } = process.env;
  const basePtyService = new BasePtyService();
  const server = new grpc.Server();

  server.addService(UnimplementedBasePtyService.definition, basePtyService);

  server.bindAsync(
    `${SERVICE_HOST}:${SERVICE_PORT}`,
    grpc.ServerCredentials.createInsecure(),
    (err, port) => {
      if (err) {
        console.log(err);
        return;
      }
      server.start();
      console.log(`Server running at: ${SERVICE_HOST}:${port}`);
    }
  );
})();
