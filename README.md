# dragoman
A secrets encryption tool that helps seal your secrets for storage alongside your source code. An example use case is to seal secrets needed by deployment infrastructure (like [Terraform](https://www.terraform.io/))

## Installation
There are a few ways to install the executable:

### Using the go installer
`$ go install github.com/meltwater/dragoman@vX.Y.Z`
- Replace `vX.Y.Z` with your desired version number
- Make sure `$GOPATH/bin` is in your `$PATH`

### Download from releases
- The executable can be downloaded from the [releases page](https://github.com/meltwater/dragoman/releases/latest) for various architectures
- Extract the executable and place it somewhere in your `$PATH`

## KMS Encrpytion and Decryption
Envelope encryption can be done using KMS to encrypt your strings locally.

### Encryption

`$ echo "Jon Snow is a Targaryen" | onstaging dragoman encrypt --kms-key-id alias/my-secret-key`
- Encrypt reads the string to encrypt from std:in. This means you can encrypt entire files like this: `$ cat myfile.txt | dragoman ...`
- `onstaging` in this situation just sets the AWS credentials needed to do KMS encryption
- `--kms-key-id` can handle multiple formats. See [the KeyId section of the kms docs](https://docs.aws.amazon.com/kms/latest/APIReference/API_Encrypt.html#API_Encrypt_RequestSyntax) for more info

### Decrypt Strings

`$ echo "ENC[KMS,...]" | onstaging dragoman decrypt`
- Decrypt reads the string provided to std:in or optionally a file via the `--input` argument
- Decrypt will output the decrypted string to std:out which can then be forwarded to a file if desired
- Decrypt will search the provided text for any encryptions and do a replace-in-place for each encryption it finds

### Decrypting Entire Files

`$ cat myfile.enc | onstaging dragoman decrypt > myfile.dec`

or

`$ onstaging dragoman decrypt -i myfile.enc > myfile.dec`

### Configuration

The KMS Encryption and Decryption crypto strategies require some environment variables to be set in order to properly configure the AWS SDK. 

**Credentials**: Details on how to configure your AWS credentials [can be found here]("github.com/aws/aws-sdk-go-v2/config").

**AWS Region**: You have 3 options for how to configure what region to use. 
1. Specify the region when calling the api using the `--aws-region` flag
2. Set the `AWS_REGION` environment variable (*preferred method*)
3. Set the `AWS_DEFAULT_REGION` environment variable

## Contributing
### Enabling git hooks
Run `$ git config core.hooksPath .githooks`