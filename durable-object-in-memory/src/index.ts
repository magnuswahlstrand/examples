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
        let id = env.LOCATION.idFromName("main");
        let locationDO = env.LOCATION.get(id);

        const location = {
            city: request.cf?.city,
            countryCode: request.cf?.country,
            postalCode: request.cf?.postalCode,
            isEUCountry: request.cf?.isEUCountry == "1"
        }
        // RPC to the Durable Object
        const resp = await locationDO.updateAndGetLocationText(location)

        return new Response(resp);
    }
}

export default handler