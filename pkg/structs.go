package pkg

type CommitObject struct {
	Treesha        string
	Parentsha      string
	Authorname     string
	Authoremail    string
	CommitterName  string
	CommitterEmail string
	Time           int64
	Timezone       string
	Message        string
}

type TreeEntry struct {
	FileMode string
	FileType string
	Sha      string
	Path     string
}

type Metadata struct {
	// Ctime time.Time i removed this because in git reset, needed to update index according to the info available on the commit tree.
	// Mtime time.Time
	Mode uint32
	Size int64
}

type PairSHAandStatus struct {
	P      PairOfSHA
	Status string
}

type Pair struct {
	Exists bool
	Sha    string
}

type PairOfSHA struct {
	P1 Pair
	P2 Pair
}

type LogContents struct {
	// branchname  string
	Parentsha   string
	Currentsha  string
	Authorname  string
	Authoremail string
	Timestamp   int64
	Gmt         string
	Operation   string
	Message     string
}
