import {DurableObject} from "cloudflare:workers";

export class Room extends DurableObject {

    constructor(ctx: DurableObjectState, env: CloudflareBindings) {
        super(ctx, env);
    }

    async fetch(request: Request): Promise<Response> {
        // Creates two ends of a WebSocket connection.
        const webSocketPair = new WebSocketPair();
        const [client, server] = Object.values(webSocketPair);
        console.log(server.extensions)

        const tags: string[] = [] // Tags be used for filtering
        this.ctx.acceptWebSocket(server, tags);

        return new Response(null, {
            status: 101,
            webSocket: client,
        });
    }

    async webSocketMessage(ws: WebSocket, message: string) {
        this.ctx.getWebSockets().forEach(ws => ws.send(message));
    }

    async webSocketClose(ws: WebSocket, code: number, reason: string, wasClean: boolean) {
        ws.close();
    }
}