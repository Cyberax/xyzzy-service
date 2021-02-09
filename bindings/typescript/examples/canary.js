"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const service_twirp_1 = require("../service.twirp");
const service_pb_1 = require("../service.pb");
let t = service_pb_1.xyzzy;
let td = service_twirp_1.createXyzzyData(process.env.XYZZY_BASE_URL || 'http://localhost:8080');
function ping() {
    return __awaiter(this, void 0, void 0, function* () {
        // Find the requests from the
        let res = yield td.ping(t.PingRequest.create());
        console.log(`Ping succeeded`);
    });
}
function main() {
    return __awaiter(this, void 0, void 0, function* () {
        yield ping();
    });
}
// @ts-ignore
(() => __awaiter(void 0, void 0, void 0, function* () {
    try {
        yield main();
    }
    catch (err) {
        console.error(err);
        process.exit(1);
    }
}))();
//# sourceMappingURL=canary.js.map