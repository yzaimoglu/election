<script>
    /** @type {import('./$types').PageData} */
    export let data;
    let districtName = data.districtReadableName
    let cityName = data.cityReadableName
    let district = data.district
</script>

<div class="container">
    {#if district.name != null}
        <h1>{cityName} - {districtName}</h1>
        <div class="alert alert-success" role="alert">
            <p>
            Oy kullanabilen: {district.eligiblevoters}
            <br>
            Gecerli oy: {district.validvotes}
            <br>
            Gecersiz oy: {district.invalidvotes}
            </p>
            <hr>
            <div class="mb-0">
                {#each district.candidates as candidate}
                    <p>
                        Kisi: {candidate.lastname}
                        <br>
                        Oy: {candidate.votes}
                    </p>
                {/each}
            </div>
        </div>
    {:else}
        <div class="alert alert-danger" role="alert">
            Böyle bir ilce yok
        </div>
    {/if}
</div>

