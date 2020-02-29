# Where were you allt his time api (wwuatt-api-go)

## Project config 

```sh
export WWUATT_ENV=local
```

## GO

### Instalation

```sh
brew install go
brew install dep
```

### Setup GO project base folder
Execute ``` go env GOPATH ```
If it is empty just add a the env var into your ```~/.zshrc``` or ```~/.bashrc``` file.
```sh
export GOPATH=~/my/path/to/code
```

## GOOGLE CLOUD ENGINE - GKE
### Install Google SDK
1. Download from the link
```https://cloud.google.com/sdk/docs```
2. Unzip the file where you want to keep it like ```~/google-cloud-sdk```
3. (Optional) Add to path:
    - Edit the file ``` ~/.zshrc``` or ``` ~/bashrc```
    - Add the following lines at the end of the file
    ```sh
    export GCLOUD_DIR=~/google-cloud-sdk # Use the path where you unzip it
    export PATH=$PATH:$GCLOUD_DIR/bin
    ```

### Install GO plugin for gcloud
```sh
gcloud components install app-engine-go
```

## DATABASE - MYSQL
### Install
```sh
brew install mysql
brew services start mysql
```

### Commnads
- First conection. Initial password: 'yourpassword'
```sh
mysql -u root -p
```

### Change root password
```sql
ALTER USER ‘root’@‘localhost' IDENTIFIED BY 'NEW_USER_PASSWORD';
```

### Create user and database
1. DB
```sql
CREATE DATABASE appdb;
```
2. User
```sql
CREATE USER 'appuser'@'%' IDENTIFIED BY 'password';
```
3. Privileges
```sql
GRANT ALL PRIVILEGES ON appdb.* TO 'appuser'@'%';
```
