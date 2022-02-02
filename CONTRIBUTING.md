# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue,
email, or any other method with the owners of this repository before making a change.

Please note we have a [code of conduct](CODE_OF_CONDUCT.md) that we ask you to follow in all your interactions with the project.

**IMPORTANT: Please do not create a Pull Request without creating an issue first.**

*Any change needs to be discussed before proceeding. Failure to do so may result in the rejection of the pull request.*

Thank you for your pull request. Please provide a description above and review
the requirements below.

## Pull Request Process

0. Check out [Pull Request Checklist](#pull-request-checklist), ensure you have fulfilled each step.
1. Check out guidelines below, the project tries to follow these, ensure you have fulfilled them as much as possible.
    * [Effective Go](https://golang.org/doc/effective_go.html)
    * [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
2. Ensure any install or build dependencies are removed before the end of the layer when doing a
   build.
3. Please ensure the [README](README.md) is up-to-date with details of changes to the command-line interface,
    this includes new environment variables, exposed ports, used file locations, and container parameters.
4. **PLEASE ENSURE YOU DO NOT INTRODUCE BREAKING CHANGES.**
5. **PLEASE ENSURE BUG FIXES AND NEW FEATURES INCLUDE TESTS.**
6. We will merge the changes once approved and notify you once the build has completed.

## Pull Request Checklist

- [x] Read the **CONTRIBUTING** document. (It's checked since you are already here.)
- [ ] Read the [**CODE OF CONDUCT**](CODE_OF_CONDUCT.md) document.
- [ ] Add tests to cover changes.
- [ ] Ensure your code follows the code style of this project.
- [ ] Ensure CI and all other PR checks are green OR
    - [ ] Code compiles correctly.
    - [ ] Created tests which fail without the change (if possible).
    - [ ] All new and existing tests passed.
- [ ] Improve and update the [README](README.md) (if necessary).

## Testing Locally

Want to test locally without opening a PR?  Follow the steps below to build a local executable

0. Verify you have the applicable go version installed with `$ go version`
1. Run `$ go mod tidy` to install needed modules
2. Run `$ go install .` to install dragoman
3. Verify `$GOPATH/bin` is in your `$PATH`

## Response Times

**Please note the below timeframes are response windows we strive to meet. Please understand we may not always be able to respond in the exact timeframes outlined below**
- New issues will be reviewed and acknowledged with a message sent to the submitter within two business days
    - ***Please ensure all of your pull requests have an associated issue.***
- The ticket will then be groomed and planned as regular sprint work and an estimated timeframe of completion will be communicated to the submitter.
- Once the ticket is complete, a final message will be sent to the submitter letting them know work is complete.

***Please feel free to ping us if you have not received a response after one week***