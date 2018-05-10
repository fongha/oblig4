# Int Elligence(5);
### Systemspesifikasjon
Systemutviklingsprosjektet vårt har vært å lage en applikasjon for værmelding basert på data fra Openweathermap. Nytteverdien av applikasjonen er at man enkelt kan søke på sin ønskede by i søkefeltet, og en vil da få opp værvarselet for de fem neste dagene. Her vil man da i tillegg til å få informasjon om temperatur og vind, også få vite hvilket type skydekke eller eventuelt hvilken type nedbør det er snakk om for de gjeldende dagene. Værvarselet dekker hver tredje time for de fem kommende dagene, og her vil man da også få en kommentar basert på temperaturen, om man da for eksempel trenger en tykk ytterjakke eller om det er sommertid med shorts og t-skjorte.

Hvis man skal gå mer ned i dybden av nytteverdien til en værmelding, er det ingen tvil om at værvarsel for ekstremvær kan være kritisk i gitte situasjoner. Dersom det blir meldt i god tid at det kommer ekstremvær, vil personer i det aktuelle området kunne forberede seg og ta forhåndsregler. Alle utrykningsetatene kan dra stor fordel av et værvarsel, da de kan planlegge bestemte oppdrag og forberede akutte utrykninger som oppstår. Militæret vil også dra stor nytte av å benytte et værvarsel, da de til enhver tid må tilpasse flere rutiner basert på værsituasjon.

For å oppsummere, ser vi at samfunnet er avhengig av værvarsel; alt fra å skulle planlegge en telttur til at kommunen skal måke gatene når snøen kommer. Det er ikke til å legge skjul på at det handler om at vi alle får muligheten til planlegge de neste dagene basert på været.







### Systemarkitektur
Vår værapplikasjon består av: En klient som kjører applikasjonen i en browser, lokal webserver på port “:8080”, og databaseserverne til Openweathermap. Kildekoden/backend er skrevet i backend, mens HTML og CSS er brukt til presentasjon/front end.

Klienten vil gjennom applikasjonens kode opprette en forbindelse til webserver, og vil lytte etter innkommende TCP-pakker på nettverksporten “:8080”. Brukeren vil skrive inn en norsk by som input, og webserveren vil da sende en http.GET request til “Openweathermap” sine databaser om tilgang til data. Data vil bli sendt fra Openweathermap i form av en nettside med rå JSON, til serveren. Hvilket API som sendes fra Openweathermap blir bestemt av inputen til brukeren på klienten. Serveren vil da gi klienten en response med denne dataen.  Behandling av denne dataen skjer i kildekoden til applikasjonen, og dataen vil bli presentert til klienten ved hjelp av HTML og CSS. 
APIet fra Openweathermap blir oppdatert kontinuerlig. 
 
Den bearbeidede dataen som vises i applikasjonen blir lagret i RAM, da vi ikke har noen ekstern lagringsplass. Siden kildekoden ikke inneholder noen globale variabler, eller individuelle tråder, vil informasjonen blir kun hentet og bearbeidet når klienten sender en request i applikasjonen (runtime).
ILLUSTRASJON LENGER NED 

















