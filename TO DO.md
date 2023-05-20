## USER REQUIREMENTS:
  
  Awareness raising: 
    - sezione di HELP in cui vengono date alcune info su CER, BLOCKCHAIN, ecc..

  Energy community:
    - l'utente deve far parte di una CER per partecipare e per ottenere i certificati necessari e le credenziali di accesso
    - "noi" interagiamo con il delegato della CER e creiamo tanti peer e credenziali quanti sono i membri

  User interface:
    - statistiche:
      - 

  Ease of Use:
    - wallet e storico transazioni visibili sull'app

  Privacy:
    - transazioni anonime (peer che partecipano alle transazioni sono criptati)
    - no associazione tra peer e utente fisico


## FUNCTIONAL REQUIREMENTS:

  Blockchain:
    - abbiamo rispettato i requisiti di una blockchain 

  Security:
    - HASH & CA 

  Smart Contracts:
    - non puoi non mandare i soldi/energia (è tutto automatico)

  Virtual currency:
    - creiamo un asset per le monete (future implementazioni, integrare con altre applicazioni tipo car sharing)

  Mobile application&web site:
    - 


## TECHNICAL REQUIREMENTS:

  Programming languages:
    - 
  
  Data analysis:
    - thingspeak 

  Simulation:
    - creiamo un simulation environment in cui facciamo vedere tot transazioni
    - la durata della simulazione sarà lunga abbastanza da ottenere sufficienti dati statistici (24/48h)
    - il numero di peer sarà maggiore di 2 => TEST 
      - scalabilità sia all'interno delle CER che all'esterno (possibilità di creare più canali e quindi più CER)

  Hardware and application protocols:
    - arduino a prova di faults 


Private Data Collection => come funzionano? le usiamo?

Quale asset transfer è da copiare?

Asta => più winner (come già fatto in auction dutch), l'asset di energia da trasferire viene diviso in più asset
=> TEST
  - offerta > domanda con avanzo di energia a fine asta => asset energetico rimanente viene distrutto e l'energia viene "venduta" alla grid
  - offerta < domanda => domanda del primo bidder soddisfatta in parte 

Una volta ricevuto l'asset energetico viene registrato e poi distrutto
Una volta ricevuto l'asset monetario viene inglobato nel wallet del peer => TEST di SICUREZZA (viene modificato correttamente l'asset?)

Come forzare il trasferimento di energia e soldi ad asta chiusa? => TEST DA FARE

Nodo delle offerte: percentuale dell'energia prodotta in eccesso viene automaticamente trasferita al nodo delle donazioni 
(il nodo non partecipa all'asta ma riceve tipo l'1% di tutta l'energia prodotta in eccesso)

## SMART CONTRACTS:

### SBE : State Based Endorsement, ovvero particolari smart contract che permettono di modificare l'endorsement per uno specifico asset (per esempio dopo aver creato un asset usando l'endorsement policy di default della chaincode specifico che d'ora in poi per agire su quello specifico asset non è necessario nuovamente l'endorsement degli altri peer, basta la conferma del proprietario)

  - #### Energy-contract 
    - è realmente necessario? Oltre al nodo delle donazioni è necessario 
      avere un asset energetico che viene creato per essere distrutto?

    - funzioni
      - createAsset
      - updateAsset
      - deleteAsset
      - transferAsset

    - non è necessario che sia salvato in una private data collection, tutti devono poter vedere quanta energia ho a disposizione da vendere
    - è fondamentale che un asset possa essere modificato solo dal proprietario
    - endorsement policy della chaincode: maggioranza dei peer nella rete => modificabile usando le SBE? (ho bisogno della maggioranza per creare l'asset e trasferirlo ma per update e delete invece no => controllare se le SBE sono specifiche per funzione dello smart contract o per asset)

  - #### Money-contract

    - funzioni
      - createAsset
      - updateAsset
      - deleteAsset
      - transferAsset

    - le informazioni relative ai soldi non devono mai essere in chiaro:
      - l'asset monetario di un peer deve essere salvato in una Private Data Collection
      - creazione e update devono essere oscuri a tutti i peer

      - se quando viene eletto il vincitore dell'asta viene anche pubblicato il denaro che è stato offerto allora non è necessario nascondere la quantità di denaro trasferito come conseguenza dell'acquisizione dell'asset energetico (basterebbe guardare la quantità di energia che è stata trasferita e moltiplicarla per il denaro offerto)

  - #### Auction (dutch)
    - perché il riferimento è solo alle organizzazioni e non al singolo peer?
    - il denaro nelle offerte (sia da parte del venditore che degli acquirenti) sempre valutato in €/unità di energia

    - quando vengono pubblicate le offerte degli acquirenti pubblichiamo il denaro offerto giusto?

  - #### Application 
    - è uno script esterno alla blockchain presente su un device (raspberry) che è immutabile (tamper-proof)








