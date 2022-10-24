# Answers to Basic questions

## *What is a term?*
A term is an election cycle

## *What makes RAFT different from Paxos?*
| Raft | Paxos|
|-|-|
Only leader can submit proposals | Any node can submit proposals
Nodes with the latest logs are eligible to become leaders | No such restrictions to become a proposer
Checks log continuity before confirming a log | No checks, allows void logs
Keeps track of the commit index | No log connectivity, updates may require additional messages


### *Can you mutate state from any node in the cluster for RAFT?*
No only the leader node can dictate changes

### *Can you mutate state from any node in the cluster for Paxos?*
Yes, can result in livelocks.

## *What is CAP theorem?*
A CA database delivers consistency and availability, but it can’t deliver fault tolerance if any two nodes in the system have a partition between them. Clearly, this is where CAP theorem and NoSQL databases collide: there are no NoSQL databases you would classify as CA under the CAP theorem. In a distributed database, there is no way to avoid system partitions. So, although CAP theorem stating a CA distributed database is possible exists, there is currently no true CA distributed database system. The modern goal of CAP theorem analysis should be for system designers to generate optimal combinations of consistency and availability for particular applications.

### *What does each letter stand for?*
* Consistency
* Availability
* Partition Tolerance

### *Which elements (if any) of CAP does RAFT aim to satisfy?*
Consistency & Partition Tolerance

## *Does RAFT perform heartbeats?*
Yes

### *Why does/doesn’t it?*
To check if a server is still up

### *How does/would it?*
The leader sends empty AppendLog RPC's its followers
