rrpBuildGoCode {
    projectKey = 'utilities'
    buildImage = 'amr-registry.caas.intel.com/rrp/ci-go-build-image:1.12.0-alpine'
    testDependencies = ['consul']
    skipBuild = true
    skipDocker = true
    testStepsInParallel = false
    protexProjectName = 'bb-utilities'

    notify = [
        slack: [ success: '#ima-build-success', failure: '#ima-build-failed' ]
    ]
}
