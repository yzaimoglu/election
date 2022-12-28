// @ts-nocheck
import { get as mvget } from '$lib/milletvekili/api'
import { get as bilgiget } from '$lib/bilgi/api'
import { capitalizeWord } from '$lib/utilities/stringUtilities';

export async function load({ params }) {
    let city;
    let cityName = params.city

    // Get the city from the database
    city = await mvget("city/"+cityName+"/");

    // Return the city
    return { city }
}