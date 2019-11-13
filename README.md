# IMT2681 Cloud Technologies Project

### Erlend Fonnes, Johan Selnes, Aksel Baardsen, Knut Jørgen Totland, Benjamin Skinstad


#### Setup
The application uses firebase and requires a credential file called **FBKEY.json**, unless other name is passed as a command line argument.
The following environment variables are required to be specified:

```
VALVE_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
RIOT_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GOOGLE_OAUTH2_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GOOGLE_OAUTH2_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxx
HMAC_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

It is intended for these to be put in an **.env** file (just like sample.env, replacing the x's), which is injected into the environment variables for the running application by [joho/godotenv/autoload](https://github.com/joho/godotenv), which is imported in cmd/root. Whichever way they are added to the environment for the application, they are required to be present with valid values for the application to run.


The application accepts the following commandline arguments:
```
 -h, --help                  Help for CTA2
 -p, --port int              Specifies which port the API should listen to (default 80)
 -d, --domain string         Specifies the domain for the redirect URI used for authentication (default "localhost")
 -f, --fbkey string          Path to the firebase key file (default "./FBKEY.json")
 -v, --verbose               Verbose logging
 -j, --jsonFormatter         JSON logging format
 -s, --shutdownTimeout int   Sets the timeout (in seconds) for graceful shutdown (default 15)
 -c, --clientTimeout int     Sets the timeout (in seconds) for the http client which makes requests to the external APIs (default 15)
```

#### Authentication
###### Configuration
OpenID Connect with Google as the provider is used to authenticate users of the application. Therefore, valid Google OAUTH2 credentials are required, like shown in the Setup section. In addition, the [Google APIs Project](https://console.developers.google.com/) needs to be configured with scope as *email*, *profile* and *openid*, although only the **openid** scope is actually used (neither email nor profile are stored in the application). To our knowledge, it is currently not possible to reduce the scope further. The project also needs "http://%s:%d/api/v1/authcallback" to be set as a **Authorised rediredt URI**, where "%s" replaced with applicable domain and "%d" with the desired port.


###### Usage
To login to the application, the user should send a GET request to /api/v1/login. This route should redirect the user to Googles OAuth consent screen, where the user needs to be signed in to a Google account and accept sending the required data to the application. The user is then redirected back to the application (/api/v1/authcallback), where a JWT token is sent back unless some error has occured. This token should be sent with every request requiring authentication as the **Authorization** header. Verification of the token is handled by the *auth middleware*.


##### Server setup is based on [gorilla/mux graceful-shutdown example](https://github.com/gorilla/mux#graceful-shutdown)