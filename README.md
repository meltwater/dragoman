# dragoman
A secrets encryption tool that helps seal your secrets for storage alongside your source code. An example use case is to seal secrets needed by deployment infrastructure (like [Terraform](https://www.terraform.io/))

# Installation
There are a few ways to install the executable:

## Using the go installer
`$ go install github.com/meltwater/dragoman@vX.Y.Z`
- Replace `vX.Y.Z` with your desired version number
- Make sure `$GOPATH/bin` is in your `$PATH`

## Download from releases
- The executable can be downloaded from the [releases page](https://github.com/meltwater/dragoman/releases/latest) for various architectures
- Extract the executable and place it somewhere in your `$PATH`

# KMS Encrpytion and Decryption
Envelope encryption can be done using KMS to encrypt your strings locally. When encrypting secrets, you will be returned a string in the format `ENC[KMS,{{YOUR_ENCRYPTED_SECRET}}]`. It's important to leave the ENC and KMS wrapping to allow for proper decryption later.

You will either need your AWS Credentials in ~/.aws/credentials or all of the following environment variables set:
- `$AWS_ACCESS_KEY`
- `$AWS_SECRET_ACCESS_KEY`
- `$AWS_REGION`

Details on how to configure your AWS credentials [can be found here]("github.com/aws/aws-sdk-go-v2/config")

## CLI Usage
```bash
# Encryption using AWS KMS
$ echo -n "Jon Snow is a Targaryen" | dragoman encrypt --kms-key-id=alias/my-secret-key

# ---

# Decryption using AWS KMS
$ echo "ENC[KMS,...]" | dragoman decrypt

# ---

# Decrypt an entire file
$ cat my_encrypted.file | dragoman decrypt > my_decrypted.file
# OR
$ dragoman decrypt -f my_encrypted.file > my_decrypted.file
```
### Notes on Encryption
- Encrypt reads the string to encrypt from std:in. This means you can encrypt entire files like this: `$ cat myfile.txt | dragoman ...`
- `--kms-key-id` can handle multiple formats. See [the KeyId section of the kms docs](https://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html#API_Encrypt_RequestSyntax) for more info

### Notes on Decryption
- Decrypt reads the string provided to std:in or optionally a file via the `--input` argument
- Decrypt will output the decrypted string to std:out which can then be forwarded to a file if desired
- Decrypt will search the provided text for any encryptions and do a replace-in-place for each encryption it finds

## Configuration

The KMS Encryption and Decryption crypto strategies require some environment variables to be set in order to properly configure the AWS SDK. 

### AWS Region
You have 3 options for how to configure what region to use. 
1. Specify the region when calling the api using the `--aws-region` flag
2. Set the `AWS_REGION` environment variable (*preferred method*)
3. Set the `AWS_DEFAULT_REGION` environment variable

# Contributing
## Submitting PRs
We will gladly take a look at PRs for improvements to this codebase. Add your changes to a branch, create a PR to `main` and once all checks successfully complete, ping Carlitos in [#carlitos-el-paraiso](https://meltwater.slack.com/archives/CB1EZMNJZ)
## Enabling git hooks
We suggest enabling git hooks so that linting and unit test automatically get run when committing and pushing changes.

To enable git hooks run `$ git config core.hooksPath .githooks` in the project root

## Builds
Builds can be done by running `go build` in the project root

## Testing
The unit tests can be run by running `go test ./{package_name}` in a specific package or for all packages by running `go test ./...` in the project root.
# Support, Questions and Feedback
The lead maintainer for this project is Team Carlito's Way.

[Email (all.carlitosway@meltwater.com)](mailto:all.carlitosway@meltwater.com) | [Slack (#carlitos-el-paraiso)](https://meltwater.slack.com/archives/CB1EZMNJZ)