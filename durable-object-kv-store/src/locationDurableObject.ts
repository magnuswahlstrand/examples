import {DurableObject} from "cloudflare:workers";

// Durable Object
export class Location extends DurableObject {

    constructor(ctx: DurableObjectState, env: unknown) {
        super(ctx, env);
    }

    async trackClientCity(city: string) {
        const counter = (await this.ctx.storage.get('counter') as Record<string, number>) || {}
        const updated = {
            ...counter,
            [city]: (counter[city] || 0) + 1
        }

        await this.ctx.storage.put('counter', updated)
        return updated
    }
}