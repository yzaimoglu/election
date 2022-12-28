<html lang="tr" class="bg-slate-900 text-gray-100">
    <body>
        <div class="mx-auto max-w-screen-lg px-3 py-6">
            <div class="flex flex-col gap-y-3 sm:flex-row sm:items-center sm:justify-between">
                <a href="/demo/astro-boilerplate">
                    <div class="flex items-center bg-gradient-to-br from-sky-500 to-cyan-400 bg-clip-text text-xl font-bold text-transparent">
                        <a class="flex gap-x-1 font medium text-gray-200" href="/">
                            SEÇİM 2023
                        </a>
                    </div>
                </a>
                <nav>
                    <ul class="flex gap-x-3 font-medium text-gray-200">
                        <li class="hover:text-white">
                            <a href="/secim">Sonuçlar</a>
                        </li>
                        <li class="hover:text-white">
                            <a href="/yardim">Yardım</a>
                        </li>
                        <li class="hover:text-white">
                            <a href="/proje">Proje</a>
                        </li>
                        <li class="text-cyan-600 hover:text-white">
                            <a href="https://intern.2023secimi.org">Giriş</a>
                        </li>
                    </ul>
                </nav>
            </div>
        </div>
        <slot></slot>
        <div class="mx-auto max-w-screen-lg px-3 py-6">
            <div class="border-t border-gray-600 pt-5">
                <div class="text-sm text-gray-200">
                    © Copyright 
                    2022
                    - 
                    Seçim 2023 ekibi
                </div>
            </div>
        </div>
    </body>
</html>