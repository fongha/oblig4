<h1>Systemarkitektur</h1><br>
Vår værapplikasjon består av: En klient som kjører applikasjonen i en browser, lokal webserver på port “:8080”, og databaseserverne til Openweathermap. Kildekoden/backend er skrevet i Go, mens HTML og CSS er brukt til presentasjon/front end.
<br><br>
Klienten vil gjennom applikasjonens kode opprette en forbindelse til webserver, og vil lytte etter innkommende TCP-pakker på nettverksporten “:8080”. Brukeren vil skrive inn en norsk by som input, og webserveren vil da sende en http.GET request til “Openweathermap” sine databaser om tilgang til API for den byen. APIet vil bli sendt fra Openweathermap i form av en nettside med rå JSON, til serveren. Hvilket API som sendes fra Openweathermap blir bestemt av brukerens input på klienten. Serveren vil da gi klienten en response med denne dataen. Behandling av denne dataen skjer i kildekoden til applikasjonen, og dataen vil bli presentert til klienten ved hjelp av HTML og CSS. 
APIet fra Openweathermap blir oppdatert kontinuerlig. <br><br>

Den bearbeidede dataen som vises i applikasjonen blir lagret i RAM, da vi ikke har noen ekstern lagringsplass. Siden kildekoden ikke inneholder noen globale variabler, eller individuelle tråder, vil informasjonen blir kun hentet og bearbeidet når klienten sender en request i applikasjonen (runtime).
<br><br>
<img src="https://i.imgur.com/1ngqdTH.jpg" alt="Systemarkitektur">
