# dragoman
A secrets encryption tool

## Usage
### KMS Encrpytion and Decryption
#### Configuration
The KMS Encryption and Decryption crypto strategies require some environment variables to be set in order to properly configure the AWS SDK. 

**Credentials**: Details on how to configure your AWS credentials [can be found here]("github.com/aws/aws-sdk-go-v2/config").

**AWS Region**: You have 3 options for how to configure what region to use. 
1. Specify the region when calling the api using the `--aws-region` flag
2. Set the `AWS_REGION` environment variable (*preferred method*)
3. Set the `AWS_DEFAULT_REGION` environment variable



## Contributing
### Enabling git hooks
Run `$ git config core.hooksPath .githooks`