export function sleep(duration: number) {
  const start = new Date().getTime();
  while (new Date().getTime() - start < duration) {}
}
