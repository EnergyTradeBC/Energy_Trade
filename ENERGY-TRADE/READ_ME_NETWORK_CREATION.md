E' stato creato il file SETUP_WHOLE_NETWORK.sh, che, se lanciato, crea tutto il network fino ad arrivare alla readiness della chaincode.
Di seguito verranno spiegati tutti i comandi e tutte le scelte, sperando siano i più chari possibile.
Nota: non è ancora stata divisa la creazione in creazione che facciamo noi (cosiddetto NETWORKUP) e in simulazione degli utenti (cosiddetto NETWORKUP2). Questo perchè ci sono
ancora delle cose su cui lavorare che sono elencate e spiegate di seguito. Quando tutto funziona, verrà fatto.

./network.sh up

Semplicemente tira su la rete con 1 orderer e 2 orgs (Org1 e Org2). La spiegazione di questa scelta viene spiegata in seguito (vicino al simbolo * nei paragrafi successivi)

./scripts/createChannel.sh genesis -c channeltest

Crea il canale "channeltest". Tale canale è configurato per eventualmente ospitare esattamente Org1 e/o Org2. Spiegazione sulla configurazione verrà sempre eseguita in
seguito (vicino al simbolo *).

./scripts/createChannel.sh join 1 -c channeltest
./scripts/createChannel.sh join 2 -c channeltest

Aggiunge Org1 e Org2 al canale "channeltest" (nota: questo viene correttamente eseguito perchè il canale si aspettava queste organizzazioni).

./scripts/addOrg.sh files 3 -c channeltest
./scripts/addOrg.sh add 3 -c channeltest

Crea i file relativi all'Org3 e successivamente crea l'Org3 stessa. Notare come non è più necessario specificare le porte in nessuno dei due comandi.
Questo perchè le stesse porte dovevano essere specificate più volte in seguito, il che avrebbe reso molti script hard-coded, perciò questa è stata la soluzione.
Il primo comando richiedeva sempre porte nella forma XX051 XX054, il secondo XX051 XX052. Visto che l'Org1 aveva XX = 7 e l'Org2 XX = 9, allora ho deciso di
generalizzare come segue: data una generica OrgY, XX sarà calcolata come Y + 8. Quindi Org3 utilizzerà porte 11051-11052-..., Org4 12051-12052-..., Org5 13051-13052-... e così via.

In molti script veniva utilizzata la funzione setGlobals() (e setGlobalsCLI()) contenuta in envVar.sh, la quale permetteva di agire come uno dei peer del network.
Funzionava solo per Org1, Org2 e Org3 e tutti i path e le porte erano hard-coded. Per questo motivo ho creato un nuovo script (actAsOrg.sh) che contiene le funzioni
actAdOrg() e actAsOrgCLI(): queste funzioni sono l'esatto equivalente di setGlobals() e setGlobalsCLI() con la differenza che funzionano per ogni Org calcolando
automaticamente path e porte (secondo la convenzione spiegata prima). Il motivo per cui ho cambiato il nome delle funzioni è per far sì che si possa comunque
continuare a importare lo script envVar.sh per tutto il resto che contiene senza che così facendo vengano importate due funzioni omonime.

./scripts/createChannel.sh updateConfig 3 -c channeltest

Questo comando è importantissimo: cambia la configurazione del canale "channeltest" in modo che si prepari ad ospitare anche l'Org3.
Se questo comando non venisse eseguito e si passasse subito al successivo (./scripts/createChannel.sh join 3 -c channeltest) per aggiungere l'Org3 al canale,
questo non darebbe errore, ma l'Org3 non sarebbe propriamente parte del canale. Questo lo si nota arrivati alla chaincode, infatti Org3 non verrebbe riconosciuta
e mancherebbe nell'elenco di tutte le Org che devono approvare la chaincode stessa.
Quello che questo comando fa, a livello più basso, è prendere la configurazione attuale del canale, modificarla per includere Org3, e poi far firmare (per approvazione)
questa modifica prima a Org1 (hard-coded) e poi a Org2(hard-coded). Quando entrambe le firme vengono raccolte, la configurazione viene a tutti gli effetti cambiata e 
Org3 può aggiungersi al canale.
Eseguendo alcuni test ho scoperto che:
Se la configurazione del canale ha Org1 e Org2, allora Org3 può entrare con le firme di Org1 e Org2 (100% di firmatari)
Se la configurazione del canale ha Org1, Org2 e Org3, allora Org4 può entrare con le firme di Org1 e Org2 (66.6% di firmatari)
Se la configurazione del canale ha Org1, Org2, Org3 e Org4, allora Org5 NON può entrare con le firme di Org1 e Org2 (50% di firmatari)
Da questo si evince che probabilmente è necessario il 50% + 1 di firme per entrare.
Sarà necessario fare modifiche per decidere in maniera dinamica chi firma la modifica.
(*) Da questo però si capisce una cosa, ovvero che il canale, quando viene creato, contiene già la configurazione per ospitare Org1 e Org2 di default. Lo script che causa ciò è
uno di quelli talmente grandi che VisualStudio si rifiuta di aprirlo. Per questo motivo innanzitutto sarebbe letteralmente impossibile scoprire come creare un canale senza
questa configurazione iniziare, e inoltre non sappiamo neanche se è possibile creare un canale vuoto, visto che per aggiungere un Org è necessario cambiare configuazione e
farla firmare da org già presenti nel canale.
Da ciò ne deriva che non ha neanche senso creare un networkup senza Org1 e Org2, visto che quando si crea poi il canale queste due Org ci devono essere comunque.
P.S.: questa cosa veniva fatta solo per l'Org3 utilzzando i seguenti file:
org3-scripts/updateChannelConfig.sh
configUpdate.sh
che però erano hard-coded e funzionanti solo per l'Org3
Ho quindi creato una versione generica sotto forma di:
org-scripts/updateChannelConfig.sh
configUpdateModified.sh
(nota: org-scripts/joinChannel.sh è al momento inutilizzato, ma può servire, quindi per favore lasciatelo e ignoratelo)

./scripts/createChannel.sh join 3 -c channeltest

Easy, dopo il casino di poco fa, ora possiamo aggiungere Org3 al canale.

./scripts/deployCC.sh packageChaincode -ccn chaincodetest -ccp ../asset-transfer-basic/chaincode-go

Facciamo il deploy della chaincode "chaincodetest" ottenuta da "asset-transfer-basic".

./scripts/deployCC.sh installChaincode 1 -ccn chaincodetest
./scripts/deployCC.sh installChaincode 2 -ccn chaincodetest
./scripts/deployCC.sh installChaincode 3 -ccn chaincodetest

Installiamo la chaincode "chaincodetest" sui peer delle 3 org.

./scripts/deployCC.sh checkCommitReadiness 1 -ccn chaincodetest -c channeltest

Verifichiamo che per ora la chaincode NON è ancora stata approvata da nessuna delle 3 org.
E' a questo punto che se non abbiamo eseguito il comando per fare l'update della configurazione del canale (./scripts/createChannel.sh updateConfig 3 -c channeltest)
non vediamo l'Org3.

./scripts/deployCC.sh approveForMyOrg 1 -ccn chaincodetest -c channeltest
./scripts/deployCC.sh approveForMyOrg 2 -ccn chaincodetest -c channeltest
./scripts/deployCC.sh approveForMyOrg 3 -ccn chaincodetest -c channeltest

Approviamo la chaincode per tutte e 3 le org.

./scripts/deployCC.sh checkCommitReadiness 1 -ccn chaincodetest -c channeltest

Verifichiamo che ora la chaincode è stata approvata da tutte e 3 le org.

(Nota finale: ora lavoro sul network down :D)