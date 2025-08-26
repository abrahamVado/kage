// kage/src/index.ts
// SUPER COMMENTS — IMPLEMENTATION ROADMAP
import { ipcMain } from 'electron';
const NS = 'kage' as const;
export function activate() {
  ipcMain.handle(`${NS}:ping`, () => ({ ok: true, purpose: "Advanced steganography: hide & encrypt arbitrary files inside media carriers." }));
}
export function deactivate() {
  ipcMain.removeHandler(`${NS}:ping`);
}