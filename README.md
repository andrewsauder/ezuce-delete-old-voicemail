# ezuce-delete-old-voicemail

Binary exeuctable that connects to your Ezuce Mongodb instance and deletes all voicemails older than the number of months specified. Binaries are available for Windows and Linux in the /dist directory.

## Arguments

The available arguments are:

- **Delete Older Than X Months**
- **Mongo URL**
- **Mongo Username**
- **Mongo Password**

An example usage, passing all arguments:
```bash
deloldvm -months=6 -url="127.0.0.1:27017" -u="myusername" -p="mypassword"
```

### Delete Older Than X Months

- Tells the program how many months of voicemail to keep. For example, if you enter 10, the program will delete all voicemail older than 10 months. It defaults to 6.

```bash
deloldvm -months=6
```

### Mongo URL

- Tells the program what database to connect to. Defaults to 127.0.0.1:27017 (default local mongo instance)

```bash
deloldvm -url="127.0.0.1:27017"
```

### Mongo Username

- Tells the program what username to connect to Mongo with. Defaults to blank for no authentication.

```bash
deloldvm -u="myusername"
```

### Mongo Password

- Tells the program what password to connect to Mongo with. Defaults to blank for no authentication.

```bash
deloldvm -p="mypassword"
```


## Logging
A log file named *delete-old-voicemail.log* is automatically created in the directory where the executable resides. This log will show database connection errors as well what voicemails are being deleted while the application is running. The total number of voicemails it deletes is appended after each run.
