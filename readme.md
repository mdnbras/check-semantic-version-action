# Utilities for GitOps

The GitOps Manager Utility is a powerful tool designed to simplify and automate common tasks related to Git repository management and continuous integration operations. With an intuitive command-line interface, it provides features such as version verification, updating GitHub variables, and validating commits in merge requests. Through this utility, streamline your GitOps workflow and enhance your team's efficiency in the development and deployment of projects.

### commits-verify

```yaml
  - name: Check Commits Pattern
    uses: mdnbras/check-semantic-version-action@v1
    with:
      command: 'commits-verify' # command
      owner: 'owner' # github owner
      repository: 'repository' # github repository
      mergeRequestId: 10 # Merge request identifier
      gbtoken: '' # access personal token
      bypass: 'YES' # YES or NO
      urlWebhook: '' # Discord webhook URL (optional)
      regexPattern: '' # regex pattern (optional)
```

### version-verify

```yaml
  - name: Check Semantic Version
    uses: mdnbras/check-semantic-version-action@v1
    with:
      command: 'version-verify' # command
      versionOld: '' # Old Semantic Version
      versionNew: '' # New Semantic Version
```
