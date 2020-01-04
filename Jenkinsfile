node {
    stage 'Init'
    step([
      $class: 'GitHubCommitStatusSetter',
      contextSource: [
        $class: 'ManuallyEnteredCommitContextSource',
        context: 'jenkins/test-build'
      ],
      statusResultSource: [
        $class: 'ConditionalStatusResultSource',
        results: [
          [
            $class: 'AnyBuildResult',
            message: 'Build container',
            state: 'PENDING'
          ]
        ]
      ]
    ])

    step([
      $class: 'GitHubCommitStatusSetter',
      contextSource: [
        $class: 'ManuallyEnteredCommitContextSource',
        context: 'jenkins/test-build'
      ],
      statusResultSource: [
        $class: 'ConditionalStatusResultSource',
        results: [
          [
            $class: 'AnyBuildResult',
            message: 'Build container',
            state: 'PENDING'
          ]
        ]
      ]
    ])

    step([
      $class: 'GitHubCommitStatusSetter',
      contextSource: [
        $class: 'ManuallyEnteredCommitContextSource',
        context: 'jenkins/test-unit'
      ],
      statusResultSource: [
        $class: 'ConditionalStatusResultSource',
        results: [
          [
            $class: 'AnyBuildResult',
            message: 'Run tests',
            state: 'PENDING'
          ]
        ]
      ]
    ])
    step([
      $class: 'GitHubCommitStatusSetter',
      contextSource: [
        $class: 'ManuallyEnteredCommitContextSource',
        context: 'jenkins/test-service-fmt'
      ],
      statusResultSource: [
        $class: 'ConditionalStatusResultSource',
        results: [
          [
            $class: 'AnyBuildResult',
            message: 'service - go fmt',
            state: 'PENDING'
          ]
        ]
      ]
    ])
    step([
      $class: 'GitHubCommitStatusSetter',
      contextSource: [
        $class: 'ManuallyEnteredCommitContextSource',
        context: 'jenkins/test-service-lint'
      ],
      statusResultSource: [
        $class: 'ConditionalStatusResultSource',
        results: [
          [
            $class: 'AnyBuildResult',
            message: 'service - golint',
            state: 'PENDING'
          ]
        ]
      ]
    ])
    step([
      $class: 'GitHubCommitStatusSetter',
      contextSource: [
        $class: 'ManuallyEnteredCommitContextSource',
        context: 'jenkins/test-service-vet'
      ],
      statusResultSource: [
        $class: 'ConditionalStatusResultSource',
        results: [
          [
            $class: 'AnyBuildResult',
            message: 'service - go vet',
            state: 'PENDING'
          ]
        ]
      ]
    ])
    step([
      $class: 'GitHubCommitStatusSetter',
      contextSource: [
        $class: 'ManuallyEnteredCommitContextSource',
        context: 'jenkins/test-service-unit'
      ],
      statusResultSource: [
        $class: 'ConditionalStatusResultSource',
        results: [
          [
            $class: 'AnyBuildResult',
            message: 'service - unit tests',
            state: 'PENDING'
          ]
        ]
      ]
    ])

    def vault070 = docker.image('vault:0.7.0');
    def postgres96 = docker.image('postgres:9.6');
    vault070.pull()
    postgres96.pull()

    // Mark the code checkout 'stage'....
    stage 'Checkout'
    // Get some code from a GitHub repository
    checkout scm

    def tests=["jenkins/run-linter", "jenkins/test-build", "jenkins/test-unit",
               "jenkins/test-service-fmt", "jenkins/test-service-lint",
               "jenkins/test-service-vet", "jenkins/test-service-unit"]
    // [].removeElement is not whitelisted
    def completedTests=[]
    try {
        // Mark the code build 'stage'....
        stage 'Build'
        def service = docker.build 'ai_decision:snapshot'
        // docker.build has no argument for different Dockerfile than default
        sh 'docker build -t ai_decision:snapshot -f Dockerfile.test .'
        def testService = docker.image 'ai_decision:snapshot'

        sh 'docker build -t ai_decision_service:snapshot -f service/Dockerfile ./service/'
        sh 'docker build -t ai_decision_service_test:snapshot -f service/Dockerfile.test ./service/'
        def testGoService = docker.image 'ai_decision_service_test:snapshot'

        step([
          $class: 'GitHubCommitStatusSetter',
          contextSource: [
            $class: 'ManuallyEnteredCommitContextSource',
            context: 'jenkins/test-build'
          ],
          statusResultSource: [
            $class: 'ConditionalStatusResultSource',
            results: [
              [
                $class: 'AnyBuildResult',
                message: 'Build OK',
                state: 'SUCCESS'
              ]
            ]
          ]
        ])
        completedTests.add('jenkins/test-build')

        stage 'Linter'
        testService.inside("") {
          sh 'pylama'
          step([
            $class: 'GitHubCommitStatusSetter',
            contextSource: [
              $class: 'ManuallyEnteredCommitContextSource',
              context: 'jenkins/run-linter'
            ],
            statusResultSource: [
              $class: 'ConditionalStatusResultSource',
              results: [
                [
                  $class: 'AnyBuildResult',
                  message: 'Style check passed.',
                  state: 'SUCCESS'
                ]
              ]
            ]
          ])
          completedTests.add('jenkins/run-linter')
        }

        stage 'Test'
        def workspace = pwd()

        // NOTE: ran as root to enable go compilation in jenkins
        testService.inside("") {
            // fail if tests don't pass
            sh 'export LOG_LEVEL=\"DEBUG\"; pytest -s --cov=/python/src/github.com/callstats-io/ai-decision/src /python/src/github.com/callstats-io/ai-decision/test'
            step([
              $class: 'GitHubCommitStatusSetter',
              contextSource: [
                $class: 'ManuallyEnteredCommitContextSource',
                context: 'jenkins/test-unit'
              ],
              statusResultSource: [
                $class: 'ConditionalStatusResultSource',
                results: [
                  [
                    $class: 'AnyBuildResult',
                    message: 'Unit tests passed.',
                    state: 'SUCCESS'
                  ]
                ]
              ]
            ])
            completedTests.add('jenkins/test-unit')
        }

        postgres96.withRun("-e POSTGRES_PASSWORD=test -e POSTGRES_USER=ai_decision -e POSTGRES_DB=ai_decision") { postgre->
          vault070.withRun("--link ${postgre.id}:postgres -e VAULT_DEV_ROOT_TOKEN_ID=def76607-c3c0-4390-939d-1c844619340a") {vault ->
            // NOTE: ran as root to enable go compilation in jenkins
            testGoService.inside("--link ${postgre.id}:postgres --link ${vault.id}:vault --user=root \
            -e ENV=test \
            -e LOG_LEVEL=ERROR \
            -e SERVICE_NAME=ai_decision \
            -e HTTP_PORT=10001 \
            -e GRPC_PORT=10000 \
            -e VAULT_ADDR=http://vault:8200 \
            -e VAULT_SKIP_VERIFY=true \
            -e VAULT_TEST_BOOTSTRAP_TOKEN=def76607-c3c0-4390-939d-1c844619340a \
            -e VAULT_POSTGRES_NAME=ai_decision \
            -e VAULT_POSTGRES_ROOT_URL='postgres://ai_decision:test@postgres:5432/ai_decision_test?sslmode=disable' \
            -e VAULT_ENABLE_POSTGRES=true \
            -e VAULT_POSTGRES_CREDS_PATH=test/postgresql/ai_decision/creds/ai_decision \
            -e POSTGRES_CONN_TMPL='postgres://%s:%s@postgres:5432/ai_decision_test' \
            -e POSTGRES_ROOT_ROLE=ai_decision \
            -e POSTGRES_READ_ONLY_ROLE=ai_decision_read_only") {
                // fail if the files have not been formatted according to the go formatting rules
                sh 'cd /go/src/github.com/callstats-io/ai-decision/service; test -z $(gofmt -l ./src/)'
                step([
                  $class: 'GitHubCommitStatusSetter',
                  contextSource: [
                    $class: 'ManuallyEnteredCommitContextSource',
                    context: 'jenkins/test-service-fmt'
                  ],
                  statusResultSource: [
                    $class: 'ConditionalStatusResultSource',
                    results: [
                      [
                        $class: 'AnyBuildResult',
                        message: 'service - go code formatted correctly',
                        state: 'SUCCESS'
                      ]
                    ]
                  ]
                ])
                completedTests.add('jenkins/test-service-fmt')

                // fail if the files don't pass the simple linter rules
                sh 'cd /go/src/github.com/callstats-io/ai-decision/service; test -z $(golint $(go list ./... | grep -v vendor | grep -v protos))'
                step([
                  $class: 'GitHubCommitStatusSetter',
                  contextSource: [
                    $class: 'ManuallyEnteredCommitContextSource',
                    context: 'jenkins/test-service-lint'
                  ],
                  statusResultSource: [
                    $class: 'ConditionalStatusResultSource',
                    results: [
                      [
                        $class: 'AnyBuildResult',
                        message: 'service - golint ok',
                        state: 'SUCCESS'
                      ]
                    ]
                  ]
                ])
                completedTests.add('jenkins/test-service-lint')

                // fail if the files don't pass the go vetting tool checks
                sh 'cd /go/src/github.com/callstats-io/ai-decision/service; test -z $(go vet $(go list ./... | grep -v vendor | grep -v protos))'
                step([
                  $class: 'GitHubCommitStatusSetter',
                  contextSource: [
                    $class: 'ManuallyEnteredCommitContextSource',
                    context: 'jenkins/test-service-vet'
                  ],
                  statusResultSource: [
                    $class: 'ConditionalStatusResultSource',
                    results: [
                      [
                        $class: 'AnyBuildResult',
                        message: 'service - go vet ok',
                        state: 'SUCCESS'
                      ]
                    ]
                  ]
                ])
                completedTests.add('jenkins/test-service-vet')

                // fail if the unit tests do not pass
                sh 'cd /go/src/github.com/callstats-io/ai-decision/service; go test $(go list ./src/...)'
                step([
                  $class: 'GitHubCommitStatusSetter',
                  contextSource: [
                    $class: 'ManuallyEnteredCommitContextSource',
                    context: 'jenkins/test-service-unit'
                  ],
                  statusResultSource: [
                    $class: 'ConditionalStatusResultSource',
                    results: [
                      [
                        $class: 'AnyBuildResult',
                        message: 'service - go test passed',
                        state: 'SUCCESS'
                      ]
                    ]
                  ]
                ])
                completedTests.add('jenkins/test-service-unit')
            }
          }
        }
      }
      catch (err) {
        // [].each is not supported. JENKINS-26481
        for(int i = 0; i < tests.size(); i++) {
            if (!completedTests.contains(tests[i])) {
                step([
                  $class: 'GitHubCommitStatusSetter',
                  contextSource: [
                    $class: 'ManuallyEnteredCommitContextSource',
                    context: tests[i]
                  ],
                  statusResultSource: [
                    $class: 'ConditionalStatusResultSource',
                    results: [
                      [
                        $class: 'AnyBuildResult',
                        message: 'Build Fail',
                        state: 'FAILURE'
                      ]
                    ]
                  ]
                ])
            }
        }
        throw err
    }
}
