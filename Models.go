package main

type Node struct {
	PartitionKey         int
	ReplicationAddresses []string // Also Replication Factor by getting length
	Store                map[string]string
	CommitLog            []string // For when a server is down
}
