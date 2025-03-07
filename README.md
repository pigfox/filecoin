# glif-test

This project provides a simple API for interacting with the Filecoin blockchain, allowing you to check wallet balances and transfer FIL.

## Setup

1.  **Environment Configuration:**
    * Save `env.example` as `.env`.
    * Enter your own values in the `.env` file:
        * `PRIVATE_KEY`: The private key of the sender address.
        * `WALLET_ADDRESS`: The address corresponding to the private key.
        * `TRANSFER_RECEIVER_ADDRESS`: The address to which you want to send FIL.

2.  **Docker Setup (Terminal 1):**
    * Start the Docker containers: `docker-compose up`
    * Stop the Docker containers: `docker-compose down`
    * **Note:** Use `docker.sh` only as a last resort (nuclear option).

3.  **Start the API (Terminal 2):**
    * Navigate to the `api` directory and run: `go run *.go`

## API Endpoints (Terminal 3)

1.  **Get Wallet Balance:**

    * **Request:**
        ```
        GET http://localhost:8080/wallet?address=0x...
        ```
    * **Response:**
        ```json
        {
            "fil": "0.689337087983363083",
            "iFil": "0"
        }
        ```

2.  **Transfer FIL:**

    * **Request:**
        ```
        POST http://localhost:8080/transaction
        ```
        * **Header:** `Content-Type: application/json`
        * **Body:**
            ```json
            {
                "sender": "0x8822D9472E96F22565cc07aD259DC03b100B18d6",
                "receiver": "0xb04d6a4949fa623629e0ED6bd4Ecb78A8C847693",
                "amount": "10000000000000000"
            }
            ```
    * **Response:**
        ```json
        {
            "tx_hash": "0xde68344fee4189e41536a9bc067bd1294a336f326cd3d4845df75fee52d756de",
            "sender": "0x8822D9472E96F22565cc07aD259DC03b100B18d6",
            "receiver": "0xb04d6a4949fa623629e0ED6bd4Ecb78A8C847693",
            "amount": "10000000000000000",
            "status": "pending",
            "timestamp": "2025-03-06T19:45:52.112824531-08:00"
        }
        ```

3.  **List Transactions:**

    * **Request:**
        ```
        GET http://localhost:8080/transactions?address=0x...
        ```
    * **Response:**
        ```json
        [
           {
               "tx_hash": "0xde68344fee4189e41536a9bc067bd1294a336f326cd3d4845df75fee52d756de",
               "sender": "0x8822D9472E96F22565cc07aD259DC03b100B18d6",
               "receiver": "0xb04d6a4949fa623629e0ED6bd4Ecb78A8C847693",
               "amount": "10000000000000000",
               "status": "success",
               "timestamp": "2025-03-06T19:45:52.112825Z"
           },
           {
               "tx_hash": "0x151ecae1d9d1744481ec1775084d644984c165fc06a324fcadea73aed3e6c07f",
               "sender": "0x8822D9472E96F22565cc07aD259DC03b100B18d6",
               "receiver": "0xb04d6a4949fa623629e0ED6bd4Ecb78A8C847693",
               "amount": "10000000000000000",
               "status": "pending",
               "timestamp": "2025-03-06T20:18:22.485153Z"
           }
        ]
        ```

4.  **Transaction Mining Status:**

    * After initiating a transfer, the transaction status is initially `"pending"`.
    * Once the transaction is mined, the status will be updated to `"success"`.
    * Example of a successful transaction list:

        ```json
        [
           {
               "tx_hash": "0xde68344fee4189e41536a9bc067bd1294a336f326cd3d4845df75fee52d756de",
               "sender": "0x8822D9472E96F22565cc07aD259DC03b100B18d6",
               "receiver": "0xb04d6a4949fa623629e0ED6bd4Ecb78A8C847693",
               "amount": "10000000000000000",
               "status": "success",
               "timestamp": "2025-03-06T19:45:52.112825Z"
           },
           {
               "tx_hash": "0x151ecae1d9d1744481ec1775084d644984c165fc06a324fcadea73aed3e6c07f",
               "sender": "0x8822D9472E96F22565cc07aD259DC03b100B18d6",
               "receiver": "0xb04d6a4949fa623629e0ED6bd4Ecb78A8C847693",
               "amount": "10000000000000000",
               "status": "success",
               "timestamp": "2025-03-06T20:18:22.485153Z"
           }
        ]
        ```

5.  **Mining Process Observation:**

    * You can observe the transaction mining process in Terminal 2, where the API is running.
    * Finally, if a mining attempt does not success within allowed transaction time it will be listed as ```"status": "timeout"```. 