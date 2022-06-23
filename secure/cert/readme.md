## MAKE SELF-SIGNED CERTIFICATE in terminal

**write command in terminal to generate privat key (key) and open key (csr):**

```console
$ openssl req -new -newkey rsa:2048 -nodes -keyout forum.key -out forum.csr
```
*need to enter country, city and personal information*

**you will see:**

>*Generating a RSA private key*
>
>*................................+++++*
>
>*..........+++++*
>
>*writing new private key to 'forum.key'*

>*You are about to be asked to enter information that will be incorporated*
>
>*into your certificate request.*
>
>*What you are about to enter is what is called a Distinguished Name or a DN.*
>
>*There are quite a few fields but you can leave some blank*
>
>*For some fields there will be a default value,*
>
>*If you enter '.', the field will be left blank.*

>*Country Name (2 letter code) [AU]:KZ*
>
>*State or Province Name (full name) [Some-State]:Akmola*
>
>*Locality Name (eg, city) []:Nur-Sultan*
>
>*Organization Name (eg, company) [Internet Widgits Pty Ltd]:Forum-Security*
>
>*Organizational Unit Name (eg, section) []:Security*
>
>*Common Name (e.g. server FQDN or YOUR name) []:localhost*
>
>*Email Address []: xxx*
>
>
>*Please enter the following 'extra' attributes*
>
>*to be sent with your certificate request*
>
>*A challenge password []: xxx*
>
>*An optional company name []:Forum*


**write command in terminal to generate certificate (crt):**

```console
$ openssl x509 -req -days 365 -in forum.csr -signkey forum.key -out forum.crt
```

**you will see:**

>*Signature ok*
>
>*subject=C = KZ, ST = Akmola, L = Nur-Sultan, O = Forum-Security, OU = Forum, CN = Faiz, emailAddress = xxx*
>
>*Getting Private key*
