Thanks for submitting a PR to DRATs.

## Checklist

Please provide the following information (and links if possible):

### [ ] What component are you testing? 

### [ ] Is the component an default component in `cf-deployment`?

### [ ] Have you created a `TestCase` and added it to the list of cases to be run?

### [ ] Have you added any new properties/information to all of the following:
* [ ] [integration_config.json](../ci/integration_config.json): Specifically, an `include_<testcase-name>` property
* [ ] [documentation in docs/](../docs/)
* [ ] [tasks in ci/](../ci/)
* [ ] [scripts in scripts/](../scripts/)

### [ ] Have you manually validated your `TestCase` against a deployed Cloud Foundry? If so, which version?

### [ ] Does this change rely on a particular version of `cf-deployment`?

### [ ] Are there any optional components of Cloud Foundry that should be enabled for this new `TestCase` to succeed?  Are their presence checked for in the `CheckDeployment` method of your `TestCase`?

### [ ] Are you available for a cross-team pair to help troubleshoot your PR?  What timezones are you based in?

### [ ] Have you submitted a pull-request to modify the `cf-deployment` [backup and restore ops files](https://github.com/cloudfoundry/cf-deployment/blob/master/operations/backup-and-restore/) to add a backup job and properties where appropriate?

## Do you have any other useful information for us?

We're on the #bbr cloudfoundry Slack channel if you need us.