<script>
    // @ts-nocheck
        /** @type {import('./$types').PageData} */
        export let data;
        let quarters = data.quarters
        let districts = data.districts
        let constituencies = data.constituencies
        let cities = data.cities
</script>
    <div class="container">
        <h1>Mahalleler</h1>
        {#each Object.values(quarters) as quarter}
            <div class="alert alert-success" role="alert">
                <h4 class="alert-heading">Mahalle: {quarter.name}</h4>
                <p>Il: {quarter.city} - {quarter.citynumber}</p>
                <p>
                <strong>Secim bölgesi ve ilce: {quarter.constituency} - {quarter.district}</strong>
                <br>
                Oy kullanabilen: {quarter.eligiblevoters}
                <br>
                Gecerli oy: {quarter.validvotes}
                <br>
                Gecersiz oy: {quarter.invalidvotes}
                </p>
                <hr>
                <div class="mb-0">
                    {#each quarter.candidates as candidate}
                        <p>
                            Kisi: {candidate.lastname}
                            <br>
                            Oy: {candidate.votes}
                        </p>
                    {/each}
                </div>
            </div>
        {/each}
        
        <h1>Ilceler</h1>
        {#each Object.values(districts) as district}
            <div class="alert alert-secondary" role="alert">
                <h4 class="alert-heading">Ilce: {district.name}</h4>
                <p>Il: {district.city} - {district.citynumber}</p>
                <p>
                <strong>Secim bölgesi: {district.constituency}</strong>
                <br>
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
        {/each}
        
        <h1>Secim bölgeleri</h1>
        {#each Object.values(constituencies) as constituency}
            <div class="alert alert-secondary" role="alert">
                <h4 class="alert-heading">Secim bölgesi: {constituency.name}</h4>
                <p>Il: {constituency.city} - {constituency.citynumber}</p>
                <p>
                Oy kullanabilen: {constituency.eligiblevoters}
                <br>
                Gecerli oy: {constituency.validvotes}
                <br>
                Gecersiz oy: {constituency.invalidvotes}
                </p>
                <hr>
                <div class="mb-0">
                    {#each constituency.candidates as candidate}
                        <p>
                            Kisi: {candidate.lastname}
                            <br>
                            Oy: {candidate.votes}
                        </p>
                    {/each}
                </div>
            </div>
        {/each}
        
        <h1>Iller</h1>
        {#each Object.values(cities) as city}
            <div class="alert alert-secondary" role="alert">
                <h4 class="alert-heading">Il: {city.name} - {city.number}</h4>
                <p>
                Oy kullanabilen: {city.eligiblevoters}
                <br>
                Gecerli oy: {city.validvotes}
                <br>
                Gecersiz oy: {city.invalidvotes}
                </p>
                <hr>
                <div class="mb-0">
                    {#each city.candidates as candidate}
                        <p>
                            Kisi: {candidate.lastname}
                            <br>
                            Oy: {candidate.votes}
                        </p>
                    {/each}
                </div>
            </div>
        {/each}
    </div>
    
    
    
    
    