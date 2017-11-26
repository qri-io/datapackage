package dsfs

// PackageFile specifies the different types of files that are
// stored in a package
type PackageFile int

const (
	// PackageFileUnknown is the default package file, which
	// should be erroneous, as there is no sensible default
	// for PackageFile
	PackageFileUnknown PackageFile = iota
	// PackageFileDataset is the maind dataset.json file
	// that contains all dataset metadata, and is the only
	// required file to constitute a dataset
	PackageFileDataset
	// PackageFileStructure isolates this dataset's structure
	// in it's own file
	PackageFileStructure
	// PackageFileAbstractStructure is the abstract verion of
	// structure
	PackageFileAbstractStructure
	// PackageFileResources lists the resource datasets
	// that went into creating a dataset
	// TODO - I think this can be removed now that Query exists
	PackageFileResources
	// PackageFileCommitMsg is isolates the user-entered
	// documentation of the changes to this dataset's history
	PackageFileCommitMsg
	// PackageFileQuery isloates the concrete query that
	// generated this dataset
	PackageFileQuery
	// PackageFileAbstractQuery is the abstract version of
	// the operation performed to create this dataset
	PackageFileAbstractQuery
)

// filenames maps PackageFile to their filename counterparts
var filenames = map[PackageFile]string{
	PackageFileUnknown:           "",
	PackageFileDataset:           "dataset.json",
	PackageFileStructure:         "structure.json",
	PackageFileAbstractStructure: "abstract_structure.json",
	PackageFileAbstractQuery:     "abstract_query.json",
	PackageFileResources:         "resources",
	PackageFileCommitMsg:         "commit.json",
	PackageFileQuery:             "query.json",
}

// String implements the io.Stringer interface for PackageFile
func (p PackageFile) String() string {
	return p.Filename()
}

// Filename gives the canonical filename for a PackageFile
func (p PackageFile) Filename() string {
	return filenames[p]
}
