<script>
// @ts-nocheck
    import { get as bilgiget } from '$lib/bilgi/api'
    import { get as mvget } from '$lib/milletvekili/api'
    import { onMount } from 'svelte'
    import _ from 'lodash'
    /** @type {import('./$types').PageData} */
    export let data;
    let candidates = data.city.candidates
    let candidatesBilgi, candidatesResults
    let show = false

    // Return the candidate names together
    function getCandidateNamesTogether() {
      let candidateNamesTogether = ""
      candidates.forEach((candidate, index) => {
            let concatenator = "+"
            if(index+1 === candidates.length) {
                concatenator = "";
            } 
            candidateNamesTogether += candidate.firstname+"-"+candidate.lastname+concatenator;
      });
      return candidateNamesTogether
    }

    // Make the GET Request to the API
    async function getCandidatesFromBilgi(input) {
      return await bilgiget("individuals/"+input+"/")
    }

    // Make the GET Request to the API
    async function getCandidatesResults() {
      return await mvget("results/"+data.city.name+"/")
    }

    // Get a specific value from a candidate bilgi
    function getValueFromCandidateBilgi(candidate, value) {
      let returnValue
      candidatesBilgi.forEach(candidateBilgi => {
        if(candidate.firstname === candidateBilgi.firstname && candidate.lastname === candidateBilgi.lastname) {
          returnValue = candidateBilgi[value]
        }
      });
      return returnValue
    }

    // Get a specific value from a candidate result
    function getValueFromCandidateResult(candidate, value) {
      let returnValue
      candidatesResults.candidates.forEach(candidateResult => {
        if(candidate.firstname === candidateResult.firstname && candidate.lastname === candidateResult.lastname) {
          returnValue = candidateResult[value]
        }
      });
      return returnValue
    }

    function beautifyNumber(number) {
      return new Intl.NumberFormat('tr-TR').format(number)
    }

    onMount(async function () {
      // Encode the candidateNamesTogether as a base64 string
      let encoded = btoa(getCandidateNamesTogether())
      getCandidatesFromBilgi(encoded).then(result => {
        candidatesBilgi = result
        console.log("Candidates Bilgi")
        console.log(candidatesBilgi)
        getCandidatesResults().then(result => {
          candidatesResults = result
          candidatesResults.candidates = _.orderBy(candidatesResults.candidates, 'percentage', 'desc')
          candidates = _.orderBy(candidates, 'votes', 'desc')
          show = true
        });
      });
    });

</script>

{#if show}
<div class="mx-auto max-w-screen-lg px-3 py-6">
  <nav class="flex" aria-label="Breadcrumb">
    <ol class="inline-flex items-center space-x-1 md:space-x-3 mb-5">
      <li class="inline-flex items-center">
        <a href="/secim" class="inline-flex items-center text-sm font-medium text-gray-700 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
          <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10,1.375c-3.17,0-5.75,2.548-5.75,5.682c0,6.685,5.259,11.276,5.483,11.469c0.152,0.132,0.382,0.132,0.534,0c0.224-0.193,5.481-4.784,5.483-11.469C15.75,3.923,13.171,1.375,10,1.375 M10,17.653c-1.064-1.024-4.929-5.127-4.929-10.596c0-2.68,2.212-4.861,4.929-4.861s4.929,2.181,4.929,4.861C14.927,12.518,11.063,16.627,10,17.653 M10,3.839c-1.815,0-3.286,1.47-3.286,3.286s1.47,3.286,3.286,3.286s3.286-1.47,3.286-3.286S11.815,3.839,10,3.839 M10,9.589c-1.359,0-2.464-1.105-2.464-2.464S8.641,4.661,10,4.661s2.464,1.105,2.464,2.464S11.359,9.589,10,9.589"></path>
          </svg>
          Türkiye Geneli
        </a>
      </li>
      <li aria-current="page">
        <div class="flex items-center">
          <svg class="w-6 h-6 text-gray-400" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"></path></svg>
          <a href={"/secim/"+data.city.name} class="ml-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ml-2 dark:text-gray-400 dark:hover:text-white">{data.city.readablename}</a>
        </div>
      </li>
    </ol>
  </nav>
  
  <!-- Milletvekilligi -->
  <div class="p-4 w-full max-w-md bg-white rounded-lg border shadow-md sm:p-8 dark:bg-gray-800 dark:border-gray-700">
    <div class="flex justify-between items-center mb-4">
        <h5 class="text-xl font-bold leading-none text-gray-900 dark:text-white">Milletvekilliği</h5>
    </div>
    <div class="flow-root">
      <ul role="milletvekili-partileri" class="divide-y divide-gray-200 dark:divide-gray-700">
        {#each candidates as candidate}
          <li class="pt-3 pb-0 mb-5 sm:pt-4">
            <div class="flex items-center space-x-4">
                <div class="flex-shrink-0">
                    <img class="w-8 h-8 rounded-full" src={"data:image/gif;base64," + getValueFromCandidateBilgi(candidate, "image")} alt={candidate.lastname}>
                </div>
                <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900 truncate dark:text-white">
                        {candidate.firstname + " " + candidate.lastname}
                    </p>
                    <p class="text-sm text-gray-500 truncate dark:text-gray-400">
                        {getValueFromCandidateBilgi(candidate, "affiliation")}
                    </p>
                    <br>
                    <p class="text-sm text-gray-500 truncate dark:text-gray-400">
                        Oy: {beautifyNumber(candidate.votes)}
                    </p>
                    <div class="text-xs">
                      <div class="flex justify-between mb-1">
                          <span class="text-sm font-medium text-cyan-800 dark:text-white"></span>
                          <span class="text-xs font-medium text-cyan-800 dark:text-white">{getValueFromCandidateResult(candidate, "percentage").toFixed(2)}%</span>
                      </div>
                      <div class="w-full bg-gray-400 rounded-full h-2.5 dark:bg-gray-700">
                          <div class={"bg-"+ getValueFromCandidateBilgi(candidate, "color").js +"-800 h-2.5 rounded-full"} style={"width: "+45+"%"}></div>
                      </div>
                  </div>
                </div>
            </div>
          </li>
        {/each}
      </ul>
    </div>
  </div>

  <!-- <!-- Cumhurbaskanligi -->
  <!-- <div class="p-4 w-full max-w-md bg-white rounded-lg border mt-5 shadow-md sm:p-8 dark:bg-gray-800 dark:border-gray-700"> -->
  <!--   <div class="flex justify-between items-center mb-4"> -->
  <!--       <h5 class="text-xl font-bold leading-none text-gray-900 dark:text-white">Cumhurbaşkanlığı</h5> -->
  <!--   </div> -->
  <!--   <div class="flow-root"> -->
  <!--     <ul role="milletvekili-partileri" class="divide-y divide-gray-200 dark:divide-gray-700"> -->
  <!--       {#each candidates as candidate} -->
  <!--         <li class="pt-3 pb-0 mb-5 sm:pt-4"> -->
  <!--           <div class="flex items-center space-x-4"> -->
  <!--               <div class="flex-shrink-0"> -->
  <!--                   <img class="w-8 h-8 rounded-full" src="/docs/images/people/profile-picture-5.jpg" alt={candidate.lastname}> -->
  <!--               </div> -->
  <!--               <div class="flex-1 min-w-0"> -->
  <!--                   <p class="text-sm font-medium text-gray-900 truncate dark:text-white"> -->
  <!--                       {candidate.firstname + " " + candidate.lastname} -->
  <!--                   </p> -->
  <!--                   <p class="text-sm text-gray-500 truncate dark:text-gray-400"> -->
  <!--                       {candidate.affiliation} -->
  <!--                   </p> -->
  <!--                   <div class="text-xs"> -->
  <!--                     <div class="flex justify-between mb-1"> -->
  <!--                         <span class="text-sm font-medium text-cyan-800 dark:text-white"></span> -->
  <!--                         <span class="text-xs font-medium text-cyan-800 dark:text-white">{45}%</span> -->
  <!--                     </div> -->
  <!--                     <div class="w-full bg-gray-400 rounded-full h-2.5 dark:bg-gray-700"> -->
  <!--                         <div class="bg-red-800 h-2.5 rounded-full" style={"width: "+45+"%"}></div> -->
  <!--                     </div> -->
  <!--                 </div> -->
  <!--               </div> -->
  <!--           </div> -->
  <!--         </li> -->
  <!--       {/each} -->
  <!--     </ul> -->
  <!--   </div> -->
  <!-- </div> -->


</div>
{/if}
