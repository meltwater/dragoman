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

# KMS Encryption and Decryption
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

# Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.md) to understand how to submit pull requests to us, and also see our [code of conduct](CODE_OF_CONDUCT.md).

All ideas for new features and bug reports will be kept in [github.com/meltwater/dragoman/issues](https://github.com/meltwater/dragoman/issues).

## Enabling git hooks
We suggest enabling git hooks so that linting and unit test automatically get run when committing and pushing changes.

To enable git hooks run `$ git config core.hooksPath .githooks` in the project root

## Authors and Acknowledgement
See the list of [all contributors](https://github.com/meltwater/dragoman/graphs/contributors).

## Inspiration
The idea stemmed from a subset of functionality of the [Secretary](https://github.com/meltwater/secretary) tool that is now deprecated. Internal teams still use this subset of functionality to handle secrets origination problems and as a result, we decided to resurrect that functionality in this tool.

# License and Copyright
This project is licensed under the [MIT License](LICENSE).