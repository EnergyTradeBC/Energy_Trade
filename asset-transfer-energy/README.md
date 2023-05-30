## COMMENTI
Non so ancora come vada formattata la cartella (chaincode) che contiene lo smart contract e di cui poi va fatto il deploy

- Prima di fare qualsiasi operazione faccio un check sull'identità con la funzione "verifyClientOrgMatchesPeerOrg"

- per il momento l'owner ID che viene utilizzato all'interno dell'asset è l'id del client stesso (Fabric Application) che ha effettuato la richiesta di creazione dell'asset => idee migliori?  

- Che convenzione usiamo per l'asset ID? In fase di creazione possiamo accedere ad informazioni presenti nel Fabric Application ma quando invece dobbiamo fare un "transferAsset", creando un nuovo asset per il compratore, possiamo usare informazioni molto più limitate