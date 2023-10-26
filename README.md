### Algorithm
I chose a kind of hashcash algorithm with sha256 because of this pros:
- simplicity of understanding and implementation
- simple and fast validation on server side
- complexity and target value of the challenge can be tuned fast.
It also has some cons:
- tuning of difficulty of this algorithm is not as precise as, for example, Bitcoin algorithm
- it may be too difficult for weak computers or too easy for powerful ones to solve.

### Blockchain
I implemented blockchain generation mechanism, it may be overkill for this task,
but my aim is to get acquainted with this technology and to see how it works.

### POW mechanism description
- client sends challenge request message to server
- server replies to client with last block from blockchain. Block contains information about target value and it's leading quantity in hash.
- client, by brute force, basing on new block generated by client, has to calculate hash that will suit target and difficulty.
- server checks new block and it's hash, if block is validated, than it is added to blockchain, and server sends to client a quote from online API resource

### Protocol
As a protocol for client-server communication i use json object with "header" and "payload" fields.
- by "header" server and client can indicate type of request and proceed with required logic
- "payload", if not empty, contains json with data.




