# hhhash

golang HTTP Headers Hashing CLI

## Description

HTTP Headers Hashing (HHHash) is a technique used to create a fingerprint of an HTTP server based on the headers it returns. HHHash employs one-way hashing to generate a hash value for the set of header keys returned by the server. See https://www.foo.be/2023/07/HTTP-Headers-Hashing_HHHash for more info.

## Usage

```bash
./hhhash https://www.circl.lu/
hhh:1:78f7ef0651bac1a5ea42ed9d22242ed8725f07815091032a34ab4e30d3c3cefc
```