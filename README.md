# easyinvesting
Easy Investing is a investments visualizer dashboard made with Python and angular

## Backend installation

1. Install docker on your machine
2. Clone the repository
   ```bash
   git clone https://github.com/JotaEspig/easyinvesting.git
    ```

3. Navigate to the backend directory
    ```bash
    cd easyinvesting/backend
    ```
4. Generate a self-signed SSL certificate (if you don't have one)
    ```bash
    mkdir -p apache/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
      -keyout apache/ssl/server.key \
      -out apache/ssl/server.crt \
      -subj "/CN=localhost"
    ```

5. Build the docker image
    ```bash
    docker build -t easyinvesting-backend .
    ```
6. Run the docker container
    ```bash
    docker run \
        --restart=always \
        -v /srv/sqlite-data:/data \
        -e EASYINVESTING_SECRET_KEY=your_secret \
        -e BRAPI_TOKEN=your_token \
        -d -p 443:443 \
        easyinvesting-backend
    ```

## Frontend installation

1. Install docker on your machine
2. Clone the repository
   ```bash
   git clone https://github.com/JotaEspig/easyinvesting.git
   ```
3. Navigate to the frontend directory
   ```bash
    cd easyinvesting/frontend
    ```
4. Generate a self-signed SSL certificate (if you don't have one)
    ```bash
    mkdir -p nginx/ssl
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
      -keyout nginx/ssl/server.key \
      -out nginx/ssl/server.crt \
      -subj "/CN=localhost"
    ```
5. Create a `.env.ts` file in the `src/app` directory with the following content:
    ```typescript
    export const ENV = {
        backendurl: "<apiUrl>", // Example http://localhost:8000
    };
    ```
   Adjust the `apiUrl` to match the backend API URL
5. Build the docker image
    ```bash
    docker build -t easyinvesting-frontend .
    ```
6. Run the docker container
    ```bash
    docker run --restart=always -d -p 443:443 easyinvesting-frontend
    ```
