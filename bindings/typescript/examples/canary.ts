import {createXyzzyData} from '../service.twirp';
import {xyzzy} from '../service.pb';

let t = xyzzy;

let td = createXyzzyData(process.env.XYZZY_BASE_URL || 'http://localhost:8080');

async function ping() {
    // Find the requests from the
    let res = await td.ping(t.PingRequest.create());
    console.log(`Ping succeeded`);
}

async function main() {
    await ping();
}

// @ts-ignore
(async () => {
    try {
        await main()
    } catch (err) {
        console.error(err);
        process.exit(1);
    }
})();
