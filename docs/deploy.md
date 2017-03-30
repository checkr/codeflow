### Deploying an application with Codeflow
Deploying an application is easy.  When commits/PRs are merged into the master branch you will see them on the Deploy tab under "Features". You can deploy them by clicking the "Deploy" button. 

#### Continuous Deployment settings
On the settings page there are 2 deployment related settings:
* Continuous Integration (CircleCI):  This enables CircleCI hooks and will give Codeflow awareness about the status of tests passing.  Deploys will not automatically go out if a test is failing when this box is checked.
* Continuous Delivery (Deploy on green):  This enables continuous deployment.  Every feature that is pushed to master branch will be deployed if 1) the docker build is successful, 2) the CircleCI tests have passed (optional see above).

### Rollback
To rollback a deployment you simply click "Rollback" on a previous release.

### Re-deploy
When scaling up the number of resource count or when changing an environment setting you must do a deploy for it to take effect.  You can deploy the same version of the code as is currently running but with the new settings by clicking "Re-deploy" on the current release.