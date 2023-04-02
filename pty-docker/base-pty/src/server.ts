import { BasePty } from "./base-pty.js";

(() => {
  const basePty = new BasePty();
  basePty.runCmd("echo Hello world");
})();
