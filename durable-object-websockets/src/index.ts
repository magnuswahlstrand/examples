import {Hono} from 'hono'
import {Room} from "./durable_object_room";

export {Room}

const app = new Hono<{ Bindings: CloudflareBindings }>()

app.get('/ws', async (c) => {
        const upgradeHeader = c.req.header('Upgrade');
        if (!upgradeHeader || upgradeHeader !== 'websocket') {
            return new Response('Durable Object expected Upgrade: websocket', {status: 426});
        }

        const roomName = "magnus"
        const id: DurableObjectId = c.env.ROOM.idFromName(roomName)
        return c.env.ROOM.get(id).fetch(c.req.raw)
    }
)

export default app