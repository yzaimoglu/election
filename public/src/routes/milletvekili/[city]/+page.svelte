<script>
    /** @type {import('./$types').PageData} */
    export let data;
    let cityName = data.cityReadableName
    let city = data.city
    let candidates = data.candidates
</script>


    <div class="mx-auto max-w-screen-lg px-3 py-6">
        <div class="flex flex-col items-center md:flex-row md:justify-between md:gap-x-24">
            <div>
                <h1 class="text-3xl font-bold">
                    Türkiye <span class="bg-gradient-to-br from-sky-500 to-cyan-400 bg-clip-text text-transparent">Seçim</span> Sonuçları 
                </h1>
                <p class="mt-6 text-xl leading-9 text-justify">
                    Bu sayfada Türkiyedeki 2023 seçimleriyle ilgili verilere ulaşabilirsin. İstersen
                    <a class="text-cyan-400 hover:underline" href="/milletekili/il">iller</a>
                    bazında sonuçlara, istersende
                    <a class="text-cyan-400 hover:underline" href="/milletvekili/ilce">ilçeler</a>
                    bazında veya
                    <a class="text-cyan-400 hover:underline" href="/milletvekili/mahalle">mahalleler</a>
                    bazında sonuçlara ulaşabilirsin.
                </p>
            </div>
        </div>
    </div>

    <div class="mx-auto max-w-screen-lg px-3 py-6">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-4">
            <a class="hover:translate-y-1" href="/secim">
                <div class="overflow-hidden rounded-md bg-slate-800">
                    <div class="aspect-square">
                        <img class="h-full object-cover object-center" src="/images/turkey_flag.webp" alt="Türkiye" loading="lazy"/>
                    </div>
                    <div class="px-3 pt-4 pb-6 text-center">
                        <h2 class="text-xl font-semibold">Türkiye Geneli</h2>
                        <div class="mt-2 text-sm">Türkiye genelindeki sonuçları görmek için buraya bas</div>
                    </div>
                </div>
            </a>
            <a class="hover:translate-y-1" href="/secim/il">
                <div class="overflow-hidden rounded-md bg-slate-800">
                    <div class="aspect-square">
                        <img class="h-25 object-cover object-center" src="/images/ankara_map.png" alt="Ankara" loading="lazy"/>
                    </div>
                    <div class="px-3 pt-4 pb-6 text-center">
                        <h2 class="text-xl font-semibold">İller</h2>
                        <div class="mt-2 text-sm">İl bazında sonuçları görmek için buraya bas</div>
                    </div>
                </div>
            </a>
            <a class="hover:translate-y-1" href="/secim/ilce">
                <div class="container mx-auto bg-cyan-600">
                    <div class="overflow-hidden rounded-md bg-slate-800 h-max">
                        <div class="aspect-square">
                            <img class="h-50 object-cover object-center" src="/images/ankara_map.png" alt="Çankaya" loading="lazy"/>
                        </div>
                        <div class="px-3 pt-4 pb-6 text-center">
                            <h2 class="text-xl font-semibold">İlçeler</h2>
                            <div class="mt-2 text-sm">İlçe bazında sonuçları görmek için buraya bas</div>
                        </div>
                    </div>
                </div>
            </a>
            <a class="hover:translate-y-1" href="/secim/mahalle">
                <div class="overflow-hidden rounded-md bg-slate-800">
                    <div class="aspect-square">
                        <img class="h-50 w-full object-cover object-center" src="/images/ankara_map.png" alt="Çukurambar Mah." loading="lazy"/>
                    </div>
                    <div class="px-3 pt-4 pb-6 text-center">
                        <h2 class="text-xl font-semibold">Mahalleler</h2>
                        <div class="mt-2 text-sm">Mahalle bazında sonuçları görmek için buraya bas</div>
                    </div>
                </div>
            </a>
        </div>
    </div>

<div class="text-3xl font-bold">
    {#if city.name != null}
        <h1>{cityName}</h1>
        <div class="alert alert-success" role="alert">
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
    {:else}
    <div class="alert alert-danger" role="alert">
        Böyle bir sehir yok
    </div>
    {/if}
</div>

