// @ts-nocheck
import { get as mvget } from '$lib/milletvekili/api'
import { get as bilgiget } from '$lib/bilgi/api'
import { capitalizeWord } from '$lib/utilities/stringUtilities';

export async function load({ params }) {
    let cities;

    // Get the cities from the database
    cities = await mvget("cities/");

    // Return the cities
    return { cities }
}