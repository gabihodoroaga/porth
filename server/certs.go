package main

const serverCert = `
-----BEGIN CERTIFICATE-----
MIIEGTCCAgGgAwIBAgIJANV+QVqklenPMA0GCSqGSIb3DQEBCwUAMDcxFjAUBgNV
BAMMDXBvcnRoIFJvb3QgQ0ExEDAOBgNVBAoMB2hvZG8ucm8xCzAJBgNVBAYTAlJP
MB4XDTE5MDExNDAwMTAzN1oXDTIwMDExNDAwMTAzN1owNjEVMBMGA1UEAwwMcG9y
dGggU2VydmVyMRAwDgYDVQQKDAdob2RvLnJvMQswCQYDVQQGEwJSTzCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAKF6C2zcILlcTk+rySLqSW4xMV7PUK4r
2PMbCW3io09Ohz1QBH/Wrho8fAxP1oKre/gbgwvQoE9uQ3QwI1m9OPEVhMWS6Keo
tLDdTuLLOhseJpOdy1iW+cQo9lfVxiX0qjVmhNZG6eq4W8hAziW6YcOOHkvpZ8JY
cOAR0PimCYN4ZbPImX8K1v+atDfaKLmcO6x9i5WoA7F5bpkoG8D1mDE8DQjfCfzP
0CVtABgAMDCr4ka9UJyUp29AAk17Si6WXxBuqT7eXc1sV/QA9iFXpcVuyetLUatd
ea7DyEscUYwWvsRm131zA2rPjQDzsd+O6bET7EkxFXiDs+SpWnkQlpcCAwEAAaMp
MCcwJQYDVR0RBB4wHIIJbG9jYWxob3N0gg93d3cuZXhhbXBsZS5jb20wDQYJKoZI
hvcNAQELBQADggIBAKDgJvVRn+5/ufrGwkkNYpgc5HDo19sntfCRi8+NsZ+EKD0/
70c9XVlXsB4BYo3zsppJ03dPDUZqSc2FUjwRK4sEwHNVka+f0/oMm8TWwwhhFbuf
Gd2OCnO4i1cOYkzsXHc2i+MwNP8fDnFMrIHoDxXxuQUri5EuS8x7zdQbux3YS8yJ
46Vu+UbvlKBMBZ5f1UvVVUfxQ8aMEBhi5NObBk41cVVvSKBhVBWVT42/K3hOJPEu
HXAP/a4z8J4mUwFFLjiK93PIkMjanJGmJbybjFoJaawhOYW1d5uWBhQ1dGHuH8iC
emUYM7FTs/AxD0yKX8AgrhCuiLtYCFxWSl/H2wfEqv3S7ZljDe+qyLvB0CMxEqu7
gqr08yBdEtCzs90zFwbBH7GjsbR8hX+2u0GR0dO1o8NsISZZdmy1s+pJqflWOe1N
xqyPouh8VgD2jh0Q2Fz/TnTuq4ZGAkpVFbjQyFrXQqPmxGTLzQBmdaauyc0gloSQ
RnC3s4MGFHRtDmBvWej+WMGMLW1uU3ClfTgvryx1OPBcgjlhhjH3bvSPcS0UPQPE
u6QoPriCuso0Ftlw6Frw0mVwy4SHfzNIdwo07iTe6weFZ5pkndCu2mgXmbhMZG1H
m4aZzYuZ9+aDs6X0iOQumh6S/7M8OxS70NmZXj7FTsSyRfnu8CIzp/Oml4Ja
-----END CERTIFICATE-----`

const serverKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoXoLbNwguVxOT6vJIupJbjExXs9QrivY8xsJbeKjT06HPVAE
f9auGjx8DE/Wgqt7+BuDC9CgT25DdDAjWb048RWExZLop6i0sN1O4ss6Gx4mk53L
WJb5xCj2V9XGJfSqNWaE1kbp6rhbyEDOJbphw44eS+lnwlhw4BHQ+KYJg3hls8iZ
fwrW/5q0N9oouZw7rH2LlagDsXlumSgbwPWYMTwNCN8J/M/QJW0AGAAwMKviRr1Q
nJSnb0ACTXtKLpZfEG6pPt5dzWxX9AD2IVelxW7J60tRq115rsPISxxRjBa+xGbX
fXMDas+NAPOx347psRPsSTEVeIOz5KlaeRCWlwIDAQABAoIBAAV/2pPCi3rEogYk
m50MzaNrGXtZJC5KYAEnkpKjfVxeGE0GRKD19sf991uT/IJGxNoWVcg5Orx5zOJQ
IWQVBbNwQC1aa4IKRN0hLGac9UWnKTktzpcdzTPZEUq8hRsV0hYvf6ask5ri3H1d
d6fhSSMX3ABJ2rbLvExlTvCo9vn7OeCmFzptwqk0Hgkh6KW+S3MeepVwTQBupXOi
KPaDYTwfLParEAreBpDbC4Pxhdf7d29wm0e0chR9Ie6kPZMhFwxTkVmCidiLAm9Q
ZJcoIW5jL+SOVnLrsQe58tqtQVVg8wgDnK6KuOe6NXO7B0QzRabhdc1HMnw8/Cuh
4qQlAIECgYEA0U0IHUKTS0FUB1qYoDnKcyMys3dZ7VsVxnTtcKYg2oF+K4r9T9wy
ccrY3DqzAoa0ekkRhaHKODXb4inSSrWtrQV7GnbG0Ddd3iobo+bK8UnWn3bH+t4c
ZlA8PpnPFkRTuNxguEC49TMC1w1gUF+t4M8q4FFkAA600vruTDIBLd8CgYEAxYFd
ZxKavz7f/tfBya+ncEnCZA5gKJrJv7yW0KXNbeprLPqq0jixY/zNcfkrBgk+IBZw
wcSvntSMre27yH2mYZf0vemYi2XJUSY+y2YNGKrHKKZDrr5T5C55ensnJsq2IVKd
5hntgKHhOHtR4+HEYLi4y+IfsQFER5v1+lExvkkCgYEAiKLhSTjNL7PWR3a9bNxN
bhzsXHzuGCX+cTCkUYYirIMc+xAhjqERzXe/WwZ3Fo8aAzrwVWzptwhyI5Np1ZwF
ZY7ObthbslJy1TZoFPf2RM8Pbcr9gqi9oY1/xt5icwboISa9fYvDM0+56uqwlcfg
m4KjWw3HWsI/Cf0G1HdQjcECgYB1sq0FspmTZJW52bu7RDlE+j+kvshhCjU2VN2P
Q29TlEIAUPUhR/W2fz2zMOiJtVJXbugNIPgDb+jR8X1Zcj+HozWPQzjLwYGiIWeE
cLFXRNZgjAyDgxqdPXDZI7DmNiEpZIGCUWsun8mGjj7zzWPou8wse/mk0vtsrS19
2YsKOQKBgQCXoM3sXH3S7dE0xGjcMA4vqPndHAGTKBbkz5qzXSfxKE6DAM+3UWJQ
1dy6Ef7XpHHNU0GgfGrPbGjiLa54eIPzDViBJXUv+0h5ru/RQ6KLJYt4kfZrBhGf
0cD8dSimkQ3+HA7ycNczN8Iwm4eJ2wpcL0nYOC/0/uWNG0lJraWYrw==
-----END RSA PRIVATE KEY-----`

const rootCA = `
-----BEGIN CERTIFICATE-----
MIIE6jCCAtICCQCfO6Jse3YEKTANBgkqhkiG9w0BAQsFADA3MRYwFAYDVQQDDA1w
b3J0aCBSb290IENBMRAwDgYDVQQKDAdob2RvLnJvMQswCQYDVQQGEwJSTzAeFw0x
OTAxMTQwMDEwMzRaFw0yMTExMDMwMDEwMzRaMDcxFjAUBgNVBAMMDXBvcnRoIFJv
b3QgQ0ExEDAOBgNVBAoMB2hvZG8ucm8xCzAJBgNVBAYTAlJPMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAxowntKuCxfobkHdGg9r9VDfDec7OEZQOK+X3
rrntxnpobDYTdXb5eXxMv5uEZ41KGp4yDHDm2drU+GsmsDEWFcXFuNGExenaXr4s
AKRiUG9x5J+PSaFj6acD8SMiY7cJx89oSy8nFYX8uAzB2D3Oa8UKSXBmCim6mAoV
PGHxKMmkYL3XZVRoIizLkClRbgmW2EGtWH9n71woCeQ/IGofqLCawz+Kj9KMiVs9
JQzB9kN//pZbz5gMYVBY+XzLZJYGfgjjbtX8oaQ0JXFuo9ra8h44xPGO5oEmPoNz
t0dOz8xwFgl6lYlIcVE7LOTuE30eZct+tGlD9itm8FjBx8YW9DTXsdRhyRPGC5Eb
H2qHaD38ZuQvGXQaRbxiX7gxPX4JnnoYOjtVBDW6J3k0PVb9s9Fu+V5AolTlmG8b
SnNf1W+H1611UPYBFypNF2XLSHm346UmTyMUdd6vAs64DMXZjqingizmx0gT4GzH
IYtkw0smMiVfmOvkdHW+baGRws7fTYSPKxg4nDdaWE+NLbPnsHIbXVioEr/MFzMo
ZHOYJKf15ys+zMzxBaGvzS3Rjr7RKCg7MWKnp+fnNFLWcsZPLb5Oav90jonckwmL
bY1AAwE2ZrL35H3s5OJOADo0pQY+odIx9atuakpb/YRKUMqepAG3qFntWGyqgW0W
fbO0iKUCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAoFf/cs2jXQjw3snmRW/zOgu0
GWIeP4cEjRMoAMjHV3e6S3N7cKs0ohGgpQZb3QDsGhSr29pmycjkMs0Mmq1PN9JX
P8YFN0slj/4NAAhU6OnS9FUV2WuMK5XP+VV8m6+qO4nwJgIivdilAx6O/tARWkyO
vJn9SqYprAuY4WFocN6J4NUjsVpOQ0WoSunOYg5uRI/AyFrvWkw7SqAOkweVO2Ww
4hQt3FCPn8B9rkBuGKOn+rFpY7o/64Ncez+DmyJf69O/HiTjHE+bGbbfWaIv2Jai
lfD9p1SEl9eq2bsOTr3iM1ADHe/nuLZMBbYmOmpvMhRssW7sSG8+EL/4Fxer3VUE
mUHK7RXlvbl7e2zezWTYgBR5PAcd6YvK3hBXcCMZ8CjqYZ5A9gFSNmWoCrjW9vgh
ZJPF7RuO50qJsOI+uArXtrhVOwnC1tAJvaO5TOjl/ViGnVeJ8bBUOFi5DFzVEWYy
jlkFoHicgDRmQxfB6TCAiJL/wHWwfXqlGajCDuvrqbMP7PEkw4HYrO87Triz/s3u
K15xFtbB5D+gBsm7mYXRvpmuYwV7VqiUtNDud5WBHjKUBYqu0JoonXIDlcAEc9Qy
R1d/7TzV+58GI6s5O1Iq+Bt7eE5QyxGAk4gdDTVNkYS5Jtr2THE6VGxbhjfqR+gj
oZoNb+94g6/2TnxR2Sg=
-----END CERTIFICATE-----`