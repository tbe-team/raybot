import{c as a,i as t,f as r,u}from"./http-Lmb0ssjr.js";/**
 * @license lucide-vue-next v0.525.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const m=a("triangle-alert",[["path",{d:"m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3",key:"wmoenq"}],["path",{d:"M12 9v4",key:"juzpu7"}],["path",{d:"M12 17h.01",key:"p32p05"}]]),n={reboot(){return t.post("/system/reboot")},stopEmergency(){return t.post("/system/stop-emergency")},getInfo(e){return t.get("/system/info",e)},getStatus(e){return t.get("/system/status",e)}},s="system-info",y="system-status";function o(){return r({mutationFn:n.reboot})}function i(){return r({mutationFn:n.stopEmergency})}function f(e){return u({queryKey:[s],queryFn:()=>n.getInfo(e==null?void 0:e.axiosOpts),refetchInterval:e==null?void 0:e.refetchInterval})}function S(e){return u({queryKey:[y],queryFn:()=>n.getStatus({...e==null?void 0:e.axiosOpts}),refetchInterval:e==null?void 0:e.refetchInterval})}export{m as T,S as a,f as b,o as c,i as u};
