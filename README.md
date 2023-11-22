# 2FA Go

This is just a a small test implementation of 2FA OTP authentication in Go using [pquerna/otp](https://github.com/pquerna/otp/).

## Improvements

Since it is a initial test, the idea is to port this to use the [Encore](https://encore.dev) framework instead of Gin :)

Other improvents include:
- Validate the entirety of the payload;
- Hash password;
- Add QR code generation;
- Add recovery phrases for recovery;
- ...