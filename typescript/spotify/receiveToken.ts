import express from 'express';
import { URLSearchParams } from 'url';
import api  from './api';
import apiConfig from '../../config/api_config.json';
import tokenConfig from '../../config/token_config.json';

const config = tokenConfig.spotify;

function server() {
    const parsedRedirectUri = new URL(config.redirect_uri);
    const app = express();
    
    app.get('/', (req, res) => {
        const queryParams = api.getAuthorizationQueryParams(config.client_id, config.redirect_uri);
        const searchParams = new URLSearchParams(queryParams);
        return res.redirect(`${apiConfig.spotify.urls.authorization}?${searchParams}`);
    })

    app.get(parsedRedirectUri.pathname, async (req, res) => {
        const code = req.query.code as string;
        const response = await api.requestAccessToken(config.client_id, config.client_secret, code, config.redirect_uri);
        res.send({"access_token": response.access_token});
    })

    const [port, hostname] = [Number(parsedRedirectUri.port), parsedRedirectUri.hostname];

    app.listen(port, hostname, () => {
        console.log(`Starting to listen on http://${hostname}:${port}`);
    });
}

server();
