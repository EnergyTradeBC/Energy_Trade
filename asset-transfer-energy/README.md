## COMMENTI
Non so ancora come vada formattata la cartella (chaincode) che contiene lo smart contract e di cui poi va fatto il deploy

- Prima di fare qualsiasi operazione faccio un check sull'identità con la funzione "verifyClientOrgMatchesPeerOrg"

- Per il momento "verifyClientOrgMatchesPeerOrg" non funziona perché ctx.GetClientIdentity().GetMSPID() da errore e blocca l'endorsement
  ==> come controlliamo che il peer sia corretto e possa lavorare lì sopra? Possiamo controllare il clientID?

- il client ID che viene controllato e usato come owner degli asset è un client ID che viene formato usando il certificato (subject + issuer)
  ==> siccome non si riesce a controllare che l'MSP del client sia lo stesso del peer possiamo aggiungere come variabile globale del peer anche il clientID
  oltre il MSP ID in modo da creare un binding 1 a 1 tra peer e client (nello smart contract si controllerà quindi che i due ID corrispondano)

- Che convenzione usiamo per l'asset ID? ==> per il momento l'asset ha nu ID hardcoded formato da <energy_> + clientID
  infatti ogni volta che viene scambiato un asset oltre a venire cambiato l'ownerID viene cambiato anche l'assetID 
  in modo che ogni client possa accedere al proprio asset usando sempre lo stesso assetID 