Distributed Consensus
Implement a version of the RAFT consensus protocol to allow fault tolerant storage of key/value data.
Basic

Read through the RAFT white paper to understand the high level concepts for distributed consensus.

https://raft.github.io/raft.pdf

Answer the following questions:

What is a term?
What makes RAFT different from Paxos?
Can you mutate state from any node in the cluster for RAFT?
Can you mutate state from any node in the cluster for Paxos?
What is CAP theorem?
What does each letter stand for?
Which elements (if any) of CAP does RAFT aim to satisfy?
Does RAFT perform heartbeats?
Why does/doesn’t it?
How does/would it?
Medium
Create a ‘Node’ data structure which:
Stores persistent and volatile state as described in the white paper
Current term
Current index
Log
… (see white paper for a full enumeration)
Implements a REST API
/api/v1/append-entries
/api/v1/request-vote

The state, behavior and return values of the APIs are all documented in the white paper (on page 4), use this for reference.

Write unit tests to ensure the correctness of your REST API and state manipulation, your future self will thank you later.
Advanced
Use the white paper and other online resources to complete your RAFT implementation, you should be able to:

Have a CLI using cobra which allows you to startup nodes
Create a cluster of three or more nodes (ignore node addition/removal)
Running as separate processes
Connecting via REST APIs
Perform ‘Set’ operations to set a key/value pair (via a REST API)
Your implementation should have logging which clearly indicates the process taken to accept/commit a write.
Perform a ‘Get’ to read consistent data
You should only be able to receive values which have been committed (e.g. those where a majority of nodes have accepted/committed the write)
Kill a follower and still perform reads/writes
Kill the leader and observe an election taking place
Followed shortly by successful reads/writes
