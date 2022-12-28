<script>
    /** @type {import('./$types').PageData} */
    export let data;
    console.log(data.city)
</script>

<div class="mx-auto max-w-screen-lg px-3 py-6">
   <div class="text-sm breadcrumbs">
       <ul>
           <li><a href="/secim">Türkiye Geneli</a></li> 
           <li><a href={"/secim/"+data.city.name}>{data.city.readablename}</a></li>
           <li>{data.district.readablename}</li>
       </ul>
   </div>
</div>