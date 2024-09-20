import Keycloak from 'keycloak-js';

const KEYCLOAK_URL=import.meta.env.VITE_KEYCLOAK_URL
const KEYCLOAK_REALM=import.meta.env.VITE_KEYCLOAK_REALM
const KEYCLOAK_CLIENTID=import.meta.env.VITE_KEYCLOAK_CLIENTID

const keycloak = new Keycloak({
    url: KEYCLOAK_URL,
    realm: KEYCLOAK_REALM,
    clientId: KEYCLOAK_CLIENTID,
});

export default keycloak;
