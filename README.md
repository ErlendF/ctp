# IMT2681 Cloud Technologies Project
###Authors
Erlend Fonnes, Johan Selnes, Aksel Baardsen, Knut Jørgen Totland, Benjamin Skinstad

##Project report
### Project description and ambitions
The original plan of this project was to create a RESTful web application that allowed users to register accounts where the information about the playtime on games they play is calculated from official API's. This application should then return the total time spent playing games. 
We also planned to have automatic deployment of the application in Docker on Openstack via the CI/CD feature in Gitlab.
Lastly, if we had time, we would expand the functionality of the application.

Most has been achieved, but we did not have time to expand the functionality as per our ambitions. 


### Reflection


### Learning outcome


### Total work hours
The total work hours spent on this project is a little over 100 hours.
To track the group's work hours we used https://toggl.com/app/timer.

##Application information and setup
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



#### API endpoints
All enpoints start with "/api/v1/", thus the prefix has been omitted from the listing bellow. For the enpoints requiring authentication, the **Authorization** header needs to contain a valid JWT, as specified in the Authentication (usage) section.


No authentication:
```
/login                              (GET): Redirects to Googles OAuth consent screen, used for the user to login.
/authcallback                       (GET): The redirect URI where the user is returned after loging in. Returnes a JWT used for authentication for the enpoints listed above.
/user/{username:[a-zA-Z0-9 ]{1,15}} (GET): Get information about a pulbic user with a username.
```


Requires authentication:
```
/user         (GET): Returns all information about the user themselves.
/user        (POST): Updates information about the user themselves.
/user      (DELETE): Deletes the user and all information related to them.
/updategames (POST): Fetches new data from the servies registered for the user.
```

To update the user information, "/user" endpoint expects the following body for the POST request (values may be replaced, although they are required to be valid):
```
{
	"username": "newUsername",
	"lol": {
		"summonerName": "LOPER",
		"summonerRegion": "EUW1"
	},
	"valve": {
		"username": "olaroa3"
	},
	"overwatch": {
		"battleTag": "Onijuan-2670",
		"platform": "pc",
		"region": "eu"
	}
}
```

For the "/updategames" path, the request body is ignored.


#### Application structure
The application is split into two main parts: *cmd* and *pkg*. *cmd* serves as the central function of the application. *pkg* contains everything that is either used by *cmd or another package in pkg*. We consider the user to be the central part of the application as all actions and information is related to or belongs to the user. Therefore, the handler only takes a UserManager as a parameter and the **handler struct in pkg/server/handler.go [embedds](https://travix.io/type-embedding-in-go-ba40dd4264df) the UserManager**, allowing the handler to use each of the functions specified in the *UserManager interface*. The handler functions themselves contain a minimum amount of logic, merely calling functions from the UserManager, thus only handling i/o and logging.


**The Models package contains interfaces, structs, constants and function which are used by several packages to simplify the internal dependency graph.** For example, every struct used by multiple packages is defined in Models. Defining the struct in either of the packages would therefore create a direct dependency between them (or be a duplication). 


Interfaces are widely used throughout the application to facilitate testing. This makes it possible to mock them, reducing the scope of the test. Interfaces are also used for the handler, where they serve to decouple the packages from eachother, preventing several direct dependencies. The **Organizer** interface is used to combine the interfaces for all the packages which provide games and stats to the application, simplifying the passing of the interfaces to the *UserManager*. Similar to how *handler* embedds the UserManager interface, the UserManager struct (which fulfils the interface) embedds the *organizer*, allowing it to call each of the functions specified in the *organizer interface* (and every interface within it).



#### Repository structure
The repository has the following main components:
 - **cmd**: Lists all possible commands for the application. Currently, there are none other than root. Main.go serves merely to start the *Run* function of cmd/root.go. Was created by Cobra during project initialization.
 - **pkg**: Contains all packages used in the application. See Application structure.
 - **.gitignore**: Specifies what files should be ignored by git.
 - **go.mod** and **go.sum**: Go modules.
 - **LICENSE**: Apache license, created by Cobra during project initialization.
 - **main.go**: Starts the application.
 - **README.md**: The file you are currently reading.


##### Server setup is based on [gorilla/mux graceful-shutdown example](https://github.com/gorilla/mux#graceful-shutdown)