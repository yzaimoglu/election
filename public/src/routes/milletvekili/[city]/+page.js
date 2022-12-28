// @ts-nocheck
import { get as mvget } from '$lib/milletvekili/api'
import { get as bilgiget } from '$lib/bilgi/api'
import { capitalizeWord } from '$lib/utilities/stringUtilities';

export async function load({ params }) {
    let city;
    let candidates;
    let encoded;
    let candidateNamesTogether = "";
    let cityName = params.city.toLowerCase();
    let cityReadableName = capitalizeWord(cityName);
    if(params !== null) {
        city = await mvget("city/"+cityName+"/");
        city.candidates.forEach((candidate, index) => {
            let concatenator = "+"
            if(index+1 === city.candidates.length) {
                concatenator = "";
            } 
            candidateNamesTogether += candidate.firstname+"-"+candidate.lastname+concatenator;
        });
        encoded = btoa(candidateNamesTogether)
        candidates = await bilgiget("individuals/"+encoded+"/")
    }
    return { cityReadableName, city, candidates }
}