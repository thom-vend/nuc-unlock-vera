# nuc-unlock-vera
Basic tool: fetch ciphered stuff on http(s) server, decrypt, run command

See nucunlocker.yml for configuration options

## Options
```
Usage of nucunlocker:
  -c string
        path to config file (default "nucunlocker.yml")
  -d string
        data to encrypt/decrypt
  -m string
        run mode (unlock/encrypt/decrypt)
  -p string
        password to encrypt/decrypt (for encrypt/decrypt mode)
  -r    retry on http 419 error indefinitely
```
