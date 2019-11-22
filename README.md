# IMT2681 Cloud Technologies Project
### Authors
Erlend Fonnes, Johan Selnes, Aksel Baardsen, Knut Jørgen Totland, Benjamin Skinstad

## Project report
### Project description
#### Original plan
The original plan of this project was to create a RESTful web application that allows users to register accounts where the information about the playtime on games they play is calculated from other APIs. The user should be able to register their accounts for four different game "providers": Blizzard, Jagex, Valve and Riot Games. This application should then show the total time spent playing games.

We also planned to have automatic deployment of the application in Docker on Openstack via the CI/CD feature in Gitlab. We also wanted to use CI feature to run tests automatically and run linting tools.

Lastly, if we had time, we would expand the core functionality of the application or add additional games (or "providers") to the service.


#### Achievements (what has and has not been achieved/changed in the final product)
We managed to let users create accounts, and save games to that account, using OpenID Connect for authentication. All useful data is saved in Firebase for persistent storage. Nearing the end of the assignment we implemented all wanted CI/CD functionality.

All goals in our original plan has been achieved, except for the fact that we did not have time to expand the core functionality outside of playtime (as per our ambitions). <!---(we managed to add some extra functionality to the jagex account display)  <-- må skrive hva/referere til det hvis vi skal ha med dette

### Reflection
#### What went well
<!--Denne sectionen trenger innvoller & peer review-->
We managed to implement wanted core functionality. Tests run without issues.

In the end we used CI/CD for both deploying and linting. Use of [spf13/cobra](https://github.com/spf13/cobra) and [gorilla/mux](https://github.com/gorilla/mux) worked splendidly.


#### What went wrong
We underestimated the workload needed to complete this project.
We used ~110 hours on this project, while 75 hours was expected.
Riot did not have time to process our application for a permanent API key, so we were only able to use a 24-hour personal API key (application process time was longer than the project timeperiod).
This Riot API key problem was fixed with a hack to allow a new and valid Riot API key to be injected into the running application. To do this a POST request should be sent to the "/riotapikey" endpoint. The API key should be sent in the body as shown bellow:
```
RGAPI-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```


#### Hard aspects
 - Managing time
 - Prioritizing important aspects
 - Distributing workload
 - Implementing enough meaningful tests to reach 75% coverage (see testing section)

### Learning outcome
During the run of this project the group members have learned how to work with authentication, cobra file-structure, go testing using mocks and interfaces, OAuth2 and authentication using Google, Gitlab CI/CD, and that proper documentation makes many hassles go away. The group also discovered problems with local deployment concerning Google's *Authorized redirect URIs* (see Authentication OAuth workaround).

### Total work hours
The total work hours spent on this project is a little over 110 hours.
To track the group's work hours we used https://toggl.com/app/timer.



## Application information and setup
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

###### OAuth2 workaround
For OAuth2, it is recommended to pass a *state* parameter with the request to prevent CSRF attacks. In our case, we very simply stored the state as a cookie and compared the state stored in the cookie with the state from the request. This was of course not foolproof, as cookie was unencrypted and could potentially be tampered with. It was however an additional security measure, which could quite easily be expanded upon (for example by storing the state serverside using something like [gorilla/sessions](https://github.com/gorilla/sessions) with a backend store, or merely encrypting the cookie).

However, to deploy the project, we ended up using SkyHigh. We then received a *floating IP*, which only accessible on the internal NTNU network. However, when setting **Authorised redirect URI** in Google Developer Console, this is not a valid **public top-level domain**. Thus, as a workaround for the project deployment, we use [xip.io](http://xip.io/) as a custom DNS server. The *redirect URI* is thus set to **http://\<floating ip\>.xip.io:\<port\>/api/v1/authcallback**, which will redirect to xip.io. This means that everything essentially functions as intended **ecxept for the state cookie**. Thus, for this deployment, we have commented out the code validating the state in *pkg/auth/auth.go*. It has been commented out, not removed, to show what it would have looked like. All other paths than */login* will function as intended with the current deployment.


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
/user      (DELETE): Deletes specified fields from the user. If none are specified, the entire user and all related information is deleted.
/updategames (POST): Fetches new data from the servies registered for the user.
```

To update the user information, "/user" endpoint expects the following body for the POST request (values may be replaced, although they are required to be valid):
```
{
	"name": "newUsername",
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

To delete specific fields, "/user" endpoint expects the following body for the DELETE request (all other values are ignored):
```
["name", "games", "lol", "valve", "overwatch", "runescape", "games"]
```
If no fields are specified in the DELETE request, the entire user and all their data is deleted.

For all other paths, the request body is ignored.


#### Application structure
The application is split into two main parts: *cmd* and *pkg*. *cmd* serves as the central function of the application. *pkg* contains everything that is either used by *cmd or another package in pkg*. We consider the user to be the central part of the application as all actions and information is related to or belongs to the user. Therefore, the handler only takes a UserManager as a parameter and the **handler struct in pkg/server/handler.go [embedds](https://travix.io/type-embedding-in-go-ba40dd4264df) the UserManager**, allowing the handler to use each of the functions specified in the *UserManager interface*. The handler functions themselves contain a minimum amount of logic, merely calling functions from the UserManager, thus only handling i/o and logging.


**The Models package contains interfaces, structs, constants and function which are used by several packages to simplify the internal dependency graph.** For example, every struct used by multiple packages is defined in Models. Defining the struct in either of the packages would therefore create a direct dependency between them (or be a duplication).


Interfaces are widely used throughout the application to facilitate testing. This makes it possible to mock them, reducing the scope of the test. Interfaces are also used for the handler, where they serve to decouple the packages from eachother, preventing several direct dependencies. The **Organizer** interface is used to combine the interfaces for all the packages which provide games and stats to the application, simplifying the passing of the interfaces to the *UserManager*. Similar to how *handler* embeds the UserManager interface, the UserManager struct (which fulfils the interface) embeds the *organizer*, allowing it to call each of the functions specified in the *organizer interface* (and every interface within it).



#### Repository structure
The repository has the following main components:
 - **cmd**: Lists all possible commands for the application. Currently, there are none other than root. Main.go serves merely to start the *Run* function of cmd/root.go. Was created by Cobra during project initialization.
 - **pkg**: Contains all packages used in the application. See Application structure.
 - **.gitignore**: Specifies what files should be ignored by git.
 - **.gitlab-ci.yml**: Runs tests, linting, checks that the project compiles and deploys it to Openstack.
 - **.golangci.yml**: Golangci-lint configuration file.
 - **Dockerfile**: Barebones Dockerization. See [documentation](https://docs.docker.com/engine/reference/builder/).
 - **docker-compose.yml**: Barebones compose with "restart: always" and importing of environment variables from .env.
 - **go.mod** and **go.sum**: Go modules.
 - **LICENSE**: Apache license, created by Cobra during project initialization.
 - **main.go**: Starts the application.
 - **README.md**: The file you are currently reading.
 - **sample.env**: Shows which variables are expected to be present in the .env file.

#### Testing
As the project contains multiple packages, to run the tests (and get code coverage), use  ```go test ./... -cover```.

All the tests are unit tests where each of the required interfaces are mocked. This is to prevent the tests from testing other packages or external APIs which should **not** be part of a unit test. To mock responses from external sources, I used [bxcodec/faker](https://github.com/bxcodec/faker) (ecxept for the jagex test) to generate test data. To perform the actual checks throughout the tests, I used [stretchr/testify](https://github.com/stretchr/testify). All of the tests are **[table driven](https://github.com/golang/go/wiki/TableDrivenTests)**. No integration nor acceptance tests were made for the project. There is therefore no test for the database package.

We were dismayed that the only metric for tests were *code coverage*. This meant that the usefullness of the test, and what they are actually testing is utterly irrelevant, as long as enough of the code is executed. In our opinion, [test coverage alone is not a good metric](https://hackernoon.com/is-test-coverage-a-good-metric-for-test-or-code-quality-92fef332c871). As such, some of the tests contain very little actual testing. Specifically the tests *pkg/server/router_test.go*, *pkg/server/server_test.go* and the tests in *pkg/models* contain very little actual testing, as there is very little to test. The functions are nearly devoid of actual logic. It is possible to performe some more extensive tests on for example the router (checking that it contains each route as expected, and only allows certain methods), but this is very impractical and time consuming. This is also the reason *pkg/db* contain no test and *pkg/auth* contain few tests (these would also be unit tests as they would involve the database and OAuth provider respectively).


#### Using various std. lib packages, in addition to:
 - [spf13/cobra](https://github.com/spf13/cobra) for project initialization
 - [joho/godotenv/autoload](https://github.com/joho/godotenv) to import .env file to os.env
 - [gorilla/mux](https://github.com/gorilla/mux) for routing. Provides more functionality and flexibility than a standard router.
 - [sirupsen/logrus](https://github.com/sirupsen/logrus) for logging. Provides more functionality and flexibility than the log std.lib.
 - [stretchr/testify](https://github.com/stretchr/testify) for testing. Makes testing easier using assert and require.
 - [bxcodec/faker](https://github.com/bxcodec/faker) for generating test data.
 - [mitchellh/mapstructure](https://github.com/mitchellh/mapstructure) for parsing response from Firebase.
 - [coreos/go-oidc](https://github.com/coreos/go-oidc) as Go OpenID Connect client.
 - [golang.org/x/oauth2](https://godoc.org/golang.org/x/oauth2) for OAuth2 authorization and authentication support.
 - [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) for JSON Web Tokens (JWT).
 - [google.golang.org/grpc](https://godoc.org/google.golang.org/grpc) for checking rpc codes.
 - Various firebase related libraries


##### Server setup is based on [gorilla/mux graceful-shutdown example](https://github.com/gorilla/mux#graceful-shutdown)
