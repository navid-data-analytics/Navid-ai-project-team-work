# ai-decision

## Initialize

You need Docker (inside docker-machine as an example) and docker-compose installed.

    git clone git@github.com:callstats-io/ai-decision.git .
    cd ai-decision
    ./scripts/setup_dev.sh

## Start the system

This describes running the system from this folder.
Anything that has to be run additionally in localfullstack is stated as "In localfullstack ...".

### Once

The recommended way to run full development environment is through [localfullstack](https://github.com/callstats-io/csio-command-center/tree/master/localfullstack)
to ensure developer environment compatibility (everyone has the same setup) and as close to production-like environment as possible.

In localfullstack
 - start CRS (look at readme for help)
 - execute these to start AI decision
 ```
    docker-compose up -d vault_ai_decision_service
    docker-compose run vaultbootstrap_ai_decision_service
    docker-compose up --build ai_decision_service
 ```

If run from this folder, then the names are different:

    docker-compose up -d vault_service
    docker-compose run vaultbootstrapdev_service
    docker-compose up --build ai_decision_service

### Each time

In local fullstack, fill the localfullstack mongo with data (execute for each appID):

    docker-compose run dashboardscripts --localAppID 1234 --appID 1234

To populate localfullstack mongo with data used now (14.09.2018), execute:

    docker-compose run dashboardscripts --localAppID 347489791 --localAppID 331463193 --localAppID 380077084 --localAppID 722018081  --localAppID 748139165 --localAppID 815616092 --localAppID 486784909  --localAppID 713087156

Run AI decision (analytics):
Make sure that the IP address in docker-compose.yml `CRS_GRPC_ADDR` corresponds to docker-machine IP

    docker-compose build ai_decision
    docker-compose run --rm ai_decision

## Access AID service postgres

Connect to local postgres:

    docker-compose run --rm postgres_ai_decision_service_shell

Or from folder (password is defined in [docker-compose.yml](./docker-compose.yml) postgres_service):

    docker-compose run --rm postgres_service_shell

## Jupyter environment in repo folder

With jupyter notebook as environment of choice (best for analysis):

    docker-compose build csiojupyter
    docker-compose up csiojupyter

## Running tests in repo folder
AID (analytics):

    docker-compose build test_ai_decision
    docker-compose run --rm test_ai_decision

AID service:

    docker-compose build test_ai_decision_service
    docker-compose run --rm test_ai_decision_service

## Running with Prometheus metrics

    docker-compose build ai_decision
    docker-compose run -p 8084:8084 ai_decision

In order to avoid writing ports when running, use:

    docker-compose up ai_decision

The submitted metrics to Prometheus are in the following URL:

    <docker-machine ip>:8084

## Upgrade go libraries
Updating our own go-common once in a while is a good idea to get bug fixes or new features.
Also, we might want to update e.g. gRPC or general go version.

- go version can be changed in `Dockerfile` and `Dockerfile.test`
- updating libraries is done with glide
    - change `glide.yaml` to have the new versions you want
    - run `glide update`, and `glide install` can never hurt as well :)

Make sure that you are in the right path! When your `$GOPATH` points to `/x` then this repo should be cloned to `/x/src/github.com/callstats-io/ai-decision`. Above commands should then be executed from that directory, otherwise funny things happen.


## Get production/test cluster data

To get insight into current messages in test/prod postgresSQL:

```
# Assuming you are in /ai-decision/ !
cd scripts

python generate_psql_csv.py --username {vault username} --password {vault password} --db_type {dbtype}

# or just

python generate_psql_csv.py

# to be prompted for relevant credentials
```

db_type is either test or prod.

## Message Control

Messages sent to postgres can be manipulated. There are 3 ways to manipulate a message:

- Deleting a message from database
- Manual date setting for suppression
- Unsuppressing specific date for message client

#### Message deletion:

Message deletion happens with aid-service start. They are provided by passing IDs of the messages from database to the deployment scripts or docker-compose (depending whether it is deployed by kubernetes or just localfullstack run).

To delete messages from local environment add delete flag to the ai_decision_service docker-compose setup.
--delete=1,2,3 will delete messages with id 1,2 and 3. The messages are only deleted on redeployment of ai_decision service.
The line in docker compose should look like this:

```
    command: '/bin/bash -c "/go/bin/ai-decision-service --server=false --migrate=init && exec /go/bin/ai-decision-service --server=true --migrate=up --delete=1,2,3"'
```


To delete messages from test/prod deployment cluster, the flag should be added in deployment_scripts/kubernetes/ai_decision/service-migrate.yml.

The same flag should be added as another element of the list that is already there:

```
    command: ["/go/bin/ai-decision-service", "--server=false", "--migrate=up", "--delete=1,2"]
```

Where 1 and 2 are IDs of the messages to be deleted.

#### Manual Suppression:

Sometimes there is a need for updating suppression date that is already there in MessageClient. Normally MessageClient remembers the date of the last sent message so we omit sending duplicate messages to the database. In some cases, we want to update this date manually. One example is changing a parameter that would send messages more frequently, but ran from scratch would send messages from the past - like changing threshold for RTT fluctuation from 50ms to 20ms would find many more messages in the last two years - which should have been suppressed. Normally they would, but if the application has not met this message since 8 months ago, there is existing 8 month window for new, 20ms messages.

To nivelate this possibility we want to manually set the suppression date to the date of deployment, hence every appID in the pipeline will suppress earlier undetected RTT fluctuations.

Another example would be deleting X number of messages, because data from the passed came out inconsistent, then setting new earlier suppression date to update the messages with correct data.

To update the date the MANUAL_DATE parameter in docker compose have to be set, as such:

```
 MANUAL_DATE: "{\"722018081\":{\"MidtermOQFluctuationImmediatelyHigh\":[\"02-03-2019\"]}}"
```

This flag sets date for appID 722018081 and message type MidtermOQFluctuationImmediatelyHigh to 02-03-2019.

For test/prod deployment clusters the same parameter has to be updated in the deployment_scripts/kubernetes/ai_decision/prod.py or test.py under "manual_date" in settings dict. The string parameter needs additional braces escaping, i.e.
```
 "manual_date": "\'{\"722018081\":{\"MidtermOQFluctuationImmediatelyHigh\":[\"02-03-2019\"]}}\'"
```

#### Unsuppressing Specific Date

For fixing a specific message, it is possible to delete it from the database then unsuppress its specific type, date and appID. To do that the UNSUPPRESS parameter was added to docker-compose under ai_decision environment variables.

```
 UNSUPPRESS: "{\"380077084\":{\"MidtermRttFluctuationImmediatelyHigh\":[\"07-03-2019\"]}}"
```

For test/prod deployment clusters the same parameter has to be updated in the deployment_scripts/kubernetes/ai_decision/prod.py or test.py under "manual_date" in settings dict. The string parameter needs additional braces escaping, i.e.

```
 "unsuppress": "\'{\"380077084\":{\"MidtermRttFluctuationImmediatelyHigh\":[\"07-03-2019\"]}}\'"
```

Will make possible to send the message for the date 07-03-2019 only if the message is of type MidtermRttFluctuationImmediatelyHigh and of appID 380077084. This is ought to be used if the suppression date should remain the same.