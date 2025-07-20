# Mastering Eventual Consistency utilizing Transaction Log Tailing with MongoDB

Struggling without Database Transactions in Microservices?

Distributed event-driven architectures runs smoothly - until they don't. When data ends up incomplete, how do you detect it and clean it up?  This article walks through a file store example of handling those edge cases while ensuring reliable event distribution using MongoDB Change Streams.

Read my [article](https://medium.com/@kinneko-de/22cf9bf9e712), where I explain the concept in detail.

## How to start

1. Execute the script `./run-sut.sh` to start the MongoDB replica set. 
2. Wait until you see `mongodb-init exited with code 0` in the terminal.
3. Your MongoDB replica set is now ready and you can connect to it.