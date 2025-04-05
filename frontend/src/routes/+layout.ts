import Keycloak, { type KeycloakInitOptions } from 'keycloak-js';
import type { LayoutLoad } from './$types';
import { browser } from '$app/environment';

export const load: LayoutLoad = async () => {
	const instance = {
		url: 'http://127.0.0.1:8080',
		realm: 'budget',
		clientId: 'app'
	};

	const keycloak = new Keycloak(instance);
	const initOptions: KeycloakInitOptions = { onLoad: 'login-required', checkLoginIframe: false };

	let keycloakPromise;
	if (browser) {
		keycloakPromise = keycloak.init(initOptions).then((auth) => {
			if (auth) {
				return keycloak;
			}
		});
	}

	return {
		keycloak: keycloakPromise
	};
};
