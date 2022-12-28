import { get } from '$lib/milletvekili/api.js'

export async function load() {
    const quarters = await get("quarters/");
    const districts = await get("districts/");
    const constituencies = await get("constituencies/");
    const cities = await get("cities/");
	return { quarters, districts, constituencies, cities };
}