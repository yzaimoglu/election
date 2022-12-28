// @ts-nocheck
import { get } from '$lib/milletvekili/api'
import { capitalizeWord } from '$lib/utilities/stringUtilities';

export async function load({ params }) {
    let district
    let cityName = params.city.toLowerCase()
    let districtName = params.district.toLowerCase()
    let districtReadableName = capitalizeWord(districtName)
    let cityReadableName = capitalizeWord(cityName)
    // Check if there is such a district
    if(params !== null) {
        district = await get("district/"+cityName+"/"+districtName+"/")
    }
    // Return the readable name and the district object
    return { districtReadableName, cityReadableName, district }
}