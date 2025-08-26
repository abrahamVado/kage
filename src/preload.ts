import { contextBridge, ipcRenderer } from 'electron';
const NS = 'kage' as const;
const api = { ping: () => ipcRenderer.invoke(`${NS}:ping`) };
contextBridge.exposeInMainWorld(NS, api);