/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run `wrangler dev src/index.ts` in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run `wrangler publish src/index.ts --name my-worker` to publish your worker
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */
import {Location} from "./locationDurableObject";

export {Location}

export interface Env {
    LOCATION: DurableObjectNamespace<Location>;
}

// Worker
const handler: ExportedHandler<Env> = {
    async fetch(request, env: Env) {
        const url = new URL(request.url)
        if (url.pathname != "/") {
            return new Response("Not found", {status: 404})
        }

        let id = env.LOCATION.idFromName("main");
        let locationDO = env.LOCATION.get(id);


        // RPC to the Durable Object
        const resp = await locationDO.trackClientCity(request?.cf?.city || "unknown");

        const sorted = '\n\n' + Object.entries(resp)
            .sort(([, a], [, b]) => b - a)
            .map(([city, count]) => `${city}: ${count}`)
            .join('\n')

        return new Response(sorted);
    }
}

export default handler