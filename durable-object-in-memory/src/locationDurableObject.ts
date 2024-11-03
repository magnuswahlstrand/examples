import {DurableObject} from "cloudflare:workers";

// Durable Object
export class Location extends DurableObject {

    constructor(ctx: DurableObjectState, env: unknown) {
        super(ctx, env);
    }

    async updateAndGetLocationText(newLocation: LocationInfo) {
        let responseText: string
        if (this.location == null) {
            responseText = templateNewLocation(newLocation)
        } else {
            responseText = templateUpdatedLocation(this.location, newLocation)
        }
        this.location = newLocation
        return responseText;
    }
}

function formatLocation(location: LocationInfo) {
    return `City: ${location.city}
Country Code: ${location.countryCode}
Postal Code: ${location.postalCode}
Is EU Country: ${location.isEUCountry ? "Yes" : "No"}
`
}

function templateNewLocation(location: LocationInfo) {
    return `This is the first request to this Durable object instance. Location was not set.

New state:
${formatLocation(location)}`
}

function templateUpdatedLocation(oldLocation: LocationInfo, newLocation: LocationInfo) {
    return `Durable object was already loaded with an in memory state. Updating state.  
    
Previous state:
${formatLocation(oldLocation)}

New state:
${formatLocation(newLocation)}`
}
