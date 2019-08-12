# pgp-verify

Simple commandline tool to verify pgp signatures tailored to be used in scripts.
In contrast to gnupg it does not require any database and will set a proper exit
code on success/failure.


## Usage

```
$ ./pgp-verify test.sh test/channel-rust-beta-date.txt.asc test/rust-key.gpg.ascii
```

On success it exit with status code 0; if the signature is not valid it will
exit with status code 1. Note that `pgp-verify` ignores the expire date of the
pgp key on purpose. If you need this to check for expire dates, use a different tool.
